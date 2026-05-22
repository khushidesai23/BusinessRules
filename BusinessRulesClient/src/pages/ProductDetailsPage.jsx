import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ProductAPI, CategoryAPI, AssignmentAPI } from '../api';
import { Save, ArrowLeft, AlertCircle, Loader2 } from 'lucide-react';

export default function ProductDetailsPage() {
    const { id } = useParams();
    const isNew = id === 'new';
    const navigate = useNavigate();

    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    const [productData, setProductData] = useState([]);
    const [productName, setProductName] = useState('');
    const [categoryId, setCategoryId] = useState('');
    const [categories, setCategories] = useState([]);

    useEffect(() => {
        if (isNew) {
            fetchCategories();
        } else {
            fetchProductDetails();
        }
    }, [id]);

    const fetchCategories = async () => {
        const res = await CategoryAPI.getAll();
        if (res.message === 'success') {
            setCategories(res.data || []);
        }
    };

    const fetchProductDetails = async () => {
        setLoading(true);
        const res = await ProductAPI.getById(id);
        if (res.message === 'success' && res.data && res.data.length > 0) {
            setCategoryId(res.data[0].categoryId);
            // Constructing form data state
            const initialData = res.data.map(item => ({
                attributeId: item.attributeId,
                attributeName: item.attributeName,
                dataType: item.dataType,
                value: item.data !== undefined ? item.data : ''
            }));
            // extract product name if attributeId == 0
            const nameField = res.data.find(d => d.attributeId === 0);
            if (nameField && nameField.data) setProductName(nameField.data);
            setProductData(initialData);
            setError('');
        } else {
            setError(res.message);
        }
        setLoading(false);
    };

    const handleCategorySelect = async (e) => {
        const selectedCatId = parseInt(e.target.value, 10);
        setCategoryId(selectedCatId);
        if (!selectedCatId) {
            setProductData([]);
            return;
        }

        setLoading(true);
        const res = await AssignmentAPI.getCategoryWiseCommonAttributes([selectedCatId]);
        if (res.message === 'success') {
            const assignedAttrs = (res.data || []).filter(attr => attr.assigned);
            setProductData(assignedAttrs.map(attr => ({
                attributeId: attr.id,
                attributeName: attr.name,
                dataType: attr.dataType,
                value: ''
            })));
            setError('');
        } else {
            setError(res.message);
        }
        setLoading(false);
    };

    const handleValueChange = (attributeId, newValue) => {
        setProductData(prev => prev.map(item =>
            item.attributeId === attributeId ? { ...item, value: newValue } : item
        ));
    };

    const handleSave = async (e) => {
        e.preventDefault();
        if (!categoryId) {
            setError('Category is required.');
            return;
        }

        setSaving(true);
        const apiData = [
            // include product name as attributeId 0
            { attributeId: 0, value: productName },
            ...productData.map(item => ({
                attributeId: item.attributeId,
                value: item.value.toString()
            }))
        ];

        const productIdStr = isNew ? "" : id;
        const res = await ProductAPI.upsert(categoryId, productIdStr, apiData);

        if (res.message === 'success') {
            navigate('/product');
        } else {
            setError(res.message);
        }
        setSaving(false);
    };

    return (
        <div className="page-container">
            <div className="flex justify-between items-center" style={{ marginBottom: '2rem' }}>
                <div className="flex items-center gap-4">
                    <button className="btn btn-secondary" onClick={() => navigate('/product')}>
                        <ArrowLeft size={18} />
                    </button>
                    <h1>{isNew ? 'Create Product' : 'Edit Product'}</h1>
                </div>
                <button
                    className="btn btn-primary"
                    onClick={handleSave}
                    disabled={saving || loading || !categoryId || !productName.trim()}
                >
                    {saving ? <Loader2 size={18} className="animate-spin" /> : <Save size={18} />}
                    {saving ? 'Saving...' : 'Save Product'}
                </button>
            </div>

            {error && (
                <div className="alert alert-error">
                    <AlertCircle size={20} />
                    <span>{error}</span>
                </div>
            )}

            <div className="card w-full" style={{ maxWidth: '1000px', margin: '0 auto' }}>
                <div style={{ borderBottom: '1px solid #e5e7eb', paddingBottom: '1.5rem', marginBottom: '2rem' }}>
                    <h2 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '1rem' }}>Product Information</h2>
                    
                    {/* Category Section */}
                    <div style={{ display: 'grid', gridTemplateColumns: isNew ? '1fr' : 'auto', gap: '1rem', marginBottom: '2rem' }}>
                        {isNew && (
                            <div className="form-group">
                                <label className="form-label" style={{ fontWeight: '600' }}>Select Category *</label>
                                <select
                                    className="form-input"
                                    value={categoryId}
                                    onChange={handleCategorySelect}
                                    style={{ borderColor: !categoryId ? '#fca5a5' : undefined }}
                                >
                                    <option value="">-- Choose Category --</option>
                                    {categories.map(cat => (
                                        <option key={cat.id} value={cat.id}>{cat.name}</option>
                                    ))}
                                </select>
                                {!categoryId && <span style={{ fontSize: '0.875rem', color: '#ef4444' }}>Category is required</span>}
                            </div>
                        )}
                        {!isNew && categoryId && (
                            <div style={{ padding: '0.75rem', backgroundColor: '#f3f4f6', borderRadius: '0.375rem', border: '1px solid #d1d5db' }}>
                                <p style={{ fontSize: '0.875rem', color: '#6b7280', marginBottom: '0.25rem' }}>Category ID</p>
                                <p style={{ fontSize: '1rem', fontWeight: '600' }}>{categoryId}</p>
                            </div>
                        )}
                    </div>

                    {/* Product Name Field */}
                    <div className="form-group">
                        <label className="form-label" style={{ fontWeight: '600' }}>Product Name *</label>
                        <input
                            type="text"
                            className="form-input"
                            placeholder="Enter product name"
                            value={productName}
                            onChange={(e) => setProductName(e.target.value)}
                            style={{ borderColor: !productName.trim() && productName !== '' ? '#fca5a5' : undefined }}
                        />
                        {!productName.trim() && <span style={{ fontSize: '0.875rem', color: '#ef4444' }}>Product name is required</span>}
                    </div>
                </div>

                {/* Attributes Section */}
                {categoryId && (
                    <div>
                        <h2 style={{ fontSize: '1.25rem', fontWeight: '600', marginBottom: '1.5rem' }}>Product Attributes</h2>
                        {loading ? (
                            <div className="flex items-center justify-center" style={{ padding: '3rem' }}>
                                <Loader2 size={32} className="animate-spin text-primary" color="#6366f1" />
                            </div>
                        ) : productData.length === 0 ? (
                            <div style={{ padding: '2rem', textAlign: 'center', backgroundColor: '#f9fafb', borderRadius: '0.375rem', border: '1px dashed #d1d5db' }}>
                                <p style={{ color: '#6b7280' }}>No attributes assigned to this category</p>
                            </div>
                        ) : (
                            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '1.5rem' }}>
                                {productData.map((field) => (
                                    <div key={field.attributeId} className="form-group" style={{ marginBottom: 0 }}>
                                        <label className="form-label" style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontWeight: '600' }}>
                                            {field.attributeName}
                                            <span className="badge badge-purple" style={{ fontSize: '0.65rem', padding: '0.25rem 0.5rem' }}>{field.dataType}</span>
                                        </label>
                                        {field.dataType === 'boolean' ? (
                                            <select
                                                className="form-input"
                                                value={field.value}
                                                onChange={(e) => handleValueChange(field.attributeId, e.target.value)}
                                            >
                                                <option value="">-- Select --</option>
                                                <option value="true">True</option>
                                                <option value="false">False</option>
                                            </select>
                                        ) : field.dataType === 'integer' || field.dataType === 'float' ? (
                                            <input
                                                type="number"
                                                step={field.dataType === 'float' ? 'any' : '1'}
                                                className="form-input"
                                                placeholder={`Enter ${field.dataType}`}
                                                value={field.value}
                                                onChange={(e) => handleValueChange(field.attributeId, e.target.value)}
                                            />
                                        ) : (
                                            <input
                                                type="text"
                                                className="form-input"
                                                placeholder={`Enter ${field.dataType}`}
                                                value={field.value}
                                                onChange={(e) => handleValueChange(field.attributeId, e.target.value)}
                                            />
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}

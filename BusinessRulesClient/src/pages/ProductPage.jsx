import React, { useState, useEffect } from 'react';
import { ProductAPI } from '../api';
import { Plus, Box, AlertCircle } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export default function ProductPage() {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        fetchProducts();
    }, []);

    const fetchProducts = async () => {
        setLoading(true);
        const res = await ProductAPI.getAll();
        if (res.message === 'success') {
            setProducts(res.data || []);
            setError('');
        } else {
            setError(res.message);
        }
        setLoading(false);
    };

    return (
        <div className="page-container">
            <div className="flex justify-between items-center" style={{ marginBottom: '2rem' }}>
                <h1>Products</h1>
                <button
                    className="btn btn-primary"
                    onClick={() => navigate('/product/new')}
                >
                    <Plus size={18} />
                    Create Product
                </button>
            </div>

            {error && (
                <div className="alert alert-error">
                    <AlertCircle size={20} />
                    <span>{error}</span>
                </div>
            )}

            <div className="card w-full">
                {loading ? (
                    <div className="flex items-center justify-center" style={{ padding: '3rem' }}>
                        <div style={{ textAlign: 'center' }}>
                            <div style={{ marginBottom: '1rem' }}>Loading products...</div>
                        </div>
                    </div>
                ) : products.length === 0 ? (
                    <div style={{ padding: '3rem', textAlign: 'center', backgroundColor: '#f9fafb', borderRadius: '0.375rem' }}>
                        <Box size={32} style={{ margin: '0 auto 1rem', opacity: 0.4 }} />
                        <p style={{ color: '#6b7280', fontSize: '1rem', marginBottom: '1rem' }}>No products found</p>
                        <button
                            className="btn btn-primary"
                            onClick={() => navigate('/product/new')}
                        >
                            <Plus size={16} />
                            Create Your First Product
                        </button>
                    </div>
                ) : (
                    <div style={{ overflowX: 'auto' }}>
                        <table className="data-table" style={{ width: '100%', borderCollapse: 'collapse' }}>
                            <thead>
                                <tr style={{ borderBottom: '2px solid #e5e7eb', backgroundColor: '#f9fafb' }}>
                                    <th style={{ padding: '1rem', textAlign: 'left', fontWeight: '600', color: '#374151' }}>Product Name</th>
                                    <th style={{ padding: '1rem', textAlign: 'left', fontWeight: '600', color: '#374151' }}>Category</th>
                                    <th style={{ padding: '1rem', textAlign: 'left', fontWeight: '600', color: '#374151' }}>Product ID</th>
                                    <th style={{ padding: '1rem', textAlign: 'center', fontWeight: '600', color: '#374151' }}>Action</th>
                                </tr>
                            </thead>
                            <tbody>
                                {products.map((prod, index) => (
                                    <tr 
                                        key={prod.id} 
                                        style={{ 
                                            borderBottom: '1px solid #e5e7eb',
                                            backgroundColor: index % 2 === 0 ? '#ffffff' : '#f9fafb',
                                            transition: 'background-color 0.2s',
                                            cursor: 'pointer'
                                        }}
                                        onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#f3f4f6'}
                                        onMouseLeave={(e) => e.currentTarget.style.backgroundColor = index % 2 === 0 ? '#ffffff' : '#f9fafb'
                                        }
                                    >
                                        <td style={{ padding: '1rem' }}>
                                            <div className="flex items-center gap-2">
                                                <Box size={18} color="#6366f1" style={{ opacity: 0.7 }} />
                                                <span style={{ fontWeight: '500', color: '#1f2937' }}>{prod.name || '(Unnamed)'}</span>
                                            </div>
                                        </td>
                                        <td style={{ padding: '1rem' }}>
                                            <span className="badge badge-purple" style={{ backgroundColor: '#eee5ff', color: '#7c3aed', padding: '0.375rem 0.75rem', borderRadius: '0.25rem', fontSize: '0.875rem' }}>
                                                {prod.categoryPath || prod.category || 'N/A'}
                                            </span>
                                        </td>
                                        <td style={{ padding: '1rem', fontSize: '0.875rem', color: '#6b7280', fontFamily: 'monospace' }}>
                                            {prod.id.substring(0, 8)}...
                                        </td>
                                        <td style={{ padding: '1rem', textAlign: 'center' }}>
                                            <button 
                                                className="btn btn-secondary" 
                                                onClick={(e) => { 
                                                    e.stopPropagation(); 
                                                    navigate(`/product/${prod.id}`); 
                                                }}
                                                style={{ padding: '0.5rem 1rem', fontSize: '0.875rem' }}
                                            >
                                                Edit
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
}

import React, { useState, useEffect } from 'react';
import { AttributeAPI } from '../api';
import { Plus, Tag, AlertCircle } from 'lucide-react';

export default function AttributePage() {
    const [attributes, setAttributes] = useState([]);
    const [newAttrName, setNewAttrName] = useState('');
    const [newAttrType, setNewAttrType] = useState('string');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [isAdding, setIsAdding] = useState(false);

    const dataTypes = ['string', 'float', 'integer', 'boolean'];

    useEffect(() => {
        fetchAttributes();
    }, []);

    const fetchAttributes = async () => {
        setLoading(true);
        const res = await AttributeAPI.getAll();
        if (res.message === 'success') {
            setAttributes(res.data || []);
            setError('');
        } else {
            setError(res.message);
        }
        setLoading(false);
    };

    const handleCreateAttribute = async (e) => {
        e.preventDefault();
        if (!newAttrName.trim()) return;

        setIsAdding(true);
        const res = await AttributeAPI.create(newAttrName, newAttrType);
        if (res.message === 'success') {
            setNewAttrName('');
            setNewAttrType('string');
            fetchAttributes();
        } else {
            setError(res.message);
        }
        setIsAdding(false);
    };

    return (
        <div className="page-container">
            <div className="flex justify-between items-center" style={{ marginBottom: '2rem' }}>
                <h1>Attributes</h1>
            </div>

            {error && (
                <div className="alert alert-error">
                    <AlertCircle size={20} />
                    <span>{error}</span>
                </div>
            )}

            <div className="flex gap-4">
                <div className="card flex-1">
                    <h3>Attribute List</h3>
                    {loading ? (
                        <p className="badge badge-neutral">Loading attributes...</p>
                    ) : (
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Name</th>
                                    <th>Data Type</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {attributes.length === 0 ? (
                                    <tr>
                                        <td colSpan="4" style={{ textAlign: 'center', padding: '2rem' }}>
                                            No attributes found. Create your first one!
                                        </td>
                                    </tr>
                                ) : (
                                    attributes.map(attr => (
                                        <tr key={attr.id}>
                                            <td>{attr.id}</td>
                                            <td>
                                                <div className="flex items-center gap-2">
                                                    <Tag size={16} color="#8b949e" />
                                                    {attr.name}
                                                </div>
                                            </td>
                                            <td>
                                                <span className="badge badge-purple">{attr.dataType}</span>
                                            </td>
                                            <td>
                                                <button className="btn btn-error btn-sm" onClick={async () => {
                                                    if (!window.confirm('Delete this attribute?')) return;
                                                    const res = await AttributeAPI.delete(attr.id);
                                                    if (res.message === 'success') fetchAttributes();
                                                    else setError(res.message);
                                                }}>Delete</button>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    )}
                </div>

                <div className="card" style={{ width: '300px', height: 'fit-content' }}>
                    <h3>Add Attribute</h3>
                    <form onSubmit={handleCreateAttribute}>
                        <div className="form-group">
                            <label className="form-label">Attribute Name</label>
                            <input
                                type="text"
                                className="form-input"
                                placeholder="e.g. Price"
                                value={newAttrName}
                                onChange={(e) => setNewAttrName(e.target.value)}
                                disabled={isAdding}
                            />
                        </div>
                        <div className="form-group">
                            <label className="form-label">Data Type</label>
                            <select
                                className="form-input"
                                value={newAttrType}
                                onChange={(e) => setNewAttrType(e.target.value)}
                                disabled={isAdding}
                            >
                                {dataTypes.map(type => (
                                    <option key={type} value={type}>{type}</option>
                                ))}
                            </select>
                        </div>
                        <button
                            type="submit"
                            className="btn btn-primary w-full"
                            disabled={isAdding || !newAttrName.trim()}
                        >
                            <Plus size={18} />
                            {isAdding ? 'Adding...' : 'Add Attribute'}
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
}

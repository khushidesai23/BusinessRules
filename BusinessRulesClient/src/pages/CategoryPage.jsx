import React, { useState, useEffect } from 'react';
import { CategoryAPI } from '../api';
import { Plus, AlertCircle } from 'lucide-react';

export default function CategoryPage() {
    const [categories, setCategories] = useState([]);
    const [newCategoryName, setNewCategoryName] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [isAdding, setIsAdding] = useState(false);

    useEffect(() => {
        fetchCategories();
    }, []);

    const fetchCategories = async () => {
        setLoading(true);
        const res = await CategoryAPI.getAll();
        if (res.message === 'success') {
            setCategories(res.data || []);
            setError('');
        } else {
            setError(res.message);
        }
        setLoading(false);
    };

    const handleCreateCategory = async (e) => {
        e.preventDefault();
        if (!newCategoryName.trim()) return;

        setIsAdding(true);
        const res = await CategoryAPI.create(newCategoryName);
        if (res.message === 'success') {
            setNewCategoryName('');
            fetchCategories();
        } else {
            setError(res.message);
        }
        setIsAdding(false);
    };

    const handleDeleteCategory = async (id) => {
        if (!window.confirm('Are you sure you want to delete this category?')) return;
        setLoading(true);
        const res = await CategoryAPI.delete(id);
        if (res.message === 'success') {
            fetchCategories();
        } else {
            setError(res.message);
        }
        setLoading(false);
    };

    return (
        <div className="page-container">
            <div className="flex justify-between items-center" style={{ marginBottom: '2rem' }}>
                <h1>Categories</h1>
            </div>

            {error && (
                <div className="alert alert-error">
                    <AlertCircle size={20} />
                    <span>{error}</span>
                </div>
            )}

            <div className="flex gap-4">
                <div className="card flex-1">
                    <h3>Category List</h3>
                    {loading ? (
                        <p className="badge badge-neutral">Loading categories...</p>
                    ) : (
                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Name</th>
                                    <th>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {categories.length === 0 ? (
                                    <tr>
                                        <td colSpan="2" style={{ textAlign: 'center', padding: '2rem' }}>
                                            No categories found. Create your first one!
                                        </td>
                                    </tr>
                                ) : (
                                    categories.map(cat => (
                                        <tr key={cat.id}>
                                            <td>{cat.id}</td>
                                            <td>{cat.name}</td>
                                            <td>
                                                <button className="btn btn-error btn-sm" onClick={() => handleDeleteCategory(cat.id)}>Delete</button>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    )}
                </div>

                <div className="card" style={{ width: '300px', height: 'fit-content' }}>
                    <h3>Add Category</h3>
                    <form onSubmit={handleCreateCategory}>
                        <div className="form-group">
                            <label className="form-label">Category Name</label>
                            <input
                                type="text"
                                className="form-input"
                                placeholder="e.g. Electronics"
                                value={newCategoryName}
                                onChange={(e) => setNewCategoryName(e.target.value)}
                                disabled={isAdding}
                            />
                        </div>
                        <button
                            type="submit"
                            className="btn btn-primary w-full"
                            disabled={isAdding || !newCategoryName.trim()}
                        >
                            <Plus size={18} />
                            {isAdding ? 'Adding...' : 'Add Category'}
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
}

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

    const handleDeleteProduct = async (id) => {
        if (!window.confirm('Are you sure you want to delete this product?')) return;
        setLoading(true);
        const res = await ProductAPI.delete(id);
        if (res.message === 'success') fetchProducts();
        else setError(res.message);
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
                <h3>Product List</h3>
                {loading ? (
                    <p className="badge badge-neutral">Loading products...</p>
                ) : (
                    <table className="data-table">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>Name</th>
                                <th>Category</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            {products.length === 0 ? (
                                <tr>
                                    <td colSpan="4" style={{ textAlign: 'center', padding: '2rem' }}>
                                        No products found. Add a new one!
                                    </td>
                                </tr>
                            ) : (
                                products.map(prod => (
                                    <tr key={prod.id} className="cursor-pointer" onClick={() => navigate(`/product/${prod.id}`)}>
                                        <td>{prod.id}</td>
                                        <td>
                                            <div className="flex items-center gap-2">
                                                <Box size={16} color="#8b949e" />
                                                {prod.name}
                                            </div>
                                        </td>
                                        <td>
                                            <span className="badge badge-purple">{prod.categoryPath}</span>
                                        </td>
                                        <td>
                                            <button className="btn btn-secondary" onClick={(e) => { e.stopPropagation(); navigate(`/product/${prod.id}`); }}>Edit</button>
                                            <button className="btn btn-error ml-2" onClick={(e) => { e.stopPropagation(); handleDeleteProduct(prod.id); }}>Delete</button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
}

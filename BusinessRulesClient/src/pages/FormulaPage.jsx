import React, { useState, useEffect } from 'react';
import { FormulaAPI } from '../api';
import { Plus, Calculator, AlertCircle } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export default function FormulaPage() {
    const [formulas, setFormulas] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        fetchFormulas();
    }, []);

    const fetchFormulas = async () => {
        setLoading(true);
        const res = await FormulaAPI.getAll();
        if (res.message === 'success') {
            setFormulas(res.data || []);
            setError('');
        } else {
            setError(res.message);
        }
        setLoading(false);
    };

    const handleDeleteFormula = async (categoryId, targetAttributeId) => {
        if (!window.confirm('Delete this formula?')) return;
        setLoading(true);
        const res = await FormulaAPI.delete(categoryId, targetAttributeId);
        if (res.message === 'success') fetchFormulas();
        else setError(res.message);
        setLoading(false);
    };

    return (
        <div className="page-container">
            <div className="flex justify-between items-center" style={{ marginBottom: '2rem' }}>
                <h1>Formulas</h1>
                <button
                    className="btn btn-primary"
                    onClick={() => navigate('/formula/new')}
                >
                    <Plus size={18} />
                    Create Formula
                </button>
            </div>

            {error && (
                <div className="alert alert-error">
                    <AlertCircle size={20} />
                    <span>{error}</span>
                </div>
            )}

            <div className="card w-full">
                <h3>Formula List</h3>
                {loading ? (
                    <p className="badge badge-neutral">Loading formulas...</p>
                ) : (
                    <table className="data-table">
                        <thead>
                            <tr>
                                <th>Category</th>
                                <th>Target Attribute</th>
                                <th>Formula</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            {formulas.length === 0 ? (
                                <tr>
                                    <td colSpan="4" style={{ textAlign: 'center', padding: '2rem' }}>
                                        No formulas found. Create one!
                                    </td>
                                </tr>
                            ) : (
                                formulas.map((form, idx) => (
                                    <tr key={idx}>
                                        <td>
                                            <span className="badge badge-purple">{form.categoryName}</span>
                                        </td>
                                        <td>{form.targetAttributeName}</td>
                                        <td>
                                            <code style={{ background: '#0d1117', padding: '0.2rem 0.5rem', borderRadius: '0.25rem' }}>
                                                {form.formula}
                                            </code>
                                        </td>
                                        <td>
                                            <button className="btn btn-secondary" onClick={() => navigate(`/formula/${form.categoryId}`)}>
                                                Edit
                                            </button>
                                            <button className="btn btn-error ml-2" onClick={() => handleDeleteFormula(form.categoryId, form.targetAttributeID || form.targetAttributeId)}>
                                                Delete
                                            </button>
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

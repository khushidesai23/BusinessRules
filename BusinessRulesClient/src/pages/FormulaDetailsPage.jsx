import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { FormulaAPI, CategoryAPI, AssignmentAPI } from '../api';
import { Save, ArrowLeft, AlertCircle, Loader2 } from 'lucide-react';

export default function FormulaDetailsPage() {
    const { id } = useParams();
    const isNew = id === 'new';
    const navigate = useNavigate();

    const [categories, setCategories] = useState([]);
    const [selectedCategories, setSelectedCategories] = useState([]);

    const [attributes, setAttributes] = useState([]);
    const [targetAttributeId, setTargetAttributeId] = useState('');

    const [formula, setFormula] = useState('');

    const [loadingCats, setLoadingCats] = useState(false);
    const [loadingAttrs, setLoadingAttrs] = useState(false);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    // Autocomplete state
    const [suggestions, setSuggestions] = useState([]);
    const [showSuggestions, setShowSuggestions] = useState(false);
    const [cursorPosition, setCursorPosition] = useState(0);
    const textareaRef = useRef(null);

    useEffect(() => {
        fetchCategories();
        // If not new, we should probably fetch the specific formula details
        // But the API only has get-all. So we would fetch all and filter?
        // For now, let's just make sure we handle creation and basic edit.
    }, []);

    const fetchCategories = async () => {
        setLoadingCats(true);
        const res = await CategoryAPI.getAll();
        if (res.message === 'success') {
            setCategories(res.data || []);
        }
        setLoadingCats(false);
    };

    useEffect(() => {
        if (selectedCategories.length > 0) {
            fetchAttributes();
        } else {
            setAttributes([]);
            setTargetAttributeId('');
        }
    }, [selectedCategories]);

    const fetchAttributes = async () => {
        setLoadingAttrs(true);
        const res = await AssignmentAPI.getCategoryWiseCommonAttributes(selectedCategories);
        if (res.message === 'success') {
            // Only show attributes that are assigned to ALL selected categories
            const assignedAttrs = (res.data || []).filter(attr => attr.assigned);
            setAttributes(assignedAttrs);
            if (!assignedAttrs.find(a => a.id.toString() === targetAttributeId?.toString())) {
                setTargetAttributeId('');
            }
        }
        setLoadingAttrs(false);
    };

    const handleCategoryToggle = (catId) => {
        setSelectedCategories(prev =>
            prev.includes(catId) ? prev.filter(c => c !== catId) : [...prev, catId]
        );
    };

    // Autocomplete Logic
    const handleFormulaChange = (e) => {
        const value = e.target.value;
        setFormula(value);

        const pos = e.target.selectionStart;
        setCursorPosition(pos);

        // Find the word currently being typed
        const textBeforeCursor = value.substring(0, pos);
        const words = textBeforeCursor.split(/[\s+\-*/()=<>!]+/);
        const currentWord = words[words.length - 1];

        if (currentWord.length > 0) {
            const matchAttrs = attributes
                .filter(a => a.name.toLowerCase().includes(currentWord.toLowerCase()))
                .map(a => a.name);

            const functions = ['IF'];
            const matchFuncs = functions.filter(f => f.toLowerCase().includes(currentWord.toLowerCase()));

            const allMatches = [...matchFuncs, ...matchAttrs];
            if (allMatches.length > 0) {
                setSuggestions(allMatches);
                setShowSuggestions(true);
            } else {
                setShowSuggestions(false);
            }
        } else {
            setShowSuggestions(false);
        }
    };

    const insertSuggestion = (suggestion) => {
        const textBeforeCursor = formula.substring(0, cursorPosition);
        const words = textBeforeCursor.split(/[\s+\-*/()=<>!]+/);
        const currentWord = words[words.length - 1];

        const newFormula =
            formula.substring(0, cursorPosition - currentWord.length) +
            suggestion +
            (suggestion === 'IF' ? '(' : ' ') +
            formula.substring(textareaRef.current.selectionEnd);

        setFormula(newFormula);
        setShowSuggestions(false);
        textareaRef.current.focus();
    };

    const handleSave = async () => {
        if (selectedCategories.length === 0) {
            setError('Please select at least one category.');
            return;
        }
        if (!targetAttributeId) {
            setError('Please select a target attribute.');
            return;
        }
        if (!formula.trim()) {
            setError('Please enter a formula.');
            return;
        }

        setSaving(true);
        let hasError = false;
        for (const catId of selectedCategories) {
            const res = await FormulaAPI.create(catId, parseInt(targetAttributeId, 10), formula);
            if (res.message !== 'success') {
                setError(`Failed for category ${catId}: ` + res.message);
                hasError = true;
                break;
            }
        }

        if (!hasError) {
            navigate('/formula');
        }
        setSaving(false);
    };

    const handleDelete = async () => {
        if (isNew) return;
        if (selectedCategories.length === 0 || !targetAttributeId) {
            setError('Select categories and target attribute to delete formula');
            return;
        }
        if (!window.confirm('Delete this formula for selected categories?')) return;
        setSaving(true);
        let hasError = false;
        for (const catId of selectedCategories) {
            const res = await FormulaAPI.delete(catId, parseInt(targetAttributeId, 10));
            if (res.message !== 'success') {
                setError(`Failed for category ${catId}: ` + res.message);
                hasError = true;
                break;
            }
        }
        if (!hasError) navigate('/formula');
        setSaving(false);
    };

    return (
        <div className="page-container">
            <div className="flex justify-between items-center" style={{ marginBottom: '2rem' }}>
                <div className="flex items-center gap-4">
                    <button className="btn btn-secondary" onClick={() => navigate('/formula')}>
                        <ArrowLeft size={18} />
                    </button>
                    <h1>{isNew ? 'Create Formula' : 'Edit Formula'}</h1>
                </div>
                <div style={{ display: 'flex', gap: '8px' }}>
                    {!isNew && (
                        <button className="btn btn-error" onClick={handleDelete} disabled={saving}>
                            Delete
                        </button>
                    )}
                    <button
                        className="btn btn-primary"
                        onClick={handleSave}
                        disabled={saving || selectedCategories.length === 0}
                    >
                        {saving ? <Loader2 size={18} className="animate-spin" /> : <Save size={18} />}
                        {saving ? 'Saving...' : 'Save Formula'}
                    </button>
                </div>
            </div>

            {error && (
                <div className="alert alert-error">
                    <AlertCircle size={20} />
                    <span>{error}</span>
                </div>
            )}

            <div className="flex gap-4">
                {/* Categories Selection */}
                <div className="card w-full" style={{ maxWidth: '300px' }}>
                    <h3>Select Categories</h3>
                    {loadingCats ? (
                        <Loader2 size={18} className="animate-spin" />
                    ) : (
                        <div className="flex flex-col gap-2 mt-4" style={{ maxHeight: '400px', overflowY: 'auto' }}>
                            {categories.map(cat => (
                                <label key={cat.id} className="checkbox-container">
                                    <input
                                        type="checkbox"
                                        checked={selectedCategories.includes(cat.id)}
                                        onChange={() => handleCategoryToggle(cat.id)}
                                    />
                                    <span>{cat.name}</span>
                                </label>
                            ))}
                        </div>
                    )}
                </div>

                {/* Formula Editor */}
                <div className="card flex-1">
                    <h3>Formula Configuration</h3>

                    <div className="form-group mt-4">
                        <label className="form-label">Target Attribute</label>
                        <select
                            className="form-input"
                            value={targetAttributeId}
                            onChange={(e) => setTargetAttributeId(e.target.value)}
                            disabled={loadingAttrs || selectedCategories.length === 0}
                        >
                            <option value="">-- Choose Target Attribute --</option>
                            {attributes.map(attr => (
                                <option key={attr.id} value={attr.id}>{attr.name} ({attr.dataType})</option>
                            ))}
                        </select>
                    </div>

                    <div className="form-group" style={{ position: 'relative' }}>
                        <label className="form-label">Formula Expression</label>
                        <textarea
                            ref={textareaRef}
                            className="form-input"
                            rows={5}
                            placeholder='e.g. IF(Attribute1="hello", "hi", "bye") or Attribute1 + Attribute2'
                            value={formula}
                            onChange={handleFormulaChange}
                            onClick={handleFormulaChange}
                            onKeyUp={handleFormulaChange}
                            style={{ fontFamily: 'monospace' }}
                        />

                        {showSuggestions && suggestions.length > 0 && (
                            <div className="autocomplete-list">
                                {suggestions.map((s, idx) => (
                                    <div
                                        key={idx}
                                        className="autocomplete-item"
                                        onClick={() => insertSuggestion(s)}
                                    >
                                        {s}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    <div className="badge badge-neutral" style={{ display: 'inline-block', marginTop: '1rem' }}>
                        Hint: Start typing an attribute name or 'IF' to see autocomplete suggestions.
                    </div>
                </div>
            </div>
        </div>
    );
}

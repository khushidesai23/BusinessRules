import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  NavLink,
  Navigate,
} from "react-router-dom";
import {
  LayoutDashboard,
  Tags,
  Link as LinkIcon,
  Box,
  Calculator,
} from "lucide-react";

import CategoryPage from "./pages/CategoryPage";
import AttributePage from "./pages/AttributePage";
import AssignmentPage from "./pages/AssignmentPage";
import ProductPage from "./pages/ProductPage";
import ProductDetailsPage from "./pages/ProductDetailsPage";
import FormulaPage from "./pages/FormulaPage";
import FormulaDetailsPage from "./pages/FormulaDetailsPage";

function App() {
  return (
    <Router>
      <div className="layout">
        <aside className="sidebar">
          <h2>Business Rules</h2>
          <nav className="flex flex-col gap-2 mt-4">
            <NavLink
              to="/category"
              className={({ isActive }) =>
                `nav-item ${isActive ? "active" : ""}`
              }
            >
              <LayoutDashboard size={20} />
              Categories
            </NavLink>
            <NavLink
              to="/attribute"
              className={({ isActive }) =>
                `nav-item ${isActive ? "active" : ""}`
              }
            >
              <Tags size={20} />
              Attributes
            </NavLink>
            <NavLink
              to="/assignment"
              className={({ isActive }) =>
                `nav-item ${isActive ? "active" : ""}`
              }
            >
              <LinkIcon size={20} />
              Assignments
            </NavLink>
            <NavLink
              to="/formula"
              className={({ isActive }) =>
                `nav-item ${isActive ? "active" : ""}`
              }
            >
              <Calculator size={20} />
              Formulas
            </NavLink>
            <NavLink
              to="/product"
              className={({ isActive }) =>
                `nav-item ${isActive ? "active" : ""}`
              }
            >
              <Box size={20} />
              Products
            </NavLink>
          </nav>
        </aside>

        <main className="main-content">
          <Routes>
            <Route path="/" element={<Navigate to="/category" replace />} />
            <Route path="/category" element={<CategoryPage />} />
            <Route path="/attribute" element={<AttributePage />} />
            <Route path="/assignment" element={<AssignmentPage />} />
            <Route path="/product" element={<ProductPage />} />
            <Route path="/product/:id" element={<ProductDetailsPage />} />
            <Route path="/formula" element={<FormulaPage />} />
            <Route path="/formula/:id" element={<FormulaDetailsPage />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;

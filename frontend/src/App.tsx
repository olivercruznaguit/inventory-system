import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import './App.css'
import Login from './pages/Login'
import { ProtectedRoute } from './components/ProtectedRoutes'
import Dashboard from './pages/Dashboard'
import AppLayout from './layouts/AppLayout'
import Products from './pages/Products'
import Categories from './pages/Categories'

function App() {
    return (
        <BrowserRouter>
            <Routes>
                {/* Public */}
                <Route path="/login" element={<Login />} />

                {/* Protected application */}
                <Route
                    element={
                        <ProtectedRoute>
                            <AppLayout />
                        </ProtectedRoute>
                    }
                >
                    <Route path="/dashboard" element={<Dashboard />} />
                    <Route path="/products" element={<Products />} />
                    <Route path="/categories" element={<Categories />} />
                </Route>

                {/* Default */}
                <Route
                    path="/"
                    element={<Navigate to="/dashboard" replace />}
                />
            </Routes>
        </BrowserRouter>
    )
}

export default App

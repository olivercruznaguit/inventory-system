import type { ReactNode } from "react"
import { useAuth } from "../hooks/useAuth"
import { Navigate } from "react-router-dom"

type ProtectedRouteProps = {
    children: ReactNode
}

export function ProtectedRoute({children}: ProtectedRouteProps){
    const { user, isLoading } = useAuth()

    if (isLoading) {
        return <div>Loading...</div>
    }

    if (!user) {
        return <Navigate to="/login" replace />
    }

    return children
}
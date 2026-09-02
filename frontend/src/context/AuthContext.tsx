import { useState, type ReactNode } from "react"
import type { User } from "../types/user"
import { login as LoginApi } from "../services/api"
import type { JwtPayload } from "../types/auth"
import { jwtDecode } from "jwt-decode"
import { AuthContext } from "./AuthContext"

type AuthProviderProps = {
    children: ReactNode
}

function getInitialAuth() {
    const storedToken = localStorage.getItem("token")

    if (!storedToken) {
        return {
            user: null,
            token: null,
        }
    }

    try {
        const payload = jwtDecode<JwtPayload>(storedToken)

        return {
            token: storedToken,
            user: {
                id: payload.user_id,
                email: payload.email,
                role: payload.role,
            },
        }
    } catch {
        localStorage.removeItem("token")

        return {
            user: null,
            token: null,
        }
    }
}

export function AuthProvider({ children }: AuthProviderProps) {
    const initialAuth = getInitialAuth()
    
    const [user, setUser] = useState<User | null>(initialAuth.user)
    const [token, setToken] = useState<string | null>(initialAuth.token)

    async function login(email: string, password: string) {
        const response = await LoginApi(email, password)

        const payload = jwtDecode<JwtPayload>(response.token)

        localStorage.setItem("token", response.token)

        setToken(response.token)

        setUser({
            id: payload.user_id,
            email: payload.email,
            role: payload.role,
        })
    }

    function logout() {
        localStorage.removeItem("token")
        setUser(null)   
        setToken(null)
    }

    return (
        <AuthContext.Provider
            value={{
                user,
                token,
                login,
                logout,
            }}
        >
            {children}
        </AuthContext.Provider>
    )
}
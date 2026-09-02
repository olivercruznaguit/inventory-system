import type { UserRole } from "./user"

export type AuthResponse = {
    token: string
}

export type LoginRequest = {
    email: string
    password: string
}

export type JwtPayload = {
    user_id: number
    email: string
    role: UserRole
    exp: number
    iat: number
    iss: string
}
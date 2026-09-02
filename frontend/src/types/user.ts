export type UserRole = "USER" | "ADMIN"

export type User = {
    id: number
    email: string
    role: UserRole
}
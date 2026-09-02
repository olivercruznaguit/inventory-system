import { useState, type SyntheticEvent } from "react"
import { useAuth } from "../hooks/useAuth"
import { useNavigate } from "react-router-dom"

export default function Login() {
    const navigate = useNavigate()
    const { login } = useAuth()
    
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")

    async function handleSubmit(e: SyntheticEvent<HTMLFormElement>) {
        e.preventDefault()

        try {
            await login(email, password)
            navigate("/dashboard")
        } catch (error) {
            console.error(error)
        }
    }

    return (
        <form onSubmit={handleSubmit}>
            <div>
                <label htmlFor="email">Email</label>
                <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                />
            </div>

            <div>
                <label htmlFor="password">Password</label>
                <input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
            </div>

            <button type="submit">
                Login
            </button>
        </form>
    )
}
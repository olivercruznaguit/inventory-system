import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { AuthProvider } from './context/AuthContext.tsx'
import { SnackbarProvider } from 'notistack'
import { ThemeProvider } from "@mui/material/styles"
import theme from "./theme/theme"

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <SnackbarProvider maxSnack={3}>
        <AuthProvider>
            <ThemeProvider theme={theme}>
                <App />
            </ThemeProvider>
        </AuthProvider>
    </SnackbarProvider>
  </StrictMode>,
)

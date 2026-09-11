import { Outlet } from "react-router-dom"
import Navbar from "../components/Navbar"
import Sidebar from "../components/Sidebar"
import { Box, Toolbar } from "@mui/material"
import { useState } from "react"

export default function AppLayout() {
    const [mobileOpen, setMobileOpen] = useState(false)

    function handleMenuClick() {
        setMobileOpen(true)
    }

    function handleDrawerClose() {
        setMobileOpen(false)
    }

    return (
       <Box>
            <Navbar onMenuClick={handleMenuClick} />

            <Box sx={{ display: "flex" }}>
                <Sidebar 
                open={mobileOpen}
                onClose={handleDrawerClose} />

                <Box
                    component="main"
                    sx={{ flexGrow: 1, p: 3 }}
                >
                    <Toolbar />
                    <Outlet />
                </Box>
            </Box>
        </Box>
    )
}
import { Outlet } from "react-router-dom"
import Navbar from "../components/Navbar"
import Sidebar from "../components/Sidebar"
import { Box, Toolbar } from "@mui/material"

export default function AppLayout() {
    return (
       <Box>
            <Navbar />

            <Box sx={{ display: "flex" }}>
                <Sidebar />

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
import { AppBar, Box, Button, Toolbar, Typography } from "@mui/material"
import { useAuth } from "../hooks/useAuth"


export default function Navbar(){
    const {user, logout} = useAuth()
    return(
        <AppBar position="fixed">
            <Toolbar>
                <Typography
                    variant="h6"
                    component="div"
                    sx={{ 
                        flexGrow: 1
                    }}
                >
                    Inventory System
                </Typography>

                <Box>
                    <Typography
                        component="span"
                        sx={{ mr: 2 }}
                    >
                        {user?.email}
                    </Typography>

                    <Button
                        color="inherit"
                        onClick={logout}
                    >
                        Logout
                    </Button>
                </Box>
            </Toolbar>
        </AppBar>
    )
}
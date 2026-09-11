import { AppBar, Box, Button, IconButton, Toolbar, Typography } from "@mui/material"
import { useAuth } from "../hooks/useAuth"
import MenuIcon from '@mui/icons-material/Menu';

type NavbarProps = {
    onMenuClick: () => void
}

export default function Navbar({ onMenuClick }: NavbarProps){
    const {user, logout} = useAuth()
    return(
        <AppBar
        position="fixed"
        elevation={0}
        sx={{
            zIndex: (theme) => theme.zIndex.drawer + 1,
            backgroundColor: "background.paper",
            color: "text.primary",
            borderBottom: "1px solid",
            borderColor: "divider",
        }}
        >
            <Toolbar>
                <IconButton
                    color="inherit"
                    onClick={onMenuClick}
                    sx={{
                        mr: 1,
                        display: { xs: "inline-flex", md: "none" },
                    }}
                    aria-label="open navigation"
                >
                    <MenuIcon />
                </IconButton>
                
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
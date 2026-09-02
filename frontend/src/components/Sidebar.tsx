import { Drawer, List, ListItemButton, ListItemText, Toolbar } from "@mui/material";
import { Link } from "react-router-dom";


const drawerWidth = 220
export default function Sidebar() {
    return (
        <Drawer
        variant="permanent"
            sx={{
                width: drawerWidth,
                flexShrink: 0,
                "& .MuiDrawer-paper": {
                    width: drawerWidth,
                    boxSizing: "border-box",
                },
            }}
        >
            <Toolbar />
            
            <List>
                <ListItemButton component={Link} to="/dashboard">
                    <ListItemText primary="Dashboard" />
                </ListItemButton>

                <ListItemButton component={Link} to="/products">
                    <ListItemText primary="Products" />
                </ListItemButton>

                <ListItemButton component={Link} to="/categories">
                    <ListItemText primary="Categories" />
                </ListItemButton>
            </List>
        </Drawer>
    )
}
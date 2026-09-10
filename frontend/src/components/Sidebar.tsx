import { Drawer, List, ListItemButton, ListItemText, Toolbar } from "@mui/material";
import { NavLink } from "react-router-dom";

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
                    backgroundColor: "background.paper",
                    borderRight: "1px solid",
                    borderColor: "divider",
                },
            }}
        >
            <Toolbar />
            
            <List>
                <NavLink
                    to="/dashboard"
                    end
                    style={{ textDecoration: "none", color: "inherit" }}
                >
                    {({ isActive }) => (
                        <ListItemButton selected={isActive}>
                            <ListItemText primary="Dashboard" />
                        </ListItemButton>
                    )}
                </NavLink>

                <NavLink
                    to="/products"
                    end
                    style={{ textDecoration: "none", color: "inherit" }}
                >
                    {({ isActive }) => (
                        <ListItemButton selected={isActive}>
                            <ListItemText primary="Products" />
                        </ListItemButton>
                    )}
                </NavLink>

                <NavLink
                    to="/categories"
                    end
                    style={{ textDecoration: "none", color: "inherit" }}
                >
                    {({ isActive }) => (
                        <ListItemButton selected={isActive}>
                            <ListItemText primary="Categories" />
                        </ListItemButton>
                    )}
                </NavLink>
            </List>
        </Drawer>
    )
}
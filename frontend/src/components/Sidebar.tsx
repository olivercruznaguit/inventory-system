import { Drawer, List, ListItemButton, ListItemText, Toolbar, useMediaQuery, useTheme } from "@mui/material";
import { NavLink } from "react-router-dom";

const drawerWidth = 220

type SidebarProps = {
    open: boolean
    onClose: () => void
}

export default function Sidebar({ open, onClose }: SidebarProps) {
    const theme = useTheme()
    const isMobile = useMediaQuery(theme.breakpoints.down("md"))

    const drawer = (
        <>
            <Toolbar />

            <List>
                <NavLink
                    to="/dashboard"
                    end
                    style={{ textDecoration: "none", color: "inherit" }}
                    onClick={isMobile ? onClose : undefined}
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
                    onClick={isMobile ? onClose : undefined}
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
                    onClick={isMobile ? onClose : undefined}
                >
                    {({ isActive }) => (
                        <ListItemButton selected={isActive}>
                            <ListItemText primary="Categories" />
                        </ListItemButton>
                    )}
                </NavLink>
            </List>
        </>
    )
    return (
        <Drawer
            variant={isMobile ? "temporary" : "permanent"}
            open={isMobile ? open : true}
            onClose={onClose}
            ModalProps={{
                keepMounted: true,
            }}
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
            {drawer}
        </Drawer>
    )
}
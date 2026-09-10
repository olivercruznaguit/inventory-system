import { createTheme } from "@mui/material/styles"

const theme = createTheme({
    palette: {
        primary: {
            main: "#D97706",
            dark: "#B45309",
        },

        secondary: {
            main: "#475569",
        },

        background: {
            default: "#F8FAFC",
            paper: "#FFFFFF",
        },

        text: {
            primary: "#1E293B",
            secondary: "#64748B",
        },

        divider: "#E2E8F0",

        success: {
            main: "#4D7C5A",
        },

        warning: {
            main: "#D97706",
        },

        error: {
            main: "#B94A48",
        },

        info: {
            main: "#4F6F8F",
        },
    },

    components: {
        MuiButton: {
            styleOverrides: {
                root: {
                    borderRadius: 8,
                    textTransform: "none",
                    fontWeight: 600,
                },
            },

            variants: [
                {
                    props: { variant: "contained" },
                    style: {
                        "&:hover": {
                            backgroundColor: "#B45309",
                        },
                    },
                },
            ],
        },

        MuiTextField: {
            styleOverrides: {
                root: {
                    "& .MuiOutlinedInput-root": {
                        borderRadius: 8,

                        "&:hover .MuiOutlinedInput-notchedOutline": {
                            borderColor: "#94A3B8",
                        },

                        "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
                            borderColor: "#D97706",
                            borderWidth: 2,
                        },
                    },
                    
                    "& input[type=number]": {
                        MozAppearance: "textfield",

                        "&::-webkit-outer-spin-button, &::-webkit-inner-spin-button": {
                            display: "none",
                        },
                    },
                },
            },
        },

        MuiDialog: {
            styleOverrides: {
                paper: {
                    borderRadius: 12,
                },
            },
        },

        MuiDialogTitle: {
            styleOverrides: {
                root: {
                    fontWeight: 700,
                    paddingBottom: 8,
                },
            },
        },

        MuiDialogActions: {
            styleOverrides: {
                root: {
                    padding: "16px 24px",
                },
            },
        },

        MuiChip: {
            styleOverrides: {
                root: {
                    borderRadius: 6,
                    fontWeight: 600,
                },
            },
        },

        MuiTableCell: {
            styleOverrides: {
                head: {
                    fontWeight: 700,
                    color: "#475569",
                    backgroundColor: "#F8FAFC",
                },
                root: {
                    borderColor: "#E2E8F0",
                },
            },
        },

        MuiTableRow: {
            styleOverrides: {
                root: {
                    "&:hover": {
                        backgroundColor: "#F8FAFC",
                    },
                },
            },
        },

        MuiCard: {
            styleOverrides: {
                root: {
                    borderRadius: 10,
                    boxShadow: "0 1px 3px rgba(15, 23, 42, 0.08)",
                },
            },
        },

        MuiAlert: {
            styleOverrides: {
                root: {
                    borderRadius: 8,
                },
                message: {
                    fontWeight: 500,
                },
            },
        },

        MuiListItemButton: {
            styleOverrides: {
                root: {
                    borderRadius: 8,
                    margin: "2px 8px",

                    "&:hover": {
                        backgroundColor: "rgba(217, 119, 6, 0.08)",
                    },

                    "&.Mui-selected": {
                        backgroundColor: "rgba(217, 119, 6, 0.12)",
                        color: "#B45309",

                        "&:hover": {
                            backgroundColor: "rgba(217, 119, 6, 0.16)",
                        },
                    },
                },
            },
        },
    },
})

export default theme
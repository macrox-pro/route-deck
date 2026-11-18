import type { ReactNode } from "react";
import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { SettingsService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export type Theme = "light" | "dark" | "system";

const THEMES: ReadonlySet<string> = new Set(["light", "dark", "system"]);
const STORAGE_KEY = "ui-theme";
const DARK_MODE_CLASS = "dark-mode";

function isTheme(value: string | null | undefined): value is Theme {
    return value !== null && value !== undefined && THEMES.has(value);
}

function resolveTheme(theme: Theme): "light" | "dark" {
    if (theme === "system") {
        return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
    }
    return theme;
}

function applyDocumentTheme(theme: Theme) {
    const root = window.document.documentElement;
    const resolved = resolveTheme(theme);
    root.classList.toggle(DARK_MODE_CLASS, resolved === "dark");
    root.style.colorScheme = resolved;
    localStorage.setItem(STORAGE_KEY, theme);
}

interface ThemeContextType {
    theme: Theme;
    setTheme: (theme: Theme) => void;
    /** Persist via SettingsService and recreate native chrome. */
    setAppearance: (theme: Theme) => Promise<void>;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const useTheme = (): ThemeContextType => {
    const context = useContext(ThemeContext);

    if (context === undefined) {
        throw new Error("useTheme must be used within a ThemeProvider");
    }

    return context;
};

interface ThemeProviderProps {
    children: ReactNode;
}

export const ThemeProvider = ({ children }: ThemeProviderProps) => {
    const [theme, setThemeState] = useState<Theme>(() => {
        if (typeof window === "undefined") {
            return "system";
        }
        const cached = localStorage.getItem(STORAGE_KEY);
        return isTheme(cached) ? cached : "system";
    });

    const setTheme = useCallback((next: Theme) => {
        setThemeState(next);
    }, []);

    const setAppearance = useCallback(async (next: Theme) => {
        const hash = typeof window !== "undefined" ? window.location.hash : "";
        await SettingsService.SetAppearance(next, hash);
        setThemeState(next);
        applyDocumentTheme(next);
    }, []);

    useEffect(() => {
        let cancelled = false;
        void (async () => {
            try {
                const settings = await SettingsService.Get();
                const fromStore = settings.appearance;
                if (!cancelled && isTheme(fromStore)) {
                    setThemeState(fromStore);
                    applyDocumentTheme(fromStore);
                }
            } catch {
                // Keep FOUC cache / default until Settings is available.
            }
        })();
        return () => {
            cancelled = true;
        };
    }, []);

    useEffect(() => {
        applyDocumentTheme(theme);

        const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
        const handleChange = () => {
            if (theme === "system") {
                applyDocumentTheme(theme);
            }
        };

        mediaQuery.addEventListener("change", handleChange);
        return () => mediaQuery.removeEventListener("change", handleChange);
    }, [theme]);

    return (
        <ThemeContext.Provider value={{ theme, setTheme, setAppearance }}>{children}</ThemeContext.Provider>
    );
};

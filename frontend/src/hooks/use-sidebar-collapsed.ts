import { useCallback, useState } from "react";

const STORAGE_KEY = "routedeck.sidebar.collapsed";

function readCollapsed(): boolean {
    try {
        return localStorage.getItem(STORAGE_KEY) === "true";
    } catch {
        return false;
    }
}

function writeCollapsed(collapsed: boolean): void {
    try {
        localStorage.setItem(STORAGE_KEY, collapsed ? "true" : "false");
    } catch {
        // WebView storage may be unavailable; keep in-memory state only.
    }
}

/** Collapsed sidebar preference. Uses localStorage (not cookies) for Wails UI chrome. */
export function useSidebarCollapsed() {
    const [collapsed, setCollapsedState] = useState(readCollapsed);

    const setCollapsed = useCallback((next: boolean) => {
        setCollapsedState(next);
        writeCollapsed(next);
    }, []);

    const toggle = useCallback(() => {
        setCollapsedState((prev) => {
            const next = !prev;
            writeCollapsed(next);
            return next;
        });
    }, []);

    return { collapsed, setCollapsed, toggle };
}

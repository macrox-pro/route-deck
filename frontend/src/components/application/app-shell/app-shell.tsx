import { useEffect, useState } from "react";
import {
    Activity,
    Grid01,
    HomeLine,
    Route,
    Settings01,
    Shield01,
} from "@untitledui/icons";
import { Outlet, useLocation } from "react-router-dom";
import { SettingsService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import { SidebarNavigationSimple } from "@/components/application/app-navigation/sidebar-navigation/sidebar-simple";
import type { NavItemType } from "@/components/application/app-navigation/config";
import { APP_VERSION_FALLBACK } from "@/config/app";
import { useSidebarCollapsed } from "@/hooks/use-sidebar-collapsed";

const navItems: NavItemType[] = [
    { label: "Dashboard", href: "/", icon: HomeLine },
    { label: "Applications", href: "/applications", icon: Grid01 },
    { label: "Proxies", href: "/proxies", icon: Route },
    { label: "Profiles", href: "/profiles", icon: Shield01 },
    { label: "Activity", href: "/activity", icon: Activity },
    { label: "Settings", href: "/settings", icon: Settings01 },
];

export function AppShell() {
    const location = useLocation();
    const activeUrl = location.pathname === "" ? "/" : location.pathname;
    const { collapsed, toggle } = useSidebarCollapsed();
    const [versionLabel, setVersionLabel] = useState(APP_VERSION_FALLBACK);

    useEffect(() => {
        let cancelled = false;
        void SettingsService.Get()
            .then((settings) => {
                if (!cancelled && settings.appVersionLabel) {
                    setVersionLabel(settings.appVersionLabel);
                }
            })
            .catch(() => {
                /* keep fallback */
            });
        return () => {
            cancelled = true;
        };
    }, []);

    return (
        <div className="flex h-full min-w-0 overflow-hidden overscroll-none bg-primary text-primary">
            <SidebarNavigationSimple
                activeUrl={activeUrl}
                items={navItems}
                showAccountCard={false}
                collapsed={collapsed}
                onToggleCollapsed={toggle}
                versionLabel={versionLabel}
            />
            <main className="flex min-h-0 min-w-0 flex-1 flex-col overflow-x-hidden overflow-y-auto overscroll-none p-6 md:p-8">
                <Outlet />
            </main>
        </div>
    );
}

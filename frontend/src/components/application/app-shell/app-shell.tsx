import {
    Activity,
    Grid01,
    HomeLine,
    Route,
    Settings01,
    Shield01,
} from "@untitledui/icons";
import { Outlet, useLocation } from "react-router-dom";
import { SidebarNavigationSimple } from "@/components/application/app-navigation/sidebar-navigation/sidebar-simple";
import type { NavItemType } from "@/components/application/app-navigation/config";
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

    return (
        <div className="flex h-full min-w-0 overflow-hidden overscroll-none bg-primary text-primary">
            <SidebarNavigationSimple
                activeUrl={activeUrl}
                items={navItems}
                showAccountCard={false}
                collapsed={collapsed}
                onToggleCollapsed={toggle}
            />
            <main className="min-w-0 flex-1 overflow-x-hidden overflow-y-auto overscroll-none p-6 md:p-8">
                <Outlet />
            </main>
        </div>
    );
}

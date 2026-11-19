import { useEffect, useMemo, useState } from "react";
import {
    Activity,
    Grid01,
    HelpCircle,
    HomeLine,
    Route,
    Settings01,
    Shield01,
} from "@untitledui/icons";
import { Events } from "@wailsio/runtime";
import { Outlet, useLocation, useNavigate } from "react-router-dom";
import {
    SettingsService,
    UtilityWindowsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { EventEnvelope } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails/models";
import { SidebarNavigationSimple } from "@/components/application/app-navigation/sidebar-navigation/sidebar-simple";
import type { NavItemType } from "@/components/application/app-navigation/config";
import { APP_VERSION_FALLBACK } from "@/config/app";
import { useSidebarCollapsed } from "@/hooks/use-sidebar-collapsed";

export function AppShell() {
    const location = useLocation();
    const navigate = useNavigate();
    const activeUrl = location.pathname === "" ? "/" : location.pathname;
    const { collapsed, toggle } = useSidebarCollapsed();
    const [versionLabel, setVersionLabel] = useState(APP_VERSION_FALLBACK);

    const navItems: NavItemType[] = useMemo(
        () => [
            { label: "Dashboard", href: "/", icon: HomeLine },
            { label: "Applications", href: "/applications", icon: Grid01 },
            { label: "Proxies", href: "/proxies", icon: Route },
            { label: "Profiles", href: "/profiles", icon: Shield01 },
            { label: "Activity", href: "/activity", icon: Activity },
            { label: "Settings", href: "/settings", icon: Settings01 },
        ],
        [],
    );

    const footerItems: NavItemType[] = useMemo(
        () => [
            {
                label: "Setup Guides",
                icon: HelpCircle,
                onPress: () => {
                    void UtilityWindowsService.OpenGuides("");
                },
            },
        ],
        [],
    );

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

    useEffect(() => {
        return Events.On("main.command", (event: { data?: EventEnvelope | EventEnvelope[] }) => {
            const raw = event?.data;
            const envelope = Array.isArray(raw) ? raw[0] : raw;
            const payload = envelope?.payload ?? {};
            const command = String(payload.command ?? "");
            switch (command) {
                case "navigate": {
                    const route = String(payload.route ?? "/");
                    navigate(route);
                    break;
                }
                case "toggle-sidebar":
                    toggle();
                    break;
                case "apply":
                case "stop":
                case "recheck":
                    navigate("/", { state: { command } });
                    break;
                case "add-application":
                    navigate("/applications", { state: { command } });
                    break;
                default:
                    break;
            }
        });
    }, [navigate, toggle]);

    return (
        <div className="flex h-full min-w-0 overflow-hidden overscroll-none bg-primary text-primary">
            <SidebarNavigationSimple
                activeUrl={activeUrl}
                items={navItems}
                footerItems={footerItems}
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

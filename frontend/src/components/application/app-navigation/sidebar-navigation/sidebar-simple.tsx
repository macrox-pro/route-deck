import type { CSSProperties, ReactNode } from "react";
import { ChevronLeftDouble, ChevronRightDouble } from "@untitledui/icons";
import { Button as AriaButton } from "react-aria-components";
import { Tooltip } from "@/components/base/tooltip/tooltip";
import { cx } from "@/utils/cx";
import { NavAccountCard } from "../base-components/nav-account-card";
import { NavButton } from "../base-components/nav-button";
import { NavItemBase } from "../base-components/nav-item";
import { NavList } from "../base-components/nav-list";
import type { NavItemType } from "../config";

const EXPANDED_WIDTH = 280;
const COLLAPSED_WIDTH = 68;

interface SidebarNavigationProps {
    /** URL of the currently active item. */
    activeUrl?: string;
    /** List of items to display. */
    items: NavItemType[];
    /** List of footer items to display. */
    footerItems?: NavItemType[];
    /** Feature card to display. */
    featureCard?: ReactNode;
    /** Whether to show the account card. */
    showAccountCard?: boolean;
    /** Whether to hide the right side border. */
    hideBorder?: boolean;
    /** Whether the sidebar is collapsed to icons. */
    collapsed?: boolean;
    /** Toggle expanded/collapsed. */
    onToggleCollapsed?: () => void;
    /** Additional CSS classes to apply to the sidebar. */
    className?: string;
}

/** Desktop-first Untitled UI simple sidebar with icon-rail collapse for the Wails shell. */
export const SidebarNavigationSimple = ({
    activeUrl,
    items,
    footerItems = [],
    featureCard,
    showAccountCard = false,
    hideBorder = false,
    collapsed = false,
    onToggleCollapsed,
    className,
}: SidebarNavigationProps) => {
    const width = collapsed ? COLLAPSED_WIDTH : EXPANDED_WIDTH;

    return (
        <aside
            style={
                {
                    "--width": `${width}px`,
                    width,
                } as CSSProperties
            }
            className={cx(
                "flex h-full shrink-0 flex-col overflow-y-auto overscroll-none bg-primary pt-5 transition-[width] duration-200 ease-linear",
                !hideBorder && "border-r border-secondary",
                className,
            )}
        >
            <div className={cx("flex items-center px-3", collapsed ? "justify-center" : "justify-between gap-2 px-5")}>
                {!collapsed && <p className="truncate text-md font-semibold text-primary">MACROX RouteDeck</p>}
                {onToggleCollapsed && (
                    <Tooltip title={collapsed ? "Expand sidebar" : "Collapse sidebar"} placement="right">
                        <AriaButton
                            aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
                            onPress={onToggleCollapsed}
                            className="flex size-9 shrink-0 cursor-pointer items-center justify-center rounded-md text-fg-quaternary outline-focus-ring hover:bg-primary_hover hover:text-fg-quaternary_hover focus-visible:outline-2 focus-visible:outline-offset-2"
                        >
                            {collapsed ? <ChevronRightDouble className="size-5" /> : <ChevronLeftDouble className="size-5" />}
                        </AriaButton>
                    </Tooltip>
                )}
            </div>

            {collapsed ? (
                <ul className="mt-5 flex flex-col gap-0.5 px-3.5">
                    {items.map((item) => (
                        <li key={item.label}>
                            <NavButton current={activeUrl === item.href} href={item.href} label={item.label} icon={item.icon} />
                        </li>
                    ))}
                </ul>
            ) : (
                <NavList activeUrl={activeUrl} items={items} />
            )}

            <div className={cx("mt-auto flex flex-col gap-3 py-5", collapsed ? "items-center px-3" : "px-4")}>
                {footerItems.length > 0 && (
                    <ul className={cx("flex flex-col", collapsed && "gap-0.5")}>
                        {footerItems.map((item) => (
                            <li key={item.label} className={collapsed ? undefined : "py-px"}>
                                {collapsed ? (
                                    <NavButton current={activeUrl === item.href} href={item.href} label={item.label} icon={item.icon} />
                                ) : (
                                    <NavItemBase badge={item.badge} icon={item.icon} href={item.href} type="link" current={item.href === activeUrl}>
                                        {item.label}
                                    </NavItemBase>
                                )}
                            </li>
                        ))}
                    </ul>
                )}

                {!collapsed && featureCard}

                {!collapsed && showAccountCard && <NavAccountCard />}
            </div>
        </aside>
    );
};

export { EXPANDED_WIDTH, COLLAPSED_WIDTH };

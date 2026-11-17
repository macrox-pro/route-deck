import type { CSSProperties, FC, HTMLAttributes, ReactNode } from "react";
import { ChevronLeftDouble, ChevronRightDouble } from "@untitledui/icons";
import { Button as AriaButton, Link as AriaLink } from "react-aria-components";
import { Badge } from "@/components/base/badges/badges";
import { Tooltip } from "@/components/base/tooltip/tooltip";
import { RouteDeckMark } from "@/components/foundations/logo/route-deck-mark";
import { APP_VERSION } from "@/config/app";
import { cx } from "@/utils/cx";
import { NavAccountCard } from "../base-components/nav-account-card";
import type { NavItemType } from "../config";

const EXPANDED_WIDTH = 280;
const COLLAPSED_WIDTH = 68;

/** Shared with collapsed rail: px-4 + size-9 centers in 68px and stays put while width animates. */
const RAIL_PAD = "px-4";

const toggleButtonClassName =
    "relative z-10 flex size-9 shrink-0 cursor-pointer items-center justify-center rounded-md text-fg-quaternary outline-focus-ring hover:bg-primary_hover hover:text-fg-quaternary_hover focus-visible:outline-2 focus-visible:outline-offset-2";

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

interface SidebarNavLinkProps {
    href?: string;
    label: string;
    icon?: FC<HTMLAttributes<HTMLOrSVGElement>>;
    current?: boolean;
    collapsed: boolean;
}

/** Icon sits in a fixed size-9 column so expand/collapse does not shift it horizontally. */
function SidebarNavLink({ href, label, icon: Icon, current, collapsed }: SidebarNavLinkProps) {
    const link = (
        <AriaLink
            href={href}
            aria-label={label}
            aria-current={current ? "page" : undefined}
            className={cx(
                "group/item relative flex h-9 w-full cursor-pointer items-center rounded-md outline-focus-ring transition duration-100 ease-linear select-none",
                "bg-primary hover:bg-primary_hover focus-visible:z-10 focus-visible:outline-2 focus-visible:outline-offset-2",
                current && "bg-secondary hover:bg-secondary_hover",
            )}
        >
            <span className="flex size-9 shrink-0 items-center justify-center">
                {Icon && (
                    <Icon
                        aria-hidden="true"
                        className={cx(
                            "size-5 text-fg-quaternary transition-inherit-all group-hover/item:text-fg-quaternary_hover",
                            current && "text-fg-quaternary_hover",
                        )}
                    />
                )}
            </span>
            <span
                className={cx(
                    "min-w-0 flex-1 truncate pr-2 text-sm font-semibold text-secondary transition duration-100 ease-linear group-hover/item:text-secondary_hover",
                    current && "text-secondary_hover",
                    collapsed && "hidden",
                )}
            >
                {label}
            </span>
        </AriaLink>
    );

    if (collapsed) {
        return (
            <Tooltip title={label} placement="right">
                {link}
            </Tooltip>
        );
    }

    return link;
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
    const toggleLabel = collapsed ? "Expand sidebar" : "Collapse sidebar";

    return (
        <aside
            style={
                {
                    "--width": `${width}px`,
                    width,
                } as CSSProperties
            }
            className={cx(
                "flex h-full shrink-0 flex-col overflow-x-hidden overflow-y-auto overscroll-none bg-primary pt-5 transition-[width] duration-200 ease-linear",
                !hideBorder && "border-r border-secondary",
                className,
            )}
        >
            {/* Mark always in the same size-9 cell; label/badge hide without moving the mark. */}
            <div className={cx("flex h-9 shrink-0 items-center gap-2", RAIL_PAD)}>
                <div className="flex size-9 shrink-0 items-center justify-center">
                    <RouteDeckMark />
                </div>
                <p className={cx("truncate text-sm font-semibold text-primary", collapsed && "hidden")}>RouteDeck</p>
                <Badge color="gray" type="pill-color" size="sm" className={cx("shrink-0", collapsed && "hidden")}>
                    {APP_VERSION}
                </Badge>
            </div>

            <ul className={cx("mt-5 flex flex-col gap-0.5", RAIL_PAD)}>
                {items.map((item) => (
                    <li key={item.label}>
                        <SidebarNavLink
                            href={item.href}
                            label={item.label}
                            icon={item.icon}
                            current={activeUrl === item.href}
                            collapsed={collapsed}
                        />
                    </li>
                ))}
            </ul>

            <div className="mt-auto flex flex-col">
                <div className={cx("flex flex-col gap-3 pt-5", RAIL_PAD)}>
                    {footerItems.length > 0 && (
                        <ul className="flex flex-col gap-0.5">
                            {footerItems.map((item) => (
                                <li key={item.label}>
                                    <SidebarNavLink
                                        href={item.href}
                                        label={item.label}
                                        icon={item.icon}
                                        current={activeUrl === item.href}
                                        collapsed={collapsed}
                                    />
                                </li>
                            ))}
                        </ul>
                    )}

                    {!collapsed && featureCard}

                    {!collapsed && showAccountCard && <NavAccountCard />}
                </div>

                {onToggleCollapsed && (
                    <div className={cx("flex pt-3 pb-5", RAIL_PAD)}>
                        <Tooltip title={toggleLabel} placement="right">
                            <AriaButton aria-label={toggleLabel} onPress={onToggleCollapsed} className={toggleButtonClassName}>
                                {collapsed ? (
                                    <ChevronRightDouble className="size-5" aria-hidden="true" />
                                ) : (
                                    <ChevronLeftDouble className="size-5" aria-hidden="true" />
                                )}
                            </AriaButton>
                        </Tooltip>
                    </div>
                )}
            </div>
        </aside>
    );
};

export { EXPANDED_WIDTH, COLLAPSED_WIDTH };

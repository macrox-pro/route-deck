import type { CSSProperties, FC, HTMLAttributes, ReactNode } from "react";
import { ChevronLeftDouble, ChevronRightDouble } from "@untitledui/icons";
import { Button as AriaButton } from "react-aria-components";
import { useNavigate } from "react-router-dom";
import { Badge } from "@/components/base/badges/badges";
import { Tooltip } from "@/components/base/tooltip/tooltip";
import { RouteDeckMark } from "@/components/foundations/logo/route-deck-mark";
import { APP_VERSION_FALLBACK } from "@/config/app";
import { cx } from "@/utils/cx";
import { NavAccountCard } from "../base-components/nav-account-card";
import type { NavItemType } from "../config";

const EXPANDED_WIDTH = 232;
const COLLAPSED_WIDTH = 68;

/** Shared with collapsed rail: px-4 + size-9 centers in 68px and stays put while width animates. */
const RAIL_PAD = "px-4";

/** Shared hit target for every sidebar control (nav, footer action, collapse). */
const SIDEBAR_CONTROL =
    "relative flex h-9 cursor-pointer items-center rounded-md border-0 bg-transparent p-0 text-left font-inherit text-fg-quaternary outline-focus-ring transition duration-100 ease-linear select-none hover:bg-primary_hover hover:text-fg-quaternary_hover focus-visible:z-10 focus-visible:outline-2 focus-visible:outline-offset-2";

const toggleButtonClassName = cx(SIDEBAR_CONTROL, "z-10 size-9 shrink-0 justify-center");

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
    /** Short chrome version label ("dev" or "v1.2.3"). */
    versionLabel?: string;
    /** Additional CSS classes to apply to the sidebar. */
    className?: string;
}

interface SidebarNavButtonProps {
    href?: string;
    onPress?: () => void;
    label: string;
    icon?: FC<HTMLAttributes<HTMLOrSVGElement>>;
    current?: boolean;
    collapsed: boolean;
}

const navItemClassName = (current: boolean | undefined) =>
    cx(
        SIDEBAR_CONTROL,
        "group/item w-full",
        current && "bg-primary text-fg-brand-secondary hover:bg-primary_hover hover:text-fg-brand-secondary",
    );

/** Icon sits in a fixed size-9 column so expand/collapse does not shift it horizontally. */
function SidebarNavButton({ href, onPress, label, icon: Icon, current, collapsed }: SidebarNavButtonProps) {
    const navigate = useNavigate();

    const body = (
        <>
            <span className="flex size-9 shrink-0 items-center justify-center">
                {Icon && (
                    <Icon
                        aria-hidden="true"
                        data-icon
                        className={cx(
                            "size-5 text-fg-quaternary transition-inherit-all group-hover/item:text-fg-quaternary_hover",
                            current && "text-fg-brand-secondary group-hover/item:text-fg-brand-secondary",
                        )}
                    />
                )}
            </span>
            <span
                className={cx(
                    "min-w-0 flex-1 truncate pr-2 text-sm font-semibold text-secondary transition duration-100 ease-linear group-hover/item:text-secondary_hover",
                    current && "text-brand-secondary group-hover/item:text-brand-secondary",
                    collapsed && "hidden",
                )}
            >
                {label}
            </span>
        </>
    );

    const item = (
        <AriaButton
            aria-label={label}
            aria-current={current ? "page" : undefined}
            onPress={() => {
                if (onPress) {
                    onPress();
                    return;
                }
                if (href) {
                    navigate(href);
                }
            }}
            className={navItemClassName(current)}
        >
            {body}
        </AriaButton>
    );

    if (collapsed) {
        return (
            <Tooltip title={label} placement="right">
                {item}
            </Tooltip>
        );
    }

    return item;
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
    versionLabel = APP_VERSION_FALLBACK,
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
                "flex h-full shrink-0 flex-col overflow-x-hidden overflow-y-auto overscroll-none bg-sidebar pt-5 transition-[width] duration-200 ease-linear",
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
                    {versionLabel}
                </Badge>
            </div>

            <ul className={cx("mt-5 flex flex-col gap-0.5", RAIL_PAD)}>
                {items.map((item) => (
                    <li key={item.label}>
                        <SidebarNavButton
                            href={item.href}
                            onPress={item.onPress}
                            label={item.label}
                            icon={item.icon}
                            current={Boolean(item.href) && activeUrl === item.href}
                            collapsed={collapsed}
                        />
                    </li>
                ))}
            </ul>

            <div
                className={cx(
                    "mt-auto flex flex-col",
                    footerItems.length > 0 && "border-t border-secondary",
                )}
            >
                <div className={cx("flex flex-col gap-3", footerItems.length > 0 ? "pt-3" : "pt-5", RAIL_PAD)}>
                    {footerItems.length > 0 && (
                        <ul className="flex flex-col gap-0.5">
                            {footerItems.map((item) => (
                                <li key={item.label}>
                                    <SidebarNavButton
                                        href={item.href}
                                        onPress={item.onPress}
                                        label={item.label}
                                        icon={item.icon}
                                        current={Boolean(item.href) && activeUrl === item.href}
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

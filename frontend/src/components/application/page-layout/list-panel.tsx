import type { FC, ReactNode } from "react";
import { Button } from "@/components/base/buttons/button";
import { Dropdown } from "@/components/base/dropdown/dropdown";
import { cx } from "@/utils/cx";

export type DropdownMenuItem = {
    id: string;
    label: string;
    icon?: FC<{ className?: string }>;
    destructive?: boolean;
    onAction: () => void;
    isDisabled?: boolean;
};

export function ListPanel({ children, className }: { children: ReactNode; className?: string }) {
    return <div className={cx("divide-y divide-secondary overflow-hidden rounded-xl bg-primary ring-1 ring-secondary", className)}>{children}</div>;
}

export function ListRow({
    leading,
    title,
    meta,
    status,
    actions,
    menu,
    className,
}: {
    leading?: ReactNode;
    title: ReactNode;
    meta?: ReactNode;
    status?: ReactNode;
    actions?: ReactNode;
    menu?: DropdownMenuItem[];
    className?: string;
}) {
    const menuItems = menu ?? [];
    const normalItems = menuItems.filter((item) => !item.destructive);
    const destructiveItems = menuItems.filter((item) => item.destructive);

    return (
        <div className={cx("flex min-h-14 items-center gap-3 px-4 py-3", className)}>
            {leading ? <div className="flex size-8 shrink-0 items-center justify-center">{leading}</div> : null}
            <div className="min-w-0 flex-1">
                <div className="truncate text-sm font-semibold text-primary">{title}</div>
                {meta ? <div className="mt-0.5 truncate text-sm text-tertiary">{meta}</div> : null}
            </div>
            {status ? <div className="shrink-0">{status}</div> : null}
            {actions ? <div className="flex shrink-0 items-center gap-2">{actions}</div> : null}
            {menuItems.length > 0 ? (
                <Dropdown.Root>
                    <Dropdown.DotsButton />
                    <Dropdown.Popover>
                        <Dropdown.Menu
                            onAction={(key) => {
                                const item = menuItems.find((entry) => entry.id === String(key));
                                item?.onAction();
                            }}
                        >
                            {normalItems.map((item) => (
                                <Dropdown.Item
                                    key={item.id}
                                    id={item.id}
                                    label={item.label}
                                    icon={item.icon}
                                    isDisabled={item.isDisabled}
                                    textValue={item.label}
                                />
                            ))}
                            {normalItems.length > 0 && destructiveItems.length > 0 ? <Dropdown.Separator /> : null}
                            {destructiveItems.map((item) => (
                                <Dropdown.Item
                                    key={item.id}
                                    id={item.id}
                                    label={item.label}
                                    icon={item.icon}
                                    isDisabled={item.isDisabled}
                                    textValue={item.label}
                                    className="**:data-[slot=label]:text-error-primary"
                                />
                            ))}
                        </Dropdown.Menu>
                    </Dropdown.Popover>
                </Dropdown.Root>
            ) : null}
        </div>
    );
}

/** Compact empty row used inside ListPanel (e.g. Dashboard proxies). */
export function ListEmptyRow({ children, action }: { children: ReactNode; action?: ReactNode }) {
    return (
        <div className="flex min-h-14 items-center justify-between gap-3 px-4 py-3">
            <p className="text-sm text-tertiary">{children}</p>
            {action ? (
                <div className="shrink-0">
                    {typeof action === "string" ? (
                        <Button color="secondary" size="sm">
                            {action}
                        </Button>
                    ) : (
                        action
                    )}
                </div>
            ) : null}
        </div>
    );
}

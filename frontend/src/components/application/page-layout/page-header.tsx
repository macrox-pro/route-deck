import type { ReactNode } from "react";
import { cx } from "@/utils/cx";

export function PageHeader({
    title,
    titleAddon,
    description,
    actions,
    className,
}: {
    title: string;
    titleAddon?: ReactNode;
    description?: string;
    actions?: ReactNode;
    className?: string;
}) {
    return (
        <header className={cx("flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between", className)}>
            <div className="min-w-0">
                <div className="flex flex-wrap items-center gap-2">
                    <h1 className="text-display-xs font-semibold text-primary">{title}</h1>
                    {titleAddon}
                </div>
                {description ? <p className="mt-1 max-w-2xl text-sm text-tertiary">{description}</p> : null}
            </div>
            {actions ? <div className="flex shrink-0 flex-wrap items-center gap-2">{actions}</div> : null}
        </header>
    );
}

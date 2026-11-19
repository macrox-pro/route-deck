import type { ReactNode } from "react";
import { Badge } from "@/components/base/badges/badges";
import { cx } from "@/utils/cx";

export function PageSection({
    title,
    count,
    action,
    children,
    className,
}: {
    title: string;
    count?: number;
    action?: ReactNode;
    children: ReactNode;
    className?: string;
}) {
    return (
        <section className={cx("flex flex-col gap-3", className)}>
            <div className="flex items-center justify-between gap-3">
                <div className="flex min-w-0 items-center gap-2">
                    <h2 className="text-sm font-semibold text-secondary">{title}</h2>
                    {typeof count === "number" ? (
                        <Badge color="gray" size="sm" type="modern">
                            {count}
                        </Badge>
                    ) : null}
                </div>
                {action ? <div className="shrink-0">{action}</div> : null}
            </div>
            <div className="min-w-0">{children}</div>
        </section>
    );
}

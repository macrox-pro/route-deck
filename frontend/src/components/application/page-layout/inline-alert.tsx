import type { ReactNode } from "react";
import { AlertCircle, AlertTriangle, InfoCircle } from "@untitledui/icons";
import { Button } from "@/components/base/buttons/button";
import { cx } from "@/utils/cx";

const TONE = {
    error: {
        role: "alert" as const,
        icon: AlertCircle,
        root: "bg-error-primary text-error-primary ring-error_subtle",
        iconClass: "text-fg-error-primary",
    },
    warning: {
        role: "status" as const,
        icon: AlertTriangle,
        root: "bg-warning-primary text-warning-primary ring-warning_subtle",
        iconClass: "text-fg-warning-primary",
    },
    info: {
        role: "status" as const,
        icon: InfoCircle,
        root: "bg-secondary_subtle text-secondary ring-secondary",
        iconClass: "text-fg-quaternary",
    },
};

export function InlineAlert({
    tone,
    title,
    children,
    action,
    className,
}: {
    tone: "error" | "warning" | "info";
    title?: string;
    children: ReactNode;
    action?: { label: string; onPress: () => void };
    className?: string;
}) {
    const styles = TONE[tone];
    const Icon = styles.icon;

    return (
        <div role={styles.role} className={cx("flex gap-3 rounded-lg px-4 py-3 ring-1", styles.root, className)}>
            <Icon aria-hidden="true" data-icon className={cx("mt-0.5 size-5 shrink-0", styles.iconClass)} />
            <div className="min-w-0 flex-1">
                {title ? <p className="text-sm font-semibold">{title}</p> : null}
                <div className={cx("text-sm", title && "mt-0.5")}>{children}</div>
            </div>
            {action ? (
                <div className="shrink-0 self-center">
                    <Button color="secondary" size="sm" onClick={action.onPress}>
                        {action.label}
                    </Button>
                </div>
            ) : null}
        </div>
    );
}

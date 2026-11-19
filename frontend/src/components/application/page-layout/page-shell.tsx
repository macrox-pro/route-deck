import type { ReactNode } from "react";
import { cx } from "@/utils/cx";

/** Shared page column: 24px between sections, content capped at max-w-5xl on every screen. */
export function PageShell({ className, children }: { className?: string; children: ReactNode }) {
    return <div className={cx("mx-auto flex w-full max-w-5xl flex-col gap-6", className)}>{children}</div>;
}

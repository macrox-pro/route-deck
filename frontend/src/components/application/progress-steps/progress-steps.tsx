import type { FC } from "react";
import { AlertTriangle, Check, Grid01, Route, Shield01, TerminalSquare } from "@untitledui/icons";
import { FeaturedIcon } from "@/components/foundations/featured-icon/featured-icon";
import { Tooltip } from "@/components/base/tooltip/tooltip";
import { checklistMessage } from "@/components/application/setup-guide/checklist-copy";
import { cx } from "@/utils/cx";

export type ProgressStepId = "bridge" | "proxy" | "application" | "profile";

export type ProgressStepStatus = "done" | "partial" | "todo";

export type ProgressStep = {
    id: ProgressStepId;
    label: string;
    hint?: string;
    status: ProgressStepStatus;
};

const STEP_ICONS: Record<ProgressStepId, FC<{ className?: string }>> = {
    bridge: TerminalSquare,
    proxy: Route,
    application: Grid01,
    profile: Shield01,
};

function currentStepId(steps: ProgressStep[]): ProgressStepId | undefined {
    const current = steps.find((step) => step.status !== "done");
    return current?.id;
}

function StepDisk({
    step,
    isCurrent,
    size,
}: {
    step: ProgressStep;
    isCurrent: boolean;
    size: "xs" | "sm" | "md";
}) {
    const Icon = STEP_ICONS[step.id];
    const iconSize = size === "xs" ? "size-3" : size === "sm" ? "size-4" : "size-5";

    if (size === "xs") {
        const mark =
            step.status === "done" ? (
                <Check aria-hidden className="size-3 text-fg-success-primary" />
            ) : step.status === "partial" ? (
                <AlertTriangle aria-hidden className="size-3 text-fg-warning-primary" />
            ) : (
                <Icon aria-hidden className={cx("size-3", isCurrent ? "text-fg-brand-primary" : "text-fg-quaternary")} />
            );
        const disk = (
            <span
                className={cx(
                    "inline-flex size-5 shrink-0 items-center justify-center rounded-full",
                    step.status === "done" && "bg-success-secondary",
                    step.status === "partial" && "bg-warning-secondary",
                    step.status === "todo" && "bg-secondary",
                    step.status === "todo" && (isCurrent ? "ring-2 ring-brand" : "ring-1 ring-secondary"),
                )}
            >
                {mark}
            </span>
        );
        if (step.status === "partial" && step.hint) {
            return (
                <Tooltip title={step.hint} placement="top">
                    <span className="inline-flex">{disk}</span>
                </Tooltip>
            );
        }
        return disk;
    }

    if (step.status === "done") {
        return (
            <FeaturedIcon icon={Check} color="success" theme="light" size={size === "sm" ? "sm" : "md"} aria-hidden />
        );
    }

    if (step.status === "partial") {
        const inner = (
            <FeaturedIcon icon={AlertTriangle} color="warning" theme="light" size={size === "sm" ? "sm" : "md"} aria-hidden />
        );
        if (step.hint) {
            return (
                <Tooltip title={step.hint} placement="top">
                    <span className="inline-flex">{inner}</span>
                </Tooltip>
            );
        }
        return inner;
    }

    const ring = isCurrent ? "ring-2 ring-brand" : "ring-1 ring-secondary";
    return (
        <span
            className={cx(
                "inline-flex items-center justify-center rounded-full bg-secondary",
                size === "sm" ? "size-8" : "size-10",
                ring,
            )}
        >
            <Icon aria-hidden data-icon className={cx(iconSize, isCurrent ? "text-fg-brand-primary" : "text-fg-quaternary")} />
        </span>
    );
}

export function ProgressSteps({
    orientation,
    steps,
    currentId: currentIdProp,
    onSelect,
    compact = false,
}: {
    orientation: "vertical" | "horizontal";
    steps: ProgressStep[];
    currentId?: ProgressStepId;
    onSelect?: (id: ProgressStepId) => void;
    /** Smaller disks and always-visible labels. For a top strip in a narrow window. */
    compact?: boolean;
}) {
    const resolvedCurrent = currentIdProp ?? currentStepId(steps);
    const subtitle = checklistMessage("checklist.subtitle");

    if (orientation === "horizontal") {
        return (
            <div className={cx("flex flex-col", compact ? "gap-2" : "gap-3")}>
                <ol
                    className={
                        compact
                            ? "grid grid-cols-1 gap-x-4 gap-y-2 min-[640px]:grid-cols-2"
                            : "flex flex-wrap items-center gap-2"
                    }
                    aria-label={checklistMessage("checklist.title")}
                >
                    {steps.map((step) => {
                        const isCurrent = step.id === resolvedCurrent;
                        const clickable = !!onSelect && !compact;
                        const item = (
                            <>
                                <StepDisk step={step} isCurrent={isCurrent} size={compact ? "xs" : "sm"} />
                                <span
                                    className={cx(
                                        compact
                                            ? "min-w-0 text-xs leading-4"
                                            : "sr-only sm:not-sr-only sm:max-w-[8rem] sm:truncate sm:text-xs",
                                        step.status === "partial"
                                            ? "text-warning-primary"
                                            : step.status === "done"
                                              ? "text-success-primary"
                                              : isCurrent
                                                ? "text-brand-secondary"
                                                : "text-tertiary",
                                    )}
                                >
                                    {step.label}
                                </span>
                            </>
                        );
                        return (
                            <li key={step.id} className="flex min-w-0 items-start gap-2">
                                {clickable ? (
                                    <button
                                        type="button"
                                        className="flex items-center gap-2 rounded-md transition duration-100 hover:bg-secondary_hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand"
                                        onClick={() => onSelect?.(step.id)}
                                    >
                                        {item}
                                    </button>
                                ) : (
                                    <span className="flex min-w-0 items-start gap-2">{item}</span>
                                )}
                            </li>
                        );
                    })}
                </ol>
                {compact ? null : <p className="text-xs text-tertiary">{subtitle}</p>}
            </div>
        );
    }

    return (
        <div className="flex flex-col gap-4">
            <ul className="flex flex-col gap-3">
                {steps.map((step) => {
                    const isCurrent = step.id === resolvedCurrent;
                    const Icon = STEP_ICONS[step.id];
                    const color =
                        step.status === "done" ? "success" : step.status === "partial" ? "warning" : isCurrent ? "brand" : "gray";
                    const body = (
                        <div className="flex min-w-0 items-start gap-3">
                            {step.status === "done" || step.status === "partial" ? (
                                <StepDisk step={step} isCurrent={isCurrent} size="md" />
                            ) : (
                                <FeaturedIcon icon={Icon} color={color} theme="light" size="md" aria-hidden />
                            )}
                            <div className="min-w-0 pt-0.5">
                                <p
                                    className={cx(
                                        "text-sm font-medium",
                                        step.status === "partial"
                                            ? "text-warning-primary"
                                            : step.status === "done"
                                              ? "text-success-primary"
                                              : isCurrent
                                                ? "text-brand-secondary"
                                                : "text-secondary",
                                    )}
                                >
                                    {step.label}
                                </p>
                                {step.hint && step.status === "partial" ? (
                                    <p className="mt-0.5 text-xs text-tertiary">{step.hint}</p>
                                ) : null}
                            </div>
                        </div>
                    );
                    return (
                        <li key={step.id}>
                            {onSelect ? (
                                <button
                                    type="button"
                                    className="w-full rounded-lg px-2 py-1 text-left transition duration-100 hover:bg-secondary_hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand"
                                    onClick={() => onSelect(step.id)}
                                >
                                    {body}
                                </button>
                            ) : (
                                body
                            )}
                        </li>
                    );
                })}
            </ul>
            <p className="text-xs text-tertiary">{subtitle}</p>
        </div>
    );
}
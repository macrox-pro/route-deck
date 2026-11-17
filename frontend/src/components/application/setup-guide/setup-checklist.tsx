import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { CheckCircle, Circle } from "@untitledui/icons";
import { SetupGuidesService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { ChecklistDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function SetupChecklist({ refreshKey }: { refreshKey?: number }) {
    const [data, setData] = useState<ChecklistDTO | null>(null);
    const [error, setError] = useState<string | null>(null);

    const reload = useCallback(async () => {
        try {
            setData(await SetupGuidesService.Checklist());
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload, refreshKey]);

    if (error) {
        return (
            <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                {error}
            </div>
        );
    }
    if (!data || data.complete) {
        return null;
    }

    const steps: { done: boolean; label: string; to: string }[] = [
        { done: data.bridgeFound, label: "ProxyBridge found", to: "/help" },
        { done: data.proxyAvailable, label: "A proxy is Available", to: "/help" },
        { done: data.applicationAdded, label: "An application is in the catalog", to: "/applications" },
        { done: data.profileApplied, label: "A profile is applied", to: "/profiles" },
    ];

    return (
        <div className="rounded-xl bg-secondary p-5 ring-1 ring-secondary" role="status">
            <h2 className="text-sm font-semibold text-primary">First-run checklist</h2>
            <p className="mt-1 text-sm text-tertiary">These checks track setup progress only. They do not mean protection is Active.</p>
            <ul className="mt-4 flex flex-col gap-2.5 text-sm">
                {steps.map((step) => (
                    <li key={step.label} className="flex items-center justify-between gap-3">
                        <span className={`flex min-w-0 items-center gap-2.5 ${step.done ? "text-success-primary" : "text-secondary"}`}>
                            {step.done ? (
                                <CheckCircle aria-hidden className="size-5 shrink-0 text-fg-success-primary" />
                            ) : (
                                <Circle aria-hidden className="size-5 shrink-0 text-fg-quaternary" />
                            )}
                            <span>
                                <span className="sr-only">{step.done ? "Done: " : "To do: "}</span>
                                {step.label}
                            </span>
                        </span>
                        <Link className="shrink-0 text-brand-secondary hover:underline" to={step.to}>
                            Open
                        </Link>
                    </li>
                ))}
            </ul>
        </div>
    );
}

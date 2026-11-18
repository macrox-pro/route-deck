import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { AlertTriangle, CheckCircle, Circle } from "@untitledui/icons";
import { Button } from "@/components/base/buttons/button";
import { Tooltip } from "@/components/base/tooltip/tooltip";
import { SetupGuidesService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { ChecklistDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import { checklistMessage, proxyStepHint, proxyStepLabel } from "./checklist-copy";

type ChecklistStepRow = {
    key: string;
    done: boolean;
    partial: boolean;
    label: string;
    hint: string;
    to: string;
};

export function SetupChecklist({ refreshKey }: { refreshKey?: number }) {
    const [data, setData] = useState<ChecklistDTO | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

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
    if (!data || data.complete || data.dismissed) {
        return null;
    }

    const steps: ChecklistStepRow[] = [
        {
            key: "bridge",
            done: data.bridgeFound,
            partial: false,
            label: checklistMessage("checklist.step.bridge"),
            hint: "",
            to: "/help",
        },
        {
            key: "proxy",
            done: data.proxyAvailable,
            partial: data.proxyPartial,
            label: proxyStepLabel(data.proxyStep),
            hint: proxyStepHint(data.proxyStep),
            to: "/help",
        },
        {
            key: "application",
            done: data.applicationAdded,
            partial: false,
            label: checklistMessage("checklist.step.application"),
            hint: "",
            to: "/applications",
        },
        {
            key: "profile",
            done: data.profileApplied,
            partial: false,
            label: checklistMessage("checklist.step.profile"),
            hint: "",
            to: "/profiles",
        },
    ];

    return (
        <div className="rounded-xl bg-secondary p-5 ring-1 ring-secondary" role="status">
            <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                    <h2 className="text-sm font-semibold text-primary">{checklistMessage("checklist.title")}</h2>
                    <p className="mt-1 text-sm text-tertiary">{checklistMessage("checklist.subtitle")}</p>
                </div>
                <Button
                    color="secondary"
                    size="sm"
                    isDisabled={busy}
                    onClick={() => {
                        setBusy(true);
                        setError(null);
                        void SetupGuidesService.DismissChecklist()
                            .then(() => reload())
                            .catch((err) => {
                                setError(err instanceof Error ? err.message : String(err));
                            })
                            .finally(() => setBusy(false));
                    }}
                >
                    {checklistMessage("checklist.dismiss")}
                </Button>
            </div>
            <ul className="mt-4 flex flex-col gap-2.5 text-sm">
                {steps.map((step) => (
                    <li key={step.key} className="flex items-center justify-between gap-3">
                        <span
                            className={`flex min-w-0 items-center gap-2.5 ${
                                step.partial ? "text-warning-primary" : step.done ? "text-success-primary" : "text-secondary"
                            }`}
                        >
                            <StepIcon step={step} />
                            <span>
                                {!step.partial && (
                                    <span className="sr-only">
                                        {step.done ? checklistMessage("checklist.sr.done") : checklistMessage("checklist.sr.todo")}
                                    </span>
                                )}
                                {step.label}
                            </span>
                        </span>
                        <Link className="shrink-0 text-brand-secondary hover:underline" to={step.to}>
                            {checklistMessage("checklist.open")}
                        </Link>
                    </li>
                ))}
            </ul>
        </div>
    );
}

function StepIcon({ step }: { step: ChecklistStepRow }) {
    if (step.partial) {
        return (
            <Tooltip title={step.hint} placement="top">
                <span className="inline-flex shrink-0">
                    <AlertTriangle aria-hidden className="size-5 text-fg-warning-primary" />
                    <span className="sr-only">
                        {checklistMessage("checklist.sr.partial")}
                        {step.label}. {step.hint}
                    </span>
                </span>
            </Tooltip>
        );
    }
    if (step.done) {
        return <CheckCircle aria-hidden className="size-5 shrink-0 text-fg-success-primary" />;
    }
    return <Circle aria-hidden className="size-5 shrink-0 text-fg-quaternary" />;
}

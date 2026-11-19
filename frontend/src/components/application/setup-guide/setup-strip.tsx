import { useCallback, useEffect, useState } from "react";
import { ProgressSteps, type ProgressStep, type ProgressStepId } from "@/components/application/progress-steps/progress-steps";
import { Button } from "@/components/base/buttons/button";
import {
    checklistMessage,
    proxyStepHint,
    proxyStepLabel,
} from "@/components/application/setup-guide/checklist-copy";
import { SetupGuidesService, UtilityWindowsService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { ChecklistDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

function stepsFromChecklist(data: ChecklistDTO): ProgressStep[] {
    return [
        {
            id: "bridge",
            label: checklistMessage("checklist.step.bridge"),
            status: data.bridgeFound ? "done" : "todo",
        },
        {
            id: "proxy",
            label: proxyStepLabel(data.proxyStep),
            hint: proxyStepHint(data.proxyStep) || undefined,
            status: data.proxyAvailable ? (data.proxyPartial ? "partial" : "done") : "todo",
        },
        {
            id: "application",
            label: checklistMessage("checklist.step.application"),
            status: data.applicationAdded ? "done" : "todo",
        },
        {
            id: "profile",
            label: checklistMessage("checklist.step.profile"),
            status: data.profileApplied ? "done" : "todo",
        },
    ];
}

export function SetupStrip({ refreshKey }: { refreshKey?: number }) {
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

    const steps = stepsFromChecklist(data);

    const onDismiss = () => {
        setBusy(true);
        setError(null);
        void SetupGuidesService.DismissChecklist()
            .then(() => reload())
            .catch((err) => {
                setError(err instanceof Error ? err.message : String(err));
            })
            .finally(() => setBusy(false));
    };

    return (
        <div className="rounded-xl bg-secondary p-5 ring-1 ring-secondary" role="status">
            <div className="mb-4 flex flex-wrap items-start justify-between gap-3">
                <div className="min-w-0">
                    <h2 className="text-sm font-semibold text-primary">{checklistMessage("checklist.title")}</h2>
                </div>
                <Button color="secondary" size="sm" isDisabled={busy} onClick={onDismiss}>
                    {checklistMessage("checklist.dismiss")}
                </Button>
            </div>
            <ProgressSteps
                orientation="horizontal"
                steps={steps}
                onSelect={(id: ProgressStepId) => void UtilityWindowsService.OpenGuides(id)}
            />
        </div>
    );
}
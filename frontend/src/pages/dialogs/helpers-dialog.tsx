import { useCallback, useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { DialogFooter } from "@/components/application/dialog-shell/dialog-footer";
import { DialogShell } from "@/components/application/dialog-shell/dialog-shell";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { Checkbox } from "@/components/base/checkbox/checkbox";
import {
    ApplicationsService,
    EditorDialogsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { HelperGroupDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function HelpersDialog() {
    const { applicationId = "" } = useParams();
    const [groups, setGroups] = useState<HelperGroupDTO[]>([]);
    const [appName, setAppName] = useState("");
    const [checked, setChecked] = useState<Record<string, boolean>>({});
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);
    const [loaded, setLoaded] = useState(false);

    // Dirty / preview vs catalog confirmed only.
    const preview = useMemo(() => {
        let adds = 0;
        let removes = 0;
        for (const group of groups) {
            const was = !!group.confirmed;
            const now = !!checked[group.bundleId];
            if (!was && now) {
                adds += 1;
            }
            if (was && !now) {
                removes += 1;
            }
        }
        return { adds, removes };
    }, [checked, groups]);

    const dirty = preview.adds > 0 || preview.removes > 0;

    const selectedCount = useMemo(
        () => groups.reduce((count, group) => count + (checked[group.bundleId] ? 1 : 0), 0),
        [checked, groups],
    );
    const allSelected = groups.length > 0 && selectedCount === groups.length;
    const noneSelected = selectedCount === 0;

    useEffect(() => {
        void EditorDialogsService.SetDirty("helpers", dirty);
    }, [dirty]);

    useEffect(() => {
        let cancelled = false;
        const load = async () => {
            setLoaded(false);
            setError(null);
            try {
                const result = await ApplicationsService.InspectHelpers(applicationId);
                if (cancelled) {
                    return;
                }
                const label = result.application.displayName || result.application.bundleId;
                setAppName(label);
                setGroups(result.groups ?? []);
                // Only already-confirmed groups start checked (owner UX 2026-11-17).
                const next: Record<string, boolean> = {};
                for (const group of result.groups ?? []) {
                    next[group.bundleId] = !!group.confirmed;
                }
                setChecked(next);
                void EditorDialogsService.SetTitle(
                    "helpers",
                    label ? `Review Helpers — ${label}` : "Review Helpers",
                );
                setLoaded(true);
            } catch (err) {
                if (!cancelled) {
                    setError(err instanceof Error ? err.message : String(err));
                    setLoaded(true);
                }
            }
        };
        if (applicationId) {
            void load();
        }
        return () => {
            cancelled = true;
        };
    }, [applicationId]);

    const onCancel = useCallback(() => {
        void EditorDialogsService.RequestClose("helpers");
    }, []);

    const onToggleSelectAll = useCallback(
        (selected: boolean) => {
            setChecked(() => {
                const next: Record<string, boolean> = {};
                for (const group of groups) {
                    next[group.bundleId] = selected;
                }
                return next;
            });
        },
        [groups],
    );

    const onSave = useCallback(async () => {
        setBusy(true);
        setError(null);
        try {
            const selected = groups
                .filter((group) => checked[group.bundleId])
                .map((group) => ({
                    bundleId: group.representative.bundleId,
                    executableName: group.representative.executableName,
                    path: group.representative.path,
                    displayName: group.displayName,
                    isUpdater: group.isUpdater,
                }));
            await ApplicationsService.ConfirmHelpers(applicationId, selected);
            await EditorDialogsService.SetDirty("helpers", false);
            await EditorDialogsService.Close("helpers");
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    }, [applicationId, checked, groups]);

    return (
        <DialogShell
            footer={
                <DialogFooter
                    leading={
                        <>
                            Adds {preview.adds} · Removes {preview.removes}
                            {appName ? ` · ${appName}` : ""}
                        </>
                    }
                    saveLabel="Confirm"
                    onCancel={onCancel}
                    onSave={() => void onSave()}
                    saveDisabled={!loaded || !dirty || busy}
                    saveLoading={busy}
                />
            }
        >
            {!loaded ? (
                <div className="flex min-h-64 items-center justify-center">
                    <LoadingIndicator label="Loading helpers…" />
                </div>
            ) : (
                <div className="flex flex-col gap-5">
                    {error && (
                        <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                            {error}
                        </div>
                    )}
                    {groups.length === 0 ? (
                        <p className="text-sm text-tertiary">No helper processes found. RouteDeck found no helper processes for this app.</p>
                    ) : (
                        <>
                            <div className="flex flex-wrap items-center gap-3">
                                <Checkbox
                                    label="Select all"
                                    isSelected={allSelected}
                                    isIndeterminate={!allSelected && !noneSelected}
                                    isDisabled={busy}
                                    onChange={onToggleSelectAll}
                                />
                                <p className="text-sm text-tertiary">
                                    {selectedCount} of {groups.length} selected
                                </p>
                            </div>
                            <ul className="flex flex-col gap-3">
                                {groups.map((group) => {
                                    const count = group.executableNames?.length || 1;
                                    const hints: string[] = [];
                                    if (group.sameDeveloper) {
                                        hints.push("Same developer");
                                    }
                                    if (group.isUpdater) {
                                        hints.push("Updater");
                                    }
                                    return (
                                        <li key={group.bundleId} className="rounded-lg px-3 py-3 ring-1 ring-secondary">
                                            <Checkbox
                                                isSelected={!!checked[group.bundleId]}
                                                onChange={(selected) =>
                                                    setChecked((prev) => ({ ...prev, [group.bundleId]: selected }))
                                                }
                                                label={`${group.displayName || group.bundleId} · ${count} process${count === 1 ? "" : "es"}`}
                                            />
                                            <p className="mt-1 pl-7 text-xs text-tertiary">{group.bundleId}</p>
                                            {group.executableNames?.length ? (
                                                <p className="mt-1 pl-7 text-xs text-tertiary">
                                                    {group.executableNames.join(", ")}
                                                </p>
                                            ) : null}
                                            {hints.length > 0 ? (
                                                <p className="mt-1 pl-7 text-xs text-tertiary">{hints.join(" · ")}</p>
                                            ) : null}
                                        </li>
                                    );
                                })}
                            </ul>
                        </>
                    )}
                </div>
            )}
        </DialogShell>
    );
}

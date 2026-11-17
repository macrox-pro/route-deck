import { useCallback, useEffect, useMemo, useState } from "react";
import { SearchLg } from "@untitledui/icons";
import { Button } from "@/components/base/buttons/button";
import { Checkbox } from "@/components/base/checkbox/checkbox";
import { Input } from "@/components/base/input/input";
import { ApplicationsService } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type {
    ApplicationDTO,
    HelperDTO,
    InstalledApplicationDTO,
} from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

// Cap concurrent Icon RPC so opening Applications does not stampede AppKit.
const ICON_FETCH_LIMIT = 2;
let iconFetchesActive = 0;
const iconFetchWaiters: Array<() => void> = [];
const iconFetchByPath = new Map<string, Promise<string | null>>();

async function acquireIconSlot(): Promise<void> {
    if (iconFetchesActive < ICON_FETCH_LIMIT) {
        iconFetchesActive += 1;
        return;
    }
    await new Promise<void>((resolve) => {
        iconFetchWaiters.push(resolve);
    });
    iconFetchesActive += 1;
}

function releaseIconSlot(): void {
    iconFetchesActive -= 1;
    const next = iconFetchWaiters.shift();
    if (next) {
        next();
    }
}

function loadAppIcon(path: string): Promise<string | null> {
    const existing = iconFetchByPath.get(path);
    if (existing) {
        return existing;
    }
    const pending = (async () => {
        await acquireIconSlot();
        try {
            // Wails v3 maps []byte to a base64 string in TypeScript bindings.
            const encoded = await ApplicationsService.Icon(path);
            if (!encoded) {
                return null;
            }
            return `data:image/png;base64,${encoded}`;
        } catch {
            return null;
        } finally {
            releaseIconSlot();
        }
    })();
    iconFetchByPath.set(path, pending);
    return pending;
}

function useAppIcon(path: string | undefined, enabled: boolean) {
    const [src, setSrc] = useState<string | null>(null);
    useEffect(() => {
        if (!enabled || !path) {
            setSrc(null);
            return;
        }
        let cancelled = false;
        void loadAppIcon(path).then((next) => {
            if (!cancelled) {
                setSrc(next);
            }
        });
        return () => {
            cancelled = true;
        };
    }, [path, enabled]);
    return src;
}

function InstalledRow({
    row,
    busy,
    onAdd,
}: {
    row: InstalledApplicationDTO;
    busy: boolean;
    onAdd: (path: string) => void;
}) {
    const icon = useAppIcon(row.path, true);
    return (
        <li className="flex items-center justify-between gap-3 px-4 py-3 text-sm">
            <div className="flex min-w-0 items-center gap-3">
                <div className="flex size-8 shrink-0 items-center justify-center overflow-hidden rounded-md bg-secondary ring-1 ring-secondary">
                    {icon ? <img src={icon} alt="" className="size-8" /> : null}
                </div>
                <div className="min-w-0">
                    <p className="truncate font-medium text-primary">{row.displayName}</p>
                    <p className="truncate text-tertiary">{row.bundleId}</p>
                </div>
            </div>
            {row.inCatalog ? (
                <span className="shrink-0 text-xs text-tertiary">In catalog</span>
            ) : (
                <Button color="secondary" size="sm" isLoading={busy} onClick={() => onAdd(row.path)}>
                    Add
                </Button>
            )}
        </li>
    );
}

export function ApplicationsPage() {
    const [apps, setApps] = useState<ApplicationDTO[]>([]);
    const [installed, setInstalled] = useState<InstalledApplicationDTO[]>([]);
    const [filter, setFilter] = useState("");
    const [showSystem, setShowSystem] = useState(false);
    const [pendingAppID, setPendingAppID] = useState<string | null>(null);
    const [helpers, setHelpers] = useState<HelperDTO[]>([]);
    const [selected, setSelected] = useState<Record<string, boolean>>({});
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const [catalog, installedRows] = await Promise.all([ApplicationsService.List(), ApplicationsService.ListInstalled()]);
            setApps(catalog);
            setInstalled(installedRows);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload]);

    const filteredInstalled = useMemo(() => {
        const q = filter.trim().toLowerCase();
        return installed.filter((row) => {
            if (!showSystem && row.isSystem) {
                return false;
            }
            if (!q) {
                return true;
            }
            return row.displayName.toLowerCase().includes(q) || row.bundleId.toLowerCase().includes(q);
        });
    }, [filter, installed, showSystem]);

    const applyAddResult = (result: Awaited<ReturnType<typeof ApplicationsService.AddFromPath>>) => {
        setPendingAppID(result.application.id);
        setHelpers(result.helpers ?? []);
        const initial: Record<string, boolean> = {};
        for (const helper of result.helpers ?? []) {
            initial[helper.path] = false;
        }
        setSelected(initial);
    };

    const onPick = async () => {
        setBusy(true);
        setError(null);
        try {
            const result = await ApplicationsService.PickAndAdd();
            if (!result?.application?.id) {
                return;
            }
            applyAddResult(result);
            await reload();
        } catch (err) {
            const message = err instanceof Error ? err.message : String(err);
            if (/cancelled|canceled/i.test(message)) {
                return;
            }
            setError(message);
        } finally {
            setBusy(false);
        }
    };

    const onAddFromPath = async (path: string) => {
        setBusy(true);
        setError(null);
        try {
            const result = await ApplicationsService.AddFromPath(path);
            applyAddResult(result);
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onConfirmHelpers = async () => {
        if (!pendingAppID) {
            return;
        }
        setBusy(true);
        setError(null);
        try {
            const chosen = helpers.filter((helper) => selected[helper.path]);
            await ApplicationsService.ConfirmHelpers(pendingAppID, chosen);
            setPendingAppID(null);
            setHelpers([]);
            setSelected({});
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onDelete = async (id: string) => {
        setBusy(true);
        setError(null);
        try {
            await ApplicationsService.Delete(id);
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between gap-4">
                <h1 className="text-display-xs font-semibold text-primary">Applications</h1>
                <Button color="primary" size="md" isLoading={busy} onClick={() => void onPick()}>
                    Add from disk…
                </Button>
            </div>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            {helpers.length > 0 && pendingAppID && (
                <section className="flex flex-col gap-3 rounded-lg p-4 ring-1 ring-secondary">
                    <h2 className="text-sm font-semibold text-secondary">Helper processes</h2>
                    <p className="text-sm text-tertiary">
                        Select helpers to include as separate process identities. Unchecked helpers are never written into rules.
                    </p>
                    {helpers.map((helper) => (
                        <Checkbox
                            key={helper.path}
                            label={`${helper.displayName || helper.executableName} (${helper.bundleId})`}
                            isSelected={!!selected[helper.path]}
                            onChange={(isSelected) => setSelected((prev) => ({ ...prev, [helper.path]: isSelected }))}
                        />
                    ))}
                    <div className="flex gap-2">
                        <Button color="primary" size="md" isLoading={busy} onClick={() => void onConfirmHelpers()}>
                            Confirm helpers
                        </Button>
                        <Button
                            color="secondary"
                            size="md"
                            onClick={() => {
                                setPendingAppID(null);
                                setHelpers([]);
                                setSelected({});
                            }}
                        >
                            Skip helpers
                        </Button>
                    </div>
                </section>
            )}

            <section className="flex flex-col gap-3">
                <h2 className="text-sm font-semibold text-secondary">Catalog</h2>
                <ul className="divide-y divide-secondary rounded-lg ring-1 ring-secondary">
                    {apps.length === 0 ? (
                        <li className="px-4 py-6 text-sm text-tertiary">No applications in the catalog yet.</li>
                    ) : (
                        apps.map((app) => (
                            <li key={app.id} className="flex items-center justify-between gap-3 px-4 py-3 text-sm">
                                <div>
                                    <p className="font-medium text-primary">{app.displayName}</p>
                                    <p className="text-tertiary">{app.bundleId}</p>
                                </div>
                                <Button color="secondary-destructive" size="sm" isLoading={busy} onClick={() => void onDelete(app.id)}>
                                    Delete
                                </Button>
                            </li>
                        ))
                    )}
                </ul>
            </section>

            <section className="flex flex-col gap-3">
                <h2 className="text-sm font-semibold text-secondary">Installed on this Mac</h2>
                <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:gap-4">
                    <Input
                        aria-label="Filter installed applications"
                        icon={SearchLg}
                        size="sm"
                        value={filter}
                        onChange={setFilter}
                        placeholder="Name or bundle id"
                        className="w-full sm:max-w-xs"
                    />
                    <Checkbox label="Show system applications" isSelected={showSystem} onChange={setShowSystem} />
                </div>
                <ul className="max-h-96 divide-y divide-secondary overflow-y-auto rounded-lg ring-1 ring-secondary">
                    {filteredInstalled.length === 0 ? (
                        <li className="px-4 py-6 text-sm text-tertiary">No installed applications match.</li>
                    ) : (
                        filteredInstalled.map((row) => (
                            <InstalledRow key={row.path} row={row} busy={busy} onAdd={(path) => void onAddFromPath(path)} />
                        ))
                    )}
                </ul>
            </section>
        </div>
    );
}

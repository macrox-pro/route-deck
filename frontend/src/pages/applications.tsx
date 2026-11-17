import { useCallback, useEffect, useMemo, useState } from "react";
import { SearchLg } from "@untitledui/icons";
import { Button } from "@/components/base/buttons/button";
import { Checkbox } from "@/components/base/checkbox/checkbox";
import { Input } from "@/components/base/input/input";
import { useSnapshot } from "@/hooks/use-snapshot";
import {
    ApplicationsService,
    ProfilesService,
    ProxiesService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type {
    ApplicationDTO,
    HelperDTO,
    InstalledApplicationDTO,
    ProcessIdentityDTO,
    ProfileDTO,
    ProxyDTO,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

type CatalogPolicyLine = {
    key: string;
    profileName: string;
    processLabel: string;
    actionLabel: string;
    inForce: boolean;
    protectionWhenInForce: string;
};

function actionLabel(action: string, proxyId: string, proxyNameById: Map<string, string>): string {
    switch (action) {
        case "proxy": {
            const proxyName = proxyNameById.get(proxyId);
            return proxyName ? `Proxy · ${proxyName}` : "Proxy";
        }
        case "direct":
            return "Direct";
        case "block":
            return "Block";
        default:
            return action || "Unknown";
    }
}

function buildCatalogPolicyLines(
    applicationID: string,
    processes: ProcessIdentityDTO[],
    profiles: ProfileDTO[],
    proxyNameById: Map<string, string>,
    appliedId: string,
    appliedConfirmed: boolean,
    protection: string,
): CatalogPolicyLine[] {
    const processByID = new Map<string, ProcessIdentityDTO>();
    const appProcessIDs = new Set<string>();
    for (const process of processes) {
        processByID.set(process.id, process);
        if (process.applicationId === applicationID) {
            appProcessIDs.add(process.id);
        }
    }
    if (appProcessIDs.size === 0) {
        return [];
    }

    const lines: CatalogPolicyLine[] = [];
    for (const profile of profiles) {
        const inForce = appliedConfirmed && appliedId === profile.id;
        for (const policy of profile.policies ?? []) {
            if (!appProcessIDs.has(policy.processId)) {
                continue;
            }
            const process = processByID.get(policy.processId);
            const processLabel = process?.executableName || process?.bundleId || policy.processId;
            lines.push({
                key: `${profile.id}:${policy.processId}:${policy.action}:${policy.proxyId}`,
                profileName: profile.name || profile.id,
                processLabel,
                actionLabel: actionLabel(policy.action, policy.proxyId, proxyNameById),
                inForce,
                protectionWhenInForce: inForce ? protection || "Unknown" : "",
            });
        }
    }
    return lines;
}

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

function AppIconSlot({ path }: { path: string | undefined }) {
    const icon = useAppIcon(path, !!path);
    return (
        <div className="flex size-8 shrink-0 items-center justify-center overflow-hidden rounded-md bg-secondary ring-1 ring-secondary">
            {icon ? <img src={icon} alt="" className="size-8" /> : null}
        </div>
    );
}

function CatalogRow({
    app,
    path,
    policyLines,
    busy,
    onDelete,
}: {
    app: ApplicationDTO;
    path: string | undefined;
    policyLines: CatalogPolicyLine[];
    busy: boolean;
    onDelete: (id: string) => void;
}) {
    return (
        <li className="flex flex-col gap-2 px-4 py-3 text-sm">
            <div className="flex items-center justify-between gap-3">
                <div className="flex min-w-0 items-center gap-3">
                    <AppIconSlot path={path} />
                    <div className="min-w-0">
                        <p className="truncate font-medium text-primary">{app.displayName}</p>
                        <p className="truncate text-tertiary">{app.bundleId}</p>
                    </div>
                </div>
                <Button color="secondary-destructive" size="sm" isLoading={busy} onClick={() => onDelete(app.id)}>
                    Delete
                </Button>
            </div>
            {policyLines.length > 0 ? (
                <ul className="flex flex-col gap-1 pl-11">
                    {policyLines.map((line) => (
                        <li key={line.key} className="truncate text-xs text-tertiary">
                            <span className="text-secondary">{line.profileName}</span>
                            {" · "}
                            {line.processLabel}
                            {" · "}
                            {line.actionLabel}
                            {" · "}
                            {line.inForce ? (
                                <>
                                    In force
                                    {line.protectionWhenInForce ? ` · ${line.protectionWhenInForce}` : ""}
                                </>
                            ) : (
                                "Not in force"
                            )}
                        </li>
                    ))}
                </ul>
            ) : null}
        </li>
    );
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
    return (
        <li className="flex items-center justify-between gap-3 px-4 py-3 text-sm">
            <div className="flex min-w-0 items-center gap-3">
                <AppIconSlot path={row.path} />
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
    const { data: snapshot } = useSnapshot();
    const [apps, setApps] = useState<ApplicationDTO[]>([]);
    const [installed, setInstalled] = useState<InstalledApplicationDTO[]>([]);
    const [profiles, setProfiles] = useState<ProfileDTO[]>([]);
    const [processes, setProcesses] = useState<ProcessIdentityDTO[]>([]);
    const [proxies, setProxies] = useState<ProxyDTO[]>([]);
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
            const [catalog, installedRows, profileRows, processRows, proxyRows] = await Promise.all([
                ApplicationsService.List(),
                ApplicationsService.ListInstalled(),
                ProfilesService.List(),
                ApplicationsService.ListProcessIdentities(),
                ProxiesService.List(),
            ]);
            setApps(catalog);
            setInstalled(installedRows);
            setProfiles(profileRows);
            setProcesses(processRows);
            setProxies(proxyRows);
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

    // ApplicationDTO has no path; Icon needs a .app path from the installed index.
    const installedPathByBundle = useMemo(() => {
        const map = new Map<string, string>();
        for (const row of installed) {
            map.set(row.bundleId, row.path);
        }
        return map;
    }, [installed]);

    const proxyNameById = useMemo(() => {
        const map = new Map<string, string>();
        for (const proxy of proxies) {
            map.set(proxy.id, proxy.name || proxy.id);
        }
        return map;
    }, [proxies]);

    const policyLinesByApp = useMemo(() => {
        const map = new Map<string, CatalogPolicyLine[]>();
        const appliedId = snapshot?.appliedId ?? "";
        const appliedConfirmed = snapshot?.appliedConfirmed ?? false;
        const protection = snapshot?.protection ?? "Unknown";
        for (const app of apps) {
            map.set(
                app.id,
                buildCatalogPolicyLines(app.id, processes, profiles, proxyNameById, appliedId, appliedConfirmed, protection),
            );
        }
        return map;
    }, [apps, processes, profiles, proxyNameById, snapshot?.appliedConfirmed, snapshot?.appliedId, snapshot?.protection]);

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
                            <CatalogRow
                                key={app.id}
                                app={app}
                                path={installedPathByBundle.get(app.bundleId)}
                                policyLines={policyLinesByApp.get(app.id) ?? []}
                                busy={busy}
                                onDelete={(id) => void onDelete(id)}
                            />
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

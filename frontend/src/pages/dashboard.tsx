import { useEffect, useMemo, useState } from "react";
import { HelpCircle } from "@untitledui/icons";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { ProtectionBadge } from "@/components/application/protection-badge/protection-badge";
import { ProxyListIcon } from "@/components/application/proxy-list-icon/proxy-list-icon";
import { SetupChecklist } from "@/components/application/setup-guide/setup-checklist";
import { SuggestionBanner } from "@/components/application/setup-guide/suggestion-banner";
import { Button } from "@/components/base/buttons/button";
import { useSnapshot } from "@/hooks/use-snapshot";
import { matchProxyGuide } from "@/utils/match-proxy-guide";
import { proxyHealthLabel } from "@/utils/proxy-health-label";
import {
    ProtectionService,
    SetupGuidesService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { GuideDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function DashboardPage() {
    const { data, error, loading, reload } = useSnapshot();
    const [actionError, setActionError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);
    const [operationID, setOperationID] = useState<string | null>(null);
    const [checklistKey, setChecklistKey] = useState(0);
    const [guides, setGuides] = useState<GuideDTO[]>([]);

    useEffect(() => {
        void SetupGuidesService.List(false)
            .then(setGuides)
            .catch(() => setGuides([]));
    }, []);

    const proxyIconByID = useMemo(() => {
        const map = new Map<string, { installedPath?: string; kind: string }>();
        for (const proxy of data?.proxies ?? []) {
            const match = matchProxyGuide(proxy.host || "127.0.0.1", proxy.port || 0, proxy.kind || "socks5", guides);
            map.set(proxy.id, { installedPath: match.installedPath, kind: match.kindFallback });
        }
        return map;
    }, [data?.proxies, guides]);

    if (loading && !data) {
        return (
            <div className="flex min-h-64 items-center justify-center">
                <LoadingIndicator label="Loading snapshot…" />
            </div>
        );
    }

    const onApply = async () => {
        setBusy(true);
        setActionError(null);
        try {
            const result = await ProtectionService.Apply();
            setOperationID(result.operationId);
            await reload();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onStop = async () => {
        setBusy(true);
        setActionError(null);
        try {
            const result = await ProtectionService.Stop();
            setOperationID(result.operationId);
            await reload();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    return (
        <div className="flex flex-col gap-6">
            <div className="flex flex-wrap items-center justify-between gap-4">
                <div>
                    <h1 className="text-display-xs font-semibold text-primary">Dashboard</h1>
                    <p className="mt-1 max-w-2xl text-sm text-tertiary">
                        Protection comes from a Go snapshot after Apply. Events only trigger a reload — they are not proof of
                        Active. Applied means the CLI confirmed a profile; it is not the same as Active.
                    </p>
                </div>
                <div className="flex flex-wrap gap-2">
                    <Button color="secondary" size="md" iconLeading={HelpCircle} href="/help">
                        Setup guides
                    </Button>
                    <Button color="primary" size="md" isLoading={busy} onClick={() => void onApply()}>
                        Apply
                    </Button>
                    <Button color="secondary" size="md" isLoading={busy} onClick={() => void onStop()}>
                        Stop
                    </Button>
                    <Button color="secondary" size="md" isLoading={loading} onClick={() => void reload()}>
                        Refresh
                    </Button>
                </div>
            </div>

            {(error || actionError) && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {actionError || error}
                </div>
            )}

            <SuggestionBanner
                onAdded={() => {
                    void reload();
                    setChecklistKey((k) => k + 1);
                }}
            />
            <SetupChecklist refreshKey={checklistKey} />

            {data?.databaseError && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="alert">
                    Database: {data.databaseError}
                </div>
            )}

            {data?.bridge === "NotInstalled" && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                    ProxyBridge is not installed. Open Setup guides → ProxyBridge for install steps, or set a custom path in Settings.
                </div>
            )}

            {data?.desiredDiffersFromApplied && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                    Desired profile differs from the applied profile. Press Apply to put the desired profile into force.
                </div>
            )}

            {data?.compatibilityWarning && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                    {data.compatibilityWarning}
                </div>
            )}

            {data?.recheckPending && (
                <div className="rounded-lg bg-secondary px-4 py-3 text-sm text-secondary ring-1 ring-secondary" role="status">
                    Recheck pending after wake or network change. Protection shows Unknown until health and bridge checks finish.
                </div>
            )}

            <p className="text-sm text-tertiary">
                Helper processes that are not confirmed are not in rules — their traffic can bypass Proxy or Block. Review
                helpers on the Applications page.
            </p>

            {operationID && <p className="text-xs text-tertiary">Last operation: {operationID}</p>}

            {data && (
                <>
                    <div className="flex flex-col gap-2">
                        <div className="flex flex-wrap items-center gap-3">
                            <ProtectionBadge protection={data.protection} />
                            <span className="text-sm text-secondary">Bridge: {data.bridge || "Unknown"}</span>
                            {data.version && <span className="text-sm text-tertiary">CLI {data.version}</span>}
                            {data.desiredSet && <span className="text-sm text-tertiary">Desired: {data.desiredId || "—"}</span>}
                            {data.appliedConfirmed && <span className="text-sm text-tertiary">Applied: {data.appliedId || "—"}</span>}
                        </div>
                        <p className="max-w-2xl text-xs text-tertiary">
                            Active: routing confirmed and every protected app fully covered. Degraded: rules in effect with a
                            gap (for example TCP OK, UDP unavailable). Inactive: nothing to protect, or Direct-only policies.
                            Unprotected: rules are not in effect. Applied ≠ Active.
                        </p>
                    </div>

                    <section>
                        <h2 className="mb-2 text-sm font-semibold text-secondary">Proxies</h2>
                        {data.proxies.length === 0 ? (
                            <p className="text-sm text-tertiary">No proxies in the catalog.</p>
                        ) : (
                            <ul className="divide-y divide-secondary rounded-lg ring-1 ring-secondary">
                                {data.proxies.map((proxy) => {
                                    const icon = proxyIconByID.get(proxy.id);
                                    return (
                                        <li key={proxy.id} className="flex items-center justify-between gap-3 px-4 py-3 text-sm">
                                            <div className="flex min-w-0 items-center gap-3">
                                                <ProxyListIcon installedPath={icon?.installedPath} kind={icon?.kind || proxy.kind || "socks5"} />
                                                <span className="truncate text-primary">{proxy.name || proxy.id}</span>
                                            </div>
                                            <span className="shrink-0 text-tertiary">{proxyHealthLabel(proxy.health)}</span>
                                        </li>
                                    );
                                })}
                            </ul>
                        )}
                    </section>
                </>
            )}
        </div>
    );
}

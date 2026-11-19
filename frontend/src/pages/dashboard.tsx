import { type ReactNode, useCallback, useEffect, useMemo, useState } from "react";
import { HelpCircle, RefreshCw01 } from "@untitledui/icons";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { ProtectionBadge } from "@/components/application/protection-badge/protection-badge";
import { ProxyHealthBadge } from "@/components/application/proxy-health-badge/proxy-health-badge";
import { ProxyListIcon } from "@/components/application/proxy-list-icon/proxy-list-icon";
import { SetupChecklist } from "@/components/application/setup-guide/setup-checklist";
import { SuggestionBanner } from "@/components/application/setup-guide/suggestion-banner";
import { Button } from "@/components/base/buttons/button";
import { useSnapshot } from "@/hooks/use-snapshot";
import { matchProxyGuide } from "@/utils/match-proxy-guide";
import {
    ProtectionService,
    SetupGuidesService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { GuideDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

type DashboardPending = "apply" | "stop" | "refresh" | null;

export function DashboardPage() {
    const { data, error, loading, reload } = useSnapshot();
    const [actionError, setActionError] = useState<string | null>(null);
    const [pending, setPending] = useState<DashboardPending>(null);
    const [operationID, setOperationID] = useState<string | null>(null);
    const [checklistKey, setChecklistKey] = useState(0);
    const [guides, setGuides] = useState<GuideDTO[]>([]);
    const actionsLocked = pending !== null;

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

    const bumpChecklist = useCallback(() => setChecklistKey((k) => k + 1), []);
    const noDesired = !!data && (!data.desiredSet || !data.desiredId);
    const canStop =
        !!data &&
        (data.bridge === "Ready" || data.bridge === "Starting" || data.bridge === "Stopping" || data.appliedConfirmed);
    const applyDisabled = noDesired || (actionsLocked && pending !== "apply");
    const stopDisabled = !canStop || (actionsLocked && pending !== "stop");
    const applyTitle = noDesired ? "Set a Desired profile on Profiles first" : undefined;
    const stopTitle = !canStop ? "Nothing to stop — protection is not running" : undefined;

    if (loading && !data) {
        return (
            <div className="flex min-h-64 items-center justify-center">
                <LoadingIndicator label="Loading snapshot…" />
            </div>
        );
    }

    const onApply = async () => {
        setPending("apply");
        setActionError(null);
        try {
            const result = await ProtectionService.Apply();
            setOperationID(result.operationId);
            await reload();
            bumpChecklist();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setPending(null);
        }
    };

    const onStop = async () => {
        setPending("stop");
        setActionError(null);
        try {
            const result = await ProtectionService.Stop();
            setOperationID(result.operationId);
            await reload();
            bumpChecklist();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setPending(null);
        }
    };

    const onRefresh = async () => {
        setPending("refresh");
        setActionError(null);
        try {
            await reload();
            bumpChecklist();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setPending(null);
        }
    };

    const warnings: ReactNode[] = [];
    if (data?.databaseError) {
        warnings.push(
            <div key="db" className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="alert">
                Database: {data.databaseError}
            </div>,
        );
    }
    if (data?.bridge === "NotInstalled") {
        warnings.push(
            <div key="bridge" className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                ProxyBridge is not installed. Open Setup guides → ProxyBridge for install steps, or set a custom path in Settings.
            </div>,
        );
    }
    if (noDesired) {
        warnings.push(
            <div
                key="desired"
                className="flex flex-wrap items-center justify-between gap-3 rounded-lg bg-secondary px-4 py-3 text-sm text-secondary ring-1 ring-secondary"
                role="status"
            >
                <span>No desired profile. Set one on Profiles to enable Apply.</span>
                <Button color="secondary" size="sm" href="/profiles">
                    Open Profiles
                </Button>
            </div>,
        );
    }
    if (data?.desiredDiffersFromApplied) {
        warnings.push(
            <div key="differs" className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                Desired profile differs from the applied profile. Press Apply to put the desired profile into force.
            </div>,
        );
    }
    if (data?.compatibilityWarning) {
        warnings.push(
            <div key="compat" className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                {data.compatibilityWarning}
            </div>,
        );
    }
    if (data?.recheckPending) {
        warnings.push(
            <div key="recheck" className="rounded-lg bg-secondary px-4 py-3 text-sm text-secondary ring-1 ring-secondary" role="status">
                Recheck pending after wake or network change. Protection shows Unknown until health and bridge checks finish.
            </div>,
        );
    }

    return (
        <div className="flex min-w-0 flex-col gap-6">
            {/* Page header: title + actions aligned like Profiles / DESIGN.md */}
            <div className="flex flex-wrap items-center justify-between gap-4">
                <div className="min-w-0">
                    <h1 className="text-display-xs font-semibold text-primary">Dashboard</h1>
                    <p className="mt-1 max-w-xl text-sm text-tertiary">
                        Protection status and proxy health for this Mac.
                    </p>
                </div>
                <div className="flex flex-wrap items-center gap-4">
                    <Button color="secondary" size="md" iconLeading={HelpCircle} href="/help" isDisabled={actionsLocked}>
                        Setup guides
                    </Button>
                    {/* Affirmative / destructive pair, separated from Setup guides (NN/g proximity). */}
                    <div className="flex items-center gap-2">
                        <span title={stopTitle} className="inline-flex">
                            <Button
                                color="secondary-destructive"
                                size="md"
                                isLoading={pending === "stop"}
                                showTextWhileLoading
                                isDisabled={stopDisabled}
                                onClick={() => void onStop()}
                            >
                                Stop
                            </Button>
                        </span>
                        <span title={applyTitle} className="inline-flex">
                            <Button
                                color="primary"
                                size="md"
                                isLoading={pending === "apply"}
                                showTextWhileLoading
                                isDisabled={applyDisabled}
                                onClick={() => void onApply()}
                            >
                                Apply
                            </Button>
                        </span>
                    </div>
                </div>
            </div>

            {(error || actionError) && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {actionError || error}
                </div>
            )}

            {/* DESIGN.md order: 1 status → 2 checklist → 3 warnings → 4 proxies */}
            {data && (
                <div className="flex flex-col gap-3 rounded-lg bg-secondary_subtle p-4 ring-1 ring-secondary">
                    <div className="flex flex-wrap items-center justify-between gap-3">
                        <div className="flex min-w-0 flex-wrap items-center gap-x-3 gap-y-2">
                            <ProtectionBadge protection={data.protection} />
                            <span className="text-sm text-secondary">Bridge: {data.bridge || "Unknown"}</span>
                            {data.version ? <span className="text-sm text-tertiary">CLI {data.version}</span> : null}
                            {data.desiredSet ? <span className="text-sm text-tertiary">Desired: {data.desiredId || "—"}</span> : null}
                            {data.appliedConfirmed ? (
                                <span className="text-sm text-tertiary">Applied: {data.appliedId || "—"}</span>
                            ) : null}
                        </div>
                        <Button
                            color="tertiary"
                            size="sm"
                            iconLeading={RefreshCw01}
                            aria-label="Refresh status"
                            isLoading={pending === "refresh"}
                            isDisabled={actionsLocked && pending !== "refresh"}
                            onClick={() => void onRefresh()}
                        />
                    </div>
                    {operationID ? <p className="text-xs text-tertiary">Last operation: {operationID}</p> : null}
                </div>
            )}

            <SuggestionBanner
                onAdded={() => {
                    void reload();
                    bumpChecklist();
                }}
            />
            <SetupChecklist refreshKey={checklistKey} />

            {warnings.length > 0 ? <div className="flex flex-col gap-3">{warnings}</div> : null}

            {data && (
                <section className="flex flex-col gap-3">
                    <div className="flex flex-wrap items-end justify-between gap-2">
                        <h2 className="text-sm font-semibold text-secondary">Proxies</h2>
                        <p className="text-xs text-tertiary">
                            Unconfirmed helpers can bypass rules — review on Applications.
                        </p>
                    </div>
                    {data.proxies.length === 0 ? (
                        <p className="rounded-lg px-4 py-6 text-sm text-tertiary ring-1 ring-secondary">No proxies in the catalog.</p>
                    ) : (
                        <ul className="divide-y divide-secondary overflow-hidden rounded-lg ring-1 ring-secondary">
                            {data.proxies.map((proxy) => {
                                const icon = proxyIconByID.get(proxy.id);
                                return (
                                    <li key={proxy.id} className="flex items-center justify-between gap-3 px-4 py-3 text-sm">
                                        <div className="flex min-w-0 items-center gap-3">
                                            <ProxyListIcon installedPath={icon?.installedPath} kind={icon?.kind || proxy.kind || "socks5"} />
                                            <span className="truncate text-primary">{proxy.name || proxy.id}</span>
                                        </div>
                                        <ProxyHealthBadge health={proxy.health} />
                                    </li>
                                );
                            })}
                        </ul>
                    )}
                </section>
            )}
        </div>
    );
}

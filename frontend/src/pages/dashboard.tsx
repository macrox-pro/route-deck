import { useCallback, useEffect, useMemo, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { RefreshCw01 } from "@untitledui/icons";
import { InlineAlert } from "@/components/application/page-layout/inline-alert";
import { ListEmptyRow, ListPanel, ListRow } from "@/components/application/page-layout/list-panel";
import { PageHeader } from "@/components/application/page-layout/page-header";
import { PageSection } from "@/components/application/page-layout/page-section";
import { PageShell } from "@/components/application/page-layout/page-shell";
import { ScreenState } from "@/components/application/page-layout/screen-state";
import { ProtectionBadge } from "@/components/application/protection-badge/protection-badge";
import { ProxyHealthBadge } from "@/components/application/proxy-health-badge/proxy-health-badge";
import { ProxyListIcon } from "@/components/application/proxy-list-icon/proxy-list-icon";
import { SetupStrip } from "@/components/application/setup-guide/setup-strip";
import { SuggestionBanner } from "@/components/application/setup-guide/suggestion-banner";
import { Button } from "@/components/base/buttons/button";
import { Tooltip } from "@/components/base/tooltip/tooltip";
import { useSnapshot } from "@/hooks/use-snapshot";
import { matchProxyGuide } from "@/utils/match-proxy-guide";
import {
    ProfilesService,
    ProtectionService,
    ProxiesService,
    SetupGuidesService,
    UtilityWindowsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { GuideDTO, ProfileDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

type DashboardPending = "apply" | "stop" | "refresh" | "recheck" | null;

export function DashboardPage() {
    const { data, error, loading, reload } = useSnapshot();
    const location = useLocation();
    const navigate = useNavigate();
    const [actionError, setActionError] = useState<string | null>(null);
    const [pending, setPending] = useState<DashboardPending>(null);
    const [operationID, setOperationID] = useState<string | null>(null);
    const [checklistKey, setChecklistKey] = useState(0);
    const [guides, setGuides] = useState<GuideDTO[]>([]);
    const [profiles, setProfiles] = useState<ProfileDTO[]>([]);
    const actionsLocked = pending !== null;

    useEffect(() => {
        void SetupGuidesService.List(false)
            .then(setGuides)
            .catch(() => setGuides([]));
        void ProfilesService.List()
            .then(setProfiles)
            .catch(() => setProfiles([]));
    }, [data?.desiredId, data?.appliedId, data?.protection]);

    const profileName = useCallback(
        (id: string | undefined | null) => {
            if (!id) {
                return "—";
            }
            return profiles.find((p) => p.id === id)?.name || id;
        },
        [profiles],
    );

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

    const onApply = useCallback(async () => {
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
    }, [bumpChecklist, reload]);

    const onStop = useCallback(async () => {
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
    }, [bumpChecklist, reload]);

    const onRefresh = useCallback(async () => {
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
    }, [bumpChecklist, reload]);

    const onRecheck = useCallback(async () => {
        setPending("recheck");
        setActionError(null);
        try {
            await ProxiesService.Recheck();
            await reload();
            bumpChecklist();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setPending(null);
        }
    }, [bumpChecklist, reload]);

    useEffect(() => {
        const command = (location.state as { command?: string } | null)?.command;
        if (!command) {
            return;
        }
        navigate(location.pathname, { replace: true, state: null });
        if (command === "apply") {
            void onApply();
        } else if (command === "stop") {
            void onStop();
        } else if (command === "recheck") {
            void onRecheck();
        }
    }, [location.pathname, location.state, navigate, onApply, onRecheck, onStop]);

    return (
        <PageShell>
            <PageHeader
                title="Dashboard"
                description="Protection status and proxy health for this Mac."
                actions={
                    data ? (
                    <>
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
                        <span title={stopTitle} className="inline-flex">
                            <Button
                                color="secondary"
                                size="md"
                                isLoading={pending === "stop"}
                                showTextWhileLoading
                                isDisabled={stopDisabled}
                                onClick={() => void onStop()}
                            >
                                Stop
                            </Button>
                        </span>
                        <Tooltip title="Refresh status">
                            <Button
                                color="secondary"
                                size="md"
                                iconLeading={RefreshCw01}
                                aria-label="Refresh status"
                                isLoading={pending === "refresh"}
                                isDisabled={actionsLocked && pending !== "refresh"}
                                onClick={() => void onRefresh()}
                            />
                        </Tooltip>
                    </>
                    ) : undefined
                }
            />

            {error || actionError ? <InlineAlert tone="error">{actionError || error}</InlineAlert> : null}

            {loading && !data ? <ScreenState state="loading" title="Loading snapshot…" /> : (
            <>
            {data ? (
                <div className="flex flex-col gap-3 rounded-xl bg-secondary_subtle p-4 ring-1 ring-secondary">
                    <div className="flex min-w-0 flex-wrap items-center gap-x-3 gap-y-2">
                        <ProtectionBadge protection={data.protection} />
                        <span className="text-sm text-secondary">Bridge: {data.bridge || "Unknown"}</span>
                        {data.version ? <span className="font-mono text-xs text-tertiary">CLI {data.version}</span> : null}
                    </div>
                    <div className="flex flex-wrap gap-x-4 gap-y-1 text-sm text-tertiary">
                        <span>
                            Desired: <span className="font-medium text-secondary">{data.desiredSet ? profileName(data.desiredId) : "None"}</span>
                        </span>
                        <span>
                            Applied:{" "}
                            <span className="font-medium text-secondary">
                                {data.appliedConfirmed ? profileName(data.appliedId) : "None"}
                            </span>
                        </span>
                    </div>
                    {operationID ? <p className="font-mono text-xs text-tertiary">Last operation: {operationID}</p> : null}
                </div>
            ) : null}

            <SuggestionBanner
                onAdded={() => {
                    void reload();
                    bumpChecklist();
                }}
            />
            <SetupStrip refreshKey={checklistKey} />

            <div className="flex flex-col gap-3">
                {data?.databaseError ? (
                    <InlineAlert tone="warning" title="Database">
                        {data.databaseError}
                    </InlineAlert>
                ) : null}
                {data?.bridge === "NotInstalled" ? (
                    <InlineAlert
                        tone="warning"
                        title="ProxyBridge is not installed"
                        action={{ label: "View bridge guide", onPress: () => void UtilityWindowsService.OpenGuides("bridge") }}
                    >
                        Install ProxyBridge, or set a custom path in Settings.
                    </InlineAlert>
                ) : null}
                {noDesired ? (
                    <InlineAlert
                        tone="info"
                        title="No profile selected"
                        action={{ label: "Open profiles", onPress: () => navigate("/profiles") }}
                    >
                        Choose a Desired profile, then use Apply above.
                    </InlineAlert>
                ) : null}
                {data?.desiredDiffersFromApplied ? (
                    <InlineAlert tone="warning" title="Desired differs from Applied">
                        Use Apply above to put the Desired profile into force.
                    </InlineAlert>
                ) : null}
                {data?.compatibilityWarning ? <InlineAlert tone="warning">{data.compatibilityWarning}</InlineAlert> : null}
                {data?.recheckPending ? (
                    <InlineAlert tone="info">
                        Recheck pending after wake or network change. Protection shows Unknown until health and bridge checks finish.
                    </InlineAlert>
                ) : null}
            </div>

            {data ? (
                <PageSection title="Proxies" count={data.proxies.length}>
                    <ListPanel>
                        {data.proxies.length === 0 ? (
                            <ListEmptyRow
                                action={
                                    <Button color="secondary" size="sm" onClick={() => navigate("/proxies")}>
                                        Go to Proxies
                                    </Button>
                                }
                            >
                                No proxies yet.
                            </ListEmptyRow>
                        ) : (
                            data.proxies.map((proxy) => {
                                const icon = proxyIconByID.get(proxy.id);
                                return (
                                    <ListRow
                                        key={proxy.id}
                                        leading={
                                            <ProxyListIcon installedPath={icon?.installedPath} kind={icon?.kind || proxy.kind || "socks5"} />
                                        }
                                        title={proxy.name || proxy.id}
                                        meta={
                                            <span className="font-mono text-xs">
                                                {proxy.host}:{proxy.port}
                                            </span>
                                        }
                                        status={<ProxyHealthBadge health={proxy.health} />}
                                    />
                                );
                            })
                        )}
                    </ListPanel>
                </PageSection>
            ) : null}
            </>
            )}
        </PageShell>
    );
}

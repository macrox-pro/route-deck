import { useCallback, useEffect, useMemo, useState } from "react";
import { Events } from "@wailsio/runtime";
import { Key01, Route, Trash01 } from "@untitledui/icons";
import { InlineAlert } from "@/components/application/page-layout/inline-alert";
import { ListPanel, ListRow } from "@/components/application/page-layout/list-panel";
import { PageHeader } from "@/components/application/page-layout/page-header";
import { PageShell } from "@/components/application/page-layout/page-shell";
import { ScreenState } from "@/components/application/page-layout/screen-state";
import { ProxyHealthBadge } from "@/components/application/proxy-health-badge/proxy-health-badge";
import { ProxyListIcon } from "@/components/application/proxy-list-icon/proxy-list-icon";
import { SuggestionBanner } from "@/components/application/setup-guide/suggestion-banner";
import { Button } from "@/components/base/buttons/button";
import { useSnapshot } from "@/hooks/use-snapshot";
import { matchProxyGuide } from "@/utils/match-proxy-guide";
import {
    EditorDialogsService,
    ProxiesService,
    SetupGuidesService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { GuideDTO, ProxyDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function ProxiesPage() {
    const { data: snapshot } = useSnapshot();
    const [rows, setRows] = useState<ProxyDTO[]>([]);
    const [guides, setGuides] = useState<GuideDTO[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [busyID, setBusyID] = useState<string | null>(null);
    const [headerBusy, setHeaderBusy] = useState(false);
    const [loading, setLoading] = useState(true);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const [proxyRows, guideRows] = await Promise.all([ProxiesService.List(), SetupGuidesService.List(false)]);
            setRows(proxyRows);
            setGuides(guideRows);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload]);

    useEffect(() => {
        const unsub = Events.On("catalog.changed", () => {
            void reload();
        });
        return () => {
            unsub();
        };
    }, [reload]);

    const healthByID = useMemo(() => {
        const map = new Map<string, string>();
        for (const proxy of snapshot?.proxies ?? []) {
            map.set(proxy.id, proxy.health);
        }
        return map;
    }, [snapshot?.proxies]);

    const iconByID = useMemo(() => {
        const map = new Map<string, { installedPath?: string; kind: string }>();
        for (const row of rows) {
            const match = matchProxyGuide(row.host, row.port, row.kind, guides);
            map.set(row.id, { installedPath: match.installedPath, kind: match.kindFallback });
        }
        return map;
    }, [guides, rows]);

    const withBusy = async (id: string, fn: () => Promise<unknown>) => {
        setBusyID(id);
        setError(null);
        try {
            await fn();
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusyID(null);
        }
    };

    return (
        <PageShell>
            <PageHeader
                title="Proxies"
                description="Local SOCKS5 or HTTP endpoints from your VPN apps."
                actions={
                    !loading && rows.length > 0 ? (
                        <>
                            <Button
                                color="secondary"
                                size="md"
                                isLoading={headerBusy}
                                onClick={() => {
                                    setHeaderBusy(true);
                                    void ProxiesService.Recheck()
                                        .then(() => reload())
                                        .catch((err) => setError(err instanceof Error ? err.message : String(err)))
                                        .finally(() => setHeaderBusy(false));
                                }}
                            >
                                Recheck health
                            </Button>
                            <Button color="primary" size="md" onClick={() => void EditorDialogsService.OpenProxy("", "")}>
                                New proxy
                            </Button>
                        </>
                    ) : undefined
                }
            />

            {error ? <InlineAlert tone="error">{error}</InlineAlert> : null}

            <SuggestionBanner onAdded={() => void reload()} />

            {loading && rows.length === 0 ? (
                <ScreenState state="loading" title="Loading proxies…" />
            ) : rows.length === 0 ? (
                <ScreenState
                    state="empty"
                    icon={Route}
                    title="No proxies yet"
                    description="Add one, or use Setup Guides in the sidebar."
                    actions={
                        <Button color="primary" size="md" onClick={() => void EditorDialogsService.OpenProxy("", "")}>
                            New proxy
                        </Button>
                    }
                />
            ) : (
                <ListPanel>
                    {rows.map((row) => {
                        const icon = iconByID.get(row.id);
                        const rowBusy = busyID === row.id;
                        return (
                            <ListRow
                                key={row.id}
                                leading={<ProxyListIcon installedPath={icon?.installedPath} kind={icon?.kind || row.kind} />}
                                title={row.name || row.id}
                                meta={
                                    <span className="font-mono text-xs">
                                        {row.kind}://{row.host}:{row.port}
                                        {row.hasPassword ? " · password set" : " · no password"}
                                    </span>
                                }
                                status={<ProxyHealthBadge health={healthByID.get(row.id)} />}
                                actions={
                                    <Button
                                        color="secondary"
                                        size="sm"
                                        isDisabled={busyID !== null}
                                        onClick={() => void EditorDialogsService.OpenProxy(row.id, "")}
                                    >
                                        Edit
                                    </Button>
                                }
                                menu={[
                                    ...(row.hasPassword
                                        ? [
                                              {
                                                  id: "clear-password",
                                                  label: "Clear password",
                                                  icon: Key01,
                                                  isDisabled: busyID !== null,
                                                  onAction: () => void withBusy(row.id, () => ProxiesService.ClearPassword(row.id)),
                                              },
                                          ]
                                        : []),
                                    {
                                        id: "delete",
                                        label: "Delete",
                                        icon: Trash01,
                                        destructive: true,
                                        isDisabled: busyID !== null || rowBusy,
                                        onAction: () => void withBusy(row.id, () => ProxiesService.Delete(row.id)),
                                    },
                                ]}
                            />
                        );
                    })}
                </ListPanel>
            )}
        </PageShell>
    );
}

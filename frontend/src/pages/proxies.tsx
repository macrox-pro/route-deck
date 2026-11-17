import { useCallback, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { Events } from "@wailsio/runtime";
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
    const [searchParams, setSearchParams] = useSearchParams();
    const { data: snapshot } = useSnapshot();
    const [rows, setRows] = useState<ProxyDTO[]>([]);
    const [guides, setGuides] = useState<GuideDTO[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const [proxyRows, guideRows] = await Promise.all([ProxiesService.List(), SetupGuidesService.List(false)]);
            setRows(proxyRows);
            setGuides(guideRows);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
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

    useEffect(() => {
        const guideID = searchParams.get("guide");
        if (!guideID) {
            return;
        }
        void EditorDialogsService.OpenProxy("", guideID);
        const next = new URLSearchParams(searchParams);
        next.delete("guide");
        setSearchParams(next, { replace: true });
    }, [searchParams, setSearchParams]);

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

    const onDelete = async (id: string, name: string) => {
        if (!window.confirm(`Delete this proxy?\n\n${name}`)) {
            return;
        }
        setBusy(true);
        setError(null);
        try {
            await ProxiesService.Delete(id);
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onClearPassword = async (id: string) => {
        setBusy(true);
        setError(null);
        try {
            await ProxiesService.ClearPassword(id);
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onRecheck = async () => {
        setBusy(true);
        setError(null);
        try {
            await ProxiesService.Recheck();
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
                <h1 className="text-display-xs font-semibold text-primary">Proxies</h1>
                <div className="flex flex-wrap gap-2">
                    <Button color="secondary" size="md" isLoading={busy} onClick={() => void onRecheck()}>
                        Recheck health
                    </Button>
                    <Button color="primary" size="md" onClick={() => void EditorDialogsService.OpenProxy("", "")}>
                        New proxy
                    </Button>
                </div>
            </div>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            <SuggestionBanner onAdded={() => void reload()} />

            <ul className="divide-y divide-secondary rounded-lg ring-1 ring-secondary">
                {rows.length === 0 ? (
                    <li className="px-4 py-6 text-sm text-tertiary">No proxies yet.</li>
                ) : (
                    rows.map((row) => {
                        const icon = iconByID.get(row.id);
                        return (
                            <li key={row.id} className="flex flex-wrap items-center justify-between gap-3 px-4 py-3 text-sm">
                                <div className="flex min-w-0 items-center gap-3">
                                    <ProxyListIcon installedPath={icon?.installedPath} kind={icon?.kind || row.kind} />
                                    <div className="min-w-0">
                                        <p className="font-medium text-primary">{row.name || row.id}</p>
                                        <p className="text-tertiary">
                                            {row.kind}://{row.host}:{row.port}
                                            {row.hasPassword ? " · password set" : " · no password"}
                                        </p>
                                        <div className="mt-1.5">
                                            <ProxyHealthBadge health={healthByID.get(row.id)} />
                                        </div>
                                    </div>
                                </div>
                                <div className="flex flex-wrap gap-2">
                                    <Button
                                        color="secondary"
                                        size="sm"
                                        onClick={() => void EditorDialogsService.OpenProxy(row.id, "")}
                                    >
                                        Edit
                                    </Button>
                                    {row.hasPassword && (
                                        <Button color="secondary" size="sm" isLoading={busy} onClick={() => void onClearPassword(row.id)}>
                                            Clear password
                                        </Button>
                                    )}
                                    <Button
                                        color="secondary-destructive"
                                        size="sm"
                                        isLoading={busy}
                                        onClick={() => void onDelete(row.id, row.name || row.id)}
                                    >
                                        Delete
                                    </Button>
                                </div>
                            </li>
                        );
                    })
                )}
            </ul>
        </div>
    );
}

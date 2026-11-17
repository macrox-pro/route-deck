import { useCallback, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { ProxyListIcon } from "@/components/application/proxy-list-icon/proxy-list-icon";
import { SuggestionBanner } from "@/components/application/setup-guide/suggestion-banner";
import { Button } from "@/components/base/buttons/button";
import { Input } from "@/components/base/input/input";
import { NativeSelect } from "@/components/base/select/select-native";
import { useSnapshot } from "@/hooks/use-snapshot";
import { matchProxyGuide } from "@/utils/match-proxy-guide";
import { proxyHealthLabel } from "@/utils/proxy-health-label";
import { ProxiesService, SetupGuidesService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { GuideDTO, ProxyDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

const guidePrefill: Record<string, { name: string; kind: string; host: string; port: string }> = {
    singbox: { name: "sing-box local", kind: "socks5", host: "127.0.0.1", port: "" },
    generic: { name: "Local proxy", kind: "socks5", host: "127.0.0.1", port: "" },
};

const emptyForm = {
    id: "",
    name: "",
    kind: "socks5",
    host: "",
    port: "1080",
    username: "",
    password: "",
};

export function ProxiesPage() {
    const [searchParams] = useSearchParams();
    const { data: snapshot } = useSnapshot();
    const [rows, setRows] = useState<ProxyDTO[]>([]);
    const [guides, setGuides] = useState<GuideDTO[]>([]);
    const [form, setForm] = useState(emptyForm);
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);
    const [guideHint, setGuideHint] = useState<string | null>(null);

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

    useEffect(() => {
        const guideID = searchParams.get("guide");
        if (!guideID) {
            return;
        }
        const draft = guidePrefill[guideID];
        if (!draft) {
            return;
        }
        setForm({
            id: "",
            name: draft.name,
            kind: draft.kind,
            host: draft.host,
            port: draft.port,
            username: "",
            password: "",
        });
        setGuideHint(`Enter the local port from your ${draft.name.replace(/ local$/i, "")} client, then save.`);
    }, [searchParams]);

    const onSave = async () => {
        setBusy(true);
        setError(null);
        try {
            await ProxiesService.Save({
                id: form.id,
                name: form.name,
                kind: form.kind,
                host: form.host,
                port: Number(form.port) || 0,
                username: form.username,
                password: form.password,
            });
            setForm(emptyForm);
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onEdit = (row: ProxyDTO) => {
        setForm({
            id: row.id,
            name: row.name,
            kind: row.kind,
            host: row.host,
            port: String(row.port),
            username: row.username,
            password: "",
        });
    };

    const onDelete = async (id: string) => {
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
                <Button color="secondary" size="md" isLoading={busy} onClick={() => void onRecheck()}>
                    Recheck health
                </Button>
            </div>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            <SuggestionBanner onAdded={() => void reload()} />

            {guideHint && (
                <div className="rounded-lg bg-secondary px-4 py-3 text-sm text-secondary ring-1 ring-secondary" role="status">
                    {guideHint}
                </div>
            )}

            <section className="grid gap-3 rounded-lg p-4 ring-1 ring-secondary md:grid-cols-2">
                <Input
                    label="Name"
                    value={form.name}
                    placeholder="Shadowrocket"
                    onChange={(value) => setForm((f) => ({ ...f, name: value }))}
                />
                <NativeSelect
                    label="Kind"
                    value={form.kind}
                    onChange={(e) => setForm((f) => ({ ...f, kind: e.target.value }))}
                    options={[
                        { label: "SOCKS5", value: "socks5" },
                        { label: "HTTP", value: "http" },
                    ]}
                />
                <Input
                    label="Host"
                    value={form.host}
                    placeholder="127.0.0.1"
                    onChange={(value) => setForm((f) => ({ ...f, host: value }))}
                />
                <Input
                    label="Port"
                    value={form.port}
                    placeholder="1082"
                    onChange={(value) => setForm((f) => ({ ...f, port: value }))}
                />
                <Input
                    label="Username"
                    value={form.username}
                    placeholder="Optional"
                    onChange={(value) => setForm((f) => ({ ...f, username: value }))}
                />
                <Input
                    label="Password"
                    type="password"
                    value={form.password}
                    placeholder={form.id ? "Leave empty to keep" : "Optional"}
                    onChange={(value) => setForm((f) => ({ ...f, password: value }))}
                />
                <div className="flex gap-2 md:col-span-2">
                    <Button color="primary" size="md" isLoading={busy} onClick={() => void onSave()}>
                        {form.id ? "Update proxy" : "Add proxy"}
                    </Button>
                    {form.id && (
                        <Button color="secondary" size="md" onClick={() => setForm(emptyForm)}>
                            Cancel edit
                        </Button>
                    )}
                </div>
            </section>

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
                                            {" · "}
                                            {proxyHealthLabel(healthByID.get(row.id))}
                                        </p>
                                    </div>
                                </div>
                                <div className="flex flex-wrap gap-2">
                                    <Button color="secondary" size="sm" onClick={() => onEdit(row)}>
                                        Edit
                                    </Button>
                                    {row.hasPassword && (
                                        <Button color="secondary" size="sm" isLoading={busy} onClick={() => void onClearPassword(row.id)}>
                                            Clear password
                                        </Button>
                                    )}
                                    <Button color="secondary-destructive" size="sm" isLoading={busy} onClick={() => void onDelete(row.id)}>
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

import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/base/buttons/button";
import { Input } from "@/components/base/input/input";
import { NativeSelect } from "@/components/base/select/select-native";
import { ProxiesService } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { ProxyDTO } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

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
    const [rows, setRows] = useState<ProxyDTO[]>([]);
    const [form, setForm] = useState(emptyForm);
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async () => {
        setError(null);
        try {
            setRows(await ProxiesService.List());
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload]);

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

            <section className="grid gap-3 rounded-lg p-4 ring-1 ring-secondary md:grid-cols-2">
                <Input label="Name" value={form.name} onChange={(value) => setForm((f) => ({ ...f, name: value }))} />
                <NativeSelect
                    label="Kind"
                    value={form.kind}
                    onChange={(e) => setForm((f) => ({ ...f, kind: e.target.value }))}
                    options={[
                        { label: "SOCKS5", value: "socks5" },
                        { label: "HTTP", value: "http" },
                    ]}
                />
                <Input label="Host" value={form.host} onChange={(value) => setForm((f) => ({ ...f, host: value }))} />
                <Input label="Port" value={form.port} onChange={(value) => setForm((f) => ({ ...f, port: value }))} />
                <Input label="Username" value={form.username} onChange={(value) => setForm((f) => ({ ...f, username: value }))} />
                <Input
                    label="Password"
                    type="password"
                    value={form.password}
                    placeholder={form.id ? "Leave empty to keep" : undefined}
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
                    rows.map((row) => (
                        <li key={row.id} className="flex flex-wrap items-center justify-between gap-3 px-4 py-3 text-sm">
                            <div>
                                <p className="font-medium text-primary">{row.name || row.id}</p>
                                <p className="text-tertiary">
                                    {row.kind}://{row.host}:{row.port}
                                    {row.hasPassword ? " · password set" : " · no password"}
                                </p>
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
                    ))
                )}
            </ul>
        </div>
    );
}

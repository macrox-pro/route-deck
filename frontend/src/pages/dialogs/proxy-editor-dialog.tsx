import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { Events } from "@wailsio/runtime";
import { DialogFooter } from "@/components/application/dialog-shell/dialog-footer";
import { DialogShell } from "@/components/application/dialog-shell/dialog-shell";
import { InlineAlert } from "@/components/application/page-layout/inline-alert";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { ButtonGroup, ButtonGroupItem } from "@/components/base/button-group/button-group";
import { Input } from "@/components/base/input/input";
import {
    EditorDialogsService,
    ProxiesService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

const guidePrefill: Record<string, { name: string; kind: string; host: string; port: string }> = {
    singbox: { name: "sing-box local", kind: "socks5", host: "127.0.0.1", port: "" },
    generic: { name: "Local proxy", kind: "socks5", host: "127.0.0.1", port: "" },
};

const emptyForm = {
    name: "",
    kind: "socks5",
    host: "",
    port: "1080",
    username: "",
    password: "",
};

export function ProxyEditorDialog() {
    const { id = "" } = useParams();
    const [searchParams] = useSearchParams();
    const [form, setForm] = useState(emptyForm);
    const [baseline, setBaseline] = useState(emptyForm);
    const [error, setError] = useState<string | null>(null);
    const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
    const [busy, setBusy] = useState(false);
    const [guideHint, setGuideHint] = useState<string | null>(null);
    const [loaded, setLoaded] = useState(false);
    const nameRef = useRef<HTMLInputElement>(null);

    const dirty = useMemo(() => JSON.stringify(form) !== JSON.stringify(baseline), [baseline, form]);

    useEffect(() => {
        void EditorDialogsService.SetDirty("proxy", dirty);
    }, [dirty]);

    useEffect(() => {
        let cancelled = false;
        const load = async () => {
            setLoaded(false);
            setError(null);
            const guideID = searchParams.get("guide") || "";
            if (id) {
                try {
                    const rows = await ProxiesService.List();
                    const row = rows.find((proxy) => proxy.id === id);
                    if (!row) {
                        if (!cancelled) {
                            setError("Proxy not found.");
                            setLoaded(true);
                        }
                        return;
                    }
                    const next = {
                        name: row.name,
                        kind: row.kind,
                        host: row.host,
                        port: String(row.port),
                        username: row.username,
                        password: "",
                    };
                    if (!cancelled) {
                        setForm(next);
                        setBaseline(next);
                        void EditorDialogsService.SetTitle(
                            "proxy",
                            next.name ? `Edit Proxy — ${next.name}` : "Edit Proxy",
                        );
                        setLoaded(true);
                    }
                } catch (err) {
                    if (!cancelled) {
                        setError(err instanceof Error ? err.message : String(err));
                        setLoaded(true);
                    }
                }
                return;
            }
            const draft = guidePrefill[guideID];
            const next = draft
                ? {
                      name: draft.name,
                      kind: draft.kind,
                      host: draft.host,
                      port: draft.port || "1080",
                      username: "",
                      password: "",
                  }
                : emptyForm;
            if (!cancelled) {
                setForm(next);
                setBaseline(emptyForm);
                setGuideHint(
                    draft
                        ? `Enter the local port from your ${draft.name.replace(/ local$/i, "")} client, then save.`
                        : null,
                );
                void EditorDialogsService.SetTitle("proxy", "New Proxy");
                setLoaded(true);
            }
        };
        void load();
        return () => {
            cancelled = true;
        };
    }, [id, searchParams]);

    useEffect(() => {
        if (loaded) {
            nameRef.current?.focus();
        }
    }, [loaded]);

    const onCancel = useCallback(() => {
        void EditorDialogsService.RequestClose("proxy");
    }, []);

    const onSave = useCallback(async (forQuit = false) => {
        const nextErrors: Record<string, string> = {};
        if (!form.name.trim()) {
            nextErrors.name = "Name is required.";
        }
        if (!form.host.trim()) {
            nextErrors.host = "Host is required.";
        }
        if (!form.port.trim() || Number(form.port) <= 0) {
            nextErrors.port = "Port is required.";
        }
        setFieldErrors(nextErrors);
        if (Object.keys(nextErrors).length > 0) {
            if (forQuit) {
                await EditorDialogsService.AckSaveForQuit("proxy", "Proxy form is incomplete.");
            }
            return;
        }
        setBusy(true);
        setError(null);
        try {
            await ProxiesService.Save({
                id,
                name: form.name.trim(),
                kind: form.kind,
                host: form.host.trim(),
                port: Number(form.port) || 0,
                username: form.username,
                password: form.password,
            });
            await EditorDialogsService.SetDirty("proxy", false);
            if (forQuit) {
                await EditorDialogsService.AckSaveForQuit("proxy", "");
            } else {
                await EditorDialogsService.Close("proxy");
            }
        } catch (err) {
            const msg = err instanceof Error ? err.message : String(err);
            setError(msg);
            if (forQuit) {
                await EditorDialogsService.AckSaveForQuit("proxy", msg);
            }
        } finally {
            setBusy(false);
        }
    }, [form, id]);

    useEffect(() => {
        return Events.On("editors.save_for_quit", () => {
            if (!dirty) {
                return;
            }
            void onSave(true);
        });
    }, [dirty, onSave]);

    const saveDisabled = !loaded || !form.name.trim() || (!dirty && !!id) || busy;

    return (
        <DialogShell
            footer={
                <DialogFooter
                    onCancel={onCancel}
                    onSave={() => void onSave()}
                    saveDisabled={saveDisabled}
                    saveLoading={busy}
                />
            }
        >
            {!loaded ? (
                <div className="flex min-h-64 items-center justify-center">
                    <LoadingIndicator label="Loading proxy…" />
                </div>
            ) : (
                <div className="flex flex-col gap-5">
                    {error ? <InlineAlert tone="error">{error}</InlineAlert> : null}
                    {guideHint ? <InlineAlert tone="info">{guideHint}</InlineAlert> : null}
                    <Input
                        ref={nameRef}
                        label="Name"
                        value={form.name}
                        hint={fieldErrors.name}
                        isInvalid={!!fieldErrors.name}
                        placeholder="Shadowrocket"
                        onChange={(value) => setForm((f) => ({ ...f, name: value }))}
                    />
                    <section className="flex flex-col gap-3">
                        <h2 className="text-sm font-semibold text-secondary">Connection</h2>
                        <div className="flex flex-col gap-2">
                            <p className="text-sm font-medium text-secondary">Kind</p>
                            <ButtonGroup
                                size="sm"
                                selectedKeys={new Set([form.kind])}
                                onSelectionChange={(keys) => {
                                    const selected = keys as Set<string> | "all";
                                    if (selected === "all") {
                                        return;
                                    }
                                    const kind = [...selected][0];
                                    if (kind === "socks5" || kind === "http") {
                                        setForm((f) => ({ ...f, kind }));
                                    }
                                }}
                                aria-label="Kind"
                            >
                                <ButtonGroupItem id="socks5">SOCKS5</ButtonGroupItem>
                                <ButtonGroupItem id="http">HTTP</ButtonGroupItem>
                            </ButtonGroup>
                        </div>
                        <div className="flex flex-wrap items-start gap-3">
                            <div className="min-w-0 flex-1">
                                <Input
                                    label="Host"
                                    value={form.host}
                                    hint={fieldErrors.host}
                                    isInvalid={!!fieldErrors.host}
                                    placeholder="127.0.0.1"
                                    className="font-mono"
                                    onChange={(value) => setForm((f) => ({ ...f, host: value }))}
                                />
                            </div>
                            <div className="w-[120px] shrink-0">
                                <Input
                                    label="Port"
                                    value={form.port}
                                    hint={fieldErrors.port}
                                    isInvalid={!!fieldErrors.port}
                                    placeholder="1082"
                                    inputMode="numeric"
                                    className="font-mono"
                                    onChange={(value) => setForm((f) => ({ ...f, port: value }))}
                                />
                            </div>
                        </div>
                    </section>
                    <section className="flex flex-col gap-3">
                        <h2 className="text-sm font-semibold text-secondary">Authentication</h2>
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
                            hint="Saved in Keychain. Leave empty to keep."
                            placeholder={id ? "Leave empty to keep" : "Optional"}
                            onChange={(value) => setForm((f) => ({ ...f, password: value }))}
                        />
                    </section>
                </div>
            )}
        </DialogShell>
    );
}

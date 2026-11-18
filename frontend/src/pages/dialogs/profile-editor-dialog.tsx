import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { DialogFooter } from "@/components/application/dialog-shell/dialog-footer";
import { DialogShell } from "@/components/application/dialog-shell/dialog-shell";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { Button } from "@/components/base/buttons/button";
import { Checkbox } from "@/components/base/checkbox/checkbox";
import { Input } from "@/components/base/input/input";
import { NativeSelect } from "@/components/base/select/select-native";
import { appBundlePath } from "@/utils/app-bundle-path";
import { loadAppIcon } from "@/utils/load-app-icon";
import {
    ApplicationsService,
    EditorDialogsService,
    ProfilesService,
    ProxiesService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type {
    ApplicationDTO,
    HelperGroupDTO,
    PolicyDTO,
    ProcessIdentityDTO,
    ProxyDTO,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

type ActionKind = "" | "proxy" | "direct" | "block";

type ProcessDraft = {
    action: ActionKind;
    proxyId: string;
    protocolScope: string;
    udpLimitation: string;
};

const NONE_DRAFT = (): ProcessDraft => ({
    action: "",
    proxyId: "",
    protocolScope: "BOTH",
    udpLimitation: "Absent",
});

const ACTION_OPTIONS = [
    { label: "None", value: "" },
    { label: "Proxy", value: "proxy" },
    { label: "Direct", value: "direct" },
    { label: "Block", value: "block" },
];

function draftFromPolicy(policy: PolicyDTO | undefined): ProcessDraft {
    if (!policy) {
        return NONE_DRAFT();
    }
    const action = (policy.action || "") as ActionKind;
    if (action !== "proxy" && action !== "direct" && action !== "block") {
        return NONE_DRAFT();
    }
    return {
        action,
        proxyId: policy.proxyId || "",
        protocolScope: policy.protocolScope || "BOTH",
        udpLimitation: policy.udpLimitation || "Absent",
    };
}

function policiesFromDrafts(drafts: Record<string, ProcessDraft>): PolicyDTO[] {
    const out: PolicyDTO[] = [];
    for (const [processId, draft] of Object.entries(drafts)) {
        if (!processId || !draft.action) {
            continue;
        }
        out.push({
            processId,
            action: draft.action,
            proxyId: draft.action === "proxy" ? draft.proxyId : "",
            protocolScope: draft.action === "proxy" ? draft.protocolScope || "BOTH" : "",
            udpLimitation: draft.udpLimitation || "Absent",
        });
    }
    return out;
}

function draftsFromPolicies(policies: PolicyDTO[]): Record<string, ProcessDraft> {
    const drafts: Record<string, ProcessDraft> = {};
    for (const policy of policies ?? []) {
        if (!policy.processId || drafts[policy.processId]) {
            continue;
        }
        drafts[policy.processId] = draftFromPolicy(policy);
    }
    return drafts;
}

function actionSummary(draft: ProcessDraft, proxies: ProxyDTO[]): string {
    if (!draft.action) {
        return "None";
    }
    if (draft.action === "proxy") {
        const proxy = proxies.find((row) => row.id === draft.proxyId);
        return proxy ? `Proxy · ${proxy.name || proxy.id}` : "Proxy";
    }
    if (draft.action === "direct") {
        return "Direct";
    }
    if (draft.action === "block") {
        return "Block";
    }
    return draft.action;
}

function useAppIcon(path: string | undefined) {
    const [src, setSrc] = useState<string | null>(null);
    useEffect(() => {
        if (!path) {
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
    }, [path]);
    return src;
}

function AppIconSlot({ path }: { path: string | undefined }) {
    const icon = useAppIcon(path);
    return (
        <div className="flex size-8 shrink-0 items-center justify-center overflow-hidden rounded-md bg-secondary ring-1 ring-secondary">
            {icon ? <img src={icon} alt="" className="size-8" /> : null}
        </div>
    );
}

function ProcessActionControls({
    draft,
    proxies,
    onChange,
}: {
    draft: ProcessDraft;
    proxies: ProxyDTO[];
    onChange: (patch: Partial<ProcessDraft>) => void;
}) {
    return (
        <div className="flex min-w-0 flex-col gap-2 sm:items-end">
            <NativeSelect
                className="min-w-0 sm:w-40"
                aria-label="Action"
                value={draft.action}
                onChange={(e) => onChange({ action: e.target.value as ActionKind })}
                options={ACTION_OPTIONS}
            />
            {draft.action === "proxy" ? (
                <>
                    <NativeSelect
                        className="min-w-0 sm:w-40"
                        aria-label="Proxy"
                        value={draft.proxyId}
                        onChange={(e) => onChange({ proxyId: e.target.value })}
                        options={[
                            { label: "Select proxy…", value: "" },
                            ...proxies.map((proxy) => ({ label: proxy.name || proxy.id, value: proxy.id })),
                        ]}
                    />
                    <NativeSelect
                        className="min-w-0 sm:w-40"
                        aria-label="Protocol scope"
                        value={draft.protocolScope || "BOTH"}
                        onChange={(e) => onChange({ protocolScope: e.target.value })}
                        options={[
                            { label: "TCP", value: "TCP" },
                            { label: "UDP", value: "UDP" },
                            { label: "BOTH", value: "BOTH" },
                        ]}
                    />
                    <Checkbox
                        label="Direct UDP limitation"
                        isSelected={draft.udpLimitation === "DirectUDP"}
                        onChange={(isSelected) => onChange({ udpLimitation: isSelected ? "DirectUDP" : "Absent" })}
                    />
                </>
            ) : null}
        </div>
    );
}

type AppRowState = {
    app: ApplicationDTO;
    main?: ProcessIdentityDTO;
    groups: HelperGroupDTO[];
    iconPath?: string;
};

export function ProfileEditorDialog() {
    const { id = "" } = useParams();
    const [name, setName] = useState("");
    const [baselineName, setBaselineName] = useState("");
    const [drafts, setDrafts] = useState<Record<string, ProcessDraft>>({});
    const [baselineDrafts, setBaselineDrafts] = useState<Record<string, ProcessDraft>>({});
    const [proxies, setProxies] = useState<ProxyDTO[]>([]);
    const [rows, setRows] = useState<AppRowState[]>([]);
    const [search, setSearch] = useState("");
    const [configuredOnly, setConfiguredOnly] = useState(false);
    const [expanded, setExpanded] = useState<Record<string, boolean>>({});
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);
    const [loaded, setLoaded] = useState(false);
    const nameRef = useRef<HTMLInputElement>(null);

    const dirty = useMemo(() => {
        if (name !== baselineName) {
            return true;
        }
        return JSON.stringify(drafts) !== JSON.stringify(baselineDrafts);
    }, [baselineDrafts, baselineName, drafts, name]);

    useEffect(() => {
        void EditorDialogsService.SetDirty("profile", dirty);
    }, [dirty]);

    useEffect(() => {
        let cancelled = false;
        const load = async () => {
            setLoaded(false);
            setError(null);
            try {
                const [proxyRows, applicationRows, processRows, profiles] = await Promise.all([
                    ProxiesService.List(),
                    ApplicationsService.List(),
                    ApplicationsService.ListProcessIdentities(),
                    id ? ProfilesService.List() : Promise.resolve([]),
                ]);
                if (cancelled) {
                    return;
                }
                setProxies(proxyRows);
                const groupResults = await Promise.all(
                    applicationRows.map(async (app) => {
                        try {
                            return await ApplicationsService.ListHelperGroups(app.id);
                        } catch {
                            return [] as HelperGroupDTO[];
                        }
                    }),
                );
                if (cancelled) {
                    return;
                }
                const nextRows: AppRowState[] = applicationRows.map((app, index) => {
                    const processes = processRows.filter((process) => process.applicationId === app.id);
                    const mains = processes.filter((process) => process.bundleId === app.bundleId);
                    const main = [...mains].sort((a, b) => a.id.localeCompare(b.id))[0];
                    const groups = (groupResults[index] ?? []).filter((group) => group.bundleId !== app.bundleId);
                    return {
                        app,
                        main,
                        groups,
                        iconPath: appBundlePath(main?.path) || appBundlePath(processes[0]?.path),
                    };
                });
                setRows(nextRows);
                if (id) {
                    const profile = profiles.find((row) => row.id === id);
                    if (!profile) {
                        setError("Profile not found.");
                        setLoaded(true);
                        return;
                    }
                    const nextDrafts = draftsFromPolicies(profile.policies ?? []);
                    setName(profile.name);
                    setBaselineName(profile.name);
                    setDrafts(nextDrafts);
                    setBaselineDrafts(nextDrafts);
                    void EditorDialogsService.SetTitle(
                        "profile",
                        profile.name ? `Edit Profile — ${profile.name}` : "Edit Profile",
                    );
                } else {
                    setName("");
                    setBaselineName("");
                    setDrafts({});
                    setBaselineDrafts({});
                    void EditorDialogsService.SetTitle("profile", "New Profile");
                }
                setLoaded(true);
            } catch (err) {
                if (!cancelled) {
                    setError(err instanceof Error ? err.message : String(err));
                    setLoaded(true);
                }
            }
        };
        void load();
        return () => {
            cancelled = true;
        };
    }, [id]);

    useEffect(() => {
        if (loaded) {
            nameRef.current?.focus();
        }
    }, [loaded]);

    const configuredCount = useMemo(() => {
        let count = 0;
        for (const row of rows) {
            const mainDraft = row.main ? drafts[row.main.id] : undefined;
            const helperConfigured = row.groups.some((group) => {
                const processId = group.processIds?.[0];
                return processId ? !!drafts[processId]?.action : false;
            });
            if (mainDraft?.action || helperConfigured) {
                count += 1;
            }
        }
        return count;
    }, [drafts, rows]);

    const visibleRows = useMemo(() => {
        const q = search.trim().toLowerCase();
        return rows.filter((row) => {
            if (q) {
                const hay = `${row.app.displayName} ${row.app.bundleId}`.toLowerCase();
                if (!hay.includes(q)) {
                    return false;
                }
            }
            if (!configuredOnly) {
                return true;
            }
            const mainDraft = row.main ? drafts[row.main.id] : undefined;
            const helperConfigured = row.groups.some((group) => {
                const processId = group.processIds?.[0];
                return processId ? !!drafts[processId]?.action : false;
            });
            return !!(mainDraft?.action || helperConfigured);
        });
    }, [configuredOnly, drafts, rows, search]);

    const onDraftChange = (processId: string, patch: Partial<ProcessDraft>) => {
        setDrafts((prev) => {
            const current = prev[processId] ?? NONE_DRAFT();
            const next = { ...current, ...patch };
            if (!next.action) {
                const { [processId]: _, ...rest } = prev;
                return rest;
            }
            return { ...prev, [processId]: next };
        });
    };

    const onApplyToHelpers = (row: AppRowState) => {
        if (!row.main || row.groups.length === 0) {
            return;
        }
        const source = drafts[row.main.id] ?? NONE_DRAFT();
        setDrafts((prev) => {
            const next = { ...prev };
            for (const group of row.groups) {
                const processId = group.processIds?.[0];
                if (!processId) {
                    continue;
                }
                if (!source.action) {
                    delete next[processId];
                    continue;
                }
                next[processId] = {
                    action: source.action,
                    proxyId: source.proxyId,
                    protocolScope: source.protocolScope,
                    udpLimitation: source.udpLimitation,
                };
            }
            return next;
        });
    };

    const onCancel = useCallback(() => {
        void EditorDialogsService.RequestClose("profile");
    }, []);

    const onSave = useCallback(async () => {
        if (!name.trim()) {
            setError("Profile name is required.");
            return;
        }
        setBusy(true);
        setError(null);
        try {
            await ProfilesService.Save({
                id,
                name: name.trim(),
                policies: policiesFromDrafts(drafts),
            });
            await EditorDialogsService.SetDirty("profile", false);
            await EditorDialogsService.Close("profile");
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    }, [drafts, id, name]);

    const saveDisabled = !loaded || !name.trim() || !dirty || busy;

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
                    <LoadingIndicator label="Loading profile…" />
                </div>
            ) : (
            <div className="flex flex-col gap-4">
                {error && (
                    <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                        {error}
                    </div>
                )}
                <Input
                    ref={nameRef}
                    label="Profile name"
                    value={name}
                    placeholder="Work"
                    onChange={setName}
                />
                <div className="flex flex-wrap items-center gap-3">
                    <Input
                        className="min-w-48 flex-1"
                        aria-label="Search applications"
                        placeholder="Search applications"
                        value={search}
                        onChange={setSearch}
                    />
                    <Checkbox label="Configured only" isSelected={configuredOnly} onChange={setConfiguredOnly} />
                    <p className="text-sm text-tertiary">
                        {configuredCount} of {rows.length}
                    </p>
                </div>
                {rows.length === 0 ? (
                    <p className="text-sm text-tertiary">
                        No applications yet. Add one on the Applications screen.
                    </p>
                ) : (
                    <ul className="divide-y divide-secondary overflow-hidden rounded-lg ring-1 ring-secondary">
                        {visibleRows.map((row) => {
                            const mainDraft = row.main ? (drafts[row.main.id] ?? NONE_DRAFT()) : NONE_DRAFT();
                            const isOpen = !!expanded[row.app.id];
                            const helperSummary =
                                row.groups.length > 0
                                    ? `Helpers · ${row.groups.length} group${row.groups.length === 1 ? "" : "s"} · ${actionSummary(mainDraft, proxies)}`
                                    : null;
                            return (
                                <li key={row.app.id} className="flex min-w-0 flex-col gap-3 px-4 py-3">
                                    <div className="flex min-w-0 flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                                        <div className="flex min-w-0 items-center gap-3">
                                            <AppIconSlot path={row.iconPath} />
                                            <div className="min-w-0">
                                                <p className="truncate text-sm font-semibold text-primary">
                                                    {row.app.displayName}
                                                </p>
                                                <p className="truncate text-xs text-tertiary">{row.app.bundleId}</p>
                                                {row.main ? (
                                                    <p className="truncate text-xs text-tertiary">
                                                        {row.main.executableName}
                                                    </p>
                                                ) : (
                                                    <p className="text-xs text-tertiary">No main process in catalog</p>
                                                )}
                                            </div>
                                        </div>
                                        {row.main ? (
                                            <ProcessActionControls
                                                draft={mainDraft}
                                                proxies={proxies}
                                                onChange={(patch) => onDraftChange(row.main!.id, patch)}
                                            />
                                        ) : null}
                                    </div>
                                    {row.groups.length > 0 ? (
                                        <div className="ml-4 border-l border-secondary pl-4">
                                            <button
                                                type="button"
                                                className="text-left text-sm text-secondary"
                                                onClick={() =>
                                                    setExpanded((prev) => ({
                                                        ...prev,
                                                        [row.app.id]: !prev[row.app.id],
                                                    }))
                                                }
                                            >
                                                {helperSummary} {isOpen ? "▾" : "▸"}
                                            </button>
                                            {isOpen ? (
                                                <ul className="mt-2 flex flex-col gap-2">
                                                    {row.groups.map((group) => {
                                                        const processId = group.processIds?.[0] || "";
                                                        const draft = processId
                                                            ? (drafts[processId] ?? NONE_DRAFT())
                                                            : NONE_DRAFT();
                                                        const count = group.executableNames?.length || 1;
                                                        return (
                                                            <li
                                                                key={group.bundleId}
                                                                className="flex min-w-0 flex-col gap-2 sm:flex-row sm:items-start sm:justify-between"
                                                            >
                                                                <div className="min-w-0 pl-2">
                                                                    <p className="truncate text-sm text-primary">
                                                                        {group.displayName || group.bundleId}
                                                                        {count > 1 ? ` · ${count} processes` : ""}
                                                                    </p>
                                                                    <p className="truncate text-xs text-tertiary">
                                                                        {group.bundleId}
                                                                    </p>
                                                                </div>
                                                                {processId ? (
                                                                    <ProcessActionControls
                                                                        draft={draft}
                                                                        proxies={proxies}
                                                                        onChange={(patch) =>
                                                                            onDraftChange(processId, patch)
                                                                        }
                                                                    />
                                                                ) : (
                                                                    <p className="text-xs text-tertiary">
                                                                        Confirm helpers first
                                                                    </p>
                                                                )}
                                                            </li>
                                                        );
                                                    })}
                                                    <li className="flex flex-wrap gap-2">
                                                        <Button
                                                            color="secondary"
                                                            size="sm"
                                                            isDisabled={!row.main}
                                                            onClick={() => onApplyToHelpers(row)}
                                                        >
                                                            Apply app action to helpers
                                                        </Button>
                                                        <Button
                                                            color="secondary"
                                                            size="sm"
                                                            onClick={() =>
                                                                void EditorDialogsService.OpenHelpers(row.app.id)
                                                            }
                                                        >
                                                            Review helpers
                                                        </Button>
                                                    </li>
                                                </ul>
                                            ) : null}
                                        </div>
                                    ) : (
                                        <div className="ml-4">
                                            <Button
                                                color="secondary"
                                                size="sm"
                                                onClick={() => void EditorDialogsService.OpenHelpers(row.app.id)}
                                            >
                                                Review helpers
                                            </Button>
                                        </div>
                                    )}
                                </li>
                            );
                        })}
                    </ul>
                )}
            </div>
            )}
        </DialogShell>
    );
}

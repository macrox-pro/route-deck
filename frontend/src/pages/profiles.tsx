import { useCallback, useEffect, useMemo, useState } from "react";
import { ProfileListMeta } from "@/components/application/profile-list-meta/profile-list-meta";
import { Button } from "@/components/base/buttons/button";
import { Checkbox } from "@/components/base/checkbox/checkbox";
import { Input } from "@/components/base/input/input";
import { NativeSelect } from "@/components/base/select/select-native";
import { useSnapshot } from "@/hooks/use-snapshot";
import { appBundlePath } from "@/utils/app-bundle-path";
import { loadAppIcon } from "@/utils/load-app-icon";
import {
    ApplicationsService,
    ProfilesService,
    ProxiesService,
    SettingsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type {
    ApplicationDTO,
    PolicyDTO,
    ProcessIdentityDTO,
    ProfileDTO,
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

/**
 * ProxyBridge matches bundle id / executable name only (path never matches).
 * Collapse catalog rows that share an application + bundle id onto one representative.
 */
function representativeByProcessID(processes: ProcessIdentityDTO[]): Map<string, string> {
    const groups = new Map<string, ProcessIdentityDTO[]>();
    for (const process of processes) {
        const groupKey = `${process.applicationId}::${process.bundleId || process.id}`;
        const list = groups.get(groupKey) ?? [];
        list.push(process);
        groups.set(groupKey, list);
    }
    const repFor = new Map<string, string>();
    for (const members of groups.values()) {
        const sorted = [...members].sort((a, b) => a.id.localeCompare(b.id));
        const representativeID = sorted[0].id;
        for (const member of sorted) {
            repFor.set(member.id, representativeID);
        }
    }
    return repFor;
}

function draftsFromPolicies(
    policies: PolicyDTO[],
    processes: ProcessIdentityDTO[],
): Record<string, ProcessDraft> {
    const repFor = representativeByProcessID(processes);
    const drafts: Record<string, ProcessDraft> = {};
    for (const policy of policies ?? []) {
        if (!policy.processId) {
            continue;
        }
        const repID = repFor.get(policy.processId) ?? policy.processId;
        if (!drafts[repID]) {
            drafts[repID] = draftFromPolicy(policy);
        }
    }
    return drafts;
}

type ProcessGroup = {
    key: string;
    representative: ProcessIdentityDTO;
    copyCount: number;
};

function groupProcessesByBundle(processes: ProcessIdentityDTO[]): ProcessGroup[] {
    const byBundle = new Map<string, ProcessIdentityDTO[]>();
    for (const process of processes) {
        const key = process.bundleId || `id:${process.id}`;
        const list = byBundle.get(key) ?? [];
        list.push(process);
        byBundle.set(key, list);
    }
    const groups: ProcessGroup[] = [];
    for (const [key, members] of byBundle) {
        const sorted = [...members].sort((a, b) => a.id.localeCompare(b.id));
        groups.push({
            key,
            representative: sorted[0],
            copyCount: sorted.length,
        });
    }
    return groups;
}

function splitAppProcesses(app: ApplicationDTO, processes: ProcessIdentityDTO[]) {
    const rows = processes.filter((process) => process.applicationId === app.id);
    const mains = rows.filter((process) => process.bundleId === app.bundleId);
    const helpers = rows.filter((process) => process.bundleId !== app.bundleId);
    const mainSorted = [...mains].sort((a, b) => a.id.localeCompare(b.id));
    const main = mainSorted[0];
    return {
        main,
        helperGroups: groupProcessesByBundle(helpers),
        iconPath: appBundlePath(main?.path) || appBundlePath(rows[0]?.path),
    };
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

function ApplicationPolicyRow({
    app,
    processes,
    drafts,
    proxies,
    onDraftChange,
    onApplyToHelpers,
}: {
    app: ApplicationDTO;
    processes: ProcessIdentityDTO[];
    drafts: Record<string, ProcessDraft>;
    proxies: ProxyDTO[];
    onDraftChange: (processId: string, patch: Partial<ProcessDraft>) => void;
    onApplyToHelpers: (appId: string) => void;
}) {
    const { main, helperGroups, iconPath } = useMemo(() => splitAppProcesses(app, processes), [app, processes]);
    const mainDraft = main ? (drafts[main.id] ?? NONE_DRAFT()) : NONE_DRAFT();

    return (
        <li className="flex min-w-0 flex-col gap-3 px-4 py-3">
            <div className="flex min-w-0 flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div className="flex min-w-0 items-center gap-3">
                    <AppIconSlot path={iconPath} />
                    <div className="min-w-0">
                        <p className="truncate text-sm font-semibold text-primary">{app.displayName}</p>
                        <p className="truncate text-xs text-tertiary">{app.bundleId}</p>
                        {main ? (
                            <p className="truncate text-xs text-tertiary">{main.executableName}</p>
                        ) : (
                            <p className="text-xs text-tertiary">No main process in catalog</p>
                        )}
                    </div>
                </div>
                {main ? (
                    <ProcessActionControls
                        draft={mainDraft}
                        proxies={proxies}
                        onChange={(patch) => onDraftChange(main.id, patch)}
                    />
                ) : null}
            </div>

            {helperGroups.length > 0 ? (
                <ul className="ml-4 flex flex-col gap-2 border-l border-secondary pl-4">
                    {helperGroups.map((group) => {
                        const helper = group.representative;
                        const draft = drafts[helper.id] ?? NONE_DRAFT();
                        const copyNote = group.copyCount > 1 ? ` · ${group.copyCount} copies` : "";
                        return (
                            <li
                                key={group.key}
                                className="flex min-w-0 flex-col gap-2 sm:flex-row sm:items-start sm:justify-between"
                            >
                                <div className="min-w-0 pl-2">
                                    <p className="truncate text-sm text-primary">
                                        {helper.executableName || helper.bundleId}
                                        {copyNote}
                                    </p>
                                    <p className="truncate text-xs text-tertiary">{helper.bundleId}</p>
                                </div>
                                <ProcessActionControls
                                    draft={draft}
                                    proxies={proxies}
                                    onChange={(patch) => onDraftChange(helper.id, patch)}
                                />
                            </li>
                        );
                    })}
                    <li>
                        <Button
                            color="secondary"
                            size="sm"
                            isDisabled={!main}
                            onClick={() => onApplyToHelpers(app.id)}
                        >
                            Apply app action to helpers
                        </Button>
                    </li>
                </ul>
            ) : null}
        </li>
    );
}

export function ProfilesPage() {
    const { data: snapshot } = useSnapshot();
    const [profiles, setProfiles] = useState<ProfileDTO[]>([]);
    const [proxies, setProxies] = useState<ProxyDTO[]>([]);
    const [applications, setApplications] = useState<ApplicationDTO[]>([]);
    const [processes, setProcesses] = useState<ProcessIdentityDTO[]>([]);
    const [desiredID, setDesiredID] = useState("");
    const [name, setName] = useState("");
    const [editID, setEditID] = useState("");
    const [drafts, setDrafts] = useState<Record<string, ProcessDraft>>({});
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const [profileRows, proxyRows, applicationRows, processRows, settings] = await Promise.all([
                ProfilesService.List(),
                ProxiesService.List(),
                ApplicationsService.List(),
                ApplicationsService.ListProcessIdentities(),
                SettingsService.Get(),
            ]);
            setProfiles(profileRows);
            setProxies(proxyRows);
            setApplications(applicationRows);
            setProcesses(processRows);
            setDesiredID(settings.desiredSet ? settings.desiredId : "");
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload]);

    const resetForm = () => {
        setEditID("");
        setName("");
        setDrafts({});
    };

    const onSave = async () => {
        setBusy(true);
        setError(null);
        try {
            await ProfilesService.Save({
                id: editID,
                name,
                policies: policiesFromDrafts(drafts),
            });
            resetForm();
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onEdit = (profile: ProfileDTO) => {
        setEditID(profile.id);
        setName(profile.name);
        setDrafts(draftsFromPolicies(profile.policies ?? [], processes));
    };

    const onCopy = async (id: string, sourceName: string) => {
        setBusy(true);
        setError(null);
        try {
            await ProfilesService.Copy(id, `${sourceName} copy`);
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
            await ProfilesService.Delete(id);
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onSetDesired = async (id: string) => {
        setBusy(true);
        setError(null);
        try {
            await ProfilesService.SetDesired(id);
            await reload();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

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

    const onApplyToHelpers = (appId: string) => {
        const app = applications.find((row) => row.id === appId);
        if (!app) {
            return;
        }
        const { main, helperGroups } = splitAppProcesses(app, processes);
        if (!main || helperGroups.length === 0) {
            return;
        }
        const source = drafts[main.id] ?? NONE_DRAFT();
        setDrafts((prev) => {
            const next = { ...prev };
            for (const group of helperGroups) {
                const processId = group.representative.id;
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

    return (
        <div className="flex min-w-0 flex-col gap-6">
            <h1 className="text-display-xs font-semibold text-primary">Profiles</h1>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            <section className="flex min-w-0 flex-col gap-4 rounded-lg p-4 ring-1 ring-secondary">
                <Input label="Profile name" value={name} placeholder="Work" onChange={setName} />

                <div className="min-w-0 overflow-hidden rounded-lg ring-1 ring-secondary">
                    {applications.length === 0 ? (
                        <p className="px-4 py-6 text-sm text-tertiary">
                            No applications in the catalog yet. Add apps on the Applications screen first.
                        </p>
                    ) : (
                        <ul className="divide-y divide-secondary">
                            {applications.map((app) => (
                                <ApplicationPolicyRow
                                    key={app.id}
                                    app={app}
                                    processes={processes}
                                    drafts={drafts}
                                    proxies={proxies}
                                    onDraftChange={onDraftChange}
                                    onApplyToHelpers={onApplyToHelpers}
                                />
                            ))}
                        </ul>
                    )}
                </div>

                <div className="flex flex-wrap gap-2">
                    <Button color="primary" size="md" isLoading={busy} onClick={() => void onSave()}>
                        {editID ? "Update profile" : "Create profile"}
                    </Button>
                    {editID && (
                        <Button color="secondary" size="md" onClick={resetForm}>
                            Cancel edit
                        </Button>
                    )}
                </div>
            </section>

            <ul className="min-w-0 divide-y divide-secondary rounded-lg ring-1 ring-secondary">
                {profiles.length === 0 ? (
                    <li className="px-4 py-6 text-sm text-tertiary">No profiles yet.</li>
                ) : (
                    profiles.map((profile) => {
                        const isDesired = desiredID === profile.id;
                        const isApplied = !!snapshot?.appliedConfirmed && snapshot.appliedId === profile.id;
                        return (
                            <li
                                key={profile.id}
                                className="flex min-w-0 flex-col gap-3 px-4 py-3 sm:flex-row sm:items-center sm:justify-between"
                            >
                                <ProfileListMeta
                                    name={profile.name}
                                    policies={profile.policies ?? []}
                                    processes={processes}
                                    isDesired={isDesired}
                                    isApplied={isApplied}
                                />
                                <div className="flex flex-wrap gap-2 sm:shrink-0 sm:justify-end">
                                    <Button
                                        color="secondary"
                                        size="sm"
                                        isLoading={busy}
                                        isDisabled={isDesired}
                                        onClick={() => void onSetDesired(profile.id)}
                                    >
                                        {isDesired ? "Desired" : "Set desired"}
                                    </Button>
                                    <Button color="secondary" size="sm" onClick={() => onEdit(profile)}>
                                        Edit
                                    </Button>
                                    <Button
                                        color="secondary"
                                        size="sm"
                                        isLoading={busy}
                                        onClick={() => void onCopy(profile.id, profile.name)}
                                    >
                                        Copy
                                    </Button>
                                    <Button
                                        color="secondary-destructive"
                                        size="sm"
                                        isLoading={busy}
                                        onClick={() => void onDelete(profile.id)}
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

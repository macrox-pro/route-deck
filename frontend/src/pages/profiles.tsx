import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/base/buttons/button";
import { Checkbox } from "@/components/base/checkbox/checkbox";
import { Input } from "@/components/base/input/input";
import { NativeSelect } from "@/components/base/select/select-native";
import {
    ApplicationsService,
    ProfilesService,
    ProxiesService,
    SettingsService,
} from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { PolicyDTO, ProcessIdentityDTO, ProfileDTO, ProxyDTO } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

type PolicyDraft = PolicyDTO;

const emptyPolicy = (): PolicyDraft => ({
    processId: "",
    action: "direct",
    proxyId: "",
    protocolScope: "BOTH",
    udpLimitation: "Absent",
});

export function ProfilesPage() {
    const [profiles, setProfiles] = useState<ProfileDTO[]>([]);
    const [proxies, setProxies] = useState<ProxyDTO[]>([]);
    const [processes, setProcesses] = useState<ProcessIdentityDTO[]>([]);
    const [desiredID, setDesiredID] = useState("");
    const [name, setName] = useState("");
    const [editID, setEditID] = useState("");
    const [policies, setPolicies] = useState<PolicyDraft[]>([emptyPolicy()]);
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const [profileRows, proxyRows, processRows, settings] = await Promise.all([
                ProfilesService.List(),
                ProxiesService.List(),
                ApplicationsService.ListProcessIdentities(),
                SettingsService.Get(),
            ]);
            setProfiles(profileRows);
            setProxies(proxyRows);
            setProcesses(processRows);
            setDesiredID(settings.desiredSet ? settings.desiredId : "");
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
            await ProfilesService.Save({
                id: editID,
                name,
                policies: policies.map((policy) => ({
                    ...policy,
                    proxyId: policy.action === "proxy" ? policy.proxyId : "",
                    protocolScope: policy.action === "proxy" ? policy.protocolScope || "BOTH" : "",
                    udpLimitation: policy.udpLimitation || "Absent",
                })),
            });
            setEditID("");
            setName("");
            setPolicies([emptyPolicy()]);
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
        setPolicies(profile.policies.length > 0 ? profile.policies.map((p) => ({ ...p })) : [emptyPolicy()]);
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

    const updatePolicy = (index: number, patch: Partial<PolicyDraft>) => {
        setPolicies((rows) => rows.map((row, i) => (i === index ? { ...row, ...patch } : row)));
    };

    return (
        <div className="flex flex-col gap-6">
            <h1 className="text-display-xs font-semibold text-primary">Profiles</h1>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            <section className="flex flex-col gap-4 rounded-lg p-4 ring-1 ring-secondary">
                <Input label="Profile name" value={name} onChange={setName} />
                {policies.map((policy, index) => (
                    <div key={index} className="grid gap-3 rounded-lg bg-secondary p-3 md:grid-cols-2">
                        <NativeSelect
                            label="Process"
                            value={policy.processId}
                            onChange={(e) => updatePolicy(index, { processId: e.target.value })}
                            options={[
                                { label: "Select process…", value: "" },
                                ...processes.map((process) => ({
                                    label: `${process.executableName} (${process.bundleId})`,
                                    value: process.id,
                                })),
                            ]}
                        />
                        <NativeSelect
                            label="Action"
                            value={policy.action}
                            onChange={(e) => updatePolicy(index, { action: e.target.value })}
                            options={[
                                { label: "Proxy", value: "proxy" },
                                { label: "Direct", value: "direct" },
                                { label: "Block", value: "block" },
                            ]}
                        />
                        {policy.action === "proxy" && (
                            <>
                                <NativeSelect
                                    label="Proxy"
                                    value={policy.proxyId}
                                    onChange={(e) => updatePolicy(index, { proxyId: e.target.value })}
                                    options={[
                                        { label: "Select proxy…", value: "" },
                                        ...proxies.map((proxy) => ({ label: proxy.name || proxy.id, value: proxy.id })),
                                    ]}
                                />
                                <NativeSelect
                                    label="Protocol scope"
                                    value={policy.protocolScope || "BOTH"}
                                    onChange={(e) => updatePolicy(index, { protocolScope: e.target.value })}
                                    options={[
                                        { label: "TCP", value: "TCP" },
                                        { label: "UDP", value: "UDP" },
                                        { label: "BOTH", value: "BOTH" },
                                    ]}
                                />
                                <Checkbox
                                    label="Direct UDP limitation"
                                    isSelected={policy.udpLimitation === "DirectUDP"}
                                    onChange={(isSelected) =>
                                        updatePolicy(index, { udpLimitation: isSelected ? "DirectUDP" : "Absent" })
                                    }
                                />
                            </>
                        )}
                        <div className="flex gap-2 md:col-span-2">
                            <Button
                                color="secondary"
                                size="sm"
                                onClick={() => setPolicies((rows) => rows.filter((_, i) => i !== index))}
                                isDisabled={policies.length === 1}
                            >
                                Remove policy
                            </Button>
                        </div>
                    </div>
                ))}
                <div className="flex flex-wrap gap-2">
                    <Button color="secondary" size="md" onClick={() => setPolicies((rows) => [...rows, emptyPolicy()])}>
                        Add policy
                    </Button>
                    <Button color="primary" size="md" isLoading={busy} onClick={() => void onSave()}>
                        {editID ? "Update profile" : "Create profile"}
                    </Button>
                    {editID && (
                        <Button
                            color="secondary"
                            size="md"
                            onClick={() => {
                                setEditID("");
                                setName("");
                                setPolicies([emptyPolicy()]);
                            }}
                        >
                            Cancel edit
                        </Button>
                    )}
                </div>
            </section>

            <ul className="divide-y divide-secondary rounded-lg ring-1 ring-secondary">
                {profiles.length === 0 ? (
                    <li className="px-4 py-6 text-sm text-tertiary">No profiles yet.</li>
                ) : (
                    profiles.map((profile) => (
                        <li key={profile.id} className="flex flex-wrap items-center justify-between gap-3 px-4 py-3 text-sm">
                            <div>
                                <p className="font-medium text-primary">
                                    {profile.name}
                                    {desiredID === profile.id ? " · desired" : ""}
                                </p>
                                <p className="text-tertiary">{profile.policies.length} policies</p>
                            </div>
                            <div className="flex flex-wrap gap-2">
                                <Button color="secondary" size="sm" isLoading={busy} onClick={() => void onSetDesired(profile.id)}>
                                    Set desired
                                </Button>
                                <Button color="secondary" size="sm" onClick={() => onEdit(profile)}>
                                    Edit
                                </Button>
                                <Button color="secondary" size="sm" isLoading={busy} onClick={() => void onCopy(profile.id, profile.name)}>
                                    Copy
                                </Button>
                                <Button color="secondary-destructive" size="sm" isLoading={busy} onClick={() => void onDelete(profile.id)}>
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

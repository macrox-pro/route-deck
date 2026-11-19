import { useCallback, useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import { ProfileListMeta } from "@/components/application/profile-list-meta/profile-list-meta";
import { Button } from "@/components/base/buttons/button";
import { useSnapshot } from "@/hooks/use-snapshot";
import {
    ApplicationsService,
    EditorDialogsService,
    ProfilesService,
    SettingsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type {
    ProcessIdentityDTO,
    ProfileDTO,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function ProfilesPage() {
    const { data: snapshot } = useSnapshot();
    const [profiles, setProfiles] = useState<ProfileDTO[]>([]);
    const [processes, setProcesses] = useState<ProcessIdentityDTO[]>([]);
    const [desiredID, setDesiredID] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const [profileRows, processRows, settings] = await Promise.all([
                ProfilesService.List(),
                ApplicationsService.ListProcessIdentities(),
                SettingsService.Get(),
            ]);
            setProfiles(profileRows);
            setProcesses(processRows);
            setDesiredID(settings.desiredSet ? settings.desiredId : "");
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
            // Native Dialog.Question lives in Go — window.confirm is unreliable in Wails.
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

    return (
        <div className="flex min-w-0 flex-col gap-6">
            <div className="flex items-center justify-between gap-4">
                <h1 className="text-display-xs font-semibold text-primary">Profiles</h1>
                <Button color="primary" size="md" onClick={() => void EditorDialogsService.OpenProfile("")}>
                    New profile
                </Button>
            </div>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

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
                                    <Button
                                        color="secondary"
                                        size="sm"
                                        onClick={() => void EditorDialogsService.OpenProfile(profile.id)}
                                    >
                                        Edit
                                    </Button>
                                    <Button
                                        color="secondary"
                                        size="sm"
                                        isLoading={busy}
                                        onClick={() => void onCopy(profile.id, profile.name)}
                                    >
                                        Duplicate
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

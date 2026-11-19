import { useCallback, useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import { Copy01, Shield01, Trash01 } from "@untitledui/icons";
import { InlineAlert } from "@/components/application/page-layout/inline-alert";
import { ListPanel, ListRow } from "@/components/application/page-layout/list-panel";
import { PageHeader } from "@/components/application/page-layout/page-header";
import { PageShell } from "@/components/application/page-layout/page-shell";
import { ScreenState } from "@/components/application/page-layout/screen-state";
import { ProfileListMeta } from "@/components/application/profile-list-meta/profile-list-meta";
import { Button } from "@/components/base/buttons/button";
import { useSnapshot } from "@/hooks/use-snapshot";
import {
    ApplicationsService,
    EditorDialogsService,
    ProfilesService,
    SettingsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { ProcessIdentityDTO, ProfileDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function ProfilesPage() {
    const { data: snapshot } = useSnapshot();
    const [profiles, setProfiles] = useState<ProfileDTO[]>([]);
    const [processes, setProcesses] = useState<ProcessIdentityDTO[]>([]);
    const [desiredID, setDesiredID] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [busyID, setBusyID] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [appCount, setAppCount] = useState(0);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const [profileRows, processRows, settings, apps] = await Promise.all([
                ProfilesService.List(),
                ApplicationsService.ListProcessIdentities(),
                SettingsService.Get(),
                ApplicationsService.List(),
            ]);
            setProfiles(profileRows);
            setProcesses(processRows);
            setDesiredID(settings.desiredSet ? settings.desiredId : "");
            setAppCount((apps ?? []).length);
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
                title="Profiles"
                description="Decide which apps use a proxy, go direct, or are blocked."
                actions={
                    !loading && profiles.length > 0 ? (
                        <Button color="primary" size="md" onClick={() => void EditorDialogsService.OpenProfile("")}>
                            New profile
                        </Button>
                    ) : undefined
                }
            />

            {error ? <InlineAlert tone="error">{error}</InlineAlert> : null}

            {loading && profiles.length === 0 ? (
                <ScreenState state="loading" title="Loading profiles…" />
            ) : profiles.length === 0 ? (
                <ScreenState
                    state="empty"
                    icon={Shield01}
                    title="No profiles yet"
                    description="Create one, set it as Desired, then Apply on the Dashboard."
                    actions={
                        <>
                            <Button color="primary" size="md" onClick={() => void EditorDialogsService.OpenProfile("")}>
                                New profile
                            </Button>
                            {appCount === 0 ? (
                                <Button color="secondary" size="md" href="/applications">
                                    Add applications
                                </Button>
                            ) : null}
                        </>
                    }
                />
            ) : (
                <ListPanel>
                    {profiles.map((profile) => {
                        const isDesired = desiredID === profile.id;
                        const isApplied = !!snapshot?.appliedConfirmed && snapshot.appliedId === profile.id;
                        const rowBusy = busyID === profile.id;
                        return (
                            <ListRow
                                key={profile.id}
                                title={
                                    <ProfileListMeta
                                        name={profile.name}
                                        policies={profile.policies ?? []}
                                        processes={processes}
                                        isDesired={isDesired}
                                        isApplied={isApplied}
                                    />
                                }
                                actions={
                                    <>
                                        <Button
                                            color="secondary"
                                            size="sm"
                                            isLoading={rowBusy}
                                            isDisabled={isDesired || busyID !== null}
                                            onClick={() => void withBusy(profile.id, () => ProfilesService.SetDesired(profile.id))}
                                        >
                                            {isDesired ? "Desired" : "Set desired"}
                                        </Button>
                                        <Button
                                            color="secondary"
                                            size="sm"
                                            isDisabled={busyID !== null}
                                            onClick={() => void EditorDialogsService.OpenProfile(profile.id)}
                                        >
                                            Edit
                                        </Button>
                                    </>
                                }
                                menu={[
                                    {
                                        id: "duplicate",
                                        label: "Duplicate",
                                        icon: Copy01,
                                        isDisabled: busyID !== null,
                                        onAction: () => void withBusy(profile.id, () => ProfilesService.Copy(profile.id, `${profile.name} copy`)),
                                    },
                                    {
                                        id: "delete",
                                        label: "Delete",
                                        icon: Trash01,
                                        destructive: true,
                                        isDisabled: busyID !== null,
                                        onAction: () => void withBusy(profile.id, () => ProfilesService.Delete(profile.id)),
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

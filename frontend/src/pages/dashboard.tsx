import { useState } from "react";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { ProtectionBadge } from "@/components/application/protection-badge/protection-badge";
import { Button } from "@/components/base/buttons/button";
import { useSnapshot } from "@/hooks/use-snapshot";
import { ProtectionService } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function DashboardPage() {
    const { data, error, loading, reload } = useSnapshot();
    const [actionError, setActionError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);
    const [operationID, setOperationID] = useState<string | null>(null);

    if (loading && !data) {
        return (
            <div className="flex min-h-64 items-center justify-center">
                <LoadingIndicator label="Loading snapshot…" />
            </div>
        );
    }

    const onApply = async () => {
        setBusy(true);
        setActionError(null);
        try {
            const result = await ProtectionService.Apply();
            setOperationID(result.operationId);
            await reload();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onStop = async () => {
        setBusy(true);
        setActionError(null);
        try {
            const result = await ProtectionService.Stop();
            setOperationID(result.operationId);
            await reload();
        } catch (err) {
            setActionError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    return (
        <div className="flex flex-col gap-6">
            <div className="flex flex-wrap items-center justify-between gap-4">
                <div>
                    <h1 className="text-display-xs font-semibold text-primary">Dashboard</h1>
                    <p className="mt-1 text-sm text-tertiary">Protection state from Go. Events trigger a snapshot reload; they are not proof of Active.</p>
                </div>
                <div className="flex flex-wrap gap-2">
                    <Button color="primary" size="md" isLoading={busy} onClick={() => void onApply()}>
                        Apply
                    </Button>
                    <Button color="secondary" size="md" isLoading={busy} onClick={() => void onStop()}>
                        Stop
                    </Button>
                    <Button color="secondary" size="md" isLoading={loading} onClick={() => void reload()}>
                        Refresh
                    </Button>
                </div>
            </div>

            {(error || actionError) && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {actionError || error}
                </div>
            )}

            {data?.databaseError && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="alert">
                    Database: {data.databaseError}
                </div>
            )}

            {data?.bridge === "NotInstalled" && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                    ProxyBridge is not installed. Install the official CLI, or set a custom path in Settings.
                </div>
            )}

            {data?.bridge === "Conflict" && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    Another ProxyBridge CLI is already running. Stop the other CLI, then Apply again. RouteDeck will not signal foreign processes.
                </div>
            )}

            {data?.desiredDiffersFromApplied && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                    Desired profile differs from the applied profile. Press Apply to put the desired profile into force.
                </div>
            )}

            {data?.compatibilityWarning && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                    {data.compatibilityWarning}
                </div>
            )}

            {data?.recheckPending && (
                <div className="rounded-lg bg-secondary px-4 py-3 text-sm text-secondary ring-1 ring-secondary" role="status">
                    Recheck pending after wake or network change. Protection shows Unknown until health and bridge checks finish.
                </div>
            )}

            {operationID && <p className="text-xs text-tertiary">Last operation: {operationID}</p>}

            {data && (
                <>
                    <div className="flex flex-wrap items-center gap-3">
                        <ProtectionBadge protection={data.protection} />
                        <span className="text-sm text-secondary">Bridge: {data.bridge || "Unknown"}</span>
                        {data.version && <span className="text-sm text-tertiary">CLI {data.version}</span>}
                        {data.desiredSet && <span className="text-sm text-tertiary">Desired: {data.desiredId || "—"}</span>}
                        {data.appliedConfirmed && <span className="text-sm text-tertiary">Applied: {data.appliedId || "—"}</span>}
                    </div>

                    <section>
                        <h2 className="mb-2 text-sm font-semibold text-secondary">Proxies</h2>
                        {data.proxies.length === 0 ? (
                            <p className="text-sm text-tertiary">No proxies in the catalog.</p>
                        ) : (
                            <ul className="divide-y divide-secondary rounded-lg ring-1 ring-secondary">
                                {data.proxies.map((proxy) => (
                                    <li key={proxy.id} className="flex items-center justify-between px-4 py-3 text-sm">
                                        <span className="text-primary">{proxy.name || proxy.id}</span>
                                        <span className="text-tertiary">{proxy.health}</span>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </section>
                </>
            )}
        </div>
    );
}

import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { ProtectionBadge } from "@/components/application/protection-badge/protection-badge";
import { Button } from "@/components/base/buttons/button";
import { useSnapshot } from "@/hooks/use-snapshot";

export function DashboardPage() {
    const { data, error, loading, reload } = useSnapshot();

    if (loading && !data) {
        return (
            <div className="flex min-h-64 items-center justify-center">
                <LoadingIndicator label="Loading snapshot…" />
            </div>
        );
    }

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between gap-4">
                <div>
                    <h1 className="text-display-xs font-semibold text-primary">Dashboard</h1>
                    <p className="mt-1 text-sm text-tertiary">Protection state from Go. Events are not the source of truth.</p>
                </div>
                <Button color="secondary" size="md" isLoading={loading} onClick={() => void reload()}>
                    Refresh
                </Button>
            </div>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            {data?.databaseError && (
                <div className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle" role="alert">
                    Database: {data.databaseError}
                </div>
            )}

            {data && (
                <>
                    <div className="flex flex-wrap items-center gap-3">
                        <ProtectionBadge protection={data.protection} />
                        <span className="text-sm text-secondary">Bridge: {data.bridge || "Unknown"}</span>
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

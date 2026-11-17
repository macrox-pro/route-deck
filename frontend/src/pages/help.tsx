import { useCallback, useEffect, useMemo, useState } from "react";
import { SetupGuideCard } from "@/components/application/setup-guide/setup-guide-card";
import { Button } from "@/components/base/buttons/button";
import { SetupGuidesService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { GuideDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function HelpPage() {
    const [guides, setGuides] = useState<GuideDTO[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async (refresh: boolean) => {
        setBusy(true);
        setError(null);
        try {
            setGuides(await SetupGuidesService.List(refresh));
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    }, []);

    useEffect(() => {
        void reload(true);
    }, [reload]);

    const bridgeGuides = useMemo(() => guides.filter((g) => g.kind === "bridge"), [guides]);
    const clientGuides = useMemo(() => guides.filter((g) => g.kind !== "bridge"), [guides]);

    return (
        <div className="flex flex-col gap-8">
            <div className="flex flex-wrap items-center justify-between gap-4">
                <div>
                    <h1 className="text-display-xs font-semibold text-primary">Setup guides</h1>
                    <p className="mt-1 max-w-2xl text-sm text-tertiary">
                        RouteDeck needs two things: ProxyBridge to enforce per-app rules, and a local SOCKS5 or HTTP proxy from your VPN app.
                        “Port open” only means something answered on 127.0.0.1 — not that the proxy is healthy or that protection is Active.
                    </p>
                </div>
                <Button color="secondary" size="md" isLoading={busy} onClick={() => void reload(true)}>
                    Refresh
                </Button>
            </div>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            <section className="flex flex-col gap-4">
                <div>
                    <h2 className="text-lg font-semibold text-primary">ProxyBridge</h2>
                    <p className="mt-1 text-sm text-tertiary">
                        Required. Install the official ProxyBridge app and approve its network extension. RouteDeck drives the CLI; it does not
                        replace ProxyBridge.
                    </p>
                </div>
                {bridgeGuides.map((guide) => (
                    <SetupGuideCard key={guide.id} guide={guide} busy={busy} onChanged={() => void reload(false)} />
                ))}
            </section>

            <section className="flex flex-col gap-4">
                <div>
                    <h2 className="text-lg font-semibold text-primary">Local proxy apps</h2>
                    <p className="mt-1 text-sm text-tertiary">
                        Your VPN or proxy app on this Mac (Shadowrocket, v2rayN, Surge, and similar). Expose a local SOCKS5 or HTTP port, then add
                        it in RouteDeck. Prefer proxy-only mode — turn TUN / system VPN / Enhanced Mode off so RouteDeck can apply per-app
                        policies.
                    </p>
                </div>
                {clientGuides.map((guide) => (
                    <SetupGuideCard key={guide.id} guide={guide} busy={busy} onChanged={() => void reload(false)} />
                ))}
            </section>
        </div>
    );
}

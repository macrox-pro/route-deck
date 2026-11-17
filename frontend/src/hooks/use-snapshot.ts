import { useCallback, useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import { SnapshotService } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { SnapshotDTO } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

const REFRESH_EVENTS = [
    "protection.state.changed",
    "proxy.health.changed",
    "proxybridge.state.changed",
    "operation.progress",
    "operation.failed",
] as const;

export function useSnapshot() {
    const [data, setData] = useState<SnapshotDTO | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    const reload = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const snapshot = await SnapshotService.Get();
            setData(snapshot);
        } catch (err) {
            setData(null);
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload]);

    useEffect(() => {
        const unsubs = REFRESH_EVENTS.map((name) =>
            Events.On(name, () => {
                void reload();
            }),
        );
        return () => {
            for (const unsub of unsubs) {
                unsub();
            }
        };
    }, [reload]);

    return { data, error, loading, reload };
}

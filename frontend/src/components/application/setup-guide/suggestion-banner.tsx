import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/base/buttons/button";
import { SetupGuidesService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { SuggestionDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export function SuggestionBanner({ onAdded }: { onAdded?: () => void }) {
    const [items, setItems] = useState<SuggestionDTO[]>([]);
    const [busy, setBusy] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const reload = useCallback(async () => {
        try {
            setItems(await SetupGuidesService.Suggestions());
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload]);

    if (items.length === 0 && !error) {
        return null;
    }

    return (
        <div className="flex flex-col gap-2">
            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}
            {items.map((item) => (
                <div
                    key={`${item.guideId}:${item.port}`}
                    className="flex flex-wrap items-center justify-between gap-3 rounded-lg bg-brand-primary px-4 py-3 text-sm text-secondary ring-1 ring-brand"
                    role="status"
                >
                    <p>
                        Found {item.title} listening on {item.host}:{item.port}. Add it to Proxies?
                    </p>
                    <div className="flex gap-2">
                        <Button
                            color="primary"
                            size="sm"
                            isLoading={busy}
                            onClick={() => {
                                setBusy(true);
                                setError(null);
                                void SetupGuidesService.AddProxy(item.guideId)
                                    .then(async () => {
                                        await reload();
                                        onAdded?.();
                                    })
                                    .catch((err) => {
                                        setError(err instanceof Error ? err.message : String(err));
                                    })
                                    .finally(() => setBusy(false));
                            }}
                        >
                            Add
                        </Button>
                        <Button
                            color="secondary"
                            size="sm"
                            isDisabled={busy}
                            onClick={() => {
                                setBusy(true);
                                setError(null);
                                void SetupGuidesService.DismissSuggestion(item.guideId, item.port)
                                    .then(() => reload())
                                    .catch((err) => {
                                        setError(err instanceof Error ? err.message : String(err));
                                    })
                                    .finally(() => setBusy(false));
                            }}
                        >
                            Dismiss
                        </Button>
                    </div>
                </div>
            ))}
        </div>
    );
}

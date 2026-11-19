import { useCallback, useEffect, useState } from "react";
import { Browser } from "@wailsio/runtime";
import { Globe01 } from "@untitledui/icons";
import { Button } from "@/components/base/buttons/button";
import { RouteDeckMark } from "@/components/foundations/logo/route-deck-mark";
import {
    SettingsService,
    UtilityWindowsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { SettingsDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

const MACROX_WEBSITE_URL = "https://macrox.pro";
const GITHUB_REPO_URL = "https://github.com/macrox-pro/route-deck";

function GitHubIcon({ className, ...props }: { className?: string } & Record<string, unknown>) {
    return (
        <svg viewBox="0 0 16 16" fill="currentColor" aria-hidden="true" className={className} {...props}>
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8Z" />
        </svg>
    );
}

export function AboutDialog() {
    const [data, setData] = useState<SettingsDTO | null>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let cancelled = false;
        void SettingsService.Get()
            .then((settings) => {
                if (!cancelled) {
                    setData(settings);
                }
            })
            .catch((err: unknown) => {
                if (!cancelled) {
                    setError(err instanceof Error ? err.message : String(err));
                }
            });
        return () => {
            cancelled = true;
        };
    }, []);

    const onClose = useCallback(() => {
        void UtilityWindowsService.Close("about");
    }, []);

    const onOpenExternal = async (url: string) => {
        try {
            await Browser.OpenURL(url);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    };

    return (
        <div className="flex h-full flex-col bg-primary text-primary">
            <div className="flex flex-1 flex-col items-center gap-4 px-8 py-10">
                <RouteDeckMark className="size-16" />
                <div className="text-center">
                    <h1 className="text-display-xs font-semibold text-primary">RouteDeck</h1>
                    <p className="mt-1 text-sm text-tertiary">MACROX RouteDeck</p>
                </div>
                {error ? (
                    <p className="text-sm text-error-primary" role="alert">
                        {error}
                    </p>
                ) : null}
                <dl className="w-full max-w-xs divide-y divide-secondary rounded-xl bg-secondary_subtle px-4 ring-1 ring-secondary">
                    <div className="flex items-baseline justify-between gap-4 py-2.5">
                        <dt className="text-sm text-tertiary">Version</dt>
                        <dd className="text-sm font-medium text-secondary">{data?.appVersionLabel || "…"}</dd>
                    </div>
                    <div className="flex items-baseline justify-between gap-4 py-2.5">
                        <dt className="text-sm text-tertiary">Commit</dt>
                        <dd className="font-mono text-xs font-medium text-secondary">{data?.appCommit || "—"}</dd>
                    </div>
                </dl>
                <div className="flex flex-wrap justify-center gap-x-5 gap-y-2">
                    <Button color="link-color" size="sm" iconLeading={Globe01} onClick={() => void onOpenExternal(MACROX_WEBSITE_URL)}>
                        MACROX
                    </Button>
                    <Button color="link-color" size="sm" iconLeading={GitHubIcon} onClick={() => void onOpenExternal(GITHUB_REPO_URL)}>
                        GitHub
                    </Button>
                </div>
            </div>
            <div className="flex justify-end border-t border-secondary px-6 py-4">
                <Button color="primary" size="md" onClick={onClose}>
                    Done
                </Button>
            </div>
        </div>
    );
}

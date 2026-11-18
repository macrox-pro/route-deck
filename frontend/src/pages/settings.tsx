import type { ReactNode } from "react";
import { useCallback, useEffect, useState } from "react";
import { Browser } from "@wailsio/runtime";
import { Globe01, Monitor01, Moon01, Sun } from "@untitledui/icons";
import { Button } from "@/components/base/buttons/button";
import { Input } from "@/components/base/input/input";
import { Select } from "@/components/base/select/select";
import { Toggle } from "@/components/base/toggle/toggle";
import { useTheme } from "@/providers/theme-provider";
import { SettingsService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { SettingsDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

const MACROX_WEBSITE_URL = "https://macrox.pro";
const GITHUB_REPO_URL = "https://github.com/macrox-pro/route-deck";

const THEME_OPTIONS = [
    {
        id: "system" as const,
        label: "System",
        supportingText: "Match macOS appearance",
        icon: Monitor01,
    },
    {
        id: "light" as const,
        label: "Light",
        supportingText: "Always use light theme",
        icon: Sun,
    },
    {
        id: "dark" as const,
        label: "Dark",
        supportingText: "Always use dark theme",
        icon: Moon01,
    },
];

function SettingsSection({
    title,
    description,
    children,
}: {
    title: string;
    description?: string;
    children: ReactNode;
}) {
    return (
        <section className="grid gap-4 border-b border-secondary py-8 last:border-b-0 md:grid-cols-[minmax(0,14rem)_minmax(0,1fr)] md:gap-10">
            <div className="min-w-0">
                <h2 className="text-md font-semibold text-primary">{title}</h2>
                {description ? <p className="mt-1 text-sm text-tertiary">{description}</p> : null}
            </div>
            <div className="min-w-0">{children}</div>
        </section>
    );
}

function SettingsRow({
    label,
    hint,
    children,
}: {
    label: string;
    hint?: string;
    children: ReactNode;
}) {
    return (
        <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between sm:gap-6">
            <div className="min-w-0 sm:max-w-xs">
                <p className="text-sm font-medium text-secondary">{label}</p>
                {hint ? <p className="mt-0.5 text-sm text-tertiary">{hint}</p> : null}
            </div>
            <div className="shrink-0 sm:pt-0.5">{children}</div>
        </div>
    );
}

function AppearanceField({ busy, onBusy, onError }: { busy: boolean; onBusy: (v: boolean) => void; onError: (msg: string | null) => void }) {
    const { theme, setAppearance } = useTheme();

    return (
        <Select
            label="Theme"
            selectedKey={theme}
            isDisabled={busy}
            onSelectionChange={(key) => {
                if (key !== "light" && key !== "dark" && key !== "system") {
                    return;
                }
                if (key === theme) {
                    return;
                }
                onBusy(true);
                onError(null);
                void setAppearance(key)
                    .catch((err: unknown) => {
                        const message = err instanceof Error ? err.message : String(err);
                        if (!/cancelled/i.test(message)) {
                            onError(message);
                        }
                    })
                    .finally(() => onBusy(false));
            }}
            items={THEME_OPTIONS}
        >
            {(item) => (
                <Select.Item id={item.id} label={item.label} supportingText={item.supportingText} icon={item.icon} />
            )}
        </Select>
    );
}

function MetaRow({ label, value }: { label: string; value: string }) {
    return (
        <div className="flex items-baseline justify-between gap-4 py-2.5">
            <dt className="text-sm text-tertiary">{label}</dt>
            <dd className="text-right text-sm font-medium text-secondary">{value}</dd>
        </div>
    );
}

/** Official GitHub mark — Untitled UI has no brand glyph for GitHub. */
function GitHubIcon({ className, ...props }: { className?: string } & Record<string, unknown>) {
    return (
        <svg viewBox="0 0 16 16" fill="currentColor" aria-hidden="true" className={className} {...props}>
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8Z" />
        </svg>
    );
}

export function SettingsPage() {
    const [data, setData] = useState<SettingsDTO | null>(null);
    const [path, setPath] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const settings = await SettingsService.Get();
            setData(settings);
            setPath(settings.proxyBridgePath || "");
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    }, []);

    useEffect(() => {
        void reload();
    }, [reload]);

    const onSavePath = async () => {
        setBusy(true);
        setError(null);
        try {
            const settings = await SettingsService.SetProxyBridgePath(path);
            setData(settings);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onToggleDebug = async (enabled: boolean) => {
        setBusy(true);
        setError(null);
        try {
            const settings = await SettingsService.SetDebugLogging(enabled);
            setData(settings);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onToggleLaunchAtLogin = async (enabled: boolean) => {
        setBusy(true);
        setError(null);
        try {
            const settings = await SettingsService.SetLaunchAtLogin(enabled);
            setData(settings);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onOpenLogs = async () => {
        setBusy(true);
        setError(null);
        try {
            await SettingsService.OpenLogsFolder();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const onOpenExternal = async (url: string) => {
        setBusy(true);
        setError(null);
        try {
            await Browser.OpenURL(url);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
        }
    };

    const bridgeVersion =
        data == null
            ? "—"
            : `${data.bridgeVersion || "—"} ${data.versionKnown ? "(known)" : "(unknown)"}`;

    return (
        <div className="mx-auto flex w-full max-w-3xl flex-col">
            <header className="pb-2">
                <h1 className="text-display-xs font-semibold text-primary">Settings</h1>
                <p className="mt-1 text-sm text-tertiary">Appearance, startup, ProxyBridge, and diagnostics for this Mac.</p>
            </header>

            {error && (
                <div className="mt-6 rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            <div className="mt-2">
                <SettingsSection title="Appearance" description="Theme follows this preference. System tracks macOS light or dark mode.">
                    <AppearanceField busy={busy} onBusy={setBusy} onError={setError} />
                </SettingsSection>

                <SettingsSection title="General" description="How RouteDeck behaves when you start the Mac.">
                    <div className="rounded-xl bg-secondary_subtle p-4 ring-1 ring-secondary">
                        <SettingsRow
                            label="Launch at Login"
                            hint="Registers with macOS login items. Does not Apply protection automatically."
                        >
                            <Toggle
                                aria-label="Launch at Login"
                                isSelected={!!data?.launchAtLogin}
                                isDisabled={busy || data == null}
                                onChange={(isSelected) => void onToggleLaunchAtLogin(isSelected)}
                            />
                        </SettingsRow>
                    </div>
                </SettingsSection>

                <SettingsSection
                    title="ProxyBridge"
                    description="RouteDeck drives the official ProxyBridge CLI. Override the binary path only when discovery is wrong."
                >
                    <div className="flex flex-col gap-5">
                        <Input
                            label="Custom path"
                            hint="Checked before documented discovery locations. Leave empty to use discovery."
                            value={path}
                            onChange={setPath}
                            placeholder={data?.discoveredProxyBridgePath || "Absolute path to ProxyBridge binary"}
                        />
                        <div className="flex flex-wrap gap-2">
                            <Button color="primary" size="md" isLoading={busy} onClick={() => void onSavePath()}>
                                Save path
                            </Button>
                        </div>
                        {data && (
                            <dl className="divide-y divide-secondary rounded-xl bg-secondary_subtle px-4 ring-1 ring-secondary">
                                <MetaRow label="Discovery state" value={data.bridgeState || "Unknown"} />
                                <MetaRow label="Detected version" value={bridgeVersion} />
                                {data.discoveredProxyBridgePath ? (
                                    <MetaRow label="Discovered path" value={data.discoveredProxyBridgePath} />
                                ) : null}
                            </dl>
                        )}
                        {data?.compatibilityWarning ? (
                            <p className="rounded-lg bg-warning-primary px-3 py-2 text-sm text-warning-primary ring-1 ring-warning_subtle">
                                {data.compatibilityWarning}
                            </p>
                        ) : null}
                    </div>
                </SettingsSection>

                <SettingsSection title="Diagnostics" description="Logging and local files for troubleshooting.">
                    <div className="flex flex-col gap-1 rounded-xl bg-secondary_subtle p-4 ring-1 ring-secondary">
                        <SettingsRow label="Debug logging" hint="Applies live through slog LevelVar without restart.">
                            <Toggle
                                aria-label="Debug logging"
                                isSelected={!!data?.debugLogging}
                                isDisabled={busy || data == null}
                                onChange={(isSelected) => void onToggleDebug(isSelected)}
                            />
                        </SettingsRow>
                        <div className="my-3 h-px w-full bg-border-secondary" role="separator" />
                        <SettingsRow label="Log files" hint="Opens the RouteDeck logs folder in Finder.">
                            <Button color="secondary" size="sm" isLoading={busy} onClick={() => void onOpenLogs()}>
                                Open folder
                            </Button>
                        </SettingsRow>
                    </div>
                </SettingsSection>

                <SettingsSection title="About" description="Build identity and project links.">
                    <div className="flex flex-col gap-4">
                        {data ? (
                            <dl className="divide-y divide-secondary rounded-xl bg-secondary_subtle px-4 ring-1 ring-secondary">
                                <MetaRow label="Version" value={data.appVersionLabel || "dev"} />
                                <MetaRow label="Commit" value={data.appCommit || "—"} />
                            </dl>
                        ) : (
                            <p className="text-sm text-tertiary">Loading build info…</p>
                        )}
                        <div className="flex flex-wrap gap-x-5 gap-y-2">
                            <Button
                                color="link-color"
                                size="sm"
                                iconLeading={Globe01}
                                isDisabled={busy}
                                onClick={() => void onOpenExternal(MACROX_WEBSITE_URL)}
                            >
                                MACROX
                            </Button>
                            <Button
                                color="link-color"
                                size="sm"
                                iconLeading={GitHubIcon}
                                isDisabled={busy}
                                onClick={() => void onOpenExternal(GITHUB_REPO_URL)}
                            >
                                GitHub
                            </Button>
                        </div>
                    </div>
                </SettingsSection>
            </div>
        </div>
    );
}

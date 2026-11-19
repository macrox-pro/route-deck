import type { ReactNode } from "react";
import { useCallback, useEffect, useState } from "react";
import { Monitor01, Moon01, Sun } from "@untitledui/icons";
import { InlineAlert } from "@/components/application/page-layout/inline-alert";
import { PageHeader } from "@/components/application/page-layout/page-header";
import { PageShell } from "@/components/application/page-layout/page-shell";
import { ScreenState } from "@/components/application/page-layout/screen-state";
import { Button } from "@/components/base/buttons/button";
import { ButtonGroup, ButtonGroupItem } from "@/components/base/button-group/button-group";
import { InputBase } from "@/components/base/input/input";
import { InputGroup } from "@/components/base/input/input-group";
import { Toggle } from "@/components/base/toggle/toggle";
import { useTheme } from "@/providers/theme-provider";
import {
    SettingsService,
    UtilityWindowsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { SettingsDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

function SettingsSection({ title, children }: { title: string; children: ReactNode }) {
    return (
        <section className="flex flex-col gap-3">
            <h2 className="text-sm font-semibold text-secondary">{title}</h2>
            {children}
        </section>
    );
}

/** Inset group: one hairline, divided rows. Same vertical rhythm as list rows (px-4 py-3). */
function SettingsGroup({ children }: { children: ReactNode }) {
    return <div className="divide-y divide-secondary overflow-hidden rounded-xl bg-primary ring-1 ring-secondary">{children}</div>;
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
        <div className="flex items-center justify-between gap-6 px-4 py-3">
            <div className="min-w-0">
                <p className="text-sm font-medium text-primary">{label}</p>
                {hint ? <p className="mt-0.5 text-sm text-tertiary">{hint}</p> : null}
            </div>
            <div className="shrink-0">{children}</div>
        </div>
    );
}

function AppearanceField({ busy, onBusy, onError }: { busy: boolean; onBusy: (v: boolean) => void; onError: (msg: string | null) => void }) {
    const { theme, setAppearance } = useTheme();

    const onThemeChange = (keys: Set<string> | "all") => {
        if (keys === "all") {
            return;
        }
        const key = [...keys][0];
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
    };

    return (
        <ButtonGroup
            size="md"
            selectedKeys={new Set([theme])}
            onSelectionChange={(keys) => onThemeChange(keys as Set<string> | "all")}
            isDisabled={busy}
            aria-label="Theme"
        >
            <ButtonGroupItem id="system" iconLeading={Monitor01}>
                System
            </ButtonGroupItem>
            <ButtonGroupItem id="light" iconLeading={Sun}>
                Light
            </ButtonGroupItem>
            <ButtonGroupItem id="dark" iconLeading={Moon01}>
                Dark
            </ButtonGroupItem>
        </ButtonGroup>
    );
}

function MetaRow({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
    return (
        <div className="flex items-baseline justify-between gap-6 px-4 py-3">
            <dt className="shrink-0 text-sm text-tertiary">{label}</dt>
            <dd className={mono ? "min-w-0 truncate text-right font-mono text-xs text-secondary" : "text-right text-sm font-medium text-secondary"}>
                {value}
            </dd>
        </div>
    );
}

export function SettingsPage() {
    const [data, setData] = useState<SettingsDTO | null>(null);
    const [path, setPath] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);
    const [loading, setLoading] = useState(true);

    const reload = useCallback(async () => {
        setError(null);
        try {
            const settings = await SettingsService.Get();
            setData(settings);
            setPath(settings.proxyBridgePath || "");
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setLoading(false);
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

    const onToggleAutoApply = async (enabled: boolean) => {
        setBusy(true);
        setError(null);
        try {
            const settings = await SettingsService.SetAutoApply(enabled);
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

    const bridgeVersion =
        data == null
            ? "—"
            : `${data.bridgeVersion || "—"} ${data.versionKnown ? "(known)" : "(unknown)"}`;

    return (
        <PageShell>
            <PageHeader title="Settings" description="Look, startup, and the ProxyBridge path on this Mac." />

            {error ? <InlineAlert tone="error">{error}</InlineAlert> : null}

            {loading && !data ? (
                <ScreenState state="loading" title="Loading settings…" />
            ) : (
            <>

            <SettingsSection title="Appearance">
                <SettingsGroup>
                    <div className="flex flex-col gap-3 px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
                        <p className="text-sm font-medium text-primary">Theme</p>
                        <AppearanceField busy={busy} onBusy={setBusy} onError={setError} />
                    </div>
                </SettingsGroup>
            </SettingsSection>

            <SettingsSection title="General">
                <SettingsGroup>
                    <SettingsRow label="Launch at login" hint="Opens RouteDeck when you log in.">
                        <Toggle
                            aria-label="Launch at login"
                            isSelected={!!data?.launchAtLogin}
                            isDisabled={busy || data == null}
                            onChange={(isSelected) => void onToggleLaunchAtLogin(isSelected)}
                        />
                    </SettingsRow>
                    <SettingsRow label="Auto Apply on launch" hint="Applies the Desired profile when RouteDeck starts.">
                        <Toggle
                            aria-label="Auto Apply on launch"
                            isSelected={!!data?.autoApply}
                            isDisabled={busy || data == null}
                            onChange={(isSelected) => void onToggleAutoApply(isSelected)}
                        />
                    </SettingsRow>
                </SettingsGroup>
            </SettingsSection>

            <SettingsSection title="ProxyBridge">
                <InputGroup
                    label="Custom path"
                    hint="Leave empty to use the discovered binary."
                    value={path}
                    onChange={setPath}
                    isDisabled={busy}
                    trailingAddon={
                        <Button color="primary" size="md" isLoading={busy} onClick={() => void onSavePath()}>
                            Save
                        </Button>
                    }
                >
                    <InputBase placeholder={data?.discoveredProxyBridgePath || "Absolute path to the ProxyBridge binary"} />
                </InputGroup>
                {data ? (
                    <SettingsGroup>
                        <MetaRow label="Discovery state" value={data.bridgeState || "Unknown"} />
                        <MetaRow label="Detected version" value={bridgeVersion} />
                        {data.discoveredProxyBridgePath ? (
                            <MetaRow label="Discovered path" value={data.discoveredProxyBridgePath} mono />
                        ) : null}
                    </SettingsGroup>
                ) : null}
                {data?.compatibilityWarning ? (
                    <p className="rounded-lg bg-warning-primary px-4 py-3 text-sm text-warning-primary ring-1 ring-warning_subtle">
                        {data.compatibilityWarning}
                    </p>
                ) : null}
            </SettingsSection>

            <SettingsSection title="Diagnostics">
                <SettingsGroup>
                    <SettingsRow label="Debug logging" hint="More detail in the log. Applies immediately.">
                        <Toggle
                            aria-label="Debug logging"
                            isSelected={!!data?.debugLogging}
                            isDisabled={busy || data == null}
                            onChange={(isSelected) => void onToggleDebug(isSelected)}
                        />
                    </SettingsRow>
                    <SettingsRow label="Log files" hint="Show the logs folder in Finder.">
                        <Button color="secondary" size="md" isLoading={busy} onClick={() => void onOpenLogs()}>
                            Open folder
                        </Button>
                    </SettingsRow>
                </SettingsGroup>
            </SettingsSection>

            <SettingsSection title="About">
                <SettingsGroup>
                    <SettingsRow label="About RouteDeck" hint="Version, commit, and links.">
                        <Button
                            color="secondary"
                            size="md"
                            isDisabled={busy}
                            onClick={() => void UtilityWindowsService.OpenAbout()}
                        >
                            Open
                        </Button>
                    </SettingsRow>
                </SettingsGroup>
            </SettingsSection>
            </>
            )}
        </PageShell>
    );
}

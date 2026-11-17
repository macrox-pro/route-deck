import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/base/buttons/button";
import { Input } from "@/components/base/input/input";
import { Toggle } from "@/components/base/toggle/toggle";
import { SettingsService } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { SettingsDTO } from "../../bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

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

    return (
        <div className="flex flex-col gap-6">
            <h1 className="text-display-xs font-semibold text-primary">Settings</h1>

            {error && (
                <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            )}

            <section className="flex flex-col gap-4 rounded-lg p-4 ring-1 ring-secondary">
                <Input
                    label="ProxyBridge path"
                    hint="Absolute path checked before documented discovery locations."
                    value={path}
                    onChange={setPath}
                    placeholder={data?.discoveredProxyBridgePath || "Not found in documented locations"}
                />
                <div className="flex gap-2">
                    <Button color="primary" size="md" isLoading={busy} onClick={() => void onSavePath()}>
                        Save path
                    </Button>
                    <Button color="secondary" size="md" isLoading={busy} onClick={() => void onOpenLogs()}>
                        Open logs folder
                    </Button>
                </div>
                <Toggle
                    label="Debug logging"
                    hint="Applies live through slog LevelVar without restart."
                    isSelected={!!data?.debugLogging}
                    isDisabled={busy}
                    onChange={(isSelected) => void onToggleDebug(isSelected)}
                />
                <Toggle
                    label="Launch at Login"
                    hint="Registers RouteDeck with the OS login items. Does not Apply protection automatically."
                    isSelected={!!data?.launchAtLogin}
                    isDisabled={busy}
                    onChange={(isSelected) => void onToggleLaunchAtLogin(isSelected)}
                />
            </section>

            {data && (
                <section className="rounded-lg p-4 text-sm ring-1 ring-secondary">
                    <h2 className="mb-2 font-semibold text-secondary">Bridge discovery</h2>
                    <p className="text-tertiary">State: {data.bridgeState || "Unknown"}</p>
                    <p className="text-tertiary">Version: {data.bridgeVersion || "—"} {data.versionKnown ? "(known)" : "(unknown)"}</p>
                    {data.compatibilityWarning && <p className="mt-2 text-warning-primary">{data.compatibilityWarning}</p>}
                </section>
            )}
        </div>
    );
}

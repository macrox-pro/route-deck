import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Badge } from "@/components/base/badges/badges";
import { Button } from "@/components/base/buttons/button";
import { ApplicationsService, SetupGuidesService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { GuideDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

async function loadIcon(path: string): Promise<string | null> {
    try {
        const encoded = await ApplicationsService.Icon(path);
        if (!encoded) {
            return null;
        }
        return `data:image/png;base64,${encoded}`;
    } catch {
        return null;
    }
}

function detectionBadge(guide: GuideDTO) {
    switch (guide.detection) {
        case "installed":
            return (
                <Badge color="success" size="sm" type="pill-color">
                    Installed
                </Badge>
            );
        case "not_installed":
            return (
                <Badge color="gray" size="sm" type="pill-color">
                    Not installed
                </Badge>
            );
        case "unavailable":
            return (
                <Badge color="warning" size="sm" type="pill-color">
                    Detection unavailable
                </Badge>
            );
        default:
            return null;
    }
}

export function SetupGuideCard({
    guide,
    busy,
    onChanged,
}: {
    guide: GuideDTO;
    busy: boolean;
    onChanged: () => void;
}) {
    const navigate = useNavigate();
    const [icon, setIcon] = useState<string | null>(null);
    const [localBusy, setLocalBusy] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!guide.installedPath) {
            setIcon(null);
            return;
        }
        let cancelled = false;
        void loadIcon(guide.installedPath).then((src) => {
            if (!cancelled) {
                setIcon(src);
            }
        });
        return () => {
            cancelled = true;
        };
    }, [guide.installedPath]);

    const onAdd = async () => {
        setError(null);
        if (guide.addMode === "form_fallback") {
            navigate(`/proxies?guide=${encodeURIComponent(guide.id)}`);
            return;
        }
        if (guide.addMode !== "one_click") {
            return;
        }
        setLocalBusy(true);
        try {
            await SetupGuidesService.AddProxy(guide.id);
            onChanged();
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setLocalBusy(false);
        }
    };

    const onDocs = async () => {
        setError(null);
        try {
            await SetupGuidesService.OpenDocs(guide.id);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    };

    return (
        <article className="flex flex-col gap-4 rounded-xl bg-primary p-5 ring-1 ring-secondary">
            <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="flex min-w-0 items-center gap-3">
                    <div className="flex size-10 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-secondary ring-1 ring-secondary">
                        {icon ? <img src={icon} alt="" className="size-10" /> : <span className="text-xs text-tertiary">App</span>}
                    </div>
                    <div className="min-w-0">
                        <h2 className="text-lg font-semibold text-primary">{guide.title}</h2>
                        <div className="mt-1 flex flex-wrap gap-2">
                            {detectionBadge(guide)}
                            {guide.portChecked > 0 && (
                                <Badge color={guide.portOpen ? "success" : "gray"} size="sm" type="pill-color">
                                    {guide.portOpen ? `Port ${guide.portChecked} open` : `Port ${guide.portChecked} closed`}
                                </Badge>
                            )}
                            {guide.alreadyAdded && (
                                <Badge color="brand" size="sm" type="pill-color">
                                    Already added
                                </Badge>
                            )}
                        </div>
                    </div>
                </div>
                <div className="flex flex-wrap gap-2">
                    {guide.docsUrl ? (
                        <Button color="secondary" size="sm" onClick={() => void onDocs()}>
                            Docs
                        </Button>
                    ) : null}
                    {guide.addMode === "one_click" || guide.addMode === "form_fallback" ? (
                        <Button
                            color="primary"
                            size="sm"
                            isDisabled={guide.alreadyAdded}
                            isLoading={busy || localBusy}
                            onClick={() => void onAdd()}
                        >
                            {guide.alreadyAdded ? "Already added" : "Add this proxy"}
                        </Button>
                    ) : null}
                </div>
            </div>

            {guide.connectHint && (
                <p className="rounded-lg bg-warning-primary px-3 py-2 text-sm text-warning-primary ring-1 ring-warning_subtle" role="status">
                    {guide.title} is installed, but port {guide.portChecked} is not listening yet. Connect the app first. You can still add the
                    proxy — health will stay Unavailable until the port opens.
                </p>
            )}

            {error && (
                <p className="rounded-lg bg-error-primary px-3 py-2 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </p>
            )}

            {guide.endpoints?.length > 0 && (
                <div className="rounded-lg bg-secondary px-4 py-3 text-sm ring-1 ring-secondary">
                    <p className="font-medium text-secondary">Local endpoint</p>
                    {guide.endpoints.map((ep) => (
                        <p key={`${ep.host}:${ep.port}`} className="mt-1 text-primary">
                            {ep.kind.toUpperCase()} {ep.host}:{ep.port}
                            {ep.note ? <span className="text-tertiary"> — {ep.note}</span> : null}
                        </p>
                    ))}
                </div>
            )}

            <ol className="flex list-decimal flex-col gap-3 pl-5 text-sm">
                {(guide.steps || []).map((step, index) => (
                    <li key={`${guide.id}-${index}`} className={step.warning ? "text-warning-primary" : "text-secondary"}>
                        <span className="font-medium text-primary">{step.title}</span>
                        <p className="mt-0.5 text-tertiary">{step.body}</p>
                    </li>
                ))}
            </ol>
        </article>
    );
}

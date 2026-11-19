import type { ReactNode } from "react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { RefreshCw01 } from "@untitledui/icons";
import { ProgressSteps, type ProgressStep, type ProgressStepId } from "@/components/application/progress-steps/progress-steps";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { SetupGuideCard } from "@/components/application/setup-guide/setup-guide-card";
import { SuggestionBanner } from "@/components/application/setup-guide/suggestion-banner";
import {
    checklistMessage,
    proxyStepHint,
    proxyStepLabel,
} from "@/components/application/setup-guide/checklist-copy";
import { Button } from "@/components/base/buttons/button";
import { Badge } from "@/components/base/badges/badges";
import { Tabs } from "@/components/application/tabs/tabs";
import { ProtectionBadge } from "@/components/application/protection-badge/protection-badge";
import {
    SetupGuidesService,
    UtilityWindowsService,
} from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { ChecklistDTO, GuideDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

function StateRow({ badge, children }: { badge: ReactNode; children: ReactNode }) {
    return (
        <div className="flex items-start gap-4 px-4 py-3">
            <div className="flex shrink-0 pt-0.5">{badge}</div>
            <p className="min-w-0 text-sm text-tertiary">{children}</p>
        </div>
    );
}

function ProtectionStates() {
    return (
        <div className="flex flex-col gap-6">
            <section className="flex flex-col gap-3">
                <h2 className="text-sm font-semibold text-secondary">Profile</h2>
                <div className="divide-y divide-secondary overflow-hidden rounded-xl bg-primary ring-1 ring-secondary">
                    <StateRow badge={<Badge color="gray" type="modern" size="md">Applied</Badge>}>
                        The profile the CLI confirmed. Applied is not a protection state.
                    </StateRow>
                </div>
            </section>
            <section className="flex flex-col gap-3">
                <h2 className="text-sm font-semibold text-secondary">Protection</h2>
                <div className="divide-y divide-secondary overflow-hidden rounded-xl bg-primary ring-1 ring-secondary">
                    <StateRow badge={<ProtectionBadge protection="Active" />}>
                        Every app set to Proxy or Block is covered.
                    </StateRow>
                    <StateRow badge={<ProtectionBadge protection="Degraded" />}>
                        Rules are in effect, with a gap such as TCP up and UDP unavailable.
                    </StateRow>
                    <StateRow badge={<ProtectionBadge protection="Inactive" />}>
                        Nothing to protect, or the applied profile is Direct only.
                    </StateRow>
                    <StateRow badge={<ProtectionBadge protection="Unprotected" />}>
                        Rules are not in effect. Stopped, a conflict, or confirm failed.
                    </StateRow>
                    <StateRow badge={<ProtectionBadge protection="Failed" />}>
                        Apply or the bridge reported a failure.
                    </StateRow>
                    <StateRow badge={<ProtectionBadge protection="Unknown" />}>
                        Still checking. Initializing uses the same gray badge.
                    </StateRow>
                </div>
            </section>
            <p className="text-sm text-tertiary">
                Confirm helper processes on Applications. Otherwise their traffic can bypass Proxy or Block.
            </p>
        </div>
    );
}

function stepsFromChecklist(data: ChecklistDTO | null): ProgressStep[] {
    if (!data) {
        return [
            { id: "bridge", label: checklistMessage("checklist.step.bridge"), status: "todo" },
            { id: "proxy", label: proxyStepLabel("missing"), status: "todo" },
            { id: "application", label: checklistMessage("checklist.step.application"), status: "todo" },
            { id: "profile", label: checklistMessage("checklist.step.profile"), status: "todo" },
        ];
    }
    return [
        {
            id: "bridge",
            label: checklistMessage("checklist.step.bridge"),
            status: data.bridgeFound ? "done" : "todo",
        },
        {
            id: "proxy",
            label: proxyStepLabel(data.proxyStep),
            hint: data.proxyPartial ? proxyStepHint(data.proxyStep) : undefined,
            status: data.proxyAvailable ? "done" : data.proxyPartial ? "partial" : "todo",
        },
        {
            id: "application",
            label: checklistMessage("checklist.step.application"),
            status: data.applicationAdded ? "done" : "todo",
        },
        {
            id: "profile",
            label: checklistMessage("checklist.step.profile"),
            status: data.profileApplied ? "done" : "todo",
        },
    ];
}

const SECTION_IDS: Record<ProgressStepId, string> = {
    bridge: "guides-bridge",
    proxy: "guides-proxy",
    application: "guides-application",
    profile: "guides-profile",
};

export function GuidesDialog() {
    const { step: stepParam } = useParams<{ step?: string }>();
    const [guides, setGuides] = useState<GuideDTO[]>([]);
    const [checklist, setChecklist] = useState<ChecklistDTO | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [busy, setBusy] = useState(false);
    const [ready, setReady] = useState(false);
    const [tab, setTab] = useState("guides");

    const reload = useCallback(async (refresh: boolean) => {
        setBusy(true);
        setError(null);
        try {
            const [guideRows, checklistRow] = await Promise.all([
                SetupGuidesService.List(refresh),
                SetupGuidesService.Checklist(),
            ]);
            setGuides(guideRows);
            setChecklist(checklistRow);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setBusy(false);
            setReady(true);
        }
    }, []);

    useEffect(() => {
        void reload(true);
    }, [reload]);

    const scrollRef = useRef<HTMLDivElement>(null);

    const bridgeGuides = useMemo(() => guides.filter((g) => g.kind === "bridge"), [guides]);
    const clientGuides = useMemo(() => guides.filter((g) => g.kind !== "bridge"), [guides]);

    const scrollToStep = useCallback((step: ProgressStepId) => {
        const id = SECTION_IDS[step];
        const root = scrollRef.current;
        const target = root?.querySelector(`#${id}`);
        if (target instanceof HTMLElement) {
            target.scrollIntoView({ behavior: "smooth", block: "start" });
        }
    }, []);

    useEffect(() => {
        const step = stepParam as ProgressStepId | undefined;
        if (step !== "bridge" && step !== "proxy") {
            return;
        }
        setTab("guides");
        const timer = window.setTimeout(() => scrollToStep(step), 100);
        return () => window.clearTimeout(timer);
    }, [stepParam, guides.length, scrollToStep]);

    return (
        <div className="flex h-full flex-col bg-primary text-primary">
            <header className="flex shrink-0 flex-wrap items-start justify-between gap-3 border-b border-secondary px-6 py-4">
                <div className="min-w-0">
                    <h1 className="text-display-xs font-semibold text-primary">Setup guides</h1>
                    <p className="mt-1 text-sm text-tertiary">
                        Install ProxyBridge, then add the local proxy your VPN app exposes.
                    </p>
                </div>
                <Button color="secondary" size="md" iconLeading={RefreshCw01} isLoading={busy} onClick={() => void reload(true)}>
                    Refresh
                </Button>
            </header>

            {error ? (
                <div className="shrink-0 px-6 pt-6">
                    <div className="rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                        {error}
                    </div>
                </div>
            ) : null}

            <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden pt-6">
                <Tabs selectedKey={tab} onSelectionChange={(key) => setTab(String(key))} className="min-h-0 flex-1">
                    <Tabs.List type="underline" size="sm" className="shrink-0 px-6">
                        <Tabs.Item id="guides">Guides</Tabs.Item>
                        <Tabs.Item id="states">Protection states</Tabs.Item>
                    </Tabs.List>
                    <Tabs.Panel id="guides" className="min-h-0 flex-1 overflow-y-auto px-6 pt-6 pb-6 outline-hidden">
                        <div ref={scrollRef}>
                            {!ready ? (
                                <div className="flex min-h-64 items-center justify-center">
                                    <LoadingIndicator label="Loading setup guides…" />
                                </div>
                            ) : (
                                <div className="flex flex-col gap-6">
                                    <section className="flex flex-col gap-4 rounded-xl bg-secondary p-4 ring-1 ring-secondary">
                                        <div className="flex flex-col gap-1">
                                            <h2 className="text-sm font-semibold text-primary">Checklist</h2>
                                            <p className="text-sm text-tertiary">
                                                Progress only. ProxyBridge and the local proxy are on this page. Add an application and
                                                apply a profile on those screens.
                                            </p>
                                        </div>
                                        <ProgressSteps orientation="horizontal" compact steps={stepsFromChecklist(checklist)} />
                                    </section>
                                    <SuggestionBanner
                                        onAdded={() => {
                                            void reload(false);
                                        }}
                                    />
                                    <section id="guides-bridge" className="flex scroll-mt-6 flex-col gap-4">
                                        <div>
                                            <h2 className="text-lg font-semibold text-primary">ProxyBridge</h2>
                                            <p className="mt-1 text-sm text-tertiary">
                                                Required. Install the official app and approve its network extension. RouteDeck drives the
                                                CLI. It does not replace ProxyBridge.
                                            </p>
                                        </div>
                                        {bridgeGuides.map((guide) => (
                                            <SetupGuideCard
                                                key={guide.id}
                                                guide={guide}
                                                busy={busy}
                                                onChanged={() => void reload(false)}
                                            />
                                        ))}
                                    </section>
                                    <section id="guides-proxy" className="flex scroll-mt-6 flex-col gap-4">
                                        <div>
                                            <h2 className="text-lg font-semibold text-primary">Local proxy apps</h2>
                                            <p className="mt-1 text-sm text-tertiary">
                                                Shadowrocket, v2rayN, Surge, and similar. Expose a local SOCKS5 or HTTP port, then add it
                                                here. Turn TUN, system VPN, and Enhanced Mode off so per-app rules can apply.
                                            </p>
                                        </div>
                                        {clientGuides.map((guide) => (
                                            <SetupGuideCard
                                                key={guide.id}
                                                guide={guide}
                                                busy={busy}
                                                onChanged={() => void reload(false)}
                                            />
                                        ))}
                                    </section>
                                </div>
                            )}
                        </div>
                    </Tabs.Panel>
                    <Tabs.Panel id="states" className="min-h-0 flex-1 overflow-y-auto px-6 pt-6 pb-6 outline-hidden">
                        <ProtectionStates />
                    </Tabs.Panel>
                </Tabs>
            </div>

            <footer className="flex shrink-0 justify-end border-t border-secondary px-6 py-4">
                <Button color="primary" size="md" onClick={() => void UtilityWindowsService.Close("guides")}>
                    Done
                </Button>
            </footer>
        </div>
    );
}
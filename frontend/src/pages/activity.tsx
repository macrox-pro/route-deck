import { useCallback, useEffect, useRef, useState } from "react";
import { Events } from "@wailsio/runtime";
import { useVirtualizer } from "@tanstack/react-virtual";
import { RefreshCcw01 } from "@untitledui/icons";
import { EmptyState } from "@/components/application/empty-state/empty-state";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { Badge } from "@/components/base/badges/badges";
import type { BadgeColors } from "@/components/base/badges/badge-types";
import { Button } from "@/components/base/buttons/button";
import { cx } from "@/utils/cx";
import { ActivityService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";
import type { ActivityRecordDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

const REFRESH_EVENTS = [
    "protection.state.changed",
    "proxy.health.changed",
    "proxybridge.state.changed",
    "operation.progress",
    "operation.failed",
] as const;

function levelColor(level: string): BadgeColors {
    switch (level.trim().toUpperCase()) {
        case "ERROR":
        case "ERR":
        case "FATAL":
            return "error";
        case "WARN":
        case "WARNING":
            return "warning";
        case "DEBUG":
        case "TRACE":
            return "gray";
        case "INFO":
            return "blue";
        default:
            return "slate";
    }
}

function formatActivityTime(iso: string): { clock: string; day: string; title: string } {
    if (!iso) {
        return { clock: "—", day: "", title: "" };
    }
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) {
        return { clock: iso, day: "", title: iso };
    }
    return {
        clock: date.toLocaleTimeString(undefined, {
            hour: "2-digit",
            minute: "2-digit",
            second: "2-digit",
            hour12: false,
        }),
        day: date.toLocaleDateString(undefined, {
            month: "short",
            day: "numeric",
        }),
        title: iso,
    };
}

function footerStatus(loadingMore: boolean, hasMore: boolean): string {
    if (loadingMore) {
        return "Loading older events…";
    }
    if (hasMore) {
        return "Scroll for older events";
    }
    return "End of loaded log";
}

export function ActivityPage() {
    const [records, setRecords] = useState<ActivityRecordDTO[]>([]);
    const [hasMore, setHasMore] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const parentRef = useRef<HTMLDivElement>(null);
    const loadMoreLock = useRef(false);

    const reload = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const result = await ActivityService.Recent("");
            setRecords(result.records ?? []);
            setHasMore(Boolean(result.hasMore));
            parentRef.current?.scrollTo({ top: 0 });
        } catch (err) {
            setRecords([]);
            setHasMore(false);
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setLoading(false);
        }
    }, []);

    const loadMore = useCallback(async () => {
        if (loadMoreLock.current || loadingMore || loading || !hasMore || records.length === 0) {
            return;
        }
        const before = records[records.length - 1]?.time;
        if (!before) {
            return;
        }
        loadMoreLock.current = true;
        setLoadingMore(true);
        setError(null);
        try {
            const result = await ActivityService.Recent(before);
            const next = result.records ?? [];
            setRecords((prev) => [...prev, ...next]);
            setHasMore(Boolean(result.hasMore));
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        } finally {
            setLoadingMore(false);
            loadMoreLock.current = false;
        }
    }, [hasMore, loading, loadingMore, records]);

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

    // Docs: https://tanstack.com/virtual/latest/docs/framework/react/react-virtual
    const rowVirtualizer = useVirtualizer({
        count: records.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 88,
        overscan: 12,
        useFlushSync: false,
        measureElement:
            typeof window !== "undefined" && !navigator.userAgent.includes("Firefox")
                ? (element) => element.getBoundingClientRect().height
                : undefined,
    });

    const virtualItems = rowVirtualizer.getVirtualItems();
    const lastVirtualItem = virtualItems[virtualItems.length - 1];

    useEffect(() => {
        if (!lastVirtualItem || !hasMore || loading || loadingMore) {
            return;
        }
        if (lastVirtualItem.index >= records.length - 4) {
            void loadMore();
        }
    }, [hasMore, lastVirtualItem, loadMore, loading, loadingMore, records.length]);

    return (
        <div className="flex min-h-0 flex-1 flex-col gap-4">
            <div className="flex shrink-0 flex-wrap items-end justify-between gap-4">
                <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-3">
                        <h1 className="text-display-xs font-semibold text-primary">Activity</h1>
                        {records.length > 0 ? (
                            <Badge color="gray" type="modern" size="sm">
                                {records.length}
                                {hasMore ? "+" : ""} events
                            </Badge>
                        ) : null}
                    </div>
                    <p className="mt-1 text-sm text-tertiary">
                        Newest RouteDeck JSON log first. Scroll loads older pages; product events refresh the newest page.
                    </p>
                </div>
                <Button
                    color="secondary"
                    size="md"
                    iconLeading={RefreshCcw01}
                    onClick={() => void reload()}
                    isDisabled={loading || loadingMore}
                    isLoading={loading && records.length > 0}
                >
                    Refresh
                </Button>
            </div>

            {error ? (
                <div className="shrink-0 rounded-lg bg-error-primary px-4 py-3 text-sm text-error-primary ring-1 ring-error_subtle" role="alert">
                    {error}
                </div>
            ) : null}

            {loading && records.length === 0 && !error ? (
                <div className="flex min-h-64 flex-1 items-center justify-center rounded-xl ring-1 ring-secondary">
                    <LoadingIndicator label="Loading activity…" />
                </div>
            ) : null}

            {!loading && !error && records.length === 0 ? (
                <div className="flex flex-1 items-center justify-center rounded-xl ring-1 ring-secondary">
                    <EmptyState size="md">
                        <EmptyState.Header>
                            <EmptyState.FeaturedIcon color="gray" theme="modern" />
                        </EmptyState.Header>
                        <EmptyState.Content>
                            <EmptyState.Title>No activity yet</EmptyState.Title>
                            <EmptyState.Description>
                                RouteDeck JSON log records appear here after the app writes to the log file.
                            </EmptyState.Description>
                        </EmptyState.Content>
                    </EmptyState>
                </div>
            ) : null}

            {records.length > 0 ? (
                <section className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-xl bg-primary ring-1 ring-secondary">
                    <div className="flex shrink-0 items-center gap-3 border-b border-secondary bg-secondary px-4 py-2 text-xs font-medium text-tertiary">
                        <span className="w-18 shrink-0">Time</span>
                        <span className="w-16 shrink-0">Level</span>
                        <span className="min-w-0 flex-1">Event</span>
                    </div>

                    <div
                        ref={parentRef}
                        className="min-h-0 flex-1 overflow-auto overscroll-contain"
                        role="log"
                        aria-live="polite"
                    >
                        <div
                            className="relative w-full"
                            style={{
                                height: `${rowVirtualizer.getTotalSize()}px`,
                            }}
                        >
                            {virtualItems.map((virtualItem) => {
                                const rec = records[virtualItem.index];
                                const time = formatActivityTime(rec.time);
                                const meta = (rec.attrs ?? []).filter((attr) => attr.key !== "err");
                                const level = rec.level || "—";

                                return (
                                    <article
                                        key={virtualItem.key}
                                        data-index={virtualItem.index}
                                        ref={rowVirtualizer.measureElement}
                                        className={cx(
                                            "absolute top-0 left-0 w-full border-b border-secondary px-4 py-2.5",
                                            "hover:bg-secondary",
                                        )}
                                        style={{
                                            transform: `translateY(${virtualItem.start}px)`,
                                        }}
                                    >
                                        <div className="flex items-start gap-3">
                                            <time
                                                className="flex w-18 shrink-0 flex-col gap-0.5 font-mono tabular-nums"
                                                dateTime={rec.time || undefined}
                                                title={time.title}
                                            >
                                                <span className="text-xs leading-4 text-primary">{time.clock}</span>
                                                {time.day ? (
                                                    <span className="text-[11px] leading-3 text-quaternary">{time.day}</span>
                                                ) : null}
                                            </time>
                                            <div className="w-16 shrink-0 pt-0.5">
                                                <Badge color={levelColor(rec.level)} type="pill-color" size="sm">
                                                    {level}
                                                </Badge>
                                            </div>
                                            <div className="min-w-0 flex-1">
                                                <p className="text-sm font-medium wrap-break-word text-primary">
                                                    {rec.msg || "(no message)"}
                                                </p>
                                                {meta.length > 0 ? (
                                                    <div className="mt-1.5 flex flex-wrap gap-1.5">
                                                        {meta.map((item) => (
                                                            <Badge
                                                                key={`${item.key}-${item.value}`}
                                                                color="gray"
                                                                type="modern"
                                                                size="sm"
                                                            >
                                                                <span className="text-quaternary">{item.key}</span>
                                                                <span className="ml-1">{item.value}</span>
                                                            </Badge>
                                                        ))}
                                                    </div>
                                                ) : null}
                                                {rec.err ? (
                                                    <p className="mt-1.5 font-mono text-xs wrap-break-word text-error-primary">
                                                        {rec.err}
                                                    </p>
                                                ) : null}
                                            </div>
                                        </div>
                                    </article>
                                );
                            })}
                        </div>
                    </div>

                    <div className="flex shrink-0 items-center border-t border-secondary px-4 py-2">
                        <p className="text-xs text-tertiary" aria-live="polite">
                            {footerStatus(loadingMore, hasMore)}
                        </p>
                    </div>
                </section>
            ) : null}
        </div>
    );
}

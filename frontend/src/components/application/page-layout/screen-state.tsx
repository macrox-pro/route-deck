import type { FC, ReactNode } from "react";
import { EmptyState } from "@/components/application/empty-state/empty-state";
import { LoadingIndicator } from "@/components/application/loading-indicator/loading-indicator";
import { Button } from "@/components/base/buttons/button";

export function ScreenState({
    state,
    icon,
    title,
    description,
    actions,
    onRetry,
}: {
    state: "loading" | "empty" | "error";
    icon?: FC<{ className?: string }>;
    title: string;
    description?: string;
    actions?: ReactNode;
    onRetry?: () => void;
}) {
    if (state === "loading") {
        return (
            <div className="flex min-h-48 items-center justify-center py-12">
                <LoadingIndicator type="line-spinner" size="md" label={title} />
            </div>
        );
    }

    if (state === "error") {
        return (
            <EmptyState size="md">
                <EmptyState.Header>
                    <EmptyState.FeaturedIcon color="error" theme="light" icon={icon} />
                </EmptyState.Header>
                <EmptyState.Content>
                    <EmptyState.Title>{title}</EmptyState.Title>
                    {description ? <EmptyState.Description>{description}</EmptyState.Description> : null}
                </EmptyState.Content>
                <EmptyState.Footer>
                    {actions}
                    {onRetry ? (
                        <Button color="secondary" size="md" onClick={onRetry}>
                            Retry
                        </Button>
                    ) : null}
                </EmptyState.Footer>
            </EmptyState>
        );
    }

    return (
        <EmptyState size="md">
            <EmptyState.Header>
                <EmptyState.FeaturedIcon color="gray" theme="modern" icon={icon} />
            </EmptyState.Header>
            <EmptyState.Content>
                <EmptyState.Title>{title}</EmptyState.Title>
                {description ? <EmptyState.Description>{description}</EmptyState.Description> : null}
            </EmptyState.Content>
            {actions ? <EmptyState.Footer>{actions}</EmptyState.Footer> : null}
        </EmptyState>
    );
}

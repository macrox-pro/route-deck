import { EmptyState } from "@/components/application/empty-state/empty-state";

export function ActivityPage() {
    return (
        <div className="flex flex-col gap-6">
            <h1 className="text-display-xs font-semibold text-primary">Activity</h1>
            <EmptyState size="md">
                <EmptyState.Header>
                    <EmptyState.FeaturedIcon color="gray" theme="modern" />
                </EmptyState.Header>
                <EmptyState.Content>
                    <EmptyState.Title>No activity yet</EmptyState.Title>
                    <EmptyState.Description>RouteDeck events land in a later task.</EmptyState.Description>
                </EmptyState.Content>
            </EmptyState>
        </div>
    );
}

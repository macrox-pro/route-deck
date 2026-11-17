import { BadgeWithDot } from "@/components/base/badges/badges";
import type { BadgeColors } from "@/components/base/badges/badge-types";

function colorForProtection(protection: string): BadgeColors {
    switch (protection) {
        case "Active":
            return "success";
        case "Degraded":
            return "warning";
        case "Unprotected":
        case "Failed":
            return "error";
        case "Unknown":
        case "Inactive":
        case "Initializing":
        default:
            return "gray";
    }
}

interface ProtectionBadgeProps {
    protection: string;
}

export function ProtectionBadge({ protection }: ProtectionBadgeProps) {
    const color = colorForProtection(protection);
    const label = protection || "Unknown";
    return (
        <BadgeWithDot color={color} size="md" type="pill-color">
            {label}
        </BadgeWithDot>
    );
}

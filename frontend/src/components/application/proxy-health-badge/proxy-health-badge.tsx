import { AlertCircle, AlertTriangle, CheckCircle, Loading02, XCircle } from "@untitledui/icons";
import { BadgeWithIcon } from "@/components/base/badges/badges";
import type { BadgeColors } from "@/components/base/badges/badge-types";
import { proxyHealthLabel } from "@/utils/proxy-health-label";

type HealthVisual = {
    color: BadgeColors;
    label: string;
    icon: typeof CheckCircle;
};

function visualFor(health: string): HealthVisual {
    switch (health) {
        case "Available":
            return { color: "success", label: "Healthy", icon: CheckCircle };
        case "Unavailable":
            return { color: "error", label: "Unreachable", icon: XCircle };
        case "AuthenticationFailed":
            return { color: "error", label: "Auth failed", icon: AlertCircle };
        case "Checking":
            return { color: "gray", label: "Checking…", icon: Loading02 };
        case "TCPAvailableUDPUnavailable":
            // Handled as a split pair below; fallback single badge if needed.
            return { color: "warning", label: "TCP OK, UDP unavailable", icon: AlertTriangle };
        case "Unknown":
        default:
            return { color: "gray", label: proxyHealthLabel(health), icon: AlertCircle };
    }
}

interface ProxyHealthBadgeProps {
    health: string | undefined | null;
    size?: "sm" | "md";
}

/**
 * Closed ProxyHealth wire tokens as pill badges with icons.
 * TCPAvailableUDPUnavailable renders as two pills: TCP (success) + UDP (warning).
 */
export function ProxyHealthBadge({ health, size = "sm" }: ProxyHealthBadgeProps) {
    const token = health || "Unknown";
    const title = proxyHealthLabel(token);

    if (token === "TCPAvailableUDPUnavailable") {
        return (
            <span className="inline-flex flex-wrap items-center gap-1.5" title={title}>
                <BadgeWithIcon color="success" size={size} type="pill-color" iconLeading={CheckCircle}>
                    TCP
                </BadgeWithIcon>
                <BadgeWithIcon color="warning" size={size} type="pill-color" iconLeading={AlertTriangle}>
                    UDP unavailable
                </BadgeWithIcon>
            </span>
        );
    }

    const visual = visualFor(token);
    return (
        <BadgeWithIcon color={visual.color} size={size} type="pill-color" iconLeading={visual.icon} className="shrink-0">
            {visual.label}
        </BadgeWithIcon>
    );
}

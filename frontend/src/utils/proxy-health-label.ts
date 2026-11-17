/** Closed English labels for ProxyHealth wire tokens (tokens themselves are unchanged). */
const LABELS: Record<string, string> = {
    Unknown: "Unknown",
    Checking: "Checking…",
    Available: "Healthy",
    Unavailable: "Unreachable",
    AuthenticationFailed: "Auth failed",
    TCPAvailableUDPUnavailable: "TCP OK, UDP unavailable",
};

export function proxyHealthLabel(health: string | undefined | null): string {
    if (!health) {
        return LABELS.Unknown;
    }
    return LABELS[health] ?? health;
}

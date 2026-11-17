import type { GuideDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

export type ProxyGuideMatch = {
    installedPath?: string;
    kindFallback: "socks5" | "http" | string;
};

/**
 * Match a catalog proxy to a Setup guide endpoint (127.0.0.1 + port + kind).
 * Ambiguous shared ports: first installed guide with installedPath wins.
 */
export function matchProxyGuide(
    host: string,
    port: number,
    kind: string,
    guides: GuideDTO[],
): ProxyGuideMatch {
    const normalizedHost = host.trim().toLowerCase();
    const kindLower = kind.trim().toLowerCase();
    let fallbackPath: string | undefined;

    for (const guide of guides) {
        if (guide.kind === "bridge") {
            continue;
        }
        for (const ep of guide.endpoints ?? []) {
            if (ep.host.trim().toLowerCase() !== normalizedHost) {
                continue;
            }
            if (ep.port !== port) {
                continue;
            }
            if (ep.kind && ep.kind.trim().toLowerCase() !== kindLower) {
                continue;
            }
            if (guide.detection === "installed" && guide.installedPath) {
                return { installedPath: guide.installedPath, kindFallback: kindLower };
            }
            if (!fallbackPath && guide.installedPath) {
                fallbackPath = guide.installedPath;
            }
        }
    }
    return { installedPath: fallbackPath, kindFallback: kindLower };
}

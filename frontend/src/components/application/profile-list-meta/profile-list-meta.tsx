import { useEffect, useMemo, useState } from "react";
import { Badge } from "@/components/base/badges/badges";
import { appBundlePath } from "@/utils/app-bundle-path";
import { loadAppIcon } from "@/utils/load-app-icon";
import type { PolicyDTO, ProcessIdentityDTO } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

function countActions(policies: PolicyDTO[]) {
    let proxy = 0;
    let direct = 0;
    let block = 0;
    for (const policy of policies ?? []) {
        switch (policy.action) {
            case "proxy":
                proxy += 1;
                break;
            case "direct":
                direct += 1;
                break;
            case "block":
                block += 1;
                break;
            default:
                break;
        }
    }
    return { proxy, direct, block };
}

function ProfileIconStack({ paths }: { paths: string[] }) {
    const pathKey = paths.join("|");
    const [srcs, setSrcs] = useState<(string | null)[]>([]);
    useEffect(() => {
        let cancelled = false;
        if (paths.length === 0) {
            setSrcs([]);
            return;
        }
        void Promise.all(paths.map((path) => loadAppIcon(path))).then((next) => {
            if (!cancelled) {
                setSrcs(next);
            }
        });
        return () => {
            cancelled = true;
        };
    }, [pathKey, paths]);

    if (paths.length === 0) {
        return (
            <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-secondary text-xs font-semibold text-tertiary ring-1 ring-secondary">
                PF
            </div>
        );
    }

    return (
        <div className="flex shrink-0 items-center">
            {srcs.map((src, index) => (
                <div
                    key={`${paths[index]}-${index}`}
                    className="flex size-10 items-center justify-center overflow-hidden rounded-lg bg-secondary ring-1 ring-secondary"
                    style={{ marginLeft: index === 0 ? 0 : -10, zIndex: srcs.length - index }}
                >
                    {src ? <img src={src} alt="" className="size-10" /> : null}
                </div>
            ))}
        </div>
    );
}

export function ProfileListMeta({
    name,
    policies,
    processes,
    isDesired,
    isApplied,
}: {
    name: string;
    policies: PolicyDTO[];
    processes: ProcessIdentityDTO[];
    isDesired: boolean;
    isApplied: boolean;
}) {
    const counts = countActions(policies);
    const iconPaths = useMemo(() => {
        const processByID = new Map(processes.map((process) => [process.id, process]));
        const seen = new Set<string>();
        const paths: string[] = [];
        for (const policy of policies ?? []) {
            const process = processByID.get(policy.processId);
            const bundle = appBundlePath(process?.path);
            if (!bundle || seen.has(bundle)) {
                continue;
            }
            seen.add(bundle);
            paths.push(bundle);
            if (paths.length >= 3) {
                break;
            }
        }
        return paths;
    }, [policies, processes]);

    const uniqueProcesses = new Set((policies ?? []).map((p) => p.processId)).size;
    const extra = Math.max(0, uniqueProcesses - iconPaths.length);
    const policyTotal = (policies ?? []).length;

    return (
        <div className="flex min-w-0 items-center gap-3">
            <div className="flex items-center gap-1.5">
                <ProfileIconStack paths={iconPaths} />
                {extra > 0 ? <span className="text-xs text-tertiary">+{extra}</span> : null}
            </div>
            <div className="min-w-0">
                <p className="truncate text-sm font-semibold text-primary">{name}</p>
                <div className="mt-1 flex min-w-0 flex-wrap items-center gap-1.5">
                    {isDesired ? (
                        <Badge color="brand" size="sm" type="pill-color">
                            Desired
                        </Badge>
                    ) : null}
                    {isApplied ? (
                        <Badge color="success" size="sm" type="pill-color">
                            Applied
                        </Badge>
                    ) : null}
                    {counts.proxy > 0 ? (
                        <Badge color="gray" size="sm" type="pill-color">
                            Proxy {counts.proxy}
                        </Badge>
                    ) : null}
                    {counts.direct > 0 ? (
                        <Badge color="gray" size="sm" type="pill-color">
                            Direct {counts.direct}
                        </Badge>
                    ) : null}
                    {counts.block > 0 ? (
                        <Badge color="gray" size="sm" type="pill-color">
                            Block {counts.block}
                        </Badge>
                    ) : null}
                    {policyTotal === 0 ? <span className="text-xs text-tertiary">No policies</span> : null}
                </div>
            </div>
        </div>
    );
}

import { useEffect, useState } from "react";
import { loadAppIcon } from "@/utils/load-app-icon";

export function ProxyListIcon({ installedPath, kind }: { installedPath?: string; kind: string }) {
    const [src, setSrc] = useState<string | null>(null);
    useEffect(() => {
        if (!installedPath) {
            setSrc(null);
            return;
        }
        let cancelled = false;
        void loadAppIcon(installedPath).then((next) => {
            if (!cancelled) {
                setSrc(next);
            }
        });
        return () => {
            cancelled = true;
        };
    }, [installedPath]);

    const kindLabel = kind.toLowerCase() === "http" ? "HTTP" : "SOCKS5";

    return (
        <div className="flex size-8 shrink-0 items-center justify-center overflow-hidden rounded-md bg-secondary text-[10px] font-semibold text-tertiary ring-1 ring-secondary">
            {src ? <img src={src} alt="" className="size-8" /> : <span>{kindLabel === "HTTP" ? "HTTP" : "S5"}</span>}
        </div>
    );
}

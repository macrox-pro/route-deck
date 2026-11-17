/** Walk a process executable path up to the containing `.app` bundle. */
export function appBundlePath(processPath: string | undefined | null): string | undefined {
    if (!processPath) {
        return undefined;
    }
    let p = processPath.replace(/\/+$/, "");
    while (p && p !== "/" && p !== ".") {
        if (/\.app$/i.test(p)) {
            return p;
        }
        const idx = p.lastIndexOf("/");
        if (idx <= 0) {
            break;
        }
        p = p.slice(0, idx);
    }
    return undefined;
}

import { ApplicationsService } from "@bindings/github.com/macrox-pro/route-deck/internal/adapter/inbound/wails";

const ICON_FETCH_LIMIT = 2;
let iconFetchesActive = 0;
const iconFetchWaiters: Array<() => void> = [];
const iconFetchByPath = new Map<string, Promise<string | null>>();

async function acquireIconSlot(): Promise<void> {
    if (iconFetchesActive < ICON_FETCH_LIMIT) {
        iconFetchesActive += 1;
        return;
    }
    await new Promise<void>((resolve) => {
        iconFetchWaiters.push(resolve);
    });
    iconFetchesActive += 1;
}

function releaseIconSlot(): void {
    iconFetchesActive -= 1;
    const next = iconFetchWaiters.shift();
    if (next) {
        next();
    }
}

/** Load a PNG data URL for an application path via ApplicationsService.Icon. */
export function loadAppIcon(path: string): Promise<string | null> {
    const existing = iconFetchByPath.get(path);
    if (existing) {
        return existing;
    }
    const pending = (async () => {
        await acquireIconSlot();
        try {
            const encoded = await ApplicationsService.Icon(path);
            if (!encoded) {
                return null;
            }
            return `data:image/png;base64,${encoded}`;
        } catch {
            return null;
        } finally {
            releaseIconSlot();
        }
    })();
    iconFetchByPath.set(path, pending);
    return pending;
}

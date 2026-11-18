/**
 * English checklist copy keyed by stable message IDs for future frontend i18n.
 * Go sends closed proxyStep tokens only; this map owns user-visible strings.
 */
export const checklistCopy = {
    "checklist.title": "First-run checklist",
    "checklist.subtitle": "These checks track setup progress only. They do not mean protection is Active.",
    "checklist.dismiss": "Dismiss",
    "checklist.open": "Open",
    "checklist.sr.done": "Done: ",
    "checklist.sr.todo": "To do: ",
    "checklist.sr.partial": "Partial: ",
    "checklist.step.bridge": "ProxyBridge found",
    "checklist.step.application": "An application is in the catalog",
    "checklist.step.profile": "A profile is applied",
    "checklist.step.proxy.missing": "A proxy is Available",
    "checklist.step.proxy.available": "A proxy is Available",
    "checklist.step.proxy.udp_limited": "A proxy is Available (UDP limited)",
    "checklist.step.proxy.udp_limited.hint":
        "This proxy accepts TCP, but UDP associate failed. RouteDeck can still protect TCP; UDP may stay unprotected or Degraded.",
    "checklist.step.proxy.checking": "A proxy is Available (checking…)",
    "checklist.step.proxy.checking.hint":
        "RouteDeck is still probing this proxy. Wait a moment, or open Proxies to recheck.",
} as const;

export type ChecklistMessageID = keyof typeof checklistCopy;

export function checklistMessage(id: ChecklistMessageID): string {
    return checklistCopy[id];
}

export function proxyStepLabel(step: string | undefined | null): string {
    switch (step) {
        case "available":
            return checklistMessage("checklist.step.proxy.available");
        case "udp_limited":
            return checklistMessage("checklist.step.proxy.udp_limited");
        case "checking":
            return checklistMessage("checklist.step.proxy.checking");
        case "missing":
        default:
            return checklistMessage("checklist.step.proxy.missing");
    }
}

export function proxyStepHint(step: string | undefined | null): string {
    switch (step) {
        case "udp_limited":
            return checklistMessage("checklist.step.proxy.udp_limited.hint");
        case "checking":
            return checklistMessage("checklist.step.proxy.checking.hint");
        default:
            return "";
    }
}

import type { ReactNode } from "react";

export function DialogShell({ children, footer }: { children: ReactNode; footer: ReactNode }) {
    return (
        <div className="flex h-dvh min-h-0 flex-col bg-primary text-primary">
            <div className="min-h-0 flex-1 overflow-y-auto px-4 py-4">{children}</div>
            {footer}
        </div>
    );
}

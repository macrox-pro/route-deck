import type { ReactNode } from "react";
import { useEffect, useMemo } from "react";
import { System } from "@wailsio/runtime";
import { Button } from "@/components/base/buttons/button";
import { cx } from "@/utils/cx";

export type DialogFooterProps = {
    cancelLabel?: string;
    saveLabel?: string;
    saveDisabled?: boolean;
    saveLoading?: boolean;
    /** Optional status line above the buttons (same px-6 rhythm). */
    leading?: ReactNode;
    onCancel: () => void;
    onSave: () => void;
};

/**
 * Platform footer button order: macOS/Linux Cancel then Save; Windows Save then Cancel.
 * Enter saves, Esc cancels, Cmd/Ctrl+W closes (cancel).
 */
export function DialogFooter({
    cancelLabel = "Cancel",
    saveLabel = "Save",
    saveDisabled = false,
    saveLoading = false,
    leading,
    onCancel,
    onSave,
}: DialogFooterProps) {
    const windowsOrder = useMemo(() => System.IsWindows(), []);

    useEffect(() => {
        const onKeyDown = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                event.preventDefault();
                onCancel();
                return;
            }
            if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "w") {
                event.preventDefault();
                onCancel();
                return;
            }
            if (event.key === "Enter" && !event.metaKey && !event.ctrlKey && !event.altKey) {
                const target = event.target as HTMLElement | null;
                if (target && (target.tagName === "TEXTAREA" || target.isContentEditable)) {
                    return;
                }
                if (saveDisabled || saveLoading) {
                    return;
                }
                event.preventDefault();
                onSave();
            }
        };
        window.addEventListener("keydown", onKeyDown);
        return () => window.removeEventListener("keydown", onKeyDown);
    }, [onCancel, onSave, saveDisabled, saveLoading]);

    const cancel = (
        <Button key="cancel" color="secondary" size="md" onClick={onCancel}>
            {cancelLabel}
        </Button>
    );
    const save = (
        <Button key="save" color="primary" size="md" isDisabled={saveDisabled} isLoading={saveLoading} onClick={onSave}>
            {saveLabel}
        </Button>
    );

    return (
        <footer
            className={cx(
                "flex shrink-0 flex-col gap-3 border-t border-secondary px-6 py-4",
                !leading && "items-end",
            )}
        >
            {leading ? <div className="w-full text-sm text-tertiary">{leading}</div> : null}
            <div className="flex w-full items-center justify-end gap-2">
                {windowsOrder ? (
                    <>
                        {save}
                        {cancel}
                    </>
                ) : (
                    <>
                        {cancel}
                        {save}
                    </>
                )}
            </div>
        </footer>
    );
}

import type { ImgHTMLAttributes } from "react";
import routeDeckIcon from "@/assets/branding/route-deck-icon.png";
import { cx } from "@/utils/cx";

/** Compact RouteDeck app mark for desktop shell chrome (sidebar header). */
export function RouteDeckMark({ className, alt = "MACROX RouteDeck", ...props }: ImgHTMLAttributes<HTMLImageElement>) {
    return (
        <img
            src={routeDeckIcon}
            alt={alt}
            width={32}
            height={32}
            draggable={false}
            className={cx("size-8 shrink-0 rounded-lg object-cover select-none", className)}
            {...props}
        />
    );
}

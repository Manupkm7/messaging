import { forwardRef, HTMLAttributes } from "react";

export const Separator = forwardRef<
  HTMLDivElement,
  HTMLAttributes<HTMLDivElement> & {
    orientation?: "horizontal" | "vertical";
  }
>(({ className, orientation = "horizontal", ...props }, ref) => {
  return (
    <div
      ref={ref}
      className={`
        shrink-0 bg-border
        ${orientation === "horizontal" ? "h-px w-full" : "h-full w-px"}
        ${className}
      `}
      {...props}
    />
  );
});
Separator.displayName = "Separator";

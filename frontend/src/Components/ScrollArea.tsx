import { forwardRef, HTMLAttributes } from "react";

export const ScrollArea = forwardRef<
  HTMLDivElement,
  HTMLAttributes<HTMLDivElement>
>(({ className, children, ...props }, ref) => {
  return (
    <div ref={ref} className={`relative overflow-auto ${className}`} {...props}>
      {children}
    </div>
  );
});
ScrollArea.displayName = "ScrollArea";

import type React from "react";

import { forwardRef } from "react";

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "secondary" | "destructive" | "outline";
}

export const Badge = forwardRef<HTMLDivElement, BadgeProps>(
  ({ className, variant = "default", ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={`
        inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold transition-colors
        ${
          variant === "default"
            ? "bg-primary text-primary-foreground hover:bg-primary/80"
            : ""
        }
        ${
          variant === "secondary"
            ? "bg-secondary text-secondary-foreground hover:bg-secondary/80"
            : ""
        }
        ${
          variant === "destructive"
            ? "bg-destructive text-destructive-foreground hover:bg-destructive/80"
            : ""
        }
        ${variant === "outline" ? "border border-input" : ""}
        ${className}
      `}
        {...props}
      />
    );
  }
);
Badge.displayName = "Badge";

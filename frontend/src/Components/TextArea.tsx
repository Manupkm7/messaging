import { forwardRef, TextareaHTMLAttributes } from "react";

export type TextareaProps = TextareaHTMLAttributes<HTMLTextAreaElement>;

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, ...props }, ref) => {
    return (
      <textarea
        className={`
          flex min-h-20 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm
          placeholder:text-muted-foreground
          focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2
          disabled:cursor-not-allowed disabled:opacity-50
          resize-none
          ${className}
        `}
        ref={ref}
        {...props}
      />
    );
  }
);
Textarea.displayName = "Textarea";

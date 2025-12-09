import {
  createContext,
  useContext,
  useState,
  forwardRef,
  useEffect,
  useRef,
  ReactNode,
  HTMLAttributes,
} from "react";

interface DropdownMenuContextValue {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const DropdownMenuContext = createContext<DropdownMenuContextValue | undefined>(
  undefined
);

export function DropdownMenu({ children }: { children: ReactNode }) {
  const [open, setOpen] = useState(false);

  return (
    <DropdownMenuContext.Provider value={{ open, onOpenChange: setOpen }}>
      <div className="relative inline-block">{children}</div>
    </DropdownMenuContext.Provider>
  );
}

export function DropdownMenuTrigger({
  children,
  asChild,
}: {
  children: ReactNode;
  asChild?: boolean;
}) {
  const context = useContext(DropdownMenuContext);
  if (!context)
    throw new Error("DropdownMenuTrigger must be used within DropdownMenu");

  const handleClick = () => context.onOpenChange(!context.open);

  if (
    asChild &&
    typeof children === "object" &&
    children !== null &&
    "props" in children
  ) {
    return <div onClick={handleClick}>{children}</div>;
  }

  return <button onClick={handleClick}>{children}</button>;
}

export const DropdownMenuContent = forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement> & { align?: "start" | "end" }
>(({ className, align = "start", children, ...props }) => {
  const context = useContext(DropdownMenuContext);
  const contentRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        contentRef.current &&
        !contentRef.current.contains(event.target as Node)
      ) {
        context?.onOpenChange(false);
      }
    };

    if (context?.open) {
      document.addEventListener("mousedown", handleClickOutside);
      return () =>
        document.removeEventListener("mousedown", handleClickOutside);
    }
  }, [context?.open, context]);

  if (!context?.open) return null;

  return (
    <div
      ref={contentRef}
      className={`
        absolute animate-in fade-in-0 zoom-in-95 z-50 mt-2 min-w-32 rounded-lg border border-border bg-popover p-1 shadow-md
        ${align === "start" ? "left-0" : "right-0"}
        
        ${className}
      `}
      {...props}
    >
      {children}
    </div>
  );
});
DropdownMenuContent.displayName = "DropdownMenuContent";

export function DropdownMenuLabel({
  className,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={`px-2 py-1.5 text-sm font-semibold ${className}`}
      {...props}
    />
  );
}

export function DropdownMenuItem({
  className,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  const context = useContext(DropdownMenuContext);

  const handleClick = () => {
    context?.onOpenChange(false);
  };

  return (
    <div
      className={`
        relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none
        hover:bg-accent hover:text-accent-foreground
        focus:bg-accent focus:text-accent-foreground
        ${className}
      `}
      onClick={handleClick}
      {...props}
    />
  );
}

export function DropdownMenuSeparator({
  className,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={`-mx-1 my-1 h-px bg-border ${className}`} {...props} />
  );
}

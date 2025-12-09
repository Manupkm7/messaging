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

interface PopoverContextValue {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const PopoverContext = createContext<PopoverContextValue | undefined>(
  undefined
);

export function Popover({
  children,
  open: controlledOpen,
  onOpenChange,
}: {
  children: ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}) {
  const [uncontrolledOpen, setUncontrolledOpen] = useState(false);
  const open = controlledOpen ?? uncontrolledOpen;
  const handleOpenChange = onOpenChange ?? setUncontrolledOpen;

  return (
    <PopoverContext.Provider value={{ open, onOpenChange: handleOpenChange }}>
      <div className="relative inline-block">{children}</div>
    </PopoverContext.Provider>
  );
}

export function PopoverTrigger({
  children,
  asChild,
}: {
  children: ReactNode;
  asChild?: boolean;
}) {
  const context = useContext(PopoverContext);
  if (!context) throw new Error("PopoverTrigger must be used within Popover");

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

export const PopoverContent = forwardRef<
  HTMLDivElement,
  HTMLAttributes<HTMLDivElement> & { align?: "start" | "end" | "center" }
>(({ className, align = "center", children, ...props }) => {
  const context = useContext(PopoverContext);
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
        absolute bottom-full mb-2 z-50 rounded-lg border border-border bg-popover p-4 shadow-md
          ${align === "start" ? "left-0" : ""}
          ${align === "end" ? "right-0" : ""}
          ${align === "center" ? "left-1/2" : ""}
        animate-in fade-in-0 zoom-in-95
        ${className}
      `}
      {...props}
    >
      {children}
    </div>
  );
});
PopoverContent.displayName = "PopoverContent";

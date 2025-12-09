import { HTMLAttributes, ReactNode } from "react";

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
}

interface CardHeaderProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
}

interface CardContentProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
}

export function Card({ children, className = "", ...props }: CardProps) {
  return (
    <div
      className={`
        bg-white 
        border border-gray-200  
        rounded-lg shadow-sm 
        transition-all duration-200 
        hover:shadow-md
        ${className}
      `}
      role="region"
      data-testid="card"
      {...props}
    >
      {children}
    </div>
  );
}

export function CardHeader({
  children,
  className = "",
  ...props
}: CardHeaderProps) {
  return (
    <div
      className={`
        px-6 py-4 
        border-b border-gray-100 dark:border-gray-700
        ${className}
      `}
      role="heading"
      aria-level={2}
      data-testid="card-header"
      {...props}
    >
      {children}
    </div>
  );
}

export function CardContent({
  children,
  className = "",
  ...props
}: CardContentProps) {
  return (
    <div
      className={`
        px-6 py-4
        ${className}
      `}
      data-testid="card-content"
      {...props}
    >
      {children}
    </div>
  );
}

import { ReactNode } from "react";

import { Link } from "react-router-dom";

interface ButtonProps {
  children: ReactNode;
  onClick?: () => void;
  className?: string;
  icon?: ReactNode;
  page?: string;
  href?: string;
  iconPosition?: "left" | "right";
  type?: "button" | "submit" | "reset";
  disabled?: boolean;
  ariaLabel?: string;
  testId?: string;
  classNameParent?: string;
}

const Button = ({
  children,
  onClick,
  className,
  icon,
  page,
  href,
  iconPosition = "left",
  type = "button",
  disabled = false,
  ariaLabel,
  testId,
  classNameParent,
}: ButtonProps) =>
  page ? (
    <Link
      to={page ? `/${page}${href ? `/${href}` : ""}` : "#"}
      className={`min-w-[150px] inline-flex items-center select-none justify-center gap-1.5 overflow-hidden font-semibold border-2 px-2 capitalize text-black transition-all duration-300 dark:text-white ${className}`}
      aria-label={ariaLabel}
      data-testid={testId}
      onClick={onClick}
    >
      {icon && iconPosition === "left" && <>{icon}</>}

      <span className="inline-flex whitespace-nowrap">{children}</span>
      {icon && iconPosition === "right" && <>{icon}</>}
    </Link>
  ) : (
    <div
      className={`flex h-full items-center font-semibold ${classNameParent}`}
    >
      <button
        type={type}
        disabled={disabled}
        onClick={onClick}
        className={`min-w-[150px] inline-flex items-center justify-center select-none gap-1.5 overflow-hidden border-2 px-2 font-medium capitalize text-black transition-all duration-300 dark:text-white ${className}`}
        aria-label={ariaLabel}
        aria-disabled={disabled}
        data-testid={testId}
      >
        {icon && iconPosition === "left" ? <>{icon}</> : null}

        <span className="inline-flex whitespace-nowrap font-semibold">
          {children}
        </span>
        {icon && iconPosition === "right" ? <>{icon}</> : null}
      </button>
    </div>
  );

export default Button;

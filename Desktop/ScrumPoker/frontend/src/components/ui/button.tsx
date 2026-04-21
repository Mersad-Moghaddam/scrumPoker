import { ButtonHTMLAttributes } from "react";

import { cn } from "../../lib/utils";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost";
};

const variants: Record<NonNullable<ButtonProps["variant"]>, string> = {
  primary:
    "bg-brand-500 text-white shadow-lg shadow-brand-500/30 hover:bg-brand-400 disabled:opacity-50 disabled:hover:bg-brand-500",
  secondary:
    "bg-surface-2 text-app-text hover:bg-surface-3 disabled:opacity-50 disabled:hover:bg-surface-2",
  ghost:
    "bg-transparent text-app-text hover:bg-surface-2 disabled:opacity-50 disabled:hover:bg-transparent"
};

export function Button({ className, variant = "primary", ...props }: ButtonProps) {
  return (
    <button
      className={cn(
        "inline-flex items-center justify-center rounded-2xl px-4 py-3 text-sm font-semibold transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-brand-400/60",
        variants[variant],
        className
      )}
      {...props}
    />
  );
}

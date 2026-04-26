import { HTMLAttributes } from "react";

import { cn } from "../../lib/utils";

export function Card({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "rounded-[28px] border border-white/10 bg-surface-1/90 p-5 shadow-glow backdrop-blur-xl",
        className
      )}
      {...props}
    />
  );
}

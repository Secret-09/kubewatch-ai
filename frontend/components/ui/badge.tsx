import * as React from "react";
import { cn } from "@/lib/utils";

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: "default" | "success" | "warning" | "danger" | "muted";
}

const variants = {
  default: "bg-slate-800 text-slate-200 border border-white/10",
  success: "bg-emerald-500/15 text-emerald-300 border border-emerald-500/20",
  warning: "bg-amber-500/15 text-amber-300 border border-amber-500/20",
  danger: "bg-rose-500/15 text-rose-300 border border-rose-500/20",
  muted: "bg-white/5 text-slate-400 border border-slate-700",
};

export function Badge({ className, variant = "default", ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-3 py-1 text-xs font-medium tracking-wide",
        variants[variant],
        className,
      )}
      {...props}
    />
  );
}

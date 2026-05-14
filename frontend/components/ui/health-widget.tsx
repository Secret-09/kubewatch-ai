import { Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { cn } from "@/lib/utils";

const statusStyles = {
  healthy: "text-emerald-300",
  degraded: "text-amber-300",
  danger: "text-rose-300",
};

interface HealthWidgetProps {
  title: string;
  status: "healthy" | "degraded" | "danger";
  value: string;
  detail: string;
}

export function HealthWidget({ title, status, value, detail }: HealthWidgetProps) {
  return (
    <Card className="border-white/10 bg-slate-950/80">
      <CardHeader>
        <span className="text-xs uppercase tracking-[0.24em] text-slate-500">{title}</span>
        <CardTitle className={cn("text-4xl tracking-tight", statusStyles[status])}>{value}</CardTitle>
      </CardHeader>
      <CardDescription>{detail}</CardDescription>
    </Card>
  );
}

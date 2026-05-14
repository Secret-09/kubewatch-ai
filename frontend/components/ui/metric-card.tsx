import { Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { cn } from "@/lib/utils";

interface MetricCardProps {
  label: string;
  value: string;
  note?: string;
  accent?: string;
}

export function MetricCard({ label, value, note, accent, className }: MetricCardProps & { className?: string }) {
  return (
    <Card className={cn("border-transparent bg-slate-950/80", className)}>
      <CardHeader>
        <span className="text-xs uppercase tracking-[0.24em] text-slate-500">{label}</span>
        <CardTitle className={cn("text-3xl tracking-tight", accent)}>{value}</CardTitle>
      </CardHeader>
      {note ? <CardDescription>{note}</CardDescription> : null}
    </Card>
  );
}

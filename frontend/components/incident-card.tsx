import { Badge } from "@/components/ui/badge";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";

interface IncidentCardProps {
  id: string;
  namespace: string;
  category: string;
  summary: string;
  severity: "critical" | "high" | "medium" | "low";
  lastSeen: string;
  remediation: string;
}

type IncidentSeverity = "critical" | "high" | "medium" | "low";

const severityMap: Record<IncidentSeverity, "default" | "danger" | "warning" | "muted"> = {
  critical: "danger",
  high: "warning",
  medium: "default",
  low: "muted",
};

export function IncidentCard({ id, namespace, category, summary, severity, lastSeen, remediation }: IncidentCardProps) {
  return (
    <Card className="border-white/10 bg-slate-950/75">
      <CardHeader className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <CardTitle className="text-lg">{summary}</CardTitle>
          <p className="text-sm text-slate-400">{namespace} • {category}</p>
        </div>
        <Badge variant={severityMap[severity]}>{severity.toUpperCase()}</Badge>
      </CardHeader>
      <CardContent>
        <div className="grid gap-4 sm:grid-cols-[1fr_auto]">
          <div>
            <p className="text-sm text-slate-300">Remediation: {remediation}</p>
          </div>
          <div className="rounded-3xl bg-slate-900/70 px-4 py-3 text-xs uppercase tracking-[0.24em] text-slate-500">
            Last seen {lastSeen}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

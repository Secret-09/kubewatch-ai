import DashboardShell from "@/components/dashboard-shell";
import { IncidentCard } from "@/components/incident-card";
import { AlertsPanel } from "@/components/alerts-panel";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

type IncidentSeverity = "critical" | "high" | "medium" | "low";

const incidents: Array<{
  id: string;
  namespace: string;
  category: string;
  summary: string;
  severity: IncidentSeverity;
  lastSeen: string;
  remediation: string;
}> = [
  {
    id: "INC-001",
    namespace: "payments",
    category: "CrashLoopBackOff",
    summary: "Checkout pod restarting repeatedly",
    severity: "critical",
    lastSeen: "2m ago",
    remediation: "Inspect container logs and scale deployment if necessary.",
  },
  {
    id: "INC-002",
    namespace: "web",
    category: "Replica mismatch",
    summary: "Frontend deployment under-provisioned",
    severity: "high",
    lastSeen: "7m ago",
    remediation: "Validate deployment replica count and check node capacity.",
  },
  {
    id: "INC-003",
    namespace: "api",
    category: "Unhealthy Pod",
    summary: "Latency spike correlated with restarts",
    severity: "medium",
    lastSeen: "15m ago",
    remediation: "Inspect pod resource limits and service dependencies.",
  },
];

export default function IncidentsPage() {
  return (
    <DashboardShell title="Incidents" description="Track active and historical incidents across your Kubernetes clusters.">
      <div className="grid gap-8 xl:grid-cols-[1.25fr_0.75fr]">
        <section className="space-y-6">
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Incident Response</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Prioritize live alerts, review incident context, and coordinate remediation actions with a single observability workflow.</p>
            </CardContent>
          </Card>
          <div className="space-y-5">
            {incidents.map((incident) => (
              <IncidentCard key={incident.id} {...incident} />
            ))}
          </div>
        </section>

        <aside className="space-y-6">
          <AlertsPanel />
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Incident Triage</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Use this panel to surface incident severity, affected namespaces, and recommended remediation steps for faster operational recovery.</p>
            </CardContent>
          </Card>
        </aside>
      </div>
    </DashboardShell>
  );
}

import DashboardShell from "@/components/dashboard-shell";
import { MetricCard } from "@/components/ui/metric-card";
import { HealthWidget } from "@/components/ui/health-widget";
import { AlertsPanel } from "@/components/alerts-panel";
import { MetricsChart } from "@/components/metrics-chart";
import { IncidentCard } from "@/components/incident-card";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ShieldCheck, Activity, ServerCog } from "lucide-react";

const workloadHealth = [
  { label: "Payments API", value: "Degraded", detail: "Replica mismatch detected", accent: "bg-amber-400/20" },
  { label: "Frontend", value: "Healthy", detail: "Stable at 12 replicas", accent: "bg-emerald-400/15" },
  { label: "User service", value: "Warning", detail: "Restart spike observed", accent: "bg-rose-400/15" },
];

type IncidentSeverity = "critical" | "high" | "medium" | "low";

const recentIncidents: Array<{
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
    summary: "Checkout service restarted 6 times in 10 minutes",
    severity: "critical",
    lastSeen: "2m ago",
    remediation: "Review logs and rollout a recovery image.",
  },
  {
    id: "INC-002",
    namespace: "api",
    category: "High restart count",
    summary: "Cache worker pod reporting frequent restarts",
    severity: "high",
    lastSeen: "8m ago",
    remediation: "Inspect memory and lifecycle probes.",
  },
];

export default function Home() {
  return (
    <DashboardShell title="Dashboard" description="Kubernetes observability for cloud-native teams, with real-time incident awareness and cluster health insights.">
      <div className="grid gap-8 xl:grid-cols-[1.25fr_0.85fr]">
        <section className="space-y-6">
          <div className="grid gap-6 sm:grid-cols-2 xl:grid-cols-3">
            <MetricCard label="Active incidents" value="8" note="Critical and warning incidents in the last hour" accent="text-emerald-300" />
            <MetricCard label="Namespaces" value="14" note="Segmentation across cluster workloads" accent="text-slate-100" />
            <MetricCard label="Cluster status" value="Degraded" note="Service-level health requires attention" accent="text-amber-300" />
          </div>

          <MetricsChart title="Incident cadence" labels={["1h", "2h", "3h", "4h", "5h", "6h"]} values={[3, 5, 6, 4, 7, 8]} accent="bg-emerald-400/80" />

          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <div className="flex items-center gap-3 text-slate-100">
                <ShieldCheck className="h-5 w-5 text-emerald-400" />
                <CardTitle>Cluster overview</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 md:grid-cols-3">
                <div className="rounded-3xl bg-slate-900/70 p-5">
                  <p className="text-sm uppercase tracking-[0.24em] text-slate-500">Nodes</p>
                  <p className="mt-3 text-3xl font-semibold text-white">12</p>
                </div>
                <div className="rounded-3xl bg-slate-900/70 p-5">
                  <p className="text-sm uppercase tracking-[0.24em] text-slate-500">Pods</p>
                  <p className="mt-3 text-3xl font-semibold text-white">184</p>
                </div>
                <div className="rounded-3xl bg-slate-900/70 p-5">
                  <p className="text-sm uppercase tracking-[0.24em] text-slate-500">Workloads</p>
                  <p className="mt-3 text-3xl font-semibold text-white">36</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>

        <aside className="space-y-6">
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <div className="flex items-center gap-3 text-slate-100">
                <Activity className="h-5 w-5 text-cyan-400" />
                <CardTitle>Realtime alerts</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <AlertsPanel />
            </CardContent>
          </Card>

          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <div className="flex items-center gap-3 text-slate-100">
                <ServerCog className="h-5 w-5 text-violet-400" />
                <CardTitle>Deployment health</CardTitle>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {workloadHealth.map((item) => (
                  <div key={item.label} className="rounded-3xl bg-slate-900/70 p-4">
                    <div className="flex items-center justify-between gap-3">
                      <div>
                        <p className="text-sm text-slate-400">{item.label}</p>
                        <p className="mt-1 text-lg font-semibold text-white">{item.value}</p>
                      </div>
                      <span className={item.accent + " inline-flex h-10 min-w-[84px] items-center justify-center rounded-full text-sm font-medium text-emerald-200"}>
                        {item.detail}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </aside>
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
        <section className="space-y-6">
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Incident priorities</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-5">
                {recentIncidents.map((incident) => (
                  <IncidentCard key={incident.id} {...incident} />
                ))}
              </div>
            </CardContent>
          </Card>
        </section>
      </div>
    </DashboardShell>
  );
}

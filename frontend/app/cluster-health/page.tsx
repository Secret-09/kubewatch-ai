import DashboardShell from "@/components/dashboard-shell";
import { HealthWidget } from "@/components/ui/health-widget";
import { MetricsChart } from "@/components/metrics-chart";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { AlertsPanel } from "@/components/alerts-panel";

type HealthStatus = "healthy" | "degraded" | "danger";

const healthWidgets: Array<{ title: string; status: HealthStatus; value: string; detail: string }> = [
  { title: "Cluster Status", status: "degraded", value: "Degraded", detail: "2 critical incidents in the last 15 minutes." },
  { title: "Node readiness", status: "healthy", value: "92%", detail: "Most nodes are healthy and schedulable." },
  { title: "Workload health", status: "degraded", value: "78%", detail: "Several deployments reporting replica drift." },
];

export default function ClusterHealthPage() {
  return (
    <DashboardShell title="Cluster Health" description="Understand service reliability and cluster health across your Kubernetes estate.">
      <div className="grid gap-8 xl:grid-cols-[1.5fr_0.9fr]">
        <section className="space-y-6">
          <div className="grid gap-6 md:grid-cols-3">
            {healthWidgets.map((widget) => (
              <HealthWidget key={widget.title} {...widget} />
            ))}
          </div>

          <MetricsChart
            title="Incident trend"
            labels={["12h", "9h", "6h", "3h", "1h", "Now"]}
            values={[4, 5, 6, 9, 7, 10]}
            accent="bg-emerald-400/80"
          />

          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Cluster overview</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="rounded-3xl bg-slate-900/70 p-5 text-slate-300">
                  <p className="text-sm uppercase tracking-[0.24em] text-slate-500">Namespaces</p>
                  <p className="mt-3 text-3xl font-semibold text-white">14</p>
                </div>
                <div className="rounded-3xl bg-slate-900/70 p-5 text-slate-300">
                  <p className="text-sm uppercase tracking-[0.24em] text-slate-500">Pods</p>
                  <p className="mt-3 text-3xl font-semibold text-white">184</p>
                </div>
                <div className="rounded-3xl bg-slate-900/70 p-5 text-slate-300">
                  <p className="text-sm uppercase tracking-[0.24em] text-slate-500">Nodes</p>
                  <p className="mt-3 text-3xl font-semibold text-white">12</p>
                </div>
                <div className="rounded-3xl bg-slate-900/70 p-5 text-slate-300">
                  <p className="text-sm uppercase tracking-[0.24em] text-slate-500">Active incidents</p>
                  <p className="mt-3 text-3xl font-semibold text-white">8</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </section>

        <aside className="space-y-6">
          <AlertsPanel />
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Operational guidance</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Monitor incident density and cluster health indicators to prioritize remediation and maintain deployment stability.</p>
            </CardContent>
          </Card>
        </aside>
      </div>
    </DashboardShell>
  );
}

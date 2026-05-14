import DashboardShell from "@/components/dashboard-shell";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

type WorkloadHealth = "healthy" | "degraded" | "warning";

const workloads: Array<{ name: string; namespace: string; replicas: string; health: WorkloadHealth }> = [
  { name: "payments-api", namespace: "payments", replicas: "6/8", health: "degraded" },
  { name: "web-frontend", namespace: "web", replicas: "12/12", health: "healthy" },
  { name: "user-service", namespace: "api", replicas: "4/5", health: "warning" },
  { name: "auth-proxy", namespace: "platform", replicas: "3/3", health: "healthy" },
];

const healthMap: Record<WorkloadHealth, "success" | "warning" | "danger"> = {
  healthy: "success",
  degraded: "warning",
  warning: "danger",
};

export default function WorkloadsPage() {
  return (
    <DashboardShell title="Workloads" description="Inspect deployment health, scaling patterns, and workload reliability.">
      <div className="grid gap-8 xl:grid-cols-[1.2fr_0.8fr]">
        <section className="space-y-6">
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Deployment health</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Review workload status across namespaces and identify deployments that require scaling or remediation.</p>
            </CardContent>
          </Card>
          <div className="space-y-4">
            {workloads.map((workload) => (
              <Card key={workload.name} className="border-white/10 bg-slate-950/75">
                <CardHeader className="flex items-center justify-between gap-4">
                  <div>
                    <CardTitle className="text-base">{workload.name}</CardTitle>
                    <p className="text-sm text-slate-500">{workload.namespace}</p>
                  </div>
                  <Badge variant={healthMap[workload.health]}>{workload.health.toUpperCase()}</Badge>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-slate-300">Replica set: {workload.replicas}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <aside className="space-y-6">
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Workload summary</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Workload health widgets help teams identify amplification risks and performance regressions in production clusters.</p>
            </CardContent>
          </Card>
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Recommended actions</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="list-disc space-y-2 pl-5 text-slate-400">
                <li>Scale deployments with failing replicas.</li>
                <li>Review node pressure and pod eviction behavior.</li>
                <li>Balance service availability with resource limits.</li>
              </ul>
            </CardContent>
          </Card>
        </aside>
      </div>
    </DashboardShell>
  );
}

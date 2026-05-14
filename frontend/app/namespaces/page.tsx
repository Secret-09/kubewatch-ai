import DashboardShell from "@/components/dashboard-shell";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

type NamespaceHealth = "healthy" | "degraded" | "warning";

const namespaces: Array<{ name: string; pods: number; health: NamespaceHealth }> = [
  { name: "payments", pods: 28, health: "degraded" },
  { name: "web", pods: 18, health: "healthy" },
  { name: "api", pods: 35, health: "degraded" },
  { name: "platform", pods: 12, health: "healthy" },
  { name: "database", pods: 9, health: "warning" },
];

const healthMap: Record<NamespaceHealth, "success" | "warning" | "danger"> = {
  healthy: "success",
  degraded: "warning",
  warning: "danger",
};

export default function NamespacesPage() {
  return (
    <DashboardShell title="Namespaces" description="Visualize namespace health and resource distribution across the cluster.">
      <div className="grid gap-8 lg:grid-cols-[1.1fr_0.9fr]">
        <section className="space-y-6">
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Namespace health</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Namespace-level visibility helps teams understand which application domains require remediation or scaling attention.</p>
            </CardContent>
          </Card>
          <div className="grid gap-4 sm:grid-cols-2">
            {namespaces.map((namespace) => (
              <Card key={namespace.name} className="border-white/10 bg-slate-950/75">
                <CardHeader className="flex items-center justify-between gap-4">
                  <CardTitle className="text-base">{namespace.name}</CardTitle>
                  <Badge variant={healthMap[namespace.health]}>{namespace.health.toUpperCase()}</Badge>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-slate-400">{namespace.pods} pods deployed</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </section>

        <aside className="space-y-6">
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Namespace trends</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Namespace-level health and pod density support cloud-native operations, isolating faults before they affect broader services.</p>
            </CardContent>
          </Card>
          <Card className="border-white/10 bg-slate-950/80">
            <CardHeader>
              <CardTitle>Usage guidance</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-slate-400">Use namespaced views to track service ownership, cluster policies, and workload segmentation in production environments.</p>
            </CardContent>
          </Card>
        </aside>
      </div>
    </DashboardShell>
  );
}

import Link from "next/link";
import { ReactNode } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

const navItems = [
  { href: "/", label: "Dashboard" },
  { href: "/incidents", label: "Incidents" },
  { href: "/cluster-health", label: "Cluster Health" },
  { href: "/namespaces", label: "Namespaces" },
  { href: "/workloads", label: "Workloads" },
];

interface DashboardShellProps {
  title: string;
  description: string;
  children: ReactNode;
}

export default function DashboardShell({ title, description, children }: DashboardShellProps) {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <div className="mx-auto flex min-h-screen max-w-[1600px] flex-col gap-10 px-6 py-8 xl:px-12">
        <header className="flex flex-col gap-5 rounded-[2rem] border border-white/10 bg-slate-900/80 p-6 shadow-[0_30px_120px_rgba(15,23,42,0.35)] lg:flex-row lg:items-center lg:justify-between">
          <div className="space-y-3">
            <p className="text-xs uppercase tracking-[0.32em] text-emerald-300">KubeWatch AI</p>
            <h1 className="text-3xl font-semibold tracking-tight text-white sm:text-4xl">{title}</h1>
            <p className="max-w-2xl text-slate-400">{description}</p>
          </div>
          <div className="flex flex-wrap items-center gap-3">
            {navItems.map((item) => (
              <Link key={item.href} href={item.href} className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-100 transition hover:border-emerald-400/40 hover:bg-white/10">
                {item.label}
              </Link>
            ))}
            <Button variant="ghost">Open Telemetry</Button>
          </div>
        </header>

        <main className="grid gap-8">{children}</main>

        <Card className="border-white/10 bg-slate-900/65 p-6 text-slate-400">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <p className="text-sm">Built for cloud-native teams managing Kubernetes observability, incident response, and real-time operational intelligence.</p>
            <div className="flex flex-wrap items-center gap-3 text-xs uppercase tracking-[0.3em] text-slate-500">
              <span>Dark mode</span>
              <span>Responsive</span>
              <span>API-first</span>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}

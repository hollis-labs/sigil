"use client";

import { useMemo } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { SigilChart } from "@/components/sigil-chart";
import type { Repo } from "@/types/repo";

const FILLS = [
  "oklch(0.623 0.214 259.815)",
  "oklch(0.6 0.118 184.704)",
  "oklch(0.828 0.189 84.429)",
  "oklch(0.769 0.188 70.08)",
  "oklch(0.646 0.222 41.116)",
  "oklch(0.55 0.02 0)",
  "oklch(0.7 0.15 200)",
  "oklch(0.65 0.2 320)",
];

interface ReposByCategoryChartProps {
  repos: Repo[];
}

export function ReposByCategoryChart({ repos }: ReposByCategoryChartProps) {
  const data = useMemo(() => {
    const counts: Record<string, number> = {};
    for (const r of repos) {
      const cat = String(r.category ?? "Other");
      counts[cat] = (counts[cat] ?? 0) + 1;
    }
    return Object.entries(counts).map(([category, count], i) => ({
      category,
      count,
      fill: FILLS[i % FILLS.length],
    }));
  }, [repos]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Repos by Category</CardTitle>
        <CardDescription>Distribution across repo types</CardDescription>
      </CardHeader>
      <CardContent>
        <SigilChart
          type="donut"
          xKey="category"
          height={250}
          showLegend
          dataKeys={["count"]}
          config={{ count: { label: "Repos", color: "var(--chart-1)" } }}
          data={data}
        />
      </CardContent>
    </Card>
  );
}

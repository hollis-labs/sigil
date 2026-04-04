"use client";

import { useMemo } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { SigilChart } from "@/components/sigil-chart";
import type { Scorecard } from "@/types/scorecard";

interface ScoreDistributionChartProps {
  scorecards: Scorecard[];
}

export function ScoreDistributionChart({ scorecards }: ScoreDistributionChartProps) {
  const data = useMemo(() => {
    const buckets = [
      { range: "0-2", count: 0 },
      { range: "2-4", count: 0 },
      { range: "4-6", count: 0 },
      { range: "6-8", count: 0 },
      { range: "8-10", count: 0 },
    ];
    for (const sc of scorecards) {
      const score = Number(sc.overall ?? 0);
      const idx = Math.min(Math.floor(score / 2), 4);
      buckets[idx].count++;
    }
    return buckets;
  }, [scorecards]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Score Distribution</CardTitle>
        <CardDescription>Repos by overall score range</CardDescription>
      </CardHeader>
      <CardContent>
        <SigilChart
          type="bar"
          xKey="range"
          height={250}
          showGrid
          dataKeys={["count"]}
          config={{ count: { label: "Repos", color: "var(--chart-1)" } }}
          data={data}
        />
      </CardContent>
    </Card>
  );
}

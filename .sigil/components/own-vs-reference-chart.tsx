"use client";

import { useMemo } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { SigilChart } from "@/components/sigil-chart";
import type { Repo } from "@/types/repo";
import type { Scorecard } from "@/types/scorecard";
import type { DimensionScore } from "@/types/dimensionscore";

interface OwnVsReferenceChartProps {
  repos: Repo[];
  scorecards: Scorecard[];
  dimensionScores: DimensionScore[];
}

export function OwnVsReferenceChart({ repos, scorecards, dimensionScores }: OwnVsReferenceChartProps) {
  const data = useMemo(() => {
    const ownRepoIds = new Set(repos.filter((r) => r.is_own).map((r) => r.id));
    const scorecardRepoMap = new Map(scorecards.map((sc) => [sc.id, sc.repo_id]));
    const dimAccum: Record<string, { ownSum: number; ownCount: number; refSum: number; refCount: number }> = {};
    for (const ds of dimensionScores) {
      const dimName = ds.dimension_name ?? "Unknown";
      const repoId = scorecardRepoMap.get(ds.scorecard_id ?? "");
      const score = Number(ds.score ?? 0);
      if (!dimAccum[dimName]) dimAccum[dimName] = { ownSum: 0, ownCount: 0, refSum: 0, refCount: 0 };
      if (repoId && ownRepoIds.has(repoId)) {
        dimAccum[dimName].ownSum += score;
        dimAccum[dimName].ownCount++;
      } else {
        dimAccum[dimName].refSum += score;
        dimAccum[dimName].refCount++;
      }
    }
    return Object.entries(dimAccum)
      .map(([dimension, v]) => ({
        dimension,
        own: v.ownCount ? Math.round((v.ownSum / v.ownCount) * 10) / 10 : 0,
        reference: v.refCount ? Math.round((v.refSum / v.refCount) * 10) / 10 : 0,
      }))
      .slice(0, 6);
  }, [repos, scorecards, dimensionScores]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Own Projects vs References</CardTitle>
        <CardDescription>Average scores by dimension</CardDescription>
      </CardHeader>
      <CardContent>
        <SigilChart
          type="bar"
          xKey="dimension"
          height={250}
          showGrid
          showLegend
          dataKeys={["own", "reference"]}
          config={{
            own: { label: "Own Projects", color: "oklch(0.623 0.214 259.815)" },
            reference: { label: "Reference Avg", color: "oklch(0.55 0.02 0)" },
          }}
          data={data}
        />
      </CardContent>
    </Card>
  );
}

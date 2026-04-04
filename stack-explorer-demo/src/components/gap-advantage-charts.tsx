"use client";

import { useMemo } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { SigilChart } from "@/components/sigil-chart";
import type { Repo } from "@/types/repo";
import type { Scorecard } from "@/types/scorecard";
import type { DimensionScore } from "@/types/dimensionscore";

interface GapAdvantageChartsProps {
  repos: Repo[];
  scorecards: Scorecard[];
  dimensionScores: DimensionScore[];
}

export function GapAdvantageCharts({ repos, scorecards, dimensionScores }: GapAdvantageChartsProps) {
  const { advantages, gaps } = useMemo(() => {
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
    const deltas = Object.entries(dimAccum).map(([dimension, v]) => {
      const ownAvg = v.ownCount ? v.ownSum / v.ownCount : 0;
      const refAvg = v.refCount ? v.refSum / v.refCount : 0;
      return { dimension, delta: Math.round((ownAvg - refAvg) * 10) / 10 };
    });
    return {
      advantages: deltas
        .filter((d) => d.delta > 0)
        .sort((a, b) => b.delta - a.delta)
        .map((d) => ({ dimension: d.dimension, delta: d.delta })),
      gaps: deltas
        .filter((d) => d.delta < 0)
        .sort((a, b) => a.delta - b.delta)
        .map((d) => ({ dimension: d.dimension, gap: Math.abs(d.delta) })),
    };
  }, [repos, scorecards, dimensionScores]);

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <Card>
        <CardHeader>
          <CardTitle>Advantages</CardTitle>
          <CardDescription>Dimensions where own projects outperform references</CardDescription>
        </CardHeader>
        <CardContent>
          <SigilChart
            type="bar"
            xKey="dimension"
            height={250}
            showGrid
            dataKeys={["delta"]}
            config={{ delta: { label: "Advantage", color: "oklch(0.75 0.18 155)" } }}
            data={advantages.length > 0 ? advantages : [{ dimension: "No data", delta: 0 }]}
          />
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Gaps</CardTitle>
          <CardDescription>Dimensions where own projects trail references</CardDescription>
        </CardHeader>
        <CardContent>
          <SigilChart
            type="bar"
            xKey="dimension"
            height={250}
            showGrid
            dataKeys={["gap"]}
            config={{ gap: { label: "Gap", color: "oklch(0.65 0.2 25)" } }}
            data={gaps.length > 0 ? gaps : [{ dimension: "No data", gap: 0 }]}
          />
        </CardContent>
      </Card>
    </div>
  );
}

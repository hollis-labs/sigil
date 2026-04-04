"use client";

import { useState, useMemo } from "react";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { GitCompare } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet";
import { DataTable } from "@/components/data-table";
import { SigilChart } from "@/components/sigil-chart";
import type { ChartConfig } from "@/components/ui/chart";

import type { Scorecard } from "@/types/scorecard";
import type { DimensionScore } from "@/types/dimensionscore";

const COMPARE_COLORS = [
  "oklch(0.623 0.214 259.815)",
  "oklch(0.6 0.118 184.704)",
  "oklch(0.828 0.189 84.429)",
  "oklch(0.769 0.188 70.08)",
  "oklch(0.646 0.222 41.116)",
];

interface ScorecardCompareProps {
  scorecards: Scorecard[];
  dimensionScores: DimensionScore[];
}

export function ScorecardCompare({ scorecards, dimensionScores }: ScorecardCompareProps) {
  const router = useRouter();
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [compareOpen, setCompareOpen] = useState(false);

  const toggleSelection = (id: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else if (next.size < 5) next.add(id);
      else toast("Select up to 5 repos to compare");
      return next;
    });
  };

  const selectedScorecards = useMemo(
    () => scorecards.filter((sc) => selectedIds.has(sc.id)),
    [scorecards, selectedIds],
  );

  const compareData = useMemo(() => {
    if (selectedScorecards.length < 2)
      return { radar: [], dimensions: [] as string[], config: {} as ChartConfig, table: [] as Record<string, unknown>[] };
    const scIds = new Set(selectedScorecards.map((sc) => sc.id));
    const relevantScores = dimensionScores.filter((ds) => scIds.has(ds.scorecard_id ?? ""));
    const dimSet = new Set(relevantScores.map((ds) => ds.dimension_name ?? "Unknown"));
    const dimensions = [...dimSet].sort();
    const scoreLookup = new Map<string, Map<string, number>>();
    for (const ds of relevantScores) {
      const scId = ds.scorecard_id ?? "";
      if (!scoreLookup.has(scId)) scoreLookup.set(scId, new Map());
      scoreLookup.get(scId)!.set(ds.dimension_name ?? "Unknown", Number(ds.score ?? 0));
    }
    const radar = dimensions.map((dim) => {
      const point: Record<string, unknown> = { dimension: dim };
      for (const sc of selectedScorecards) {
        point[sc.repo_name ?? sc.id] = scoreLookup.get(sc.id)?.get(dim) ?? 0;
      }
      return point;
    });
    const config: ChartConfig = {};
    selectedScorecards.forEach((sc, i) => {
      config[sc.repo_name ?? sc.id] = {
        label: sc.repo_name ?? sc.id,
        color: COMPARE_COLORS[i % COMPARE_COLORS.length],
      };
    });
    const table = dimensions.map((dim) => {
      const row: Record<string, unknown> = { dimension: dim };
      for (const sc of selectedScorecards) {
        row[sc.repo_name ?? sc.id] = scoreLookup.get(sc.id)?.get(dim) ?? null;
      }
      return row;
    });
    return { radar, dimensions, config, table };
  }, [selectedScorecards, dimensionScores]);

  return (
    <>
      <Card>
        <CardHeader>
          <div className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>All Scorecards</CardTitle>
              <CardDescription>Select repos to compare across dimensions</CardDescription>
            </div>
            <Button variant="default" disabled={selectedIds.size < 2} onClick={() => setCompareOpen(true)}>
              <GitCompare className="mr-2 h-4 w-4" />
              Compare{selectedIds.size > 0 ? ` (${selectedIds.size})` : ""}
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex flex-row gap-4 items-center mb-4">
            <p className="text-sm text-muted-foreground">Score key:</p>
            <Badge variant="outline" className="border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">8-10 Strong</Badge>
            <Badge variant="outline" className="border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400">5-7 Average</Badge>
            <Badge variant="destructive">0-4 Weak</Badge>
            <Badge variant="default">— Unscored</Badge>
          </div>
          <DataTable
            data={scorecards}
            columns={[
              {
                id: "select",
                header: () => null,
                cell: ({ row }: { row: { original: { id: string } } }) => (
                  <div onClick={(e: React.MouseEvent) => e.stopPropagation()}>
                    <Checkbox
                      checked={selectedIds.has(row.original.id)}
                      onCheckedChange={() => toggleSelection(row.original.id)}
                    />
                  </div>
                ),
              },
              { accessorKey: "repo_name", header: "Repo" },
              { accessorKey: "overall", header: "Overall" },
              {
                accessorKey: "lens_name",
                header: "Lens",
                cell: ({ row }: { row: { getValue: (k: string) => unknown } }) => {
                  const v = String(row.getValue("lens_name") ?? "");
                  const variants: Record<string, string> = { "Leadership": "warning", "Security": "danger", "Community": "success", "Maintainability": "info", "Engineering": "info", "Overall": "default", "Innovation": "success", "Adoption": "warning", "Product": "default" };
                  const variant = variants[v] ?? variants[v.toLowerCase()] ?? "default";
                  return <Badge variant={variant === "danger" ? "destructive" : variant === "success" || variant === "warning" || variant === "info" ? "outline" : "secondary"} className={
                    variant === "success" ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-400" :
                    variant === "warning" ? "border-amber-500/30 bg-amber-500/10 text-amber-400" :
                    variant === "info" ? "border-blue-500/30 bg-blue-500/10 text-blue-400" :
                    variant === "danger" ? "" : ""
                  }>{v}</Badge>;
                },
              },
              { accessorKey: "scored_at", header: "Scored" },
            ]}
            emptyMessage="No matching results found."
            onRowClick={(row) => router.push(`/repo-detail/${(row as { id: string }).id}`)}
          />
        </CardContent>
      </Card>
      <Sheet open={compareOpen} onOpenChange={setCompareOpen}>
        <SheetContent
          side="right"
          className="overflow-y-auto"
          style={{ maxWidth: Math.min(672 + Math.max(0, selectedScorecards.length - 3) * 120, 912) }}
        >
          <SheetHeader>
            <SheetTitle>Compare Scorecards</SheetTitle>
            <SheetDescription>
              {selectedScorecards.map((sc) => sc.repo_name).join(" vs ")}
            </SheetDescription>
          </SheetHeader>
          <div className="flex flex-col gap-6 p-4">
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Overall Scores</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-row gap-3 justify-center">
                  {selectedScorecards.map((sc, i) => (
                    <div
                      key={sc.id}
                      className="flex flex-col items-center gap-1 rounded-lg px-5 py-3"
                      style={{ backgroundColor: COMPARE_COLORS[i % COMPARE_COLORS.length] }}
                    >
                      <div className="text-3xl font-bold text-white">
                        {Number(sc.overall ?? 0).toFixed(1)}
                      </div>
                      <div className="text-xs font-medium text-white/80">{sc.repo_name}</div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
            {compareData.radar.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Dimension Radar</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-center">
                    <SigilChart
                      type="radar"
                      xKey="dimension"
                      height={300}
                      showLegend
                      dataKeys={selectedScorecards.map((sc) => sc.repo_name ?? sc.id)}
                      config={compareData.config}
                      data={compareData.radar}
                    />
                  </div>
                </CardContent>
              </Card>
            )}
            {compareData.table.length > 0 && (
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Dimension Scores</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                      <thead>
                        <tr className="border-b">
                          <th className="text-left py-2 pr-4 font-medium text-muted-foreground">Dimension</th>
                          {selectedScorecards.map((sc, i) => (
                            <th key={sc.id} className="text-center py-2 px-3 font-medium" style={{ color: COMPARE_COLORS[i % COMPARE_COLORS.length] }}>
                              {sc.repo_name}
                            </th>
                          ))}
                          {selectedScorecards.length === 2 && (
                            <th className="text-center py-2 px-3 font-medium text-muted-foreground">Delta</th>
                          )}
                        </tr>
                      </thead>
                      <tbody>
                        {compareData.table.map((row) => {
                          const scores = selectedScorecards.map(
                            (sc) => row[sc.repo_name ?? sc.id] as number | null,
                          );
                          const delta =
                            selectedScorecards.length === 2 && scores[0] != null && scores[1] != null
                              ? Number(scores[0]) - Number(scores[1])
                              : null;
                          return (
                            <tr key={row.dimension as string} className="border-b last:border-0">
                              <td className="py-2 pr-4 text-muted-foreground">{row.dimension as string}</td>
                              {scores.map((score, i) => (
                                <td key={i} className="text-center py-2 px-3 font-mono">
                                  {score != null ? Number(score).toFixed(1) : "—"}
                                </td>
                              ))}
                              {delta !== null && (
                                <td className="text-center py-2 px-3">
                                  <Badge
                                    variant="outline"
                                    className={
                                      delta > 0 ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-400" :
                                      delta < 0 ? "border-red-500/30 bg-red-500/10 text-red-400" : ""
                                    }
                                  >
                                    {delta > 0 ? "+" : ""}{delta.toFixed(1)}
                                  </Badge>
                                </td>
                              )}
                            </tr>
                          );
                        })}
                      </tbody>
                    </table>
                  </div>
                </CardContent>
              </Card>
            )}
          </div>
        </SheetContent>
      </Sheet>
    </>
  );
}

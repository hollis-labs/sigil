"use client";

import { useState, useEffect, useMemo } from "react";
import { Loader2 } from "lucide-react";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { fetchList } from "@/lib/api";
import type { Scorecard } from "@/types/scorecard";
import type { Dimension } from "@/types/dimension";
import type { DimensionScore } from "@/types/dimensionscore";

interface Lens {
  id: string;
  name: string;
  dimension_count: number;
}

function scoreColor(score: number | null): string {
  if (score === null) return "bg-muted text-muted-foreground";
  if (score >= 8) return "bg-emerald-500/20 text-emerald-400 font-medium";
  if (score >= 5) return "bg-amber-500/20 text-amber-400 font-medium";
  return "bg-red-500/20 text-red-400 font-medium";
}

export function ScorecardHeatmap() {
  const [lenses, setLenses] = useState<Lens[]>([]);
  const [selectedLens, setSelectedLens] = useState<string>("");
  const [scorecards, setScorecards] = useState<Scorecard[]>([]);
  const [dimensions, setDimensions] = useState<Dimension[]>([]);
  const [scores, setScores] = useState<DimensionScore[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Load lenses on mount
  useEffect(() => {
    fetchList<Lens>("lenses").then((data) => {
      setLenses(data);
      if (data.length > 0) setSelectedLens(data[0].id);
    }).catch(() => setLenses([]));
  }, []);

  // Load dimensions once
  useEffect(() => {
    fetchList<Dimension>("dimensions").then(setDimensions).catch(() => setDimensions([]));
  }, []);

  // Load scorecards + scores when lens changes
  useEffect(() => {
    if (!selectedLens) return;
    setIsLoading(true);
    Promise.all([
      fetchList<Scorecard>("scorecards", { "filter[lens_id]": selectedLens }),
      fetchList<DimensionScore>("scores"),
    ])
      .then(([sc, ds]) => {
        setScorecards(sc);
        setScores(ds);
      })
      .catch(() => {
        setScorecards([]);
        setScores([]);
      })
      .finally(() => setIsLoading(false));
  }, [selectedLens]);

  // Build score lookup: scorecard_id -> dimension_id -> score
  const scoreLookup = useMemo(() => {
    const map = new Map<string, Map<string, number>>();
    for (const s of scores) {
      if (!s.scorecard_id || !s.dimension_id) continue;
      if (!map.has(s.scorecard_id)) map.set(s.scorecard_id, new Map());
      map.get(s.scorecard_id)!.set(s.dimension_id, Number(s.score) || 0);
    }
    return map;
  }, [scores]);

  // Filter dimensions to those that have scores in the current scorecard set
  const activeDimensions = useMemo(() => {
    const scorecardIds = new Set(scorecards.map((sc) => sc.id));
    const usedDimIds = new Set<string>();
    for (const s of scores) {
      if (s.scorecard_id && scorecardIds.has(s.scorecard_id) && s.dimension_id) {
        usedDimIds.add(s.dimension_id);
      }
    }
    return dimensions.filter((d) => usedDimIds.has(d.id));
  }, [dimensions, scores, scorecards]);

  // Sort scorecards by overall descending
  const sortedScorecards = useMemo(
    () => [...scorecards].sort((a, b) => (Number(b.overall) || 0) - (Number(a.overall) || 0)),
    [scorecards],
  );

  if (!lenses.length) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-3">
        <span className="text-sm text-muted-foreground">Lens:</span>
        <Select value={selectedLens} onValueChange={(v) => v && setSelectedLens(v)}>
          <SelectTrigger className="w-56">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {lenses.map((l) => (
              <SelectItem key={l.id} value={l.id}>
                {l.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <span className="text-xs text-muted-foreground">
          {sortedScorecards.length} repos &middot; {activeDimensions.length} dimensions
        </span>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
        </div>
      ) : sortedScorecards.length === 0 ? (
        <p className="text-sm text-muted-foreground py-4">No scorecards for this lens.</p>
      ) : (
        <div className="overflow-x-auto rounded-md border border-border">
          <TooltipProvider delay={200}>
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border bg-muted/50">
                  <th className="sticky left-0 z-10 bg-muted/50 px-3 py-2 text-left font-medium text-foreground">
                    Repo
                  </th>
                  <th className="px-2 py-2 text-center font-medium text-foreground whitespace-nowrap">
                    Overall
                  </th>
                  {activeDimensions.map((dim) => (
                    <th key={dim.id} className="px-2 py-2 text-center font-medium text-foreground">
                      <Tooltip>
                        <TooltipTrigger render={<span className="cursor-help whitespace-nowrap text-xs" />}>
                          {dim.name.length > 12 ? dim.name.slice(0, 11) + "…" : dim.name}
                        </TooltipTrigger>
                        <TooltipContent side="top">
                          <p className="font-medium">{dim.name}</p>
                          {dim.description && (
                            <p className="text-xs text-muted-foreground max-w-xs">{dim.description}</p>
                          )}
                        </TooltipContent>
                      </Tooltip>
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {sortedScorecards.map((sc) => {
                  const dimScores = scoreLookup.get(sc.id);
                  const overall = Number(sc.overall) || 0;
                  return (
                    <tr key={sc.id} className="border-b border-border last:border-0 hover:bg-muted/30">
                      <td className="sticky left-0 z-10 bg-background px-3 py-1.5 font-medium text-foreground whitespace-nowrap">
                        {sc.repo_name}
                      </td>
                      <td className={`px-2 py-1.5 text-center tabular-nums ${scoreColor(overall)}`}>
                        {overall.toFixed(1)}
                      </td>
                      {activeDimensions.map((dim) => {
                        const val = dimScores?.get(dim.id) ?? null;
                        return (
                          <td
                            key={dim.id}
                            className={`px-2 py-1.5 text-center tabular-nums ${scoreColor(val)}`}
                          >
                            {val !== null ? val.toFixed(0) : "—"}
                          </td>
                        );
                      })}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </TooltipProvider>
        </div>
      )}
    </div>
  );
}

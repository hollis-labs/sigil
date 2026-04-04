// hooks/use-dimensionscore.ts — fetches from Stack Explorer API
"use client";

import { useState, useEffect, useCallback } from "react";
import type { DimensionScore } from "@/types/dimensionscore";
import { fetchList } from "@/lib/api";

export function useDimensionScore() {
  const [data, setData] = useState<DimensionScore[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<DimensionScore>("scores").then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

export function useDimensionScoreById(id: string) {
  const [data, setData] = useState<DimensionScore[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<DimensionScore>("scores", { "filter[scorecard_id]": id }).then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [id, version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

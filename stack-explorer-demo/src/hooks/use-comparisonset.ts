// hooks/use-comparisonset.ts — fetches from Stack Explorer API
"use client";

import { useState, useEffect, useCallback } from "react";
import type { ComparisonSet } from "@/types/comparisonset";
import { fetchList, fetchItem } from "@/lib/api";

export function useComparisonSet() {
  const [data, setData] = useState<ComparisonSet[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<ComparisonSet>("comparison-sets").then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

export function useComparisonSetById(id: string) {
  const [data, setData] = useState<ComparisonSet | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchItem<ComparisonSet>("comparison-sets", id).then(setData).catch(() => setData(undefined)).finally(() => setIsLoading(false));
  }, [id, version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

// hooks/use-pattern.ts — fetches from Stack Explorer API
"use client";

import { useState, useEffect, useCallback } from "react";
import type { Pattern } from "@/types/pattern";
import { fetchList, fetchItem } from "@/lib/api";

export function usePattern() {
  const [data, setData] = useState<Pattern[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<Pattern>("patterns").then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

export function usePatternById(id: string) {
  const [data, setData] = useState<Pattern | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchItem<Pattern>("patterns", id).then(setData).catch(() => setData(undefined)).finally(() => setIsLoading(false));
  }, [id, version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

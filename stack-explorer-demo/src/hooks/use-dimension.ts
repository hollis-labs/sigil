// hooks/use-dimension.ts — fetches from Stack Explorer API
"use client";

import { useState, useEffect, useCallback } from "react";
import type { Dimension } from "@/types/dimension";
import { fetchList, fetchItem } from "@/lib/api";

export function useDimension() {
  const [data, setData] = useState<Dimension[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<Dimension>("dimensions").then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

export function useDimensionById(id: string) {
  const [data, setData] = useState<Dimension | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchItem<Dimension>("dimensions", id).then(setData).catch(() => setData(undefined)).finally(() => setIsLoading(false));
  }, [id, version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

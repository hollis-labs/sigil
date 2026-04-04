// hooks/use-repo.ts — fetches from Stack Explorer API
"use client";

import { useState, useEffect, useCallback } from "react";
import type { Repo } from "@/types/repo";
import { fetchList, fetchItem } from "@/lib/api";

export function useRepo() {
  const [data, setData] = useState<Repo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<Repo>("repos").then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

export function useRepoById(id: string) {
  const [data, setData] = useState<Repo | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchItem<Repo>("repos", id).then(setData).catch(() => setData(undefined)).finally(() => setIsLoading(false));
  }, [id, version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

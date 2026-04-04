// hooks/use-snapshot.ts — fetches from Stack Explorer API
"use client";

import { useState, useEffect, useCallback } from "react";
import type { Snapshot } from "@/types/snapshot";
import { fetchList } from "@/lib/api";

export function useSnapshot() {
  const [data, setData] = useState<Snapshot[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<Snapshot>("snapshots").then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

export function useSnapshotById(id: string) {
  const [data, setData] = useState<Snapshot[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<Snapshot>("snapshots", { "filter[repo_id]": id }).then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [id, version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

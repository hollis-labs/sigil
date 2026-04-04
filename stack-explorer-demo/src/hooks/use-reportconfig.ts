// hooks/use-reportconfig.ts — fetches from Stack Explorer API
"use client";

import { useState, useEffect, useCallback } from "react";
import type { ReportConfig } from "@/types/reportconfig";
import { fetchList, fetchItem } from "@/lib/api";

export function useReportConfig() {
  const [data, setData] = useState<ReportConfig[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchList<ReportConfig>("reports").then(setData).catch(() => setData([])).finally(() => setIsLoading(false));
  }, [version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

export function useReportConfigById(id: string) {
  const [data, setData] = useState<ReportConfig | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const [version, setVersion] = useState(0);

  useEffect(() => {
    setIsLoading(true);
    fetchItem<ReportConfig>("reports", id).then(setData).catch(() => setData(undefined)).finally(() => setIsLoading(false));
  }, [id, version]);

  const refetch = useCallback(() => setVersion((v) => v + 1), []);

  return { data, isLoading, refetch };
}

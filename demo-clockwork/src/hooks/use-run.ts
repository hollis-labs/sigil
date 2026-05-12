// hooks/use-run.ts — mock data hook (demo-only)
// Phase 3 stub: returns an empty list so the page renders the empty-state.
// Replace with a real fetch when wiring to the Clockwork MCP / API.

import { useState, useEffect } from "react";

// Loose shape — the demo UI accesses these fields on each row.
export interface Run {
  id: string;
  task_id?: string;
  agent?: string;
  status?: string;
  duration?: number;
  tokens?: number;
  cost?: number;
  created_at?: string;
  priority?: string;
}

export function useRun() {
  const [data, setData] = useState<Run[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Mock: no data; demo can swap this for a real fetch later.
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useRunById(_id: string) {
  return { data: null as Run | null, isLoading: false };
}

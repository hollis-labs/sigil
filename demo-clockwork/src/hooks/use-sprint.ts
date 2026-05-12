// hooks/use-sprint.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Sprint {
  id: string;
  [key: string]: unknown;
}

export function useSprint() {
  const [data, setData] = useState<Sprint[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useSprintById(_id: string) {
  return { data: null as Sprint | null, isLoading: false };
}

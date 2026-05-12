// hooks/use-plan.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Plan {
  id: string;
  [key: string]: unknown;
}

export function usePlan() {
  const [data, setData] = useState<Plan[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function usePlanById(_id: string) {
  return { data: null as Plan | null, isLoading: false };
}

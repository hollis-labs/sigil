// hooks/use-model.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Model {
  id: string;
  [key: string]: unknown;
}

export function useModel() {
  const [data, setData] = useState<Model[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useModelById(_id: string) {
  return { data: null as Model | null, isLoading: false };
}

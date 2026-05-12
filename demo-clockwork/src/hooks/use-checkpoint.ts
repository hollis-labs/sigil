// hooks/use-checkpoint.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Checkpoint {
  id: string;
  [key: string]: unknown;
}

export function useCheckpoint() {
  const [data, setData] = useState<Checkpoint[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useCheckpointById(_id: string) {
  return { data: null as Checkpoint | null, isLoading: false };
}

// hooks/use-artifact.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Artifact {
  id: string;
  [key: string]: unknown;
}

export function useArtifact() {
  const [data, setData] = useState<Artifact[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useArtifactById(_id: string) {
  return { data: null as Artifact | null, isLoading: false };
}

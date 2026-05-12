// hooks/use-project.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Project {
  id: string;
  [key: string]: unknown;
}

export function useProject() {
  const [data, setData] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useProjectById(_id: string) {
  return { data: null as Project | null, isLoading: false };
}

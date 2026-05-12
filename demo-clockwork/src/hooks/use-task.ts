// hooks/use-task.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Task {
  id: string;
  title?: string;
  status?: string;
  priority?: string;
  agent?: string;
  collection?: string;
  created_at?: string;
  updated_at?: string;
  done?: boolean;
  [key: string]: unknown;
}

export function useTask() {
  const [data, setData] = useState<Task[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useTaskById(_id: string) {
  return { data: null as Task | null, isLoading: false };
}

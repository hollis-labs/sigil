// hooks/use-template.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Template {
  id: string;
  [key: string]: unknown;
}

export function useTemplate() {
  const [data, setData] = useState<Template[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useTemplateById(_id: string) {
  return { data: null as Template | null, isLoading: false };
}

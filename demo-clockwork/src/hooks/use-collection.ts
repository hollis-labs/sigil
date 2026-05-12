// hooks/use-collection.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Collection {
  id: string;
  [key: string]: unknown;
}

export function useCollection() {
  const [data, setData] = useState<Collection[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useCollectionById(_id: string) {
  return { data: null as Collection | null, isLoading: false };
}

// hooks/use-setting.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Setting {
  id: string;
  [key: string]: unknown;
}

export function useSetting() {
  const [data, setData] = useState<Setting[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useSettingById(_id: string) {
  return { data: null as Setting | null, isLoading: false };
}

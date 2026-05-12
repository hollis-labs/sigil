// hooks/use-comment.ts — mock data hook (demo-only)
import { useState, useEffect } from "react";

export interface Comment {
  id: string;
  [key: string]: unknown;
}

export function useComment() {
  const [data, setData] = useState<Comment[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(false);
  }, []);

  return { data, isLoading, refetch: () => {} };
}

export function useCommentById(_id: string) {
  return { data: null as Comment | null, isLoading: false };
}

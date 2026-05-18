import { useCallback, useState } from 'react'

export interface DetailState<T> {
  /** Id of the open item, or null when the dialog is closed. */
  id: string | null
  data: T | null
  error: string | null
  loading: boolean
  /** Open the dialog for `id` and fetch its detail. */
  open: (id: string) => void
  close: () => void
}

/**
 * Fetch-on-open detail state for the row-click → `DetailDialog` flow shared
 * by every list page. `fetcher` must be referentially stable (a method off
 * the const `apiClient` is).
 */
export function useDetail<T>(fetcher: (id: string) => Promise<T>): DetailState<T> {
  const [id, setId] = useState<string | null>(null)
  const [data, setData] = useState<T | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const open = useCallback(
    (next: string) => {
      setId(next)
      setData(null)
      setError(null)
      setLoading(true)
      fetcher(next)
        .then((d) => setData(d))
        .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)))
        .finally(() => setLoading(false))
    },
    [fetcher],
  )

  const close = useCallback(() => {
    setId(null)
    setData(null)
    setError(null)
  }, [])

  return { id, data, error, loading, open, close }
}

import { useMemo, useState } from 'react'

export interface SearchState<T> {
  query: string
  setQuery: (q: string) => void
  /** `items` narrowed to rows whose `text` contains the (trimmed) query. */
  filtered: T[]
}

/** Case-insensitive substring filtering for a list page's filter-bar query. */
export function useSearch<T>(items: T[], text: (item: T) => string): SearchState<T> {
  const [query, setQuery] = useState('')
  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return items
    return items.filter((item) => text(item).toLowerCase().includes(q))
  }, [items, query, text])
  return { query, setQuery, filtered }
}

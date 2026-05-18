import { useMemo, type ReactNode } from 'react'
import { RefreshCw } from 'lucide-react'
import {
  Button,
  DetailDialog,
  EmptyState,
  JsonViewer,
  OperationsTablePage,
  usePoll,
  type ColumnDef,
  type SortState,
} from '@hollis-labs/sysop-ui'
import { useDetail } from './use-detail'
import { useSearch } from './use-search'

interface EntityListPageProps<TItem, TDetail> {
  /** Page title shown in the header. */
  title: string
  /** Lower-case singular noun — "page", "component", … — for empty/meta copy. */
  noun: string
  /** Fetch the list. Typically `() => api.listX()`. */
  fetchList: () => Promise<TItem[]>
  columns: ColumnDef<TItem>[]
  getRowId: (item: TItem) => string
  initialSort?: SortState
  /** Concatenated text a row is matched against by the filter-bar query. */
  searchText: (item: TItem) => string
  searchPlaceholder: string
  /** Fetch one row's detail for the dialog body. */
  fetchDetail: (id: string) => Promise<TDetail>
  /** Dialog body for a loaded detail. Defaults to a JSON view. */
  renderDetail?: (data: TDetail, id: string) => ReactNode
}

/**
 * The list-page template every Sigil entity screen is built from: a polled
 * list rendered through the kit's `OperationsTablePage`, a filter-bar query,
 * and a fetch-on-open `DetailDialog`. A concrete page is just column defs +
 * the two API calls.
 */
export function EntityListPage<TItem, TDetail>({
  title,
  noun,
  fetchList,
  columns,
  getRowId,
  initialSort,
  searchText,
  searchPlaceholder,
  fetchDetail,
  renderDetail,
}: EntityListPageProps<TItem, TDetail>) {
  // Fetch-once on mount; refreshed by the header button (the project is files
  // on disk, so there is nothing to poll on an interval).
  const { data, error, isLoading, refetch } = usePoll<TItem[]>(fetchList, 0)
  const items = useMemo(() => data ?? [], [data])
  const { query, setQuery, filtered } = useSearch(items, searchText)
  const detail = useDetail(fetchDetail)

  const emptyState = error ? (
    <EmptyState
      variant="error"
      title={`Could not load ${noun}s`}
      description={error instanceof Error ? error.message : String(error)}
    />
  ) : query ? (
    <EmptyState
      variant="no-results"
      title="No matches"
      description={`No ${noun}s match "${query}".`}
    />
  ) : (
    <EmptyState
      variant="empty"
      title={`No ${noun}s yet`}
      description={`This Sigil project has no ${noun}s defined.`}
    />
  )

  return (
    <>
      <OperationsTablePage<TItem>
        title={title}
        headerActions={
          <Button variant="ghost" size="sm" onClick={() => void refetch()}>
            <RefreshCw className="h-3.5 w-3.5" />
            Refresh
          </Button>
        }
        summaryCards={[
          { label: noun + 's', value: items.length },
          ...(query ? [{ label: 'Shown', value: filtered.length }] : []),
        ]}
        searchQuery={query}
        onSearchChange={setQuery}
        searchPlaceholder={searchPlaceholder}
        items={filtered}
        columns={columns}
        getRowId={getRowId}
        initialSort={initialSort}
        loading={isLoading && items.length === 0}
        onRowOpen={(id) => detail.open(id)}
        rowAriaLabel={(item) => `Open ${noun} ${getRowId(item)}`}
        emptyState={emptyState}
      />

      <DetailDialog
        open={detail.id !== null}
        onClose={detail.close}
        title={detail.id ?? ''}
        meta={<span className="font-mono">{noun}</span>}
      >
        <div className="p-4">
          {detail.loading && <p className="text-[13px] text-text-subtle">Loading…</p>}
          {detail.error && (
            <p className="text-[13px] text-status-blocked">{detail.error}</p>
          )}
          {detail.data !== null &&
            (renderDetail ? (
              renderDetail(detail.data, detail.id ?? '')
            ) : (
              <JsonViewer value={detail.data} className="max-h-[60vh]" />
            ))}
        </div>
      </DetailDialog>
    </>
  )
}

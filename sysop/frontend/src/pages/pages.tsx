import { Pill, type ColumnDef } from '@hollis-labs/sysop-ui'
import { apiClient, type JsonDoc, type PageSummary } from '../api/client'
import { EntityListPage } from '../lib/entity-page'

const dash = <span className="text-text-subtle">—</span>

const columns: ColumnDef<PageSummary>[] = [
  {
    key: 'id',
    header: 'ID',
    width: 'fill',
    sortValue: (p) => p.id,
    cell: (p) => <span className="font-mono text-text">{p.id}</span>,
  },
  {
    key: 'title',
    header: 'Title',
    sortValue: (p) => p.title ?? '',
    cell: (p) => p.title || dash,
  },
  {
    key: 'module',
    header: 'Module',
    sortValue: (p) => p.module ?? '',
    cell: (p) =>
      p.module ? <span className="font-mono text-[12px] text-text-muted">{p.module}</span> : dash,
  },
  {
    key: 'status',
    header: 'Status',
    align: 'right',
    cell: (p) =>
      p.error ? (
        <Pill tone="danger">parse error</Pill>
      ) : p.overlay ? (
        <Pill tone="info">{p.overlay}</Pill>
      ) : (
        <Pill tone="success">ok</Pill>
      ),
  },
]

/** Pages screen — every page config in the Sigil project. */
export function PagesPage() {
  return (
    <EntityListPage<PageSummary, JsonDoc>
      title="Pages"
      noun="page"
      fetchList={apiClient.listPages}
      columns={columns}
      getRowId={(p) => p.id}
      initialSort={{ key: 'id', dir: 'asc' }}
      searchText={(p) => `${p.id} ${p.title ?? ''} ${p.module ?? ''}`}
      searchPlaceholder="Filter pages…"
      fetchDetail={apiClient.getPage}
    />
  )
}

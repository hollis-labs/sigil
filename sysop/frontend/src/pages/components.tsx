import { Pill, type ColumnDef } from '@hollis-labs/sysop-ui'
import { apiClient, type ComponentDetail, type ComponentSummary } from '../api/client'
import { EntityListPage } from '../lib/entity-page'

const columns: ColumnDef<ComponentSummary>[] = [
  {
    key: 'type',
    header: 'Type',
    sortValue: (c) => c.type,
    cell: (c) => <span className="font-mono text-text">{c.type}</span>,
  },
  {
    key: 'category',
    header: 'Category',
    sortValue: (c) => c.category ?? '',
    cell: (c) =>
      c.category ? <Pill>{c.category}</Pill> : <span className="text-text-subtle">—</span>,
  },
  {
    key: 'description',
    header: 'Description',
    width: 'fill',
    sortValue: (c) => c.description ?? '',
    cell: (c) => <span className="text-text-muted">{c.description || '—'}</span>,
  },
]

/** Components screen — the component registry (builtin + custom). */
export function ComponentsPage() {
  return (
    <EntityListPage<ComponentSummary, ComponentDetail>
      title="Components"
      noun="component"
      fetchList={apiClient.listComponents}
      columns={columns}
      getRowId={(c) => c.type}
      initialSort={{ key: 'type', dir: 'asc' }}
      searchText={(c) => `${c.type} ${c.category ?? ''} ${c.description ?? ''}`}
      searchPlaceholder="Filter components…"
      fetchDetail={apiClient.getComponent}
    />
  )
}

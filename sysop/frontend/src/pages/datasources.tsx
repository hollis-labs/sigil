import { Pill, type ColumnDef } from '@hollis-labs/sysop-ui'
import { apiClient, type DataSourceSummary, type JsonDoc } from '../api/client'
import { EntityListPage } from '../lib/entity-page'

const columns: ColumnDef<DataSourceSummary>[] = [
  {
    key: 'alias',
    header: 'Alias',
    sortValue: (d) => d.alias,
    cell: (d) => <span className="font-mono text-text">{d.alias}</span>,
  },
  {
    key: 'description',
    header: 'Description',
    width: 'fill',
    sortValue: (d) => d.description ?? '',
    cell: (d) => <span className="text-text-muted">{d.description || '—'}</span>,
  },
  {
    key: 'capabilities',
    header: 'Capabilities',
    cell: (d) =>
      d.capabilities && d.capabilities.length > 0 ? (
        <div className="flex flex-wrap gap-1">
          {d.capabilities.map((c) => (
            <Pill key={c}>{c}</Pill>
          ))}
        </div>
      ) : (
        <span className="text-text-subtle">—</span>
      ),
  },
  {
    key: 'fields',
    header: 'Fields',
    align: 'right',
    sortValue: (d) => d.fieldCount ?? 0,
    cell: (d) => <span className="font-mono tabular-nums">{d.fieldCount ?? 0}</span>,
  },
]

/** Datasources screen — the datasource manifests bound to the project. */
export function DataSourcesPage() {
  return (
    <EntityListPage<DataSourceSummary, JsonDoc>
      title="Datasources"
      noun="datasource"
      fetchList={apiClient.listDataSources}
      columns={columns}
      getRowId={(d) => d.alias}
      initialSort={{ key: 'alias', dir: 'asc' }}
      searchText={(d) => `${d.alias} ${d.description ?? ''} ${(d.capabilities ?? []).join(' ')}`}
      searchPlaceholder="Filter datasources…"
      fetchDetail={apiClient.getDataSource}
    />
  )
}

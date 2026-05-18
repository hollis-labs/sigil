import { Pill, type ColumnDef } from '@hollis-labs/sysop-ui'
import { apiClient, type JsonDoc, type ThemeSummary } from '../api/client'
import { EntityListPage } from '../lib/entity-page'

const columns: ColumnDef<ThemeSummary>[] = [
  {
    key: 'name',
    header: 'Name',
    sortValue: (t) => t.name,
    cell: (t) => <span className="font-mono text-text">{t.name}</span>,
  },
  {
    key: 'extends',
    header: 'Extends',
    sortValue: (t) => t.extends ?? '',
    cell: (t) =>
      t.extends ? <Pill tone="info">{t.extends}</Pill> : <span className="text-text-subtle">—</span>,
  },
  {
    key: 'description',
    header: 'Description',
    width: 'fill',
    sortValue: (t) => t.description ?? '',
    cell: (t) => <span className="text-text-muted">{t.description || '—'}</span>,
  },
]

/** Themes screen — the theme token sets available to the project. */
export function ThemesPage() {
  return (
    <EntityListPage<ThemeSummary, JsonDoc>
      title="Themes"
      noun="theme"
      fetchList={apiClient.listThemes}
      columns={columns}
      getRowId={(t) => t.name}
      initialSort={{ key: 'name', dir: 'asc' }}
      searchText={(t) => `${t.name} ${t.extends ?? ''} ${t.description ?? ''}`}
      searchPlaceholder="Filter themes…"
      fetchDetail={apiClient.getTheme}
    />
  )
}

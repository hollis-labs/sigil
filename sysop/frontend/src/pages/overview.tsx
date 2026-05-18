import {
  EmptyState,
  JsonViewer,
  MetaList,
  Metric,
  PageHeader,
  usePoll,
  type MetaItem,
} from '@hollis-labs/sysop-ui'
import { apiClient, type JsonDoc } from '../api/client'

interface OverviewData {
  project: JsonDoc
  pages: number
  components: number
  datasources: number
  themes: number
}

/** Load the project config and the four entity counts in one pass. */
async function loadOverview(): Promise<OverviewData> {
  const [project, pages, components, datasources, themes] = await Promise.all([
    apiClient.getProject(),
    apiClient.listPages(),
    apiClient.listComponents(),
    apiClient.listDataSources(),
    apiClient.listThemes(),
  ])
  return {
    project,
    pages: pages.length,
    components: components.length,
    datasources: datasources.length,
    themes: themes.length,
  }
}

function asString(v: unknown): string | undefined {
  return typeof v === 'string' ? v : undefined
}

/** Overview screen — project config and entity counts at a glance. */
export function OverviewPage() {
  const { data, error, isLoading } = usePoll<OverviewData>(loadOverview, 0)

  if (error) {
    return (
      <div className="flex h-full flex-col">
        <PageHeader title="Overview" />
        <EmptyState
          variant="error"
          title="Could not reach Sigil"
          description={error instanceof Error ? error.message : String(error)}
        />
      </div>
    )
  }

  const project = data?.project ?? {}
  const defaults = (project.defaults ?? {}) as JsonDoc
  const count = (n: number | undefined) => (isLoading ? '…' : (n ?? 0))

  const meta: MetaItem[] = [
    { label: 'Project', value: asString(project.name) ?? '—' },
    { label: 'Spec version', value: asString(project.version) ?? '—' },
    { label: 'Default theme', value: asString(defaults.theme) ?? '—' },
    { label: 'Default renderer', value: asString(defaults.renderer) ?? '—' },
  ]

  return (
    <div className="flex h-full flex-col">
      <PageHeader title="Overview" />
      <div className="min-h-0 flex-1 overflow-auto">
        <section className="grid grid-cols-2 gap-6 border-b border-border-strong px-5 py-5 sm:grid-cols-4">
          <Metric label="Pages" value={count(data?.pages)} />
          <Metric label="Components" value={count(data?.components)} />
          <Metric label="Datasources" value={count(data?.datasources)} />
          <Metric label="Themes" value={count(data?.themes)} />
        </section>

        <section className="border-b border-border-strong px-5 py-4">
          <p className="mb-3 text-[10px] font-semibold uppercase tracking-[.18em] text-text-subtle">
            Project
          </p>
          <MetaList items={meta} columns={4} />
        </section>

        <section className="px-5 py-4">
          <p className="mb-3 text-[10px] font-semibold uppercase tracking-[.18em] text-text-subtle">
            sigil.yaml
          </p>
          <JsonViewer value={project} className="max-h-[40vh]" />
        </section>
      </div>
    </div>
  )
}

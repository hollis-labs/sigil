import { createApiClient } from '@hollis-labs/sysop-ui'

// Same-origin: the Go binary serves both this SPA and Sigil's JSON API, so
// an empty baseUrl resolves every request against the current origin.
const http = createApiClient({ baseUrl: '' })

/** A free-form YAML document decoded to JSON — page / theme / project bodies. */
export type JsonDoc = Record<string, unknown>

export interface HealthInfo {
  status: string
  service?: string
}

export interface PageSummary {
  id: string
  title?: string
  overlay?: string
  module?: string
  /** Set when the page file failed to parse. */
  error?: string
}

export interface ComponentSummary {
  type: string
  category?: string
  description?: string
}

export interface ComponentDetail extends ComponentSummary {
  props?: unknown
  actions?: unknown
  slots?: unknown
  shortcuts?: unknown
}

export interface DataSourceSummary {
  alias: string
  description?: string
  capabilities?: string[]
  fieldCount?: number
}

export interface ThemeSummary {
  name: string
  description?: string
  extends?: string
}

const enc = encodeURIComponent

/**
 * Concrete API client for Sigil's JSON surface (internal/server.API). One
 * method per endpoint; list calls unwrap the single-key envelope the API
 * returns so callers get a plain array.
 */
export const apiClient = {
  getHealth: () => http.get<HealthInfo>('/api/health'),

  getProject: () => http.get<JsonDoc>('/api/project'),

  listPages: () =>
    http.get<{ pages: PageSummary[] }>('/api/pages').then((r) => r.pages),
  getPage: (id: string) => http.get<JsonDoc>(`/api/pages/${enc(id)}`),

  listComponents: () =>
    http
      .get<{ components: ComponentSummary[] }>('/api/components')
      .then((r) => r.components),
  getComponent: (type: string) =>
    http.get<ComponentDetail>(`/api/components/${enc(type)}`),

  listDataSources: () =>
    http
      .get<{ datasources: DataSourceSummary[] }>('/api/datasources')
      .then((r) => r.datasources),
  getDataSource: (alias: string) =>
    http.get<JsonDoc>(`/api/datasources/${enc(alias)}`),

  listThemes: () =>
    http.get<{ themes: ThemeSummary[] }>('/api/themes').then((r) => r.themes),
  getTheme: (name: string) => http.get<JsonDoc>(`/api/themes/${enc(name)}`),
}

export type AppApiClient = typeof apiClient

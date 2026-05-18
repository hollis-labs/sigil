import { createApiContext } from '@hollis-labs/sysop-ui'
import { apiClient } from './client'

// Typed { ApiProvider, useApi } bound to this app's concrete client.
// Pages import `apiClient` directly; `useApi()` is here for components that
// want the context-injected client (e.g. tests passing a stub `client`).
export const { ApiProvider, useApi } = createApiContext(apiClient)

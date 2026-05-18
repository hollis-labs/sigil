import { useState, type ReactNode } from 'react'
import { Boxes, Database, FileText, Hexagon, LayoutDashboard, Palette } from 'lucide-react'
import { NavRail, ThemeSwitcher, type NavRailItem } from '@hollis-labs/sysop-ui'
import { OverviewPage } from './pages/overview'
import { PagesPage } from './pages/pages'
import { ComponentsPage } from './pages/components'
import { DataSourcesPage } from './pages/datasources'
import { ThemesPage } from './pages/themes'

interface Screen {
  label: string
  icon: ReactNode
  render: () => ReactNode
}

// One entry per screen — the nav rail and the active-route switch are both
// derived from this map, so adding a page is a single entry.
const SCREENS = {
  overview: {
    label: 'Overview',
    icon: <LayoutDashboard className="h-4 w-4" />,
    render: () => <OverviewPage />,
  },
  pages: {
    label: 'Pages',
    icon: <FileText className="h-4 w-4" />,
    render: () => <PagesPage />,
  },
  components: {
    label: 'Components',
    icon: <Boxes className="h-4 w-4" />,
    render: () => <ComponentsPage />,
  },
  datasources: {
    label: 'Datasources',
    icon: <Database className="h-4 w-4" />,
    render: () => <DataSourcesPage />,
  },
  themes: {
    label: 'Themes',
    icon: <Palette className="h-4 w-4" />,
    render: () => <ThemesPage />,
  },
} satisfies Record<string, Screen>

type RouteKey = keyof typeof SCREENS

/**
 * App shell — the icon nav rail on the left and the active screen. Each
 * screen renders its own kit page header, so the shell is just the rail.
 */
export function App() {
  const [route, setRoute] = useState<RouteKey>('overview')

  const nav: NavRailItem[] = (Object.keys(SCREENS) as RouteKey[]).map((key) => ({
    key,
    label: SCREENS[key].label,
    icon: SCREENS[key].icon,
    active: route === key,
    onSelect: () => setRoute(key),
  }))

  return (
    <div className="flex h-screen bg-bg text-text">
      <NavRail
        items={nav}
        logo={<Hexagon className="h-4 w-4" />}
        logoLabel="Sigil Sysop"
        footerExtra={<ThemeSwitcher />}
      />
      <main className="flex min-w-0 flex-1 flex-col">{SCREENS[route].render()}</main>
    </div>
  )
}

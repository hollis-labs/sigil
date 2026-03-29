import Link from "next/link";

const pages = [
  { href: "/component-showcase", title: "Component Showcase", description: "Cards, switches, dropdowns, combobox, accordion, tooltips" },
  { href: "/sigil-pages", title: "Pages", description: "Browse and manage Sigil page configurations" },
  { href: "/sigil-components", title: "Components", description: "Browse the 49 built-in component types" },
  { href: "/sigil-themes", title: "Themes", description: "Manage and preview theme tokens" },
  { href: "/sigil-datasources", title: "Data Sources", description: "Configure data source manifests" },
  { href: "/sigil-page-editor", title: "Page Editor", description: "Visual page configuration editor" },
  { href: "/sigil-preview", title: "Preview", description: "Live preview of page configs" },
];

const forgePages = [
  { href: "/forge-shell", title: "App Shell", description: "Sidebar nav, search, user menu" },
  { href: "/forge-overview", title: "Overview", description: "Dashboard with stat cards, table, progress" },
  { href: "/forge-deployments", title: "Deployments", description: "Data table with search, filters, badges" },
  { href: "/forge-deployment-detail", title: "Deployment Detail", description: "Tabs, cards, dropdown actions, detail hooks" },
  { href: "/forge-logs", title: "Logs", description: "Data table with level + time filters" },
  { href: "/forge-assets", title: "Assets", description: "Data table with provider + type filters" },
  { href: "/forge-jobs", title: "Scheduled Jobs", description: "Data table with search, badges" },
  { href: "/forge-workspaces", title: "Workspaces", description: "Data-driven card grid + list table" },
  { href: "/forge-monitoring", title: "Monitoring", description: "Stat cards, progress bars, semantic badges" },
  { href: "/forge-users", title: "Users", description: "Data table with search, role filters" },
  { href: "/forge-integrations", title: "Integrations", description: "Data-driven card grid with status badges" },
  { href: "/forge-settings", title: "Settings", description: "Tabs with forms, inputs, switches, danger zone" },
  { href: "/forge-command-palette", title: "Command Palette", description: "Cmd-K style command dialog with groups" },
  { href: "/forge-deploy-create", title: "New Deployment", description: "Create form with radio-group, textarea, breadcrumb" },
  { href: "/forge-asset-detail", title: "Asset Detail", description: "Breadcrumb, collapsible sections, popover" },
  { href: "/forge-monitoring-detail", title: "Monitoring Detail", description: "Toggle-group, collapsible, alerts, progress" },
];

export default function Home() {
  return (
    <div className="min-h-screen bg-background text-foreground font-sans">
      <div className="max-w-4xl mx-auto p-8">
        <div className="mb-8">
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight">Sigil Demo</h1>
          <p className="text-muted-foreground mt-2">
            Generated UI pages rendered with real shadcn/ui components.
          </p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {pages.map((page) => (
            <Link
              key={page.href}
              href={page.href}
              className="group block rounded-lg border border-border p-6 transition-colors hover:bg-muted"
            >
              <h2 className="text-lg font-semibold group-hover:underline">{page.title}</h2>
              <p className="text-sm text-muted-foreground mt-1">{page.description}</p>
            </Link>
          ))}
        </div>

        <div className="mt-12 mb-8">
          <h2 className="scroll-m-20 text-2xl font-semibold tracking-tight">Forge — Golden Template</h2>
          <p className="text-muted-foreground mt-1">
            12-page infra dashboard (Vercel/Linear style) generated entirely from Sigil YAML.
          </p>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          {forgePages.map((page) => (
            <Link
              key={page.href}
              href={page.href}
              className="group block rounded-lg border border-border p-6 transition-colors hover:bg-muted"
            >
              <h2 className="text-base font-semibold group-hover:underline">{page.title}</h2>
              <p className="text-xs text-muted-foreground mt-1">{page.description}</p>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}

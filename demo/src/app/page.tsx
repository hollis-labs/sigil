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
      </div>
    </div>
  );
}

"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import {
  AlertTriangle,
  BarChart3,
  Compass,
  FileText,
  GitCompare,
  GitFork,
  Layers,
  LayoutDashboard,
  Menu,
  Moon,
  Plus,
  Search,
  Settings,
  Shapes,
  Sun,
} from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Toaster } from "@/components/ui/sonner";
import { LensProvider, useLensContext } from "@/lib/lens-context";

// ── Navigation ────────────────────────────────────────────

const navItems = [
  { label: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { label: "Repos", href: "/repos", icon: GitFork },
  { label: "Scorecards", href: "/scorecards", icon: BarChart3 },
  { label: "Gap Analysis", href: "/gap-analysis", icon: GitCompare },
  { label: "Reports", href: "/reports", icon: FileText },
  { label: "Patterns", href: "/patterns", icon: Shapes },
  { label: "Findings", href: "/findings", icon: AlertTriangle },
  { label: "Dimensions", href: "/dimensions", icon: Layers },
  { label: "Settings", href: "/settings", icon: Settings },
];

function NavItem({
  item,
  pathname,
  onClick,
}: {
  item: (typeof navItems)[number];
  pathname: string;
  onClick?: () => void;
}) {
  const Icon = item.icon;
  const isActive = pathname === item.href || pathname.startsWith(item.href + "/");
  return (
    <Link
      href={item.href}
      onClick={onClick}
      className={`flex items-center gap-2.5 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
        isActive
          ? "bg-primary/10 text-primary"
          : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
      }`}
    >
      <Icon className="h-4 w-4" />
      {item.label}
    </Link>
  );
}

function SidebarContent({
  pathname,
  onNavigate,
}: {
  pathname: string;
  onNavigate?: () => void;
}) {
  return (
    <div className="flex flex-col gap-2 p-4 flex-1">
      <Link href="/dashboard" className="flex items-center gap-2 py-1" onClick={onNavigate}>
        <Compass className="h-5 w-5 text-primary" />
        <span className="text-lg font-semibold tracking-tight">Stack Explorer</span>
      </Link>
      <Separator />
      <nav className="flex flex-col gap-0.5">
        {navItems.map((item) => (
          <NavItem key={item.href} item={item} pathname={pathname} onClick={onNavigate} />
        ))}
      </nav>
      <div className="flex-1" />
      <div className="rounded-md border border-border bg-muted/30 px-3 py-2">
        <div className="text-[10px] uppercase tracking-wider text-muted-foreground mb-1">Database</div>
        <div className="text-xs text-foreground">112 repos &middot; 47 scored</div>
      </div>
    </div>
  );
}

// ── Layout ────────────────────────────────────────────────

export default function ExplorerLayout({ children }: { children: React.ReactNode }) {
  return (
    <LensProvider>
      <ExplorerShell>{children}</ExplorerShell>
    </LensProvider>
  );
}

function ExplorerShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { theme, setTheme } = useTheme();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [cmdOpen, setCmdOpen] = useState(false);
  const { currentLensId, setLensId, lenses: lensOptions } = useLensContext();

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setCmdOpen((prev) => !prev);
      }
    }
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, []);

  return (
    <>
      <div className="flex h-screen bg-background text-foreground font-sans">
        {/* Desktop Sidebar */}
        <aside className="hidden md:flex w-[220px] shrink-0 border-r border-border overflow-y-auto bg-muted/20 flex-col">
          <SidebarContent pathname={pathname} />
        </aside>

        {/* Mobile Sidebar */}
        <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
          <SheetContent side="left" className="w-[220px] p-0">
            <SheetHeader className="sr-only">
              <SheetTitle>Navigation</SheetTitle>
            </SheetHeader>
            <SidebarContent pathname={pathname} onNavigate={() => setMobileOpen(false)} />
          </SheetContent>
        </Sheet>

        {/* Main */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Topbar */}
          <div className="flex items-center justify-between px-4 py-2 border-b border-border">
            <div className="flex items-center gap-3">
              <Button
                variant="ghost"
                size="icon"
                className="md:hidden"
                onClick={() => setMobileOpen(true)}
              >
                <Menu className="h-4 w-4" />
              </Button>
              <button
                type="button"
                className="relative flex items-center w-56 rounded-md border border-input bg-transparent px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent transition-colors"
                onClick={() => setCmdOpen(true)}
              >
                <Search className="mr-2 h-3.5 w-3.5" />
                <span className="flex-1 text-left text-xs">Search repos, patterns...</span>
                <kbd className="ml-auto pointer-events-none hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
                  <span className="text-xs">⌘</span>K
                </kbd>
              </button>
            </div>
            <div className="flex items-center gap-2">
              {/* Global Lens Selector */}
              <div className="flex items-center gap-2">
                <span className="text-[10px] uppercase tracking-wider text-muted-foreground hidden lg:block">Lens</span>
                <Select value={currentLensId} onValueChange={(v) => v && setLensId(v)}>
                  <SelectTrigger className="h-8 w-[150px] text-xs">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {lensOptions.map((lens) => (
                      <SelectItem key={lens.id} value={lens.id}>
                        {lens.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <Separator orientation="vertical" className="h-6" />
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8"
                onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
              >
                <Sun className="h-3.5 w-3.5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                <Moon className="absolute h-3.5 w-3.5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                <span className="sr-only">Toggle theme</span>
              </Button>
              <DropdownMenu>
                <DropdownMenuTrigger className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium h-8 w-8 hover:bg-accent hover:text-accent-foreground">
                  <Plus className="h-4 w-4" />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onSelect={() => router.push("/repos?action=create")}>
                    <GitFork className="mr-2 h-4 w-4" />
                    Add Repo
                  </DropdownMenuItem>
                  <DropdownMenuItem onSelect={() => router.push("/reports?action=create")}>
                    <FileText className="mr-2 h-4 w-4" />
                    Generate Report
                  </DropdownMenuItem>
                  <DropdownMenuItem onSelect={() => router.push("/gap-analysis?action=create")}>
                    <GitCompare className="mr-2 h-4 w-4" />
                    New Comparison
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>

          {/* Page content */}
          <main className="flex-1 overflow-y-auto">
            {children}
          </main>
        </div>

        {/* Global Command Palette */}
        <CommandDialog open={cmdOpen} onOpenChange={setCmdOpen}>
          <CommandInput placeholder="Search repos, patterns, findings..." />
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            <CommandGroup heading="Navigation">
              {navItems.map((item) => {
                const Icon = item.icon;
                return (
                  <CommandItem
                    key={item.href}
                    onSelect={() => {
                      router.push(item.href);
                      setCmdOpen(false);
                    }}
                  >
                    <Icon className="mr-2 h-4 w-4" />
                    {item.label}
                  </CommandItem>
                );
              })}
            </CommandGroup>
            <CommandGroup heading="Quick Actions">
              <CommandItem onSelect={() => { router.push("/repos?action=create"); setCmdOpen(false); }}>
                <Plus className="mr-2 h-4 w-4" />
                Add Repo
              </CommandItem>
              <CommandItem onSelect={() => { router.push("/reports?action=create"); setCmdOpen(false); }}>
                <FileText className="mr-2 h-4 w-4" />
                Generate Report
              </CommandItem>
              <CommandItem onSelect={() => { router.push("/gap-analysis?action=create"); setCmdOpen(false); }}>
                <GitCompare className="mr-2 h-4 w-4" />
                New Comparison
              </CommandItem>
            </CommandGroup>
          </CommandList>
        </CommandDialog>

        <Toaster />
      </div>
    </>
  );
}

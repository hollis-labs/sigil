"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import {
  Activity,
  Bell,
  BookOpen,
  Boxes,
  Clock,
  Hexagon,
  Key,
  LayoutDashboard,
  Menu,
  MessageSquare,
  Moon,
  Plus,
  Puzzle,
  Rocket,
  ScrollText,
  Server,
  Settings,
  Sun,
  UserPlus,
  Users,
  Search,
} from "lucide-react";
import { useTheme } from "next-themes";
import { toast } from "sonner";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { CommandDialog, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Separator } from "@/components/ui/separator";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

const mainNav = [
  { label: "Overview", href: "/overview", icon: LayoutDashboard },
  { label: "Deployments", href: "/deployments", icon: Rocket },
  { label: "Logs", href: "/logs", icon: ScrollText },
  { label: "Monitoring", href: "/monitoring", icon: Activity },
  { label: "Assets", href: "/assets", icon: Server },
  { label: "Jobs", href: "/jobs", icon: Clock },
];

const secondaryNav = [
  { label: "Workspaces", href: "/workspaces", icon: Boxes },
  { label: "Users", href: "/users", icon: Users },
  { label: "Integrations", href: "/integrations", icon: Puzzle },
  { label: "Settings", href: "/settings", icon: Settings },
];

function NavItem({
  item,
  pathname,
  onClick,
}: {
  item: { label: string; href: string; icon: React.ComponentType<{ className?: string }> };
  pathname: string;
  onClick?: () => void;
}) {
  const Icon = item.icon;
  const isActive = pathname === item.href || pathname.startsWith(item.href + "/");
  return (
    <Link
      href={item.href}
      onClick={onClick}
      className={`flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
        isActive
          ? "bg-accent text-accent-foreground"
          : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
      }`}
    >
      <Icon className="h-4 w-4" />
      {item.label}
    </Link>
  );
}

function SidebarContent({ pathname, onNavigate }: { pathname: string; onNavigate?: () => void }) {
  return (
    <div className="flex flex-col gap-2 p-4 flex-1">
      <Link href="/overview" className="flex items-center gap-2 py-1" onClick={onNavigate}>
        <Hexagon className="h-5 w-5" />
        <span className="text-lg font-semibold tracking-tight">Forge</span>
      </Link>
      <Separator />
      <nav className="flex flex-col gap-1">
        {mainNav.map((item) => (
          <NavItem key={item.href} item={item} pathname={pathname} onClick={onNavigate} />
        ))}
      </nav>
      <Separator />
      <nav className="flex flex-col gap-1">
        {secondaryNav.map((item) => (
          <NavItem key={item.href} item={item} pathname={pathname} onClick={onNavigate} />
        ))}
      </nav>
      <div className="flex-1" />
      <div className="flex items-center gap-2">
        <Avatar className="h-8 w-8">
          <AvatarImage src="" alt="User" />
          <AvatarFallback>CP</AvatarFallback>
        </Avatar>
        <span className="text-sm text-muted-foreground">chrispian</span>
      </div>
    </div>
  );
}

export default function ForgeLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { theme, setTheme } = useTheme();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [cmdOpen, setCmdOpen] = useState(false);

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
    <div className="flex h-screen bg-background text-foreground font-sans">
      {/* Desktop Sidebar */}
      <aside className="hidden md:flex w-[240px] shrink-0 border-r border-border overflow-y-auto bg-muted/30 flex-col">
        <SidebarContent pathname={pathname} />
      </aside>

      {/* Mobile Sidebar */}
      <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
        <SheetContent side="left" className="w-[240px] p-0">
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
          <div className="flex items-center gap-2">
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
              className="relative flex items-center w-64 rounded-md border border-input bg-transparent px-3 py-2 text-sm text-muted-foreground hover:bg-accent transition-colors"
              onClick={() => setCmdOpen(true)}
            >
              <Search className="mr-2 h-4 w-4" />
              <span className="flex-1 text-left">Search...</span>
              <kbd className="ml-auto pointer-events-none hidden h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100 sm:flex">
                <span className="text-xs">⌘</span>K
              </kbd>
            </button>
          </div>
          <div className="flex items-center gap-1">
            <Popover>
              <PopoverTrigger render={<Button variant="ghost" size="icon" className="relative" />}>
                <Bell className="h-4 w-4" />
                <Badge variant="destructive" className="absolute -top-1 -right-1 h-4 w-4 p-0 text-[10px] flex items-center justify-center">
                  3
                </Badge>
              </PopoverTrigger>
              <PopoverContent className="w-80" align="end">
                <div className="flex items-center justify-between mb-3">
                  <h4 className="text-sm font-semibold">Notifications</h4>
                  <Button variant="ghost" size="sm" className="text-xs h-auto py-1">Mark all read</Button>
                </div>
                <Separator className="mb-2" />
                <div className="space-y-2">
                  <div className="flex gap-3 rounded-md p-2 hover:bg-accent cursor-pointer">
                    <div className="h-2 w-2 rounded-full bg-blue-500 mt-1.5 shrink-0" />
                    <div>
                      <p className="text-sm font-medium">Deployment succeeded</p>
                      <p className="text-xs text-muted-foreground">api-v2.4.1 deployed to production &middot; 5m ago</p>
                    </div>
                  </div>
                  <div className="flex gap-3 rounded-md p-2 hover:bg-accent cursor-pointer">
                    <div className="h-2 w-2 rounded-full bg-red-500 mt-1.5 shrink-0" />
                    <div>
                      <p className="text-sm font-medium">Build failed</p>
                      <p className="text-xs text-muted-foreground">worker-update on preview &middot; 2h ago</p>
                    </div>
                  </div>
                  <div className="flex gap-3 rounded-md p-2 hover:bg-accent cursor-pointer">
                    <div className="h-2 w-2 rounded-full bg-yellow-500 mt-1.5 shrink-0" />
                    <div>
                      <p className="text-sm font-medium">P95 latency alert</p>
                      <p className="text-xs text-muted-foreground">API latency exceeded 200ms &middot; 4h ago</p>
                    </div>
                  </div>
                </div>
              </PopoverContent>
            </Popover>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            >
              <Sun className="h-4 w-4 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
              <Moon className="absolute h-4 w-4 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
              <span className="sr-only">Toggle theme</span>
            </Button>
            <Button variant="ghost" size="sm" className="hidden sm:flex">
              <MessageSquare className="mr-2 h-4 w-4" />
              Feedback
            </Button>
            <Button variant="ghost" size="sm" className="hidden sm:flex">
              <BookOpen className="mr-2 h-4 w-4" />
              Docs
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium h-9 px-3 hover:bg-accent hover:text-accent-foreground">
                +
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <Link href="/deploy-create">
                  <DropdownMenuItem>New Deployment</DropdownMenuItem>
                </Link>
                <DropdownMenuItem onSelect={() => toast("Asset creation coming soon")}>New Asset</DropdownMenuItem>
                <DropdownMenuItem onSelect={() => toast("Job creation coming soon")}>New Job</DropdownMenuItem>
                <DropdownMenuItem onSelect={() => toast("Workspace creation coming soon")}>New Workspace</DropdownMenuItem>
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
        <CommandInput placeholder="Type a command or search..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup heading="Navigation">
            <CommandItem onSelect={() => { router.push("/overview"); setCmdOpen(false); }}>
              <LayoutDashboard className="mr-2 h-4 w-4" />
              Overview
            </CommandItem>
            <CommandItem onSelect={() => { router.push("/deployments"); setCmdOpen(false); }}>
              <Rocket className="mr-2 h-4 w-4" />
              Deployments
            </CommandItem>
            <CommandItem onSelect={() => { router.push("/logs"); setCmdOpen(false); }}>
              <ScrollText className="mr-2 h-4 w-4" />
              Logs
            </CommandItem>
            <CommandItem onSelect={() => { router.push("/monitoring"); setCmdOpen(false); }}>
              <Activity className="mr-2 h-4 w-4" />
              Monitoring
            </CommandItem>
            <CommandItem onSelect={() => { router.push("/assets"); setCmdOpen(false); }}>
              <Server className="mr-2 h-4 w-4" />
              Assets
            </CommandItem>
            <CommandItem onSelect={() => { router.push("/settings"); setCmdOpen(false); }}>
              <Settings className="mr-2 h-4 w-4" />
              Settings
            </CommandItem>
          </CommandGroup>
          <CommandGroup heading="Actions">
            <CommandItem onSelect={() => { router.push("/deploy-create"); setCmdOpen(false); }}>
              <Plus className="mr-2 h-4 w-4" />
              New Deployment
            </CommandItem>
            <CommandItem onSelect={() => { router.push("/users"); setCmdOpen(false); }}>
              <UserPlus className="mr-2 h-4 w-4" />
              Invite User
            </CommandItem>
          </CommandGroup>
          <CommandGroup heading="Settings">
            <CommandItem onSelect={() => { router.push("/integrations"); setCmdOpen(false); }}>
              <Puzzle className="mr-2 h-4 w-4" />
              Integrations
            </CommandItem>
            <CommandItem onSelect={() => { toast("API Keys coming soon"); setCmdOpen(false); }}>
              <Key className="mr-2 h-4 w-4" />
              API Keys
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </div>
  );
}

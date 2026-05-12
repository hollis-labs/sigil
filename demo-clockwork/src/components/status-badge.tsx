// status-badge.tsx
// Sigil custom component — display a status string as a styled badge.
// Permissive: accepts any string status; common values (open, in_progress,
// in_review, blocked, done, running, completed, failed, queued, cancelled)
// get colored variants. Unknown statuses render as a neutral outline.

import { Badge } from "@/components/ui/badge";

export interface StatusBadgeProps {
  status?: string;
  size?: "sm" | "default" | "lg";
  showLabel?: boolean;
  dot?: boolean;
}

const variantClass: Record<string, string> = {
  open: "border-blue-500/30 bg-blue-500/10 text-blue-400",
  queued: "border-blue-500/30 bg-blue-500/10 text-blue-400",
  in_progress: "border-amber-500/30 bg-amber-500/10 text-amber-400",
  running: "border-amber-500/30 bg-amber-500/10 text-amber-400",
  in_review: "border-cyan-500/30 bg-cyan-500/10 text-cyan-400",
  blocked: "border-red-500/30 bg-red-500/10 text-red-400",
  failed: "border-red-500/30 bg-red-500/10 text-red-400",
  done: "border-emerald-500/30 bg-emerald-500/10 text-emerald-400",
  completed: "border-emerald-500/30 bg-emerald-500/10 text-emerald-400",
  cancelled: "border-zinc-500/30 bg-zinc-500/10 text-zinc-400",
};

const sizeClass: Record<string, string> = {
  sm: "text-[10px] px-1.5 py-0.5",
  default: "text-xs px-2 py-0.5",
  lg: "text-sm px-3 py-1",
};

export function StatusBadge({ status, size = "default", showLabel = true, dot = false }: StatusBadgeProps) {
  const s = (status ?? "").toLowerCase();
  const cls = variantClass[s] ?? "border-zinc-700 bg-zinc-900 text-zinc-300";
  if (dot && !showLabel) {
    return <span className={`inline-block h-2 w-2 rounded-full ${cls}`} aria-label={status} />;
  }
  return (
    <Badge variant="outline" className={`${cls} ${sizeClass[size] ?? sizeClass.default}`}>
      {dot && <span className="mr-1 inline-block h-1.5 w-1.5 rounded-full bg-current" />}
      {showLabel ? status || "—" : null}
    </Badge>
  );
}

export default StatusBadge;

// run-card.tsx
// Sigil custom component — compact card showing a single agent run.
// Accepts flat-prop shape emitted by the renderer (task_id, agent, status,
// duration, tokens, cost, timestamp, priority). All props are optional so
// the component works in a list where each item supplies only some fields.

import { Clock, Cpu, DollarSign } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

export interface RunCardProps {
  // Object form (preserved for hand-written callers).
  run?: Record<string, unknown>;
  // Flat form (emitted by the Sigil renderer).
  task_id?: string;
  agent?: string;
  status?: string;
  duration?: number | string;
  tokens?: number;
  cost?: number;
  timestamp?: string;
  priority?: string;
  // Display toggles.
  compact?: boolean;
  showMetrics?: boolean;
  showTiming?: boolean;
}

function formatDuration(ms: number | string | undefined): string {
  if (ms === undefined || ms === null || ms === "") return "";
  const n = typeof ms === "string" ? Number(ms) : ms;
  if (!Number.isFinite(n)) return String(ms);
  if (n < 1000) return `${n}ms`;
  if (n < 60000) return `${(n / 1000).toFixed(1)}s`;
  const minutes = Math.floor(n / 60000);
  const seconds = Math.floor((n % 60000) / 1000);
  return `${minutes}m ${seconds}s`;
}

export function RunCard(props: RunCardProps) {
  const data = (props.run ?? props) as Record<string, unknown>;
  const agent = (data.agent as string | undefined) ?? "";
  const status = (data.status as string | undefined) ?? "";
  const duration = data.duration as number | string | undefined;
  const tokens = data.tokens as number | undefined;
  const cost = data.cost as number | undefined;
  const timestamp = data.timestamp as string | undefined;
  const showMetrics = props.showMetrics ?? true;
  const showTiming = props.showTiming ?? true;

  if (props.compact) {
    return (
      <div className="flex items-center gap-3 rounded-lg border border-zinc-800 bg-zinc-950 px-3 py-2">
        {status && <Badge variant="outline" className="text-xs">{status}</Badge>}
        <div className="min-w-0 flex-1">
          <p className="truncate text-sm font-medium">{agent || "Run"}</p>
        </div>
        {showTiming && duration !== undefined && (
          <span className="flex items-center gap-1 text-xs text-muted-foreground">
            <Clock className="h-3 w-3" />
            {formatDuration(duration)}
          </span>
        )}
      </div>
    );
  }

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <CardTitle className="truncate text-sm font-medium">{agent || "Run"}</CardTitle>
          {status && <Badge variant="outline" className="text-xs">{status}</Badge>}
        </div>
      </CardHeader>
      <CardContent className="pb-3">
        {showTiming && (duration !== undefined || timestamp) && (
          <div className="mb-2 flex items-center gap-4 text-xs text-muted-foreground">
            {duration !== undefined && (
              <span className="flex items-center gap-1">
                <Clock className="h-3.5 w-3.5" />
                {formatDuration(duration)}
              </span>
            )}
            {timestamp && <span>{timestamp}</span>}
          </div>
        )}

        {showMetrics && (tokens !== undefined || cost !== undefined) && (
          <>
            <Separator className="mb-2" />
            <div className="flex items-center gap-4">
              {tokens !== undefined && (
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <Cpu className="h-3.5 w-3.5" />
                  {tokens.toLocaleString()} tokens
                </span>
              )}
              {cost !== undefined && (
                <span className="flex items-center gap-1 text-xs text-muted-foreground">
                  <DollarSign className="h-3.5 w-3.5" />
                  ${cost.toFixed(4)}
                </span>
              )}
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}

export default RunCard;

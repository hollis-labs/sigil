// activity-panel.tsx
// Sigil custom component — minimal activity feed placeholder.
// Renderer emits flat props (datasource, taskField) which we accept and
// ignore in this demo-grade implementation. Real impl can replace the
// placeholder body when mock-data wiring exists.

import { ScrollArea } from "@/components/ui/scroll-area";

export interface ActivityPanelProps {
  // Flat props the renderer emits.
  datasource?: string;
  taskField?: string;
  // Object form (legacy / hand-written callers). Permissive `any[]` so
  // typed entity arrays (Run[], etc.) pass without manual conversion.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  items?: any[];
  maxItems?: number;
  showTokens?: boolean;
  showCost?: boolean;
  autoScroll?: boolean;
}

export function ActivityPanel({ items = [], maxItems = 50 }: ActivityPanelProps) {
  const shown = items.slice(0, maxItems);
  return (
    <ScrollArea className="h-64 rounded-md border">
      <div className="divide-y">
        {shown.length === 0 ? (
          <p className="p-3 text-sm text-muted-foreground">No activity yet.</p>
        ) : (
          shown.map((it, i) => (
            <div key={i} className="flex items-start gap-3 p-2 text-xs">
              <span className="text-muted-foreground">{String(it.timestamp ?? "")}</span>
              <span className="font-medium">{String(it.type ?? "")}</span>
              <span className="flex-1">{String(it.summary ?? "")}</span>
            </div>
          ))
        )}
      </div>
    </ScrollArea>
  );
}

export default ActivityPanel;

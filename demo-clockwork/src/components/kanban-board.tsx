// kanban-board.tsx
// Sigil custom component — minimal Kanban board placeholder.
// Renderer emits flat props (datasource, groupField, columns). Demo-grade.

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export interface KanbanBoardProps {
  datasource?: string;
  groupField?: string;
  // Permissive: accept arbitrary extra fields on column descriptors so
  // YAMLs that use `status`, `title`, etc. instead of `value`/`label`
  // still typecheck. The component reads label/value/id/title/status
  // defensively below.
  columns?: Array<Record<string, unknown>>;
  items?: Array<Record<string, unknown>>;
  // Permissive flat-prop pass-through for YAMLs that pass extra fields
  // (cards, inbox, draggable, …) the renderer emits unchanged.
  cards?: unknown;
  inbox?: unknown;
  draggable?: boolean;
}

export function KanbanBoard({ columns = [], items = [], groupField = "status" }: KanbanBoardProps) {
  const cols = Array.isArray(columns) ? columns : [];
  return (
    <div className="grid grid-cols-1 gap-3 md:grid-cols-3 lg:grid-cols-4">
      {cols.map((col, i) => {
        const colValue = String(
          col.value ?? col.status ?? col.id ?? col.label ?? col.title ?? "",
        );
        const colLabel = String(
          col.title ?? col.label ?? col.value ?? col.id ?? "",
        );
        const colItems = items.filter(
          (it) => String(it[groupField] ?? "") === colValue,
        );
        return (
          <Card key={i}>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm font-medium">
                {colLabel} ({colItems.length})
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {colItems.map((it, j) => (
                <div key={j} className="rounded-md border p-2 text-xs">
                  {String(it.title ?? it.name ?? it.id ?? "")}
                </div>
              ))}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

export default KanbanBoard;

// checkpoint-card.tsx
// Sigil custom component — compact card showing a single checkpoint.
// Accepts both the renderer's flat props (datasource, showDetail, variant)
// and direct checkpoint object form. All props optional so list-mapped
// instances work when only some fields are supplied.

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export interface CheckpointCardProps {
  // Object form.
  checkpoint?: Record<string, unknown>;
  // Flat props the renderer may emit.
  datasource?: string;
  showDetail?: boolean;
  variant?: string;
  // Common checkpoint fields (flat-prop fallback).
  id?: string;
  kind?: string;
  status?: string;
  task_id?: string;
  created_at?: string;
  expires_at?: string;
  requester?: string;
  context?: string;
  payload?: string;
  notes?: string;
}

export function CheckpointCard(props: CheckpointCardProps) {
  const data = (props.checkpoint ?? props) as Record<string, unknown>;
  const id = (data.id as string | undefined) ?? "";
  const kind = (data.kind as string | undefined) ?? "";
  const status = (data.status as string | undefined) ?? "";
  const requester = (data.requester as string | undefined) ?? "";
  const createdAt = (data.created_at as string | undefined) ?? "";

  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between">
          <CardTitle className="truncate text-sm font-medium">{kind || "Checkpoint"}</CardTitle>
          {status && <Badge variant="outline" className="text-xs">{status}</Badge>}
        </div>
        {id && <p className="text-xs text-muted-foreground">#{id.slice(0, 8)}</p>}
      </CardHeader>
      <CardContent className="space-y-1 text-xs text-muted-foreground">
        {requester && <p>By: {requester}</p>}
        {createdAt && <p>Created: {createdAt}</p>}
        {props.showDetail !== false && Boolean(data.context) && (
          <p className="line-clamp-2">{String(data.context)}</p>
        )}
      </CardContent>
    </Card>
  );
}

export default CheckpointCard;

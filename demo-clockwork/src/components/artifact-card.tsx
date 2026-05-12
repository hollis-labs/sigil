// artifact-card.tsx
// Sigil custom component — compact card listing artifacts for a task.
// Renderer emits flat props (datasource, taskField). Demo-grade placeholder.

import { FileText } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export interface ArtifactCardProps {
  datasource?: string;
  taskField?: string;
  artifact?: Record<string, unknown>;
  artifacts?: Array<Record<string, unknown>>;
}

export function ArtifactCard({ artifacts = [] }: ArtifactCardProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium">Artifacts</CardTitle>
      </CardHeader>
      <CardContent>
        {artifacts.length === 0 ? (
          <p className="text-sm text-muted-foreground">No artifacts.</p>
        ) : (
          <ul className="space-y-1 text-sm">
            {artifacts.map((a, i) => (
              <li key={i} className="flex items-center gap-2">
                <FileText className="h-3.5 w-3.5 text-muted-foreground" />
                <span>{String(a.name ?? a.path ?? "")}</span>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}

export default ArtifactCard;

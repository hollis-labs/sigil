// scope-manager.tsx
// Sigil custom component — minimal scope-manager dialog placeholder.
// Demo-grade: accepts the renderer's flat props (datasource, mode) and
// renders an inline button + optional dialog. Real implementation can
// follow once the underlying data is wired.

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

export interface ScopeManagerProps {
  datasource?: string;
  mode?: string;
  multiSelect?: boolean;
  open?: boolean;
  onAssign?: (assignment: Record<string, unknown>) => void;
  onClose?: () => void;
}

export function ScopeManager({ mode = "assign" }: ScopeManagerProps) {
  const [open, setOpen] = useState(false);
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger render={<Button variant="outline" size="sm">Manage scope</Button>} />
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Scope ({mode})</DialogTitle>
        </DialogHeader>
        <p className="text-sm text-muted-foreground">
          Scope manager placeholder — wire mock data here.
        </p>
      </DialogContent>
    </Dialog>
  );
}

export default ScopeManager;

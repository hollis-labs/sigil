"use client";

import { useState } from "react";
import { toast } from "sonner";
import { postItem } from "@/lib/api";
import { BarChart3 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface Dimension {
  id: string;
  name: string;
  category?: string;
}

interface Lens {
  id: string;
  name: string;
}

interface RepoScoringDialogProps {
  repoId: string;
  repoName?: string;
  dimensions: Dimension[];
  lenses: Lens[];
}

export function RepoScoringDialog({ repoId, repoName, dimensions, lenses }: RepoScoringDialogProps) {
  const [open, setOpen] = useState(false);
  const [lensId, setLensId] = useState("");
  const [scores, setScores] = useState<Record<string, string>>({});

  const openDialog = () => {
    const initial: Record<string, string> = {};
    for (const dim of dimensions) initial[dim.id] = "";
    setScores(initial);
    setLensId(lenses[0]?.id ?? "");
    setOpen(true);
  };

  const handleScore = async () => {
    try {
      const sc = await postItem<{ id: string }>("scorecards", { repo_id: repoId, lens_id: lensId });
      const scoreEntries = Object.entries(scores).filter(([, v]) => v !== "");
      for (const [dimId, val] of scoreEntries) {
        await postItem("scores", { scorecard_id: sc.id, dimension_id: dimId, score: Number(val) });
      }
      toast.success(`Scorecard created with ${scoreEntries.length} scores`);
      setOpen(false);
      window.location.reload();
    } catch (e) {
      toast.error(e instanceof Error ? e.message : "Failed to create scorecard");
    }
  };

  return (
    <>
      <Button variant="default" onClick={openDialog}>
        <BarChart3 className="mr-2 h-4 w-4" />
        Score
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-lg max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Score {repoName ?? "Repo"}</DialogTitle>
          </DialogHeader>
          <div className="flex flex-col gap-4 py-4">
            <div className="flex flex-col gap-1.5">
              <Label>Lens</Label>
              <Select value={lensId} onValueChange={(v) => v && setLensId(v)}>
                <SelectTrigger>
                  <SelectValue placeholder="Select lens..." />
                </SelectTrigger>
                <SelectContent>
                  {lenses.map((l) => (
                    <SelectItem key={l.id} value={l.id}>
                      {l.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="text-xs text-muted-foreground uppercase tracking-wider">
              Dimension Scores (0-10)
            </div>
            {dimensions.map((dim) => (
              <div key={dim.id} className="flex flex-row items-center gap-3">
                <Label className="w-40 text-sm shrink-0">{dim.name}</Label>
                <Input
                  type="number"
                  min={0}
                  max={10}
                  step={0.1}
                  placeholder="—"
                  className="w-20"
                  value={scores[dim.id] ?? ""}
                  onChange={(e) => setScores((prev) => ({ ...prev, [dim.id]: e.target.value }))}
                />
                <span className="text-xs text-muted-foreground">{dim.category}</span>
              </div>
            ))}
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleScore}>Create Scorecard</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

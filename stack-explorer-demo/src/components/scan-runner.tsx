"use client";

import { useState, useRef, useCallback, useEffect } from "react";
import { toast } from "sonner";
import { postItem, fetchItem } from "@/lib/api";
import { CheckCircle2, Clock, Loader2, Play, XCircle } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/data-table";
import { useScan } from "@/hooks/use-scan";
import type { Scan } from "@/types/scan";

interface ScanRunnerProps {
  repos: { id: string }[];
}

export function ScanRunner({ repos }: ScanRunnerProps) {
  const [scanning, setScanning] = useState(false);
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const { data: scans, isLoading: scansLoading, refetch: refetchScans } = useScan();

  const stopPolling = useCallback(() => {
    if (pollRef.current) {
      clearInterval(pollRef.current);
      pollRef.current = null;
    }
  }, []);

  useEffect(() => stopPolling, [stopPolling]);

  const runScan = async () => {
    if (!repos.length) {
      toast.error("No repos to scan");
      return;
    }
    setScanning(true);
    try {
      const results = await Promise.allSettled(
        repos.map((r) => postItem<Scan>("scans", { repo_id: r.id, blueprint: "se-repo-scan" })),
      );
      const succeeded = results.filter((r) => r.status === "fulfilled").length;
      const failed = results.length - succeeded;
      toast.success(
        `Queued ${succeeded} scan${succeeded !== 1 ? "s" : ""}${failed ? ` (${failed} failed)` : ""}`,
      );
      refetchScans();
      const lastOk = [...results].reverse().find((r) => r.status === "fulfilled");
      if (lastOk && lastOk.status === "fulfilled") {
        const scanId = lastOk.value.id;
        pollRef.current = setInterval(async () => {
          try {
            const updated = await fetchItem<Scan>("scans", scanId);
            if (updated.status === "completed" || updated.status === "failed") {
              stopPolling();
              setScanning(false);
              if (updated.status === "completed") toast.success("Scans completed");
              else toast.error(updated.error_message || "Scan failed");
              refetchScans();
            }
          } catch {
            stopPolling();
            setScanning(false);
          }
        }, 3000);
      } else {
        setScanning(false);
      }
    } catch (e) {
      toast.error(e instanceof Error ? e.message : "Failed to start scans");
      setScanning(false);
    }
  };

  return (
    <>
      <Button variant="default" disabled={scanning} onClick={runScan}>
        {scanning ? (
          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
        ) : (
          <Play className="mr-2 h-4 w-4" />
        )}
        {scanning ? "Scanning…" : "Run Scan"}
      </Button>
      <Card>
        <CardHeader>
          <CardTitle>Scan History</CardTitle>
          <CardDescription>Recent scan jobs</CardDescription>
        </CardHeader>
        <CardContent>
          {scansLoading ? (
            <div className="flex items-center justify-center py-6">
              <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
            </div>
          ) : (
            <DataTable
              data={(scans ?? []).slice(0, 10)}
              columns={[
                {
                  accessorKey: "status",
                  header: "Status",
                  cell: ({ row }: { row: { getValue: (k: string) => unknown } }) => {
                    const status = String(row.getValue("status") ?? "");
                    const icons: Record<string, React.ReactNode> = {
                      pending: <Clock className="h-4 w-4 text-muted-foreground" />,
                      running: <Loader2 className="h-4 w-4 animate-spin text-blue-400" />,
                      completed: <CheckCircle2 className="h-4 w-4 text-emerald-400" />,
                      failed: <XCircle className="h-4 w-4 text-red-400" />,
                    };
                    return (
                      <div className="flex items-center gap-2">
                        {icons[status] ?? null}
                        <Badge
                          variant={
                            status === "failed"
                              ? "destructive"
                              : status === "completed"
                                ? "outline"
                                : "secondary"
                          }
                          className={
                            status === "completed"
                              ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-400"
                              : status === "running"
                                ? "border-blue-500/30 bg-blue-500/10 text-blue-400"
                                : ""
                          }
                        >
                          {status}
                        </Badge>
                      </div>
                    );
                  },
                },
                { accessorKey: "repo_name", header: "Repo" },
                { accessorKey: "blueprint", header: "Blueprint" },
                {
                  accessorKey: "created_at",
                  header: "Started",
                  cell: ({ row }: { row: { getValue: (k: string) => unknown } }) => {
                    const v = row.getValue("created_at");
                    if (!v) return "—";
                    const d = new Date(String(v));
                    return d.toLocaleString(undefined, {
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    });
                  },
                },
                {
                  accessorKey: "finished_at",
                  header: "Finished",
                  cell: ({ row }: { row: { getValue: (k: string) => unknown } }) => {
                    const v = row.getValue("finished_at");
                    if (!v) return "—";
                    const d = new Date(String(v));
                    return d.toLocaleString(undefined, {
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    });
                  },
                },
              ]}
              emptyMessage="No scans yet. Click Run Scan to start one."
            />
          )}
        </CardContent>
      </Card>
    </>
  );
}

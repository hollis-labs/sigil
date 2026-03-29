// hooks/use-scheduledjob.ts — mock data for demo
import type { ScheduledJob } from "@/types/scheduledjob";

const mockScheduledJobs: ScheduledJob[] = [
  {
    id: "job-001",
    name: "db-backup-prod",
    schedule: "0 */6 * * *",
    status: "active",
    lastRun: "2026-03-29T06:00:00Z",
    nextRun: "2026-03-29T12:00:00Z",
    runtime: "shell",
    successRate: "98.5%",
  },
  {
    id: "job-002",
    name: "cache-purge",
    schedule: "*/30 * * * *",
    status: "active",
    lastRun: "2026-03-29T10:00:00Z",
    nextRun: "2026-03-29T10:30:00Z",
    runtime: "node",
    successRate: "100%",
  },
  {
    id: "job-003",
    name: "ssl-cert-renew",
    schedule: "0 0 1 * *",
    status: "failed",
    lastRun: "2026-03-01T00:00:00Z",
    nextRun: "2026-04-01T00:00:00Z",
    runtime: "shell",
    successRate: "95.0%",
  },
  {
    id: "job-004",
    name: "log-rotate",
    schedule: "0 2 * * *",
    status: "active",
    lastRun: "2026-03-29T02:00:00Z",
    nextRun: "2026-03-30T02:00:00Z",
    runtime: "python",
    successRate: "99.8%",
  },
  {
    id: "job-005",
    name: "health-check",
    schedule: "*/5 * * * *",
    status: "active",
    lastRun: "2026-03-29T10:30:00Z",
    nextRun: "2026-03-29T10:35:00Z",
    runtime: "go",
    successRate: "100%",
  },
];

export function useScheduledJob() {
  return { data: mockScheduledJobs, isLoading: false };
}

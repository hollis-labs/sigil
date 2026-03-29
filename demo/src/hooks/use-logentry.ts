// hooks/use-logentry.ts — mock data for demo
import type { LogEntry } from "@/types/logentry";

const mockLogEntries: LogEntry[] = [
  {
    id: "log-001",
    timestamp: "2026-03-29T10:32:15Z",
    level: "error",
    service: "db-replica",
    message: "Connection refused to database replica db-ro.forge.dev:5432",
    source: "worker-02",
  },
  {
    id: "log-002",
    timestamp: "2026-03-29T10:31:44Z",
    level: "warn",
    service: "worker-pool",
    message: "Worker pool response time exceeded 300ms threshold",
    source: "api-gateway",
  },
  {
    id: "log-003",
    timestamp: "2026-03-29T10:30:02Z",
    level: "info",
    service: "deployer",
    message: "Deployment api-v2.4.1 completed successfully",
    source: "ci-runner",
  },
  {
    id: "log-004",
    timestamp: "2026-03-29T10:28:30Z",
    level: "info",
    service: "cert-manager",
    message: "SSL certificate renewed for *.forge.dev",
    source: "cert-manager",
  },
  {
    id: "log-005",
    timestamp: "2026-03-29T10:25:18Z",
    level: "error",
    service: "worker-03",
    message: "Out of memory on worker-03, process killed",
    source: "systemd",
  },
  {
    id: "log-006",
    timestamp: "2026-03-29T10:22:05Z",
    level: "debug",
    service: "cache",
    message: "Cache invalidation triggered for /api/users",
    source: "redis-cache",
  },
  {
    id: "log-007",
    timestamp: "2026-03-29T10:20:51Z",
    level: "warn",
    service: "api-gateway",
    message: "Rate limit approaching for api-gateway (85%)",
    source: "rate-limiter",
  },
  {
    id: "log-008",
    timestamp: "2026-03-29T10:18:00Z",
    level: "info",
    service: "backup",
    message: "Backup completed: db-primary 2.3GB",
    source: "cron-runner",
  },
];

export function useLogEntry() {
  return { data: mockLogEntries, isLoading: false };
}

// hooks/use-integration.ts — mock data for demo
import type { Integration } from "@/types/integration";

const mockIntegrations: Integration[] = [
  {
    id: "int-001",
    name: "GitHub",
    provider: "github",
    status: "connected",
    connectedAt: "2025-06-15T08:30:00Z",
  },
  {
    id: "int-002",
    name: "Slack",
    provider: "slack",
    status: "connected",
    connectedAt: "2025-07-01T14:00:00Z",
  },
  {
    id: "int-003",
    name: "Datadog",
    provider: "datadog",
    status: "error",
    connectedAt: "2025-08-20T09:15:00Z",
  },
  {
    id: "int-004",
    name: "PagerDuty",
    provider: "pagerduty",
    status: "connected",
    connectedAt: "2025-09-05T11:00:00Z",
  },
  {
    id: "int-005",
    name: "AWS",
    provider: "aws",
    status: "connected",
    connectedAt: "2025-06-15T08:45:00Z",
  },
  {
    id: "int-006",
    name: "Cloudflare",
    provider: "cloudflare",
    status: "disconnected",
    connectedAt: "2025-10-12T16:30:00Z",
  },
];

export function useIntegration() {
  return { data: mockIntegrations, isLoading: false };
}

// hooks/use-deploymentevent.ts — mock data for demo
import type { DeploymentEvent } from "@/types/deploymentevent";

const mockEvents: DeploymentEvent[] = [
  { id: "evt-1", title: "Deployment started", timestamp: "2026-03-29T10:15:00Z", status: "info" },
  { id: "evt-2", title: "Dependencies installed", timestamp: "2026-03-29T10:15:32Z", status: "success" },
  { id: "evt-3", title: "Build completed", timestamp: "2026-03-29T10:16:48Z", status: "success" },
  { id: "evt-4", title: "Health check passed", timestamp: "2026-03-29T10:17:12Z", status: "success" },
  { id: "evt-5", title: "Traffic shifted to new deployment", timestamp: "2026-03-29T10:17:34Z", status: "success" },
];

export function useDeploymentEvent() {
  return { data: mockEvents, isLoading: false };
}

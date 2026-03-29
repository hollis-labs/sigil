// hooks/use-workspace.ts — mock data for demo
import type { Workspace } from "@/types/workspace";

const mockWorkspaces: Workspace[] = [
  {
    id: "ws-001",
    name: "Forge Production",
    slug: "forge-prod",
    owner: "chrispian",
    memberCount: 8,
    deploymentCount: 142,
    plan: "pro",
    createdAt: "2025-06-15T08:00:00Z",
  },
  {
    id: "ws-002",
    name: "Staging Lab",
    slug: "staging-lab",
    owner: "alice",
    memberCount: 5,
    deploymentCount: 67,
    plan: "pro",
    createdAt: "2025-09-01T12:00:00Z",
  },
  {
    id: "ws-003",
    name: "Personal",
    slug: "personal",
    owner: "chrispian",
    memberCount: 1,
    deploymentCount: 12,
    plan: "free",
    createdAt: "2025-04-20T10:00:00Z",
  },
  {
    id: "ws-004",
    name: "Enterprise Client",
    slug: "ent-client",
    owner: "bob",
    memberCount: 15,
    deploymentCount: 283,
    plan: "enterprise",
    createdAt: "2025-01-10T09:00:00Z",
  },
];

export function useWorkspace() {
  return { data: mockWorkspaces, isLoading: false };
}

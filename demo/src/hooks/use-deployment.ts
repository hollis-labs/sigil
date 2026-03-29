// hooks/use-deployment.ts — mock data for demo
import type { Deployment } from "@/types/deployment";

const mockDeployments: Deployment[] = [
  {
    id: "dep-001",
    name: "api-v2.4.1",
    status: "ready",
    environment: "production",
    branch: "main",
    commit: "a1b2c3d",
    author: "chrispian",
    duration: "2m 34s",
    createdAt: "2026-03-29T10:15:00Z",
  },
  {
    id: "dep-002",
    name: "auth-hotfix",
    status: "building",
    environment: "staging",
    branch: "fix/auth-token",
    commit: "e4f5g6h",
    author: "alice",
    duration: "1m 12s",
    createdAt: "2026-03-29T09:42:00Z",
  },
  {
    id: "dep-003",
    name: "cdn-config",
    status: "ready",
    environment: "production",
    branch: "main",
    commit: "i7j8k9l",
    author: "bob",
    duration: "45s",
    createdAt: "2026-03-28T16:30:00Z",
  },
  {
    id: "dep-004",
    name: "worker-update",
    status: "failed",
    environment: "preview",
    branch: "feat/worker-pool",
    commit: "m0n1o2p",
    author: "chrispian",
    duration: "3m 01s",
    createdAt: "2026-03-29T08:05:00Z",
  },
  {
    id: "dep-005",
    name: "db-migration",
    status: "deploying",
    environment: "staging",
    branch: "chore/migrate",
    commit: "q3r4s5t",
    author: "alice",
    duration: "4m 22s",
    createdAt: "2026-03-29T10:30:00Z",
  },
];

export function useDeployment() {
  return { data: mockDeployments, isLoading: false };
}

export function useDeploymentById(id: string) {
  return { data: mockDeployments.find((d) => d.id === id) ?? mockDeployments[0], isLoading: false };
}

// hooks/use-asset.ts — mock data for demo
import type { Asset } from "@/types/asset";

const mockAssets: Asset[] = [
  {
    id: "asset-001",
    name: "api-server-01",
    type: "server",
    provider: "aws",
    status: "healthy",
    region: "us-east-1",
    cost: "$142/mo",
    updatedAt: "2026-03-29T10:15:00Z",
  },
  {
    id: "asset-002",
    name: "db-primary",
    type: "database",
    provider: "aws",
    status: "healthy",
    region: "us-east-1",
    cost: "$380/mo",
    updatedAt: "2026-03-29T09:00:00Z",
  },
  {
    id: "asset-003",
    name: "db-replica",
    type: "database",
    provider: "aws",
    status: "down",
    region: "us-west-2",
    cost: "$380/mo",
    updatedAt: "2026-03-29T10:32:00Z",
  },
  {
    id: "asset-004",
    name: "cdn-global",
    type: "cdn",
    provider: "cloudflare",
    status: "healthy",
    region: "global",
    cost: "$89/mo",
    updatedAt: "2026-03-28T14:00:00Z",
  },
  {
    id: "asset-005",
    name: "redis-cache",
    type: "database",
    provider: "hetzner",
    status: "healthy",
    region: "eu-central",
    cost: "$45/mo",
    updatedAt: "2026-03-29T08:30:00Z",
  },
  {
    id: "asset-006",
    name: "lb-main",
    type: "load-balancer",
    provider: "aws",
    status: "healthy",
    region: "us-east-1",
    cost: "$18/mo",
    updatedAt: "2026-03-29T10:00:00Z",
  },
];

export function useAsset() {
  return { data: mockAssets, isLoading: false };
}

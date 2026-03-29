// hooks/use-user.ts — mock data for demo
import type { User } from "@/types/user";

const mockUsers: User[] = [
  {
    id: "user-001",
    name: "Chrispian",
    email: "chrispian@forge.dev",
    role: "owner",
    status: "active",
    lastActive: "2026-03-29T10:30:00Z",
    teams: "Platform, Infrastructure",
  },
  {
    id: "user-002",
    name: "Alice Chen",
    email: "alice@forge.dev",
    role: "admin",
    status: "active",
    lastActive: "2026-03-29T09:45:00Z",
    teams: "Backend, DevOps",
  },
  {
    id: "user-003",
    name: "Bob Martinez",
    email: "bob@forge.dev",
    role: "member",
    status: "active",
    lastActive: "2026-03-28T17:20:00Z",
    teams: "Frontend",
  },
  {
    id: "user-004",
    name: "Carol Wu",
    email: "carol@forge.dev",
    role: "member",
    status: "invited",
    lastActive: "2026-03-25T14:00:00Z",
    teams: "Design",
  },
  {
    id: "user-005",
    name: "Dave Kim",
    email: "dave@forge.dev",
    role: "viewer",
    status: "suspended",
    lastActive: "2026-03-10T11:00:00Z",
    teams: "Analytics",
  },
];

export function useUser() {
  return { data: mockUsers, isLoading: false };
}

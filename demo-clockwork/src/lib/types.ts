// lib/types.ts — minimal shared types used by hand-written demo components.
// Sigil emits per-datasource types under /types; these are the leftover
// hand-written shapes that the .sigil/components/*.tsx files import.

export type TaskStatus = string;
export type Priority = string;

export interface Comment {
  id: string;
  task_id: string;
  parent_id?: string;
  author: string;
  author_avatar?: string;
  body: string;
  edited_at?: string;
  created_at: string;
}

export interface Artifact {
  id: string;
  run_id: string;
  task_id: string;
  filename: string;
  path: string;
  mime_type: string;
  size_bytes: number;
  content?: string;
  created_at: string;
}

export interface ActivityItem {
  id: string;
  type: string;
  message: string;
  timestamp: string;
  agent?: string;
  task_id?: string;
  tokens?: number;
  cost?: number;
  status?: string;
}

// Loose type aliases — the demo's mock data doesn't enforce these so the
// shape stays permissive. Real apps narrow these via Sigil-generated
// types in @/types.
export interface Run {
  id: string;
  task_id?: string;
  agent?: string;
  model?: string;
  status?: string;
  duration_ms?: number;
  tokens_used?: number;
  total_tokens?: number;
  prompt_tokens?: number;
  completion_tokens?: number;
  cost?: number;
  total_cost?: number;
  prompt?: string;
  result?: string;
  error?: string;
  started_at?: string;
  created_at?: string;
  updated_at?: string;
  logs?: Array<{ level: string; message: string; timestamp: string }>;
}

export interface Checkpoint {
  id: string;
  [key: string]: unknown;
}

export interface Task {
  id: string;
  title?: string;
  status?: string;
  done?: boolean;
  [key: string]: unknown;
}

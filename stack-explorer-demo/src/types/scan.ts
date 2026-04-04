// types/scan.ts

export interface Scan {
  id: string;
  repo_id?: string;
  repo_name?: string;
  blueprint?: string;
  status?: string;
  result_json?: string;
  error_message?: string;
  started_at?: string;
  finished_at?: string;
  created_at?: string;
  updated_at?: string;
}

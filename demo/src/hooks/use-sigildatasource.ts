// hooks/use-sigildatasource.ts — mock data for demo
import type { SigilDataSource } from "@/types/sigildatasource";

const mockDataSources: SigilDataSource[] = [
  { alias: "SigilComponent", description: "Component registry", fieldCount: 4, capabilities: "list, search, filter" },
  { alias: "SigilPage", description: "Page configurations", fieldCount: 4, capabilities: "list, read, create, update, delete" },
  { alias: "SigilTheme", description: "Theme definitions", fieldCount: 3, capabilities: "list, create, read, update, delete" },
  { alias: "SigilDataSource", description: "Datasource manifests", fieldCount: 4, capabilities: "list, create, read, update, delete" },
];

export function useSigilDataSource() {
  return { data: mockDataSources, isLoading: false };
}

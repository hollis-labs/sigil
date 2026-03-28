// hooks/use-sigiltheme.ts — mock data for demo
import type { SigilTheme } from "@/types/sigiltheme";

const mockThemes: SigilTheme[] = [
  { name: "light", description: "Clean, modern light theme — base for all variants" },
  { name: "dark", description: "Dark variant", extends: "light" },
  { name: "default", description: "Default theme — aliases light", extends: "light" },
];

export function useSigilTheme() {
  return { data: mockThemes, isLoading: false };
}

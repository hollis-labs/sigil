// hooks/use-sigilpage.ts — mock data for demo
import type { SigilPage } from "@/types/sigilpage";

const mockPages: SigilPage[] = [
  { id: "sigil-pages", title: "Pages", overlay: "page", componentCount: 8 },
  { id: "sigil-components", title: "Component Browser", overlay: "page", componentCount: 12 },
  { id: "sigil-themes", title: "Themes", overlay: "page", componentCount: 10 },
  { id: "sigil-datasources", title: "Data Sources", overlay: "page", componentCount: 6 },
  { id: "sigil-page-editor", title: "Page Editor", overlay: "page", componentCount: 15 },
  { id: "sigil-preview", title: "Preview", overlay: "page", componentCount: 5 },
];

export function useSigilPage() {
  return { data: mockPages, isLoading: false };
}

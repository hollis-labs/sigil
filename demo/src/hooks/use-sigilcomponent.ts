// hooks/use-sigilcomponent.ts — mock data for demo
import type { SigilComponent } from "@/types/sigilcomponent";

const mockComponents: SigilComponent[] = [
  { type: "button", category: "primitives", description: "Clickable button", propsCount: 5 },
  { type: "input", category: "primitives", description: "Text input field", propsCount: 4 },
  { type: "badge", category: "primitives", description: "Status badge", propsCount: 2 },
  { type: "data-table", category: "data", description: "Sortable, filterable data table", propsCount: 11 },
  { type: "card", category: "layouts", description: "Card container", propsCount: 5 },
  { type: "tabs", category: "layouts", description: "Tab navigation", propsCount: 3 },
  { type: "modal", category: "composites", description: "Dialog overlay", propsCount: 3 },
  { type: "form", category: "forms", description: "Form with fields", propsCount: 3 },
  { type: "heading", category: "primitives", description: "Section heading", propsCount: 2 },
  { type: "select", category: "primitives", description: "Dropdown select", propsCount: 4 },
  { type: "nav-menu", category: "navigation", description: "Navigation menu", propsCount: 3 },
  { type: "chart", category: "data", description: "Data visualization", propsCount: 3 },
];

export function useSigilComponent() {
  return { data: mockComponents, isLoading: false };
}

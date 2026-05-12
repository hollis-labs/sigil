// filter-bar.tsx
// Sigil custom component — flexible filter bar.
// Accepts the renderer's flat-prop shape (filters: [{type, field, label, options}, ...]).

import { useState } from "react";
import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export interface FilterDef {
  type: string;
  field: string;
  label?: string;
  placeholder?: string;
  options?: { label: string; value: string }[];
}

export interface FilterBarProps {
  filters?: FilterDef[];
  // Optional callback fired on any filter change with the full filter state.
  onFilterChange?: (state: Record<string, string>) => void;
  // Optional datasource alias (ignored — page wires data via mock hooks).
  datasource?: string;
}

export function FilterBar({ filters = [], onFilterChange }: FilterBarProps) {
  const [state, setState] = useState<Record<string, string>>({});

  const update = (field: string, value: string) => {
    const next = { ...state, [field]: value };
    setState(next);
    onFilterChange?.(next);
  };

  return (
    <div className="flex flex-wrap items-center gap-2">
      {filters.map((f, i) => {
        if (f.type === "select" && Array.isArray(f.options)) {
          return (
            <div key={i} className="flex items-center gap-2">
              {f.label && (
                <span className="text-xs text-muted-foreground">{f.label}</span>
              )}
              <Select
                value={state[f.field] ?? ""}
                onValueChange={(v) => update(f.field, v ?? "")}
              >
                <SelectTrigger className="h-8 w-[140px] text-xs">
                  <SelectValue placeholder="—" />
                </SelectTrigger>
                <SelectContent>
                  {f.options.map((opt) => (
                    <SelectItem key={opt.value || "__empty__"} value={opt.value || "__empty__"}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          );
        }
        if (f.type === "search") {
          return (
            <div key={i} className="relative">
              <Search className="absolute left-2 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
              <Input
                type="search"
                placeholder={f.placeholder ?? f.label ?? "Search..."}
                value={state[f.field] ?? ""}
                onChange={(e) => update(f.field, e.target.value)}
                className="h-8 w-[200px] pl-7 text-xs"
              />
            </div>
          );
        }
        // date-range / unsupported — render a label only.
        return (
          <span key={i} className="text-xs text-muted-foreground">
            {f.label ?? f.field}
          </span>
        );
      })}
    </div>
  );
}

export default FilterBar;

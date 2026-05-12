// summary-cards.tsx
// Sigil custom component — grid of summary KPI cards.
// Accepts the flat-prop shape emitted by the renderer (cards: [...]) where
// each card may carry arbitrary extra fields beyond the documented shape
// (field, filter, etc.). Unknown fields are tolerated and ignored.

import { TrendingDown, TrendingUp, Minus } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

export interface SummaryCardData {
  label: string;
  value?: string | number;
  trend?: "up" | "down" | "neutral";
  trendLabel?: string;
  variant?: string;
  icon?: React.ReactNode | string;
  // Permissive: accept any extra fields (field, filter, id, …) without typing.
  [key: string]: unknown;
}

export interface SummaryCardsProps {
  cards: SummaryCardData[];
  compact?: boolean;
  // Optional datasource alias (ignored — the page wires data via mock hooks).
  datasource?: string;
}

const variantClass: Record<string, string> = {
  default: "border-zinc-800 bg-zinc-950",
  primary: "border-blue-500/30 bg-blue-500/5",
  success: "border-emerald-500/30 bg-emerald-500/5",
  warning: "border-amber-500/30 bg-amber-500/5",
  danger: "border-red-500/30 bg-red-500/5",
  destructive: "border-red-500/30 bg-red-500/5",
  info: "border-cyan-500/30 bg-cyan-500/5",
};

const trendIcons: Record<string, React.ReactNode> = {
  up: <TrendingUp className="h-3.5 w-3.5 text-emerald-400" />,
  down: <TrendingDown className="h-3.5 w-3.5 text-red-400" />,
  neutral: <Minus className="h-3.5 w-3.5 text-zinc-400" />,
};

export function SummaryCards({ cards, compact = false }: SummaryCardsProps) {
  const list = Array.isArray(cards) ? cards : [];
  return (
    <div className={`grid gap-3 ${compact ? "grid-cols-2 md:grid-cols-4" : "grid-cols-1 sm:grid-cols-2 lg:grid-cols-4"}`}>
      {list.map((card, i) => {
        const variantKey = (typeof card.variant === "string" ? card.variant : "default").toLowerCase();
        const cls = variantClass[variantKey] ?? variantClass.default;
        return (
          <Card key={i} className={cls}>
            <CardContent className="p-3">
              <p className="text-xs uppercase tracking-wider text-muted-foreground">{card.label}</p>
              <p className="mt-1 text-lg font-semibold">{String(card.value ?? "—")}</p>
              {card.trend && (
                <div className="mt-1 flex items-center gap-1 text-xs text-muted-foreground">
                  {trendIcons[card.trend]}
                  {card.trendLabel && <span>{card.trendLabel}</span>}
                </div>
              )}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

export default SummaryCards;

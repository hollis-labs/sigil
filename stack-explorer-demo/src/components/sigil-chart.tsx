"use client";

import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  Line,
  LineChart,
  Pie,
  PieChart,
  PolarAngleAxis,
  PolarGrid,
  Radar,
  RadarChart,
  XAxis,
  YAxis,
} from "recharts";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
  type ChartConfig,
} from "@/components/ui/chart";

interface SigilChartProps {
  type: "line" | "bar" | "area" | "pie" | "donut" | "radar";
  data: Record<string, unknown>[];
  config: ChartConfig;
  xKey?: string;
  dataKeys: string[];
  height?: number;
  showGrid?: boolean;
  showLegend?: boolean;
  stacked?: boolean;
}

export function SigilChart({
  type,
  data,
  config,
  xKey = "label",
  dataKeys,
  height = 250,
  showGrid = true,
  showLegend = false,
  stacked = false,
}: SigilChartProps) {
  if (type === "radar") {
    return (
      <ChartContainer config={config} style={{ height }}>
        <RadarChart data={data} margin={{ top: 8, right: 8, bottom: 8, left: 8 }}>
          <ChartTooltip cursor={false} content={<ChartTooltipContent />} />
          <PolarAngleAxis dataKey={xKey} tick={{ fontSize: 12 }} />
          <PolarGrid />
          {dataKeys.map((key) => (
            <Radar
              key={key}
              dataKey={key}
              stroke={`var(--color-${key})`}
              fill={`var(--color-${key})`}
              fillOpacity={0.15}
              strokeWidth={2}
            />
          ))}
          {showLegend && <ChartLegend content={<ChartLegendContent />} />}
        </RadarChart>
      </ChartContainer>
    );
  }

  if (type === "pie" || type === "donut") {
    const pieKey = dataKeys[0] ?? "value";
    return (
      <ChartContainer config={config} className="mx-auto aspect-square max-h-[250px]">
        <PieChart>
          <ChartTooltip cursor={false} content={<ChartTooltipContent hideLabel />} />
          <Pie
            data={data}
            dataKey={pieKey}
            nameKey={xKey}
            innerRadius={type === "donut" ? "55%" : 0}
            strokeWidth={2}
          />
          {showLegend && <ChartLegend content={<ChartLegendContent />} />}
        </PieChart>
      </ChartContainer>
    );
  }

  const ChartType = type === "bar" ? BarChart : type === "area" ? AreaChart : LineChart;

  return (
    <ChartContainer config={config} style={{ height }}>
      <ChartType accessibilityLayer data={data} margin={{ left: 0, right: 12 }}>
        {showGrid && <CartesianGrid vertical={false} />}
        <XAxis
          dataKey={xKey}
          tickLine={false}
          axisLine={false}
          tickMargin={8}
          tickFormatter={(v: string) => v.length > 5 ? v.slice(0, 3) : v}
        />
        <YAxis tickLine={false} axisLine={false} tickMargin={8} width={40} />
        <ChartTooltip cursor={false} content={<ChartTooltipContent />} />
        {showLegend && <ChartLegend content={<ChartLegendContent />} />}
        {dataKeys.map((key) => {
          const color = `var(--color-${key})`;
          if (type === "bar") {
            return (
              <Bar
                key={key}
                dataKey={key}
                fill={color}
                radius={[4, 4, 0, 0]}
                stackId={stacked ? "stack" : undefined}
              />
            );
          }
          if (type === "area") {
            return (
              <Area
                key={key}
                dataKey={key}
                type="natural"
                fill={color}
                fillOpacity={0.2}
                stroke={color}
                strokeWidth={2}
                stackId={stacked ? "stack" : undefined}
              />
            );
          }
          return (
            <Line
              key={key}
              dataKey={key}
              type="natural"
              stroke={color}
              strokeWidth={2}
              dot={false}
            />
          );
        })}
      </ChartType>
    </ChartContainer>
  );
}

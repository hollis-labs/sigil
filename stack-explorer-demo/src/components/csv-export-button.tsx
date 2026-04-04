"use client";

import { toast } from "sonner";
import { Download } from "lucide-react";
import { Button } from "@/components/ui/button";

interface CsvColumn {
  key: string;
  header: string;
}

interface CsvExportButtonProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: any[];
  columns: CsvColumn[];
  filename?: string;
}

export function CsvExportButton({ data, columns, filename = "export.csv" }: CsvExportButtonProps) {
  const exportCsv = () => {
    if (!data.length) {
      toast("No data to export");
      return;
    }
    const headers = columns.map((c) => c.header);
    const csvRows = [
      headers.join(","),
      ...data.map((row) =>
        columns
          .map((c) => `"${String(row[c.key] ?? "").replace(/"/g, '""')}"`)
          .join(","),
      ),
    ];
    const blob = new Blob([csvRows.join("\n")], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
    toast.success(`Exported ${data.length} rows`);
  };

  return (
    <Button variant="outline" onClick={exportCsv}>
      <Download className="mr-2 h-4 w-4" />
      Export CSV
    </Button>
  );
}

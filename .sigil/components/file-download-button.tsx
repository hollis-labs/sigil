"use client";

import { toast } from "sonner";
import { Database, Download } from "lucide-react";
import { Button } from "@/components/ui/button";

const API_BASE = process.env.NEXT_PUBLIC_SE_API_URL ?? "http://localhost:8081";

interface FileDownloadButtonProps {
  url: string;
  filename: string;
  label?: string;
  icon?: string;
}

export function FileDownloadButton({ url, filename, label = "Download", icon = "download" }: FileDownloadButtonProps) {
  const handleDownload = async () => {
    try {
      const fullUrl = url.startsWith("http") ? url : `${API_BASE}${url}`;
      const res = await fetch(fullUrl);
      if (!res.ok) {
        const errText = await res.text().catch(() => "");
        throw new Error(errText || `Download failed (${res.status})`);
      }
      const disposition = res.headers.get("Content-Disposition");
      const match = disposition?.match(/filename="?([^"]+)"?/);
      const resolvedFilename = match?.[1] ?? filename;
      const blob = await res.blob();
      const a = document.createElement("a");
      a.href = URL.createObjectURL(blob);
      a.download = resolvedFilename;
      a.click();
      URL.revokeObjectURL(a.href);
      toast.success(`Downloaded ${resolvedFilename}`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Download failed");
    }
  };

  const Icon = icon === "database" ? Database : Download;

  return (
    <Button variant="outline" onClick={handleDownload}>
      <Icon className="mr-2 h-4 w-4" />
      {label}
    </Button>
  );
}

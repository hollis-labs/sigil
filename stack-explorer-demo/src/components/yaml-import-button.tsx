"use client";

import { useRef } from "react";
import { toast } from "sonner";
import { Upload } from "lucide-react";
import { Button } from "@/components/ui/button";

const API_BASE = process.env.NEXT_PUBLIC_SE_API_URL ?? "http://localhost:8081";

interface YamlImportButtonProps {
  endpoint: string;
  label?: string;
  onSuccess?: () => void;
}

export function YamlImportButton({ endpoint, label = "Import YAML", onSuccess }: YamlImportButtonProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleImport = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    try {
      const text = await file.text();
      const url = endpoint.startsWith("http") ? endpoint : `${API_BASE}${endpoint}`;
      const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/x-yaml" },
        body: text,
      });
      if (!res.ok) {
        const errText = await res.text().catch(() => "");
        throw new Error(errText || `Import failed (${res.status})`);
      }
      const result = await res.json();
      toast.success(`Imported ${result.imported ?? ""} items`);
      onSuccess?.();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Import failed");
    }
    e.target.value = "";
  };

  return (
    <>
      <input ref={fileInputRef} type="file" accept=".yaml,.yml" className="hidden" onChange={handleImport} />
      <Button variant="outline" onClick={() => fileInputRef.current?.click()}>
        <Upload className="mr-2 h-4 w-4" />
        {label}
      </Button>
    </>
  );
}

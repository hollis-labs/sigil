// lib/api-context.tsx
// Sigil custom provider — Clockwork API client context. Provides a
// browser-side fetch client for the demo's hand-written mock hooks. Wired
// into the app via `providers:` in .sigil/app.yaml; renderer copies this
// file into the output's lib/ directory and wraps the app with <ApiProvider>.

import { createContext, useContext, type ReactNode } from "react";
import { ClockworkApiClient } from "@/lib/clockwork-api";

const ApiContext = createContext<ClockworkApiClient | null>(null);

export function useApi(): ClockworkApiClient {
  const ctx = useContext(ApiContext);
  if (!ctx) {
    throw new Error("useApi() must be used within an ApiProvider");
  }
  return ctx;
}

export function ApiProvider({ children }: { children: ReactNode }) {
  const client = new ClockworkApiClient();
  return (
    <ApiContext.Provider value={client}>
      {children}
    </ApiContext.Provider>
  );
}

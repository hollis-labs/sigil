package reactshadcn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrispian/sigil/internal/config"
	"github.com/chrispian/sigil/internal/renderer"
)

// pluralize adds an "s" to a word, handling common English patterns.
func pluralize(s string) string {
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		c := s[len(s)-2]
		if c != 'a' && c != 'e' && c != 'i' && c != 'o' && c != 'u' {
			return s[:len(s)-1] + "ies"
		}
	}
	return s + "s"
}

// renderProviders generates React Context + Provider files for each provider
// declared on any module.
//
// Built-in providers are generated from a datasource-driven template (the
// historical behavior). Custom providers (those with Source.Component set)
// have their source files copied from .sigil/providers/<file> into the
// output's lib/ directory; the renderer does not synthesize their bodies.
//
// Phase 3.5: providers live per-module. This function walks every module's
// Providers slice and deduplicates by provider ID — two modules declaring the
// same provider (e.g. `lens` on both `se` and `clockwork`) produce a single
// set of `lib/*-context.tsx` (built-ins) or `lib/<basename>` (custom)
// files. Per-module wrap order is handled by the layout renderer, not here.
//
// sigilDir is the path to the .sigil directory; empty means custom-provider
// file copying is skipped (used by tests that don't exercise source copying).
func renderProviders(modules []config.ModuleConfig, sigilDir string) ([]renderer.OutputFile, error) {
	var files []renderer.OutputFile
	seen := map[string]bool{} // provider ID → already emitted
	for _, mod := range modules {
		for _, p := range mod.Providers {
			if seen[p.ID] {
				continue
			}
			seen[p.ID] = true
			if p.IsCustom() {
				copied, err := copyCustomProviderFiles(p, sigilDir)
				if err != nil {
					return nil, err
				}
				files = append(files, copied...)
				continue
			}
			file, err := renderProvider(p)
			if err != nil {
				return nil, err
			}
			files = append(files, file)
		}
	}
	return files, nil
}

// renderProvidersForTests is a thin wrapper used by the existing test cases
// that exercise the per-provider rendering primitives directly. New tests
// should call renderProviders with module slices.
func renderProvidersForTests(providers []config.ProviderConfig, sigilDir string) ([]renderer.OutputFile, error) {
	return renderProviders([]config.ModuleConfig{{ID: "test", Providers: providers}}, sigilDir)
}

// copyCustomProviderFiles copies a custom provider's source + sidecars from
// the .sigil/ tree to the output. Paths in ProviderSource are interpreted
// relative to sigilDir; the output path mirrors the relative path so that
// `providers/active-runs-context.tsx` becomes `lib/active-runs-context.tsx`
// (the directory prefix is rewritten to `lib/` regardless of where it lived
// under .sigil/, keeping import paths predictable as `@/lib/<file>`).
func copyCustomProviderFiles(p config.ProviderConfig, sigilDir string) ([]renderer.OutputFile, error) {
	if sigilDir == "" {
		return nil, fmt.Errorf("provider %q: sigil dir not configured; cannot copy custom provider source", p.ID)
	}
	var out []renderer.OutputFile

	// Main component file.
	mainFile, err := readAndRewriteProviderSource(sigilDir, p.Source.Component)
	if err != nil {
		return nil, fmt.Errorf("provider %q: %w", p.ID, err)
	}
	out = append(out, mainFile)

	// Includes (e.g., SSE bridge, types).
	for _, inc := range p.Source.Includes {
		incFile, err := readAndRewriteProviderSource(sigilDir, inc)
		if err != nil {
			return nil, fmt.Errorf("provider %q include %q: %w", p.ID, inc, err)
		}
		out = append(out, incFile)
	}
	return out, nil
}

// readAndRewriteProviderSource reads a provider source file from sigilDir and
// returns an OutputFile whose path is the file's basename under `lib/`. This
// keeps all generated `@/lib/<file>` imports working regardless of the
// author's chosen organization under .sigil/.
func readAndRewriteProviderSource(sigilDir, relPath string) (renderer.OutputFile, error) {
	srcPath := filepath.Join(sigilDir, relPath)
	data, err := os.ReadFile(srcPath) //nolint:gosec // Path is from validated app config
	if err != nil {
		return renderer.OutputFile{}, fmt.Errorf("reading source %s: %w", relPath, err)
	}
	outPath := "lib/" + filepath.Base(relPath)
	return renderer.OutputFile{
		Path:    outPath,
		Content: data,
	}, nil
}

// providerWrapName returns the React component name to use when wrapping the
// layout/app tree with this provider — the configured export for custom
// providers, or the synthesized PascalCase(ID)+"Provider" for built-ins.
func providerWrapName(p *config.ProviderConfig) string {
	if p.IsCustom() {
		if p.Source.Export != "" {
			return p.Source.Export
		}
		return pascalCase(p.ID) + "Provider"
	}
	return pascalCase(p.ID) + "Provider"
}

// providerImportPath returns the @/lib/... import path (no extension) for the
// provider's exported component. For built-ins this is the generated
// `lib/<id>-context.tsx`; for custom providers it's `lib/<basename of source>`.
func providerImportPath(p *config.ProviderConfig) string {
	if p.IsCustom() {
		return "@/lib/" + stripExt(filepath.Base(p.Source.Component))
	}
	return "@/lib/" + toKebabCase(p.ID) + "-context"
}

// providerHookImport returns the (hookName, hookImportPath) pair to import the
// provider's React context hook in modules that consume the provider via
// useXxxContext(). Built-in providers always export use<ID>Context from the
// same generated file. Custom providers may or may not expose a hook —
// callers should only emit the import when the renderer actually consumes the
// hook (e.g. position: topbar). Returns ("", "") when there's no hook to
// import (custom providers that don't need topbar UI).
func providerHookImport(p *config.ProviderConfig) (name, path string) {
	if p.IsCustom() {
		// Custom providers only need a hook import if the renderer is going
		// to call into them (Position: topbar). Today's topbar UI is
		// built-in only; defer custom-provider topbar support to a later
		// phase.
		return "", ""
	}
	return "use" + pascalCase(p.ID) + "Context", providerImportPath(p)
}

// providerMounts returns the (componentName, importPath) pairs for any mount
// components a custom provider wants rendered inside its scope alongside
// {children}. Each mount name must correspond to an include in
// Source.Includes — the mount is imported from `@/lib/<basename of include>`.
// Built-ins return no mounts.
func providerMounts(p *config.ProviderConfig) []providerMountImport {
	if !p.IsCustom() || len(p.Mounts) == 0 {
		return nil
	}
	// Build basename → import path map from includes (and the main source).
	candidates := map[string]string{}
	for _, inc := range p.Source.Includes {
		base := stripExt(filepath.Base(inc))
		candidates[base] = "@/lib/" + base
	}
	// Main source may also export a mount component (rare but allowed).
	mainBase := stripExt(filepath.Base(p.Source.Component))
	if _, ok := candidates[mainBase]; !ok {
		candidates[mainBase] = "@/lib/" + mainBase
	}

	var out []providerMountImport
	for _, m := range p.Mounts {
		// Find an include whose basename (without extension) is a plausible
		// match for the mount component name. Authors typically name files
		// like `sse-active-runs-bridge.tsx` exporting `SSEActiveRunsBridge`,
		// so we kebab-case the mount name and look for that file.
		fileBase := kebabCaseFromPascal(m)
		path, ok := candidates[fileBase]
		if !ok {
			// Fallback: if no kebab match, take the first include we haven't
			// already assigned. This keeps the schema permissive when authors
			// have unusual file names.
			for _, p := range candidates {
				if !pathAssigned(out, p) {
					path = p
					break
				}
			}
		}
		out = append(out, providerMountImport{Name: m, Path: path})
	}
	return out
}

func pathAssigned(used []providerMountImport, p string) bool {
	for _, u := range used {
		if u.Path == p {
			return true
		}
	}
	return false
}

// providerMountImport names a custom-provider sibling component to render
// alongside {children} inside the provider scope.
type providerMountImport struct {
	Name string // React component name (e.g. SSEActiveRunsBridge)
	Path string // import path (e.g. @/lib/sse-active-runs-bridge)
}

// kebabCaseFromPascal converts PascalCase to kebab-case (SSEActiveRunsBridge
// → sse-active-runs-bridge). Runs of uppercase letters are kept together
// (acronym handling), with a hyphen before the start of each "word".
func kebabCaseFromPascal(s string) string {
	if s == "" {
		return s
	}
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		isUpper := c >= 'A' && c <= 'Z'
		if isUpper && i > 0 {
			prev := s[i-1]
			prevUpper := prev >= 'A' && prev <= 'Z'
			// Insert hyphen at lower→Upper boundary, or at the end of an
			// acronym run (Upper→Upper followed by lowercase: "SSEActive"
			// → "sse-active").
			if !prevUpper {
				b = append(b, '-')
			} else if i+1 < len(s) {
				next := s[i+1]
				if next >= 'a' && next <= 'z' {
					b = append(b, '-')
				}
			}
		}
		if isUpper {
			b = append(b, c+('a'-'A'))
		} else {
			b = append(b, c)
		}
	}
	return string(b)
}

// stripExt returns the filename with any extension stripped (e.g.
// "active-runs-context.tsx" → "active-runs-context").
func stripExt(name string) string {
	if idx := strings.LastIndexByte(name, '.'); idx > 0 {
		return name[:idx]
	}
	return name
}

func renderProvider(p config.ProviderConfig) (renderer.OutputFile, error) { //nolint:unparam // error kept for interface consistency
	id := p.ID
	typeName := toPascalCase(p.Datasource)
	providerName := toPascalCase(id) + "Provider"
	hookName := "use" + toPascalCase(id) + "Context"
	trackField := p.TrackField
	if trackField == "" {
		trackField = "id"
	}
	labelField := p.LabelField
	if labelField == "" {
		labelField = "name"
	}
	defaultVal := p.Default
	if defaultVal == "" {
		defaultVal = ""
	}

	resource := pluralize(strings.ToLower(p.Datasource))

	content := fmt.Sprintf(`// Generated by Sigil — do not edit manually
"use client";

import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import { fetchList } from "@/lib/api";

interface %s {
  %s: string;
  %s: string;
  [key: string]: unknown;
}

interface %sContextValue {
  current%sId: string;
  current%sName: string;
  set%sId: (id: string) => void;
  %s: %s[];
  isLoading: boolean;
}

const %sContext = createContext<%sContextValue>({
  current%sId: %q,
  current%sName: "",
  set%sId: () => {},
  %s: [],
  isLoading: true,
});

export function %s() {
  return useContext(%sContext);
}

export function %s({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<%s[]>([]);
  const [currentId, setCurrentId] = useState(%q);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchList<%s>(%q)
      .then((data) => {
        setItems(data);
        if (data.length > 0 && !currentId) setCurrentId(String(data[0].%s));
      })
      .catch(() => setItems([]))
      .finally(() => setIsLoading(false));
  }, []);

  const currentName = items.find((item) => String(item.%s) === currentId)?.%s ?? "";

  return (
    <%sContext value={{
      current%sId: currentId,
      current%sName: currentName,
      set%sId: setCurrentId,
      %s: items,
      isLoading,
    }}>
      {children}
    </%sContext>
  );
}
`,
		// Interface fields
		typeName, trackField, labelField,
		// Context value interface
		toPascalCase(id), toPascalCase(id), toPascalCase(id),
		toPascalCase(id),
		pluralize(toCamelCase(p.Datasource)), typeName,
		// Context default
		toPascalCase(id), toPascalCase(id),
		toPascalCase(id), defaultVal,
		toPascalCase(id),
		toPascalCase(id),
		pluralize(toCamelCase(p.Datasource)),
		// Hook
		hookName, toPascalCase(id),
		// Provider function
		providerName, typeName,
		defaultVal,
		// useEffect
		typeName, resource,
		trackField,
		// currentName
		trackField, labelField,
		// JSX
		toPascalCase(id),
		toPascalCase(id), toPascalCase(id), toPascalCase(id),
		pluralize(toCamelCase(p.Datasource)),
		toPascalCase(id),
	)

	fileName := fmt.Sprintf("lib/%s-context.tsx", toKebabCase(id))
	return renderer.OutputFile{
		Path:    fileName,
		Content: []byte(content),
	}, nil
}

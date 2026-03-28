// Package preview generates standalone HTML previews of Sigil page configs.
package preview

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/chrispian/sigil/internal/config"
	"github.com/chrispian/sigil/internal/renderer"
)

// PreviewConfig holds options for preview generation.
type PreviewConfig struct {
	Page        *config.Page
	Theme       *renderer.ThemeConfig
	DataSources map[string]*renderer.DataSourceManifest
}

// GenerateHTML produces a self-contained HTML string for a page config.
func GenerateHTML(cfg *PreviewConfig) ([]byte, error) {
	var buf bytes.Buffer

	page := cfg.Page

	buf.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	buf.WriteString(fmt.Sprintf("  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"))
	buf.WriteString(fmt.Sprintf("  <title>%s — Sigil Preview</title>\n", page.Title))
	buf.WriteString("  <script src=\"https://cdn.tailwindcss.com\"></script>\n")
	buf.WriteString("  <style>\n")

	// Inline theme CSS
	writeThemeCSS(&buf, cfg.Theme)
	writeComponentCSS(&buf)

	buf.WriteString("  </style>\n")
	buf.WriteString("</head>\n")
	buf.WriteString("<body class=\"bg-[rgb(var(--sigil-background))] text-[rgb(var(--sigil-text))] font-[var(--sigil-font-sans)] min-h-screen\">\n")

	// Page wrapper based on overlay type
	switch page.Overlay {
	case "modal":
		buf.WriteString("  <div class=\"fixed inset-0 bg-black/50 flex items-center justify-center p-8\">\n")
		buf.WriteString("    <div class=\"bg-[rgb(var(--sigil-surface))] rounded-lg shadow-xl max-w-2xl w-full max-h-[80vh] overflow-auto p-6\">\n")
		renderPreviewComponent(&buf, &page.Layout, 3, cfg)
		buf.WriteString("    </div>\n")
		buf.WriteString("  </div>\n")
	case "sheet", "drawer":
		buf.WriteString("  <div class=\"fixed inset-0 flex\">\n")
		buf.WriteString("    <div class=\"flex-1 bg-black/30\"></div>\n")
		buf.WriteString("    <div class=\"w-[480px] bg-[rgb(var(--sigil-surface))] shadow-xl overflow-auto p-6\">\n")
		renderPreviewComponent(&buf, &page.Layout, 3, cfg)
		buf.WriteString("    </div>\n")
		buf.WriteString("  </div>\n")
	default: // page, fullscreen
		buf.WriteString("  <div class=\"max-w-6xl mx-auto p-6\">\n")
		renderPreviewComponent(&buf, &page.Layout, 2, cfg)
		buf.WriteString("  </div>\n")
	}

	buf.WriteString("</body>\n</html>\n")

	return buf.Bytes(), nil
}

func writeThemeCSS(buf *bytes.Buffer, theme *renderer.ThemeConfig) {
	buf.WriteString("    :root {\n")

	if theme == nil {
		// Sensible defaults when no theme is loaded
		writeDefaultTokens(buf)
		buf.WriteString("    }\n")
		return
	}

	// Categories with optional prefix to avoid key collisions (e.g. radius.lg vs shadows.lg)
	type catEntry struct {
		name   string
		prefix string
	}
	categories := []catEntry{
		{"colors", ""},
		{"typography", ""},
		{"spacing", ""},
		{"radius", ""},
		{"shadows", "shadow-"},
	}
	for _, cat := range categories {
		tokens, ok := theme.Tokens[cat.name]
		if !ok || len(tokens) == 0 {
			continue
		}
		keys := sortedKeys(tokens)
		for _, key := range keys {
			val := fmt.Sprintf("%v", tokens[key])
			fmt.Fprintf(buf, "      --sigil-%s%s: %s;\n", cat.prefix, key, val)
		}
	}

	buf.WriteString("    }\n")
}

func writeDefaultTokens(buf *bytes.Buffer) {
	defaults := map[string]string{
		"background":  "255 255 255",
		"surface":     "249 250 251",
		"text":        "17 24 39",
		"text-muted":  "107 114 128",
		"accent":      "79 70 229",
		"border":      "229 231 235",
		"font-sans":   "'Inter', system-ui, sans-serif",
		"font-mono":   "'JetBrains Mono', monospace",
		"radius-sm":   "0.25rem",
		"radius-md":   "0.375rem",
		"radius-lg":   "0.5rem",
		"radius-full": "9999px",
	}
	keys := make([]string, 0, len(defaults))
	for k := range defaults {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("      --sigil-%s: %s;\n", k, defaults[k]))
	}
}

func writeComponentCSS(buf *bytes.Buffer) {
	buf.WriteString(`
    /* ── Buttons (shadcn-style) ─────────────────────────── */
    .sigil-btn {
      display: inline-flex; align-items: center; justify-content: center; gap: 0.5rem;
      white-space: nowrap; font-size: 0.875rem; font-weight: 500; line-height: 1.25rem;
      padding: 0.5rem 1rem; border: none; cursor: pointer;
      border-radius: calc(var(--sigil-md, 0.375rem));
      transition: background-color 150ms cubic-bezier(0.4,0,0.2,1), color 150ms, box-shadow 150ms;
      outline: none;
    }
    .sigil-btn:focus-visible {
      box-shadow: 0 0 0 2px rgb(var(--sigil-background)), 0 0 0 4px rgb(var(--sigil-ring, var(--sigil-primary)));
    }
    .sigil-btn:disabled { pointer-events: none; opacity: 0.5; }
    .sigil-btn-primary {
      background: rgb(var(--sigil-primary)); color: rgb(var(--sigil-primary-foreground, 255 255 255));
    }
    .sigil-btn-primary:hover { background: rgb(var(--sigil-primary) / 0.9); }
    .sigil-btn-secondary {
      background: rgb(var(--sigil-secondary)); color: rgb(var(--sigil-secondary-foreground));
    }
    .sigil-btn-secondary:hover { background: rgb(var(--sigil-secondary) / 0.8); }
    .sigil-btn-destructive {
      background: rgb(var(--sigil-danger)); color: rgb(var(--sigil-danger-foreground, 255 255 255));
    }
    .sigil-btn-destructive:hover { background: rgb(var(--sigil-danger) / 0.9); }
    .sigil-btn-outline {
      background: transparent; color: rgb(var(--sigil-foreground, var(--sigil-text)));
      border: 1px solid rgb(var(--sigil-input, var(--sigil-border)));
    }
    .sigil-btn-outline:hover { background: rgb(var(--sigil-accent)); color: rgb(var(--sigil-accent-foreground, var(--sigil-text))); }
    .sigil-btn-ghost {
      background: transparent; color: rgb(var(--sigil-foreground, var(--sigil-text)));
    }
    .sigil-btn-ghost:hover { background: rgb(var(--sigil-accent)); color: rgb(var(--sigil-accent-foreground, var(--sigil-text))); }
    .sigil-btn-link {
      background: transparent; color: rgb(var(--sigil-primary));
      padding: 0; text-decoration: underline; text-underline-offset: 4px;
    }
    .sigil-btn-link:hover { text-decoration-thickness: 2px; }
    .sigil-btn-sm { height: 2rem; padding: 0.25rem 0.75rem; font-size: 0.75rem; border-radius: calc(var(--sigil-sm, 0.125rem)); }
    .sigil-btn-lg { height: 2.75rem; padding: 0.5rem 2rem; font-size: 1rem; border-radius: calc(var(--sigil-md, 0.375rem)); }
    .sigil-btn-icon { height: 2.25rem; width: 2.25rem; padding: 0; }

    /* ── Badge (shadcn-style) ───────────────────────────── */
    .sigil-badge {
      display: inline-flex; align-items: center;
      padding: 0.125rem 0.625rem; border-radius: 9999px;
      font-size: 0.75rem; font-weight: 600; line-height: 1rem;
      border: 1px solid transparent;
      transition: background-color 150ms, color 150ms;
    }
    .sigil-badge-default {
      background: rgb(var(--sigil-primary)); color: rgb(var(--sigil-primary-foreground, 255 255 255));
    }
    .sigil-badge-secondary {
      background: rgb(var(--sigil-secondary)); color: rgb(var(--sigil-secondary-foreground));
    }
    .sigil-badge-success { background: rgb(var(--sigil-success) / 0.12); color: rgb(var(--sigil-success)); border-color: rgb(var(--sigil-success) / 0.2); }
    .sigil-badge-warning { background: rgb(var(--sigil-warning) / 0.12); color: rgb(var(--sigil-warning)); border-color: rgb(var(--sigil-warning) / 0.2); }
    .sigil-badge-danger { background: rgb(var(--sigil-danger) / 0.12); color: rgb(var(--sigil-danger)); border-color: rgb(var(--sigil-danger) / 0.2); }
    .sigil-badge-error { background: rgb(var(--sigil-danger) / 0.12); color: rgb(var(--sigil-danger)); border-color: rgb(var(--sigil-danger) / 0.2); }
    .sigil-badge-info { background: rgb(var(--sigil-info) / 0.12); color: rgb(var(--sigil-info)); border-color: rgb(var(--sigil-info) / 0.2); }
    .sigil-badge-outline {
      background: transparent; color: rgb(var(--sigil-foreground, var(--sigil-text)));
      border-color: rgb(var(--sigil-border));
    }
    .sigil-badge-destructive {
      background: rgb(var(--sigil-danger)); color: rgb(var(--sigil-danger-foreground, 255 255 255));
    }

    /* ── Input (shadcn-style) ───────────────────────────── */
    .sigil-input {
      display: flex; width: 100%; height: 2.5rem;
      padding: 0.5rem 0.75rem;
      border-radius: calc(var(--sigil-md, 0.375rem));
      border: 1px solid rgb(var(--sigil-input, var(--sigil-border)));
      background: transparent;
      color: rgb(var(--sigil-foreground, var(--sigil-text)));
      font-size: 0.875rem; line-height: 1.25rem;
      outline: none;
      transition: border-color 150ms, box-shadow 150ms;
    }
    .sigil-input::placeholder { color: rgb(var(--sigil-muted-foreground, var(--sigil-text-muted))); }
    .sigil-input:focus {
      border-color: rgb(var(--sigil-ring, var(--sigil-primary)));
      box-shadow: 0 0 0 2px rgb(var(--sigil-background)), 0 0 0 4px rgb(var(--sigil-ring, var(--sigil-primary)));
    }
    .sigil-input:disabled { opacity: 0.5; cursor: not-allowed; }
    select.sigil-input { appearance: none; padding-right: 2rem; background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M3 5l3 3 3-3'/%3E%3C/svg%3E"); background-repeat: no-repeat; background-position: right 0.75rem center; }
    textarea.sigil-input { height: auto; min-height: 5rem; resize: vertical; }

    /* ── Label ───────────────────────────────────────────── */
    .sigil-label {
      font-size: 0.875rem; font-weight: 500; line-height: 1;
      color: rgb(var(--sigil-foreground, var(--sigil-text)));
    }

    /* ── Separator ──────────────────────────────────────── */
    .sigil-separator {
      border: none; border-top: 1px solid rgb(var(--sigil-border));
      margin: 0; flex-shrink: 0;
    }

    /* ── Card ────────────────────────────────────────────── */
    .sigil-card {
      background: rgb(var(--sigil-card, var(--sigil-surface)));
      color: rgb(var(--sigil-card-foreground, var(--sigil-text)));
      border: 1px solid rgb(var(--sigil-border));
      border-radius: calc(var(--sigil-lg, 0.5rem));
      box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.05);
      overflow: hidden;
    }
    .sigil-card-header { padding: 1.5rem 1.5rem 0; display: flex; flex-direction: column; gap: 0.375rem; }
    .sigil-card-title { font-size: 1.5rem; font-weight: 600; line-height: 1; letter-spacing: -0.025em; }
    .sigil-card-description { font-size: 0.875rem; color: rgb(var(--sigil-muted-foreground, var(--sigil-text-muted))); }
    .sigil-card-content { padding: 1.5rem; }
    .sigil-card-footer { padding: 0 1.5rem 1.5rem; display: flex; align-items: center; }

    /* ── Form ────────────────────────────────────────────── */
    .sigil-form { max-width: 32rem; }

    /* ── Data Table (shadcn-style) ──────────────────────── */
    .sigil-data-table { width: 100%; overflow: auto; border: 1px solid rgb(var(--sigil-border)); border-radius: calc(var(--sigil-lg, 0.5rem)); }
    .sigil-data-table table { width: 100%; border-collapse: collapse; caption-side: bottom; font-size: 0.875rem; }
    .sigil-data-table th {
      padding: 0.75rem 1rem; text-align: left;
      font-weight: 500; color: rgb(var(--sigil-muted-foreground, var(--sigil-text-muted)));
      border-bottom: 1px solid rgb(var(--sigil-border));
      height: 2.5rem; white-space: nowrap;
    }
    .sigil-data-table td {
      padding: 0.75rem 1rem; text-align: left;
      border-bottom: 1px solid rgb(var(--sigil-border));
      vertical-align: middle;
    }
    .sigil-data-table tbody tr { transition: background-color 150ms; }
    .sigil-data-table tbody tr:hover { background: rgb(var(--sigil-muted) / 0.5); }
    .sigil-data-table tbody tr:last-child td { border-bottom: none; }

    /* ── Alert (shadcn-style) ───────────────────────────── */
    .sigil-alert {
      position: relative; width: 100%;
      padding: 1rem 1rem 1rem 1rem;
      border-radius: calc(var(--sigil-lg, 0.5rem));
      border: 1px solid rgb(var(--sigil-border));
      font-size: 0.875rem; line-height: 1.5;
    }
    .sigil-alert-info { border-color: rgb(var(--sigil-info) / 0.3); color: rgb(var(--sigil-info)); }
    .sigil-alert-warning { border-color: rgb(var(--sigil-warning) / 0.3); color: rgb(var(--sigil-warning)); }
    .sigil-alert-error { border-color: rgb(var(--sigil-danger)); color: rgb(var(--sigil-danger)); }
    .sigil-alert-success { border-color: rgb(var(--sigil-success) / 0.3); color: rgb(var(--sigil-success)); }

    /* ── Progress ────────────────────────────────────────── */
    .sigil-progress-wrap {
      position: relative; width: 100%; height: 0.75rem; overflow: hidden;
      border-radius: 9999px; background: rgb(var(--sigil-secondary, var(--sigil-surface-2)));
    }
    .sigil-progress-bar {
      height: 100%; border-radius: 9999px;
      background: rgb(var(--sigil-primary));
      transition: width 300ms ease;
    }

    /* ── Avatar ──────────────────────────────────────────── */
    .sigil-avatar {
      display: inline-flex; align-items: center; justify-content: center;
      border-radius: 9999px; overflow: hidden; flex-shrink: 0;
      background: rgb(var(--sigil-muted, var(--sigil-surface-2)));
      color: rgb(var(--sigil-foreground, var(--sigil-text)));
      font-weight: 600; font-size: 0.75rem;
    }
    .sigil-avatar-sm { width: 2rem; height: 2rem; }
    .sigil-avatar-md { width: 2.5rem; height: 2.5rem; }
    .sigil-avatar-lg { width: 3rem; height: 3rem; }

    /* ── Modal / Dialog (shadcn-style) ──────────────────── */
    .sigil-modal {
      background: rgb(var(--sigil-background));
      border: 1px solid rgb(var(--sigil-border));
      border-radius: calc(var(--sigil-lg, 0.5rem));
      box-shadow: 0 25px 50px -12px rgb(0 0 0 / 0.25);
      overflow: hidden;
    }
    .sigil-modal-header {
      display: flex; flex-direction: column; gap: 0.375rem;
      padding: 1.5rem 1.5rem 0;
    }
    .sigil-modal-content { padding: 1.5rem; }

    /* ── Sheet / Drawer ─────────────────────────────────── */
    .sigil-sheet {
      background: rgb(var(--sigil-background));
      border-left: 1px solid rgb(var(--sigil-border));
      box-shadow: -10px 0 30px -5px rgb(0 0 0 / 0.1);
    }
    .sigil-sheet-header {
      display: flex; flex-direction: column; gap: 0.375rem;
      padding: 1.5rem 1.5rem 0;
    }
    .sigil-sheet-content { padding: 1.5rem; }

    /* ── Tabs (shadcn-style) ─────────────────────────────── */
    .sigil-tabs-list {
      display: inline-flex; align-items: center;
      padding: 0.25rem; gap: 0.125rem;
      background: rgb(var(--sigil-muted, var(--sigil-surface-2)));
      border-radius: calc(var(--sigil-md, 0.375rem));
    }
    .sigil-tab-trigger {
      display: inline-flex; align-items: center; justify-content: center;
      padding: 0.375rem 0.75rem; border: none; cursor: pointer;
      font-size: 0.875rem; font-weight: 500; line-height: 1.25rem;
      background: transparent; color: rgb(var(--sigil-muted-foreground, var(--sigil-text-muted)));
      border-radius: calc(var(--sigil-sm, 0.125rem));
      white-space: nowrap;
      transition: background-color 150ms, color 150ms, box-shadow 150ms;
    }
    .sigil-tab-trigger:hover { color: rgb(var(--sigil-foreground, var(--sigil-text))); }
    .sigil-tab-trigger-active {
      background: rgb(var(--sigil-background));
      color: rgb(var(--sigil-foreground, var(--sigil-text)));
      box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
    }
`)
}

func renderPreviewComponent(buf *bytes.Buffer, c *config.Component, depth int, cfg *PreviewConfig) {
	prefix := strings.Repeat("  ", depth)

	switch c.Type {
	case "rows":
		gap := getPropString(c.Props, "gap", "4")
		padding := getPropString(c.Props, "padding", "0")
		align := getPropString(c.Props, "align", "")
		classes := fmt.Sprintf("flex flex-col gap-%s", gap)
		if padding != "0" {
			classes += fmt.Sprintf(" p-%s", padding)
		}
		if align == "center" {
			classes += " items-center"
		}
		fmt.Fprintf(buf, "%s<div class=%q>\n", prefix, classes)
		for i := range c.Children {
			renderPreviewComponent(buf, &c.Children[i], depth+1, cfg)
		}
		fmt.Fprintf(buf, "%s</div>\n", prefix)

	case "columns":
		gap := getPropString(c.Props, "gap", "4")
		justify := getPropString(c.Props, "justify", "")
		align := getPropString(c.Props, "align", "")
		classes := fmt.Sprintf("flex flex-row gap-%s", gap)
		if justify == "between" {
			classes += " justify-between"
		} else if justify == "center" {
			classes += " justify-center"
		} else if justify == "end" {
			classes += " justify-end"
		}
		if align == "center" {
			classes += " items-center"
		}
		fmt.Fprintf(buf, "%s<div class=%q>\n", prefix, classes)
		for i := range c.Children {
			renderPreviewComponent(buf, &c.Children[i], depth+1, cfg)
		}
		fmt.Fprintf(buf, "%s</div>\n", prefix)

	case "grid":
		cols := getPropString(c.Props, "columns", "3")
		gap := getPropString(c.Props, "gap", "4")
		classes := fmt.Sprintf("grid grid-cols-%s gap-%s", cols, gap)
		fmt.Fprintf(buf, "%s<div class=%q>\n", prefix, classes)
		for i := range c.Children {
			renderPreviewComponent(buf, &c.Children[i], depth+1, cfg)
		}
		fmt.Fprintf(buf, "%s</div>\n", prefix)

	case "heading":
		level := getPropString(c.Props, "level", "2")
		text := getPropString(c.Props, "text", "")
		tag := "h" + level
		classes := headingClasses(level)
		fmt.Fprintf(buf, "%s<%s class=%q>%s</%s>\n", prefix, tag, classes, text, tag)

	case "text":
		content := getPropString(c.Props, "text", "")
		variant := getPropString(c.Props, "variant", "body")
		classes := textClasses(variant)
		fmt.Fprintf(buf, "%s<p class=%q>%s</p>\n", prefix, classes, content)

	case "button":
		label := getPropString(c.Props, "label", "Button")
		variant := getPropString(c.Props, "variant", "primary")
		icon := getPropString(c.Props, "icon", "")
		classes := fmt.Sprintf("sigil-btn sigil-btn-%s", variant)
		iconHTML := ""
		if icon != "" {
			iconHTML = fmt.Sprintf(`<span class="text-sm">%s</span> `, iconSymbol(icon))
		}
		fmt.Fprintf(buf, "%s<button class=%q>%s%s</button>\n", prefix, classes, iconHTML, label)

	case "search-bar":
		placeholder := getPropString(c.Props, "placeholder", "Search...")
		fmt.Fprintf(buf, "%s<input type=\"search\" class=\"sigil-input w-full\" placeholder=%q/>\n", prefix, placeholder)

	case "data-table":
		renderPreviewDataTable(buf, c, prefix, cfg)

	case "form":
		renderPreviewForm(buf, c, prefix, depth, cfg)

	case "badge":
		text := getPropString(c.Props, "text", "")
		variant := getPropString(c.Props, "variant", "default")
		classes := fmt.Sprintf("sigil-badge sigil-badge-%s", variant)
		fmt.Fprintf(buf, "%s<span class=%q>%s</span>\n", prefix, classes, text)

	case "input":
		name := getPropString(c.Props, "name", "")
		label := getPropString(c.Props, "label", "")
		inputType := getPropString(c.Props, "type", "text")
		placeholder := getPropString(c.Props, "placeholder", "")
		if label != "" {
			fmt.Fprintf(buf, "%s<label class=\"sigil-label\">%s</label>\n", prefix, label)
		}
		fmt.Fprintf(buf, "%s<input type=%q name=%q class=\"sigil-input\" placeholder=%q/>\n", prefix, inputType, name, placeholder)

	case "select":
		renderPreviewSelect(buf, c, prefix)

	case "modal", "sheet":
		renderPreviewOverlay(buf, c, prefix, depth, cfg)

	case "separator":
		fmt.Fprintf(buf, "%s<hr class=\"sigil-separator\"/>\n", prefix)

	case "progress":
		value := getPropString(c.Props, "value", "0")
		fmt.Fprintf(buf, "%s<div class=\"sigil-progress-wrap\">\n", prefix)
		fmt.Fprintf(buf, "%s  <div class=\"sigil-progress-bar\" style=\"width: %s%%\"></div>\n", prefix, value)
		fmt.Fprintf(buf, "%s</div>\n", prefix)

	case "avatar":
		alt := getPropString(c.Props, "alt", "User")
		size := getPropString(c.Props, "size", "md")
		classes := fmt.Sprintf("sigil-avatar sigil-avatar-%s", size)
		fmt.Fprintf(buf, "%s<div class=%q>%s</div>\n", prefix, classes, alt[:1])

	case "alert":
		message := getPropString(c.Props, "message", "")
		variant := getPropString(c.Props, "variant", "info")
		classes := fmt.Sprintf("sigil-alert sigil-alert-%s", variant)
		fmt.Fprintf(buf, "%s<div class=%q>%s</div>\n", prefix, classes, message)

	case "icon":
		name := getPropString(c.Props, "name", "")
		fmt.Fprintf(buf, "%s<span title=%q>%s</span>\n", prefix, name, iconSymbol(name))

	case "label":
		text := getPropString(c.Props, "text", "")
		fmt.Fprintf(buf, "%s<label class=\"sigil-label\">%s</label>\n", prefix, text)

	case "spacer":
		size := getPropString(c.Props, "size", "4")
		fmt.Fprintf(buf, "%s<div class=\"h-%s\"></div>\n", prefix, size)

	case "tabs":
		renderPreviewTabs(buf, c, prefix, depth, cfg)

	default:
		idAttr := ""
		if c.ID != "" {
			idAttr = fmt.Sprintf(` id=%q`, c.ID)
		}
		fmt.Fprintf(buf, "%s<div data-component=%q%s class=\"border border-dashed border-[rgb(var(--sigil-border))] rounded p-4\">\n", prefix, c.Type, idAttr)
		fmt.Fprintf(buf, "%s  <span class=\"text-xs text-[rgb(var(--sigil-text-muted))]\">[%s]</span>\n", prefix, c.Type)
		for i := range c.Children {
			renderPreviewComponent(buf, &c.Children[i], depth+1, cfg)
		}
		fmt.Fprintf(buf, "%s</div>\n", prefix)
	}
}

func renderPreviewDataTable(buf *bytes.Buffer, c *config.Component, prefix string, cfg *PreviewConfig) {
	datasource := getPropString(c.Props, "datasource", "")

	fmt.Fprintf(buf, "%s<div class=\"sigil-data-table\">\n", prefix)
	fmt.Fprintf(buf, "%s  <table>\n", prefix)
	fmt.Fprintf(buf, "%s    <thead>\n", prefix)
	fmt.Fprintf(buf, "%s      <tr>\n", prefix)

	// Extract columns
	var columns []colDef

	if cols, ok := c.Props["columns"]; ok {
		if colSlice, ok := cols.([]interface{}); ok {
			for _, col := range colSlice {
				if colMap, ok := col.(map[string]interface{}); ok {
					columns = append(columns, colDef{
						Label: fmt.Sprintf("%v", colMap["label"]),
						Field: fmt.Sprintf("%v", colMap["field"]),
					})
				}
			}
		}
	}

	for _, col := range columns {
		fmt.Fprintf(buf, "%s        <th>%s</th>\n", prefix, col.Label)
	}

	fmt.Fprintf(buf, "%s      </tr>\n", prefix)
	fmt.Fprintf(buf, "%s    </thead>\n", prefix)
	fmt.Fprintf(buf, "%s    <tbody>\n", prefix)

	// Generate mock data rows
	rows := generateMockRows(datasource, columns, cfg)
	for _, row := range rows {
		fmt.Fprintf(buf, "%s      <tr>\n", prefix)
		for _, col := range columns {
			val := row[col.Field]
			fmt.Fprintf(buf, "%s        <td>%s</td>\n", prefix, val)
		}
		fmt.Fprintf(buf, "%s      </tr>\n", prefix)
	}

	fmt.Fprintf(buf, "%s    </tbody>\n", prefix)
	fmt.Fprintf(buf, "%s  </table>\n", prefix)
	fmt.Fprintf(buf, "%s</div>\n", prefix)
}

func renderPreviewForm(buf *bytes.Buffer, c *config.Component, prefix string, depth int, cfg *PreviewConfig) {
	fmt.Fprintf(buf, "%s<form class=\"sigil-form flex flex-col gap-4\" onsubmit=\"event.preventDefault()\">\n", prefix)

	if fields, ok := c.Props["fields"]; ok {
		if fieldSlice, ok := fields.([]interface{}); ok {
			for _, f := range fieldSlice {
				if fm, ok := f.(map[string]interface{}); ok {
					name := fmt.Sprintf("%v", fm["name"])
					label := fmt.Sprintf("%v", fm["label"])
					fieldType := fmt.Sprintf("%v", fm["type"])
					placeholder := ""
					if p, ok := fm["placeholder"]; ok {
						placeholder = fmt.Sprintf("%v", p)
					}

					fmt.Fprintf(buf, "%s  <div class=\"flex flex-col gap-1\">\n", prefix)
					fmt.Fprintf(buf, "%s    <label class=\"sigil-label\">%s</label>\n", prefix, label)

					switch fieldType {
					case "textarea":
						rows := "3"
						if r, ok := fm["rows"]; ok {
							rows = fmt.Sprintf("%v", r)
						}
						fmt.Fprintf(buf, "%s    <textarea name=%q class=\"sigil-input\" placeholder=%q rows=%q></textarea>\n",
							prefix, name, placeholder, rows)
					case "select":
						fmt.Fprintf(buf, "%s    <select name=%q class=\"sigil-input\">\n", prefix, name)
						if opts, ok := fm["options"]; ok {
							if optSlice, ok := opts.([]interface{}); ok {
								for _, opt := range optSlice {
									if om, ok := opt.(map[string]interface{}); ok {
										val := fmt.Sprintf("%v", om["value"])
										lab := fmt.Sprintf("%v", om["label"])
										fmt.Fprintf(buf, "%s      <option value=%q>%s</option>\n", prefix, val, lab)
									}
								}
							}
						}
						fmt.Fprintf(buf, "%s    </select>\n", prefix)
					case "checkbox":
						fmt.Fprintf(buf, "%s    <input type=\"checkbox\" name=%q class=\"accent-[rgb(var(--sigil-accent))]\"/>\n", prefix, name)
					default:
						fmt.Fprintf(buf, "%s    <input type=%q name=%q class=\"sigil-input\" placeholder=%q/>\n",
							prefix, fieldType, name, placeholder)
					}

					fmt.Fprintf(buf, "%s  </div>\n", prefix)
				}
			}
		}
	}

	submitLabel := "Submit"
	if submit, ok := c.Props["submit"]; ok {
		if sm, ok := submit.(map[string]interface{}); ok {
			if l, ok := sm["label"]; ok {
				submitLabel = fmt.Sprintf("%v", l)
			}
		}
	}
	fmt.Fprintf(buf, "%s  <button type=\"submit\" class=\"sigil-btn sigil-btn-primary\">%s</button>\n", prefix, submitLabel)
	fmt.Fprintf(buf, "%s</form>\n", prefix)
}

func renderPreviewSelect(buf *bytes.Buffer, c *config.Component, prefix string) {
	name := getPropString(c.Props, "name", "")
	label := getPropString(c.Props, "label", "")
	if label != "" {
		fmt.Fprintf(buf, "%s<label class=\"sigil-label\">%s</label>\n", prefix, label)
	}
	fmt.Fprintf(buf, "%s<select name=%q class=\"sigil-input\">\n", prefix, name)
	if opts, ok := c.Props["options"]; ok {
		if optSlice, ok := opts.([]interface{}); ok {
			for _, opt := range optSlice {
				if om, ok := opt.(map[string]interface{}); ok {
					val := fmt.Sprintf("%v", om["value"])
					lab := fmt.Sprintf("%v", om["label"])
					fmt.Fprintf(buf, "%s  <option value=%q>%s</option>\n", prefix, val, lab)
				}
			}
		}
	}
	fmt.Fprintf(buf, "%s</select>\n", prefix)
}

func renderPreviewOverlay(buf *bytes.Buffer, c *config.Component, prefix string, depth int, cfg *PreviewConfig) {
	title := getPropString(c.Props, "title", "")
	size := getPropString(c.Props, "size", "md")
	classes := fmt.Sprintf("sigil-%s sigil-%s-%s", c.Type, c.Type, size)

	fmt.Fprintf(buf, "%s<div class=%q>\n", prefix, classes)
	if title != "" {
		fmt.Fprintf(buf, "%s  <div class=\"sigil-%s-header\">\n", prefix, c.Type)
		fmt.Fprintf(buf, "%s    <h3 class=\"font-semibold\">%s</h3>\n", prefix, title)
		fmt.Fprintf(buf, "%s  </div>\n", prefix)
	}
	fmt.Fprintf(buf, "%s  <div class=\"sigil-%s-content\">\n", prefix, c.Type)
	for i := range c.Children {
		renderPreviewComponent(buf, &c.Children[i], depth+2, cfg)
	}
	fmt.Fprintf(buf, "%s  </div>\n", prefix)
	fmt.Fprintf(buf, "%s</div>\n", prefix)
}

func renderPreviewTabs(buf *bytes.Buffer, c *config.Component, prefix string, depth int, cfg *PreviewConfig) {
	fmt.Fprintf(buf, "%s<div>\n", prefix)
	fmt.Fprintf(buf, "%s  <div class=\"sigil-tabs-list\">\n", prefix)
	for i, child := range c.Children {
		label := getPropString(child.Props, "label", child.ID)
		activeClass := "sigil-tab-trigger"
		if i == 0 {
			activeClass += " sigil-tab-trigger-active"
		}
		fmt.Fprintf(buf, "%s    <button class=%q>%s</button>\n", prefix, activeClass, label)
	}
	fmt.Fprintf(buf, "%s  </div>\n", prefix)

	// Only render the first tab's children
	if len(c.Children) > 0 {
		fmt.Fprintf(buf, "%s  <div class=\"mt-4\">\n", prefix)
		for i := range c.Children[0].Children {
			renderPreviewComponent(buf, &c.Children[0].Children[i], depth+2, cfg)
		}
		fmt.Fprintf(buf, "%s  </div>\n", prefix)
	}
	fmt.Fprintf(buf, "%s</div>\n", prefix)
}

// Mock data generation

type colDef struct {
	Label string
	Field string
}

func generateMockRows(datasource string, columns []colDef, cfg *PreviewConfig) []map[string]string {
	rowCount := 7

	// Try to use datasource field definitions for smarter mock data
	var ds *renderer.DataSourceManifest
	if cfg != nil && cfg.DataSources != nil {
		ds = cfg.DataSources[datasource]
	}

	rows := make([]map[string]string, rowCount)
	for i := range rows {
		rows[i] = make(map[string]string)
		for _, col := range columns {
			rows[i][col.Field] = mockValue(col.Field, col.Label, i, ds)
		}
	}
	return rows
}

func mockValue(field, label string, rowIndex int, ds *renderer.DataSourceManifest) string {
	// Check datasource field definition for type hints
	if ds != nil {
		for _, f := range ds.Fields {
			if f.Name == field {
				return mockValueForType(f.Type, f.Name, f.Values, rowIndex)
			}
		}
	}

	// Infer from field name
	lower := strings.ToLower(field)
	switch {
	case lower == "id":
		return fmt.Sprintf("%d", rowIndex+1)
	case strings.Contains(lower, "name") || strings.Contains(lower, "title"):
		names := []string{"Dashboard Overview", "User Management", "Settings Panel", "Report Builder", "Data Explorer", "Theme Editor", "Component Library"}
		return names[rowIndex%len(names)]
	case strings.Contains(lower, "email"):
		emails := []string{"alice@example.com", "bob@example.com", "carol@example.com", "dave@example.com", "eve@example.com", "frank@example.com", "grace@example.com"}
		return emails[rowIndex%len(emails)]
	case strings.Contains(lower, "status"):
		statuses := []string{"active", "pending", "completed", "draft", "active", "archived", "active"}
		return statuses[rowIndex%len(statuses)]
	case strings.Contains(lower, "date") || strings.Contains(lower, "created") || strings.Contains(lower, "updated"):
		dates := []string{"2026-03-01", "2026-02-28", "2026-02-25", "2026-02-20", "2026-02-15", "2026-02-10", "2026-02-05"}
		return dates[rowIndex%len(dates)]
	case strings.Contains(lower, "count") || strings.Contains(lower, "total") || strings.Contains(lower, "amount"):
		values := []int{142, 89, 2340, 67, 1023, 456, 78}
		return fmt.Sprintf("%d", values[rowIndex%len(values)])
	case strings.Contains(lower, "role") || strings.Contains(lower, "type"):
		roles := []string{"Admin", "Editor", "Viewer", "Manager", "Developer", "Designer", "Analyst"}
		return roles[rowIndex%len(roles)]
	case strings.Contains(lower, "description"):
		descs := []string{"Primary dashboard view", "User account management", "Application configuration", "Custom report builder", "Data analysis tools", "Theme customization", "UI component browser"}
		return descs[rowIndex%len(descs)]
	default:
		return fmt.Sprintf("%s-%d", label, rowIndex+1)
	}
}

func mockValueForType(fieldType, name string, values []string, rowIndex int) string {
	if len(values) > 0 {
		return values[rowIndex%len(values)]
	}

	switch fieldType {
	case "string":
		return mockValue(name, name, rowIndex, nil)
	case "integer", "number":
		nums := []int{1, 15, 42, 108, 256, 512, 1024}
		return fmt.Sprintf("%d", nums[rowIndex%len(nums)])
	case "boolean":
		if rowIndex%2 == 0 {
			return "true"
		}
		return "false"
	case "date":
		dates := []string{"2026-03-01", "2026-02-28", "2026-02-25", "2026-02-20", "2026-02-15", "2026-02-10", "2026-02-05"}
		return dates[rowIndex%len(dates)]
	case "datetime":
		datetimes := []string{"2026-03-01 09:00", "2026-02-28 14:30", "2026-02-25 11:15", "2026-02-20 08:45", "2026-02-15 16:00", "2026-02-10 10:30", "2026-02-05 13:00"}
		return datetimes[rowIndex%len(datetimes)]
	case "email":
		emails := []string{"alice@example.com", "bob@example.com", "carol@example.com", "dave@example.com", "eve@example.com", "frank@example.com", "grace@example.com"}
		return emails[rowIndex%len(emails)]
	case "url":
		return fmt.Sprintf("https://example.com/item/%d", rowIndex+1)
	default:
		return fmt.Sprintf("%s-%d", name, rowIndex+1)
	}
}

// Helpers (duplicated from gotempl to avoid coupling)

func getPropString(props map[string]interface{}, key, defaultVal string) string {
	if v, ok := props[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return defaultVal
}

func headingClasses(level string) string {
	switch level {
	case "1":
		return "text-2xl font-bold text-[rgb(var(--sigil-text))]"
	case "2":
		return "text-xl font-semibold text-[rgb(var(--sigil-text))]"
	case "3":
		return "text-lg font-semibold text-[rgb(var(--sigil-text))]"
	case "4":
		return "text-base font-medium text-[rgb(var(--sigil-text))]"
	default:
		return "text-xl font-semibold text-[rgb(var(--sigil-text))]"
	}
}

func textClasses(variant string) string {
	switch variant {
	case "muted":
		return "text-sm text-[rgb(var(--sigil-text-muted))]"
	case "small":
		return "text-xs text-[rgb(var(--sigil-text))]"
	default:
		return "text-base text-[rgb(var(--sigil-text))]"
	}
}

func iconSymbol(name string) string {
	icons := map[string]string{
		"plus": "+", "minus": "-", "edit": "&#9998;", "delete": "&#128465;",
		"search": "&#128269;", "close": "&#10005;", "check": "&#10003;",
		"arrow-left": "&#8592;", "arrow-right": "&#8594;", "settings": "&#9881;",
		"user": "&#128100;", "home": "&#8962;", "star": "&#9733;", "heart": "&#9829;",
		"download": "&#8595;", "upload": "&#8593;", "refresh": "&#8635;",
		"filter": "&#9783;", "sort": "&#8645;", "menu": "&#9776;",
	}
	if sym, ok := icons[name]; ok {
		return sym
	}
	return "&#9679;" // bullet as fallback
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

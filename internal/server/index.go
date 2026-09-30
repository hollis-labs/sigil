package server

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hollis-labs/sigil/internal/components"
	"github.com/hollis-labs/sigil/internal/config"
	"gopkg.in/yaml.v3"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	pages := s.ListPages()
	projectName := s.projectName()
	theme := s.loadTheme()

	// Validate each page
	registry := components.NewDefaultRegistry()
	type pageRow struct {
		PageInfo
		Status     string // "pass", "warn", "fail", "error"
		ErrorCount int
		WarnCount  int
	}

	var rows []pageRow
	for _, p := range pages {
		row := pageRow{PageInfo: p}
		if p.Error != "" {
			row.Status = "error"
		} else {
			page, err := s.loadPage(p.ID)
			if err != nil {
				row.Status = "error"
			} else {
				result := config.Validate(page, registry)
				row.ErrorCount = len(result.Errors)
				row.WarnCount = len(result.Warnings)
				if !result.Valid {
					row.Status = "fail"
				} else if len(result.Warnings) > 0 {
					row.Status = "warn"
				} else {
					row.Status = "pass"
				}
			}
		}
		rows = append(rows, row)
	}

	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	buf.WriteString("  <meta charset=\"UTF-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	buf.WriteString(fmt.Sprintf("  <title>%s — Sigil Dev Server</title>\n", projectName))
	buf.WriteString("  <script src=\"https://cdn.tailwindcss.com\"></script>\n")
	buf.WriteString("  <style>\n")

	// Theme tokens
	buf.WriteString("    :root {\n")
	if theme != nil {
		for _, cat := range []string{"colors", "typography", "radius"} {
			if tokens, ok := theme.Tokens[cat]; ok {
				for k, v := range tokens {
					buf.WriteString(fmt.Sprintf("      --sigil-%s: %v;\n", k, v))
				}
			}
		}
	} else {
		buf.WriteString("      --sigil-background: 9 9 11;\n")
		buf.WriteString("      --sigil-surface: 24 24 27;\n")
		buf.WriteString("      --sigil-border: 63 63 70;\n")
		buf.WriteString("      --sigil-text: 244 244 245;\n")
		buf.WriteString("      --sigil-text-muted: 161 161 170;\n")
		buf.WriteString("      --sigil-accent: 79 70 229;\n")
	}
	buf.WriteString("    }\n")
	buf.WriteString("  </style>\n")
	buf.WriteString("</head>\n")
	buf.WriteString("<body class=\"min-h-screen\" style=\"background:rgb(var(--sigil-background));color:rgb(var(--sigil-text));font-family:system-ui,-apple-system,sans-serif\">\n")

	// Header
	buf.WriteString("  <div style=\"border-bottom:1px solid rgb(var(--sigil-border));padding:20px 32px;display:flex;align-items:center;gap:16px\">\n")
	buf.WriteString(fmt.Sprintf("    <h1 style=\"font-size:1.25rem;font-weight:700;margin:0\">%s</h1>\n", projectName))
	buf.WriteString("    <span style=\"font-size:0.75rem;padding:2px 8px;border-radius:9999px;background:rgb(var(--sigil-accent));color:white\">Sigil Dev Server</span>\n")
	buf.WriteString(fmt.Sprintf("    <span style=\"flex:1\"></span><span style=\"color:rgb(var(--sigil-text-muted));font-size:0.8rem\">%d pages</span>\n", len(rows)))
	buf.WriteString("  </div>\n")

	// Page table
	buf.WriteString("  <div style=\"padding:24px 32px\">\n")
	buf.WriteString("    <table style=\"width:100%;border-collapse:collapse\">\n")
	buf.WriteString("      <thead>\n")
	buf.WriteString("        <tr style=\"border-bottom:1px solid rgb(var(--sigil-border))\">\n")
	for _, col := range []string{"Status", "ID", "Title", "Overlay", "Module", "Components"} {
		buf.WriteString(fmt.Sprintf("          <th style=\"text-align:left;padding:8px 12px;font-size:0.7rem;text-transform:uppercase;letter-spacing:0.05em;color:rgb(var(--sigil-text-muted));font-weight:600\">%s</th>\n", col))
	}
	buf.WriteString("        </tr>\n")
	buf.WriteString("      </thead>\n")
	buf.WriteString("      <tbody>\n")

	for _, row := range rows {
		buf.WriteString("        <tr style=\"border-bottom:1px solid rgba(var(--sigil-border),0.5);transition:background 150ms\" onmouseover=\"this.style.background='rgb(var(--sigil-surface))'\" onmouseout=\"this.style.background='transparent'\">\n")

		// Status badge
		var statusHTML string
		switch row.Status {
		case "pass":
			statusHTML = `<span style="color:#22c55e;font-size:0.8rem" title="Valid">&#10003;</span>`
		case "warn":
			statusHTML = fmt.Sprintf(`<span style="color:#eab308;font-size:0.75rem" title="%d warnings">&#9888;</span>`, row.WarnCount)
		case "fail":
			statusHTML = fmt.Sprintf(`<span style="color:#ef4444;font-size:0.8rem" title="%d errors">&#10007;</span>`, row.ErrorCount)
		case "error":
			statusHTML = `<span style="color:#ef4444;font-size:0.75rem" title="Parse error">&#9888;</span>`
		}
		buf.WriteString(fmt.Sprintf("          <td style=\"padding:10px 12px;text-align:center\">%s</td>\n", statusHTML))

		// ID (linked)
		if row.Error == "" {
			buf.WriteString(fmt.Sprintf("          <td style=\"padding:10px 12px\"><a href=\"/pages/%s\" style=\"color:rgb(var(--sigil-accent));text-decoration:none;font-weight:500\">%s</a></td>\n", row.ID, row.ID))
		} else {
			buf.WriteString(fmt.Sprintf("          <td style=\"padding:10px 12px;color:rgb(var(--sigil-text-muted))\">%s</td>\n", row.ID))
		}

		buf.WriteString(fmt.Sprintf("          <td style=\"padding:10px 12px\">%s</td>\n", row.Title))
		buf.WriteString(fmt.Sprintf("          <td style=\"padding:10px 12px\"><span style=\"font-size:0.75rem;padding:1px 6px;border-radius:4px;background:rgb(var(--sigil-surface));color:rgb(var(--sigil-text-muted))\">%s</span></td>\n", row.Overlay))
		buf.WriteString(fmt.Sprintf("          <td style=\"padding:10px 12px;color:rgb(var(--sigil-text-muted))\">%s</td>\n", row.Module))
		buf.WriteString(fmt.Sprintf("          <td style=\"padding:10px 12px;color:rgb(var(--sigil-text-muted))\">%d</td>\n", row.ComponentCount))
		buf.WriteString("        </tr>\n")
	}

	buf.WriteString("      </tbody>\n")
	buf.WriteString("    </table>\n")

	if len(rows) == 0 {
		buf.WriteString("    <p style=\"color:rgb(var(--sigil-text-muted));text-align:center;padding:40px;font-size:0.9rem\">No pages found. Create one with <code style=\"background:rgb(var(--sigil-surface));padding:2px 6px;border-radius:4px\">sigil new page my-page</code></p>\n")
	}

	buf.WriteString("  </div>\n")

	// Reload script
	buf.WriteString(`<script>
const es = new EventSource('/events');
es.addEventListener('reload', () => location.reload());
es.addEventListener('error', () => setTimeout(() => location.reload(), 2000));
</script>`)
	buf.WriteString("\n</body>\n</html>\n")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}

func (s *Server) projectName() string {
	path := filepath.Join(s.SigilDir, "sigil.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return "Sigil Project"
	}
	var proj struct {
		Name string `yaml:"name"`
	}
	if yaml.Unmarshal(data, &proj) != nil || proj.Name == "" {
		return "Sigil Project"
	}
	return proj.Name
}

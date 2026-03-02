package gotempl

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/chrispian/sigil/internal/config"
	"github.com/chrispian/sigil/internal/renderer"
)

// renderPage generates a .templ file for a page config.
func renderPage(ctx *renderer.RenderContext) ([]renderer.OutputFile, error) {
	page := ctx.Page
	fileName := toSnakeCase(page.ID) + ".templ"

	var buf bytes.Buffer
	data := pageTemplateData{
		Page:     page,
		FuncName: toPascalCase(page.ID),
		GoModule: ctx.GoModule,
	}

	tmpl, err := template.New("page").Funcs(templateFuncs).Parse(pageTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing page template: %w", err)
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing page template for %q: %w", page.ID, err)
	}

	return []renderer.OutputFile{
		{Path: "pages/" + fileName, Content: buf.Bytes()},
	}, nil
}

type pageTemplateData struct {
	Page     *config.Page
	FuncName string
	GoModule string
}

var templateFuncs = template.FuncMap{
	"renderComponent": renderComponentToString,
	"indent":          indentString,
}

func renderComponentToString(c config.Component, depth int) string {
	var buf bytes.Buffer
	renderComponent(&buf, &c, depth)
	return buf.String()
}

func renderComponent(buf *bytes.Buffer, c *config.Component, depth int) {
	prefix := strings.Repeat("        ", 1) + strings.Repeat("    ", depth)

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
			renderComponent(buf, &c.Children[i], depth+1)
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
			renderComponent(buf, &c.Children[i], depth+1)
		}
		fmt.Fprintf(buf, "%s</div>\n", prefix)

	case "grid":
		cols := getPropString(c.Props, "columns", "3")
		gap := getPropString(c.Props, "gap", "4")
		classes := fmt.Sprintf("grid grid-cols-%s gap-%s", cols, gap)
		fmt.Fprintf(buf, "%s<div class=%q>\n", prefix, classes)
		for i := range c.Children {
			renderComponent(buf, &c.Children[i], depth+1)
		}
		fmt.Fprintf(buf, "%s</div>\n", prefix)

	case "heading":
		level := getPropString(c.Props, "level", "2")
		text := getPropString(c.Props, "text", "")
		tag := "h" + level
		classes := headingClasses(level)
		fmt.Fprintf(buf, "%s<%s class=%q>%s</%s>\n", prefix, tag, classes, text, tag)

	case "text":
		content := getPropString(c.Props, "content", "")
		variant := getPropString(c.Props, "variant", "body")
		classes := textClasses(variant)
		fmt.Fprintf(buf, "%s<p class=%q>%s</p>\n", prefix, classes, content)

	case "button":
		label := getPropString(c.Props, "label", "Button")
		variant := getPropString(c.Props, "variant", "primary")
		icon := getPropString(c.Props, "icon", "")
		classes := fmt.Sprintf("sigil-btn sigil-btn-%s", variant)
		htmx := renderHTMXAttrs(c.Actions)
		iconHTML := ""
		if icon != "" {
			iconHTML = fmt.Sprintf(`<span class="icon icon-%s"></span> `, icon)
		}
		fmt.Fprintf(buf, "%s<button class=%q%s>%s<span>%s</span></button>\n", prefix, classes, htmx, iconHTML, label)

	case "search-bar":
		placeholder := getPropString(c.Props, "placeholder", "Search...")
		datasource := getPropString(c.Props, "datasource", "")
		debounce := getPropString(c.Props, "debounce", "300")
		hxTarget := ""
		if datasource != "" {
			hxTarget = fmt.Sprintf(` hx-get="/api/%s" hx-trigger="input changed delay:%sms" hx-target="#table-%s"`,
				strings.ToLower(datasource)+"s", debounce, strings.ToLower(datasource)+"s")
		}
		fmt.Fprintf(buf, "%s<input type=\"search\" class=\"sigil-input w-full\" placeholder=%q%s/>\n", prefix, placeholder, hxTarget)

	case "data-table":
		renderDataTable(buf, c, prefix)

	case "form":
		renderForm(buf, c, prefix, depth)

	case "badge":
		value := getPropString(c.Props, "value", "")
		variant := getPropString(c.Props, "variant", "default")
		classes := fmt.Sprintf("sigil-badge sigil-badge-%s", variant)
		fmt.Fprintf(buf, "%s<span class=%q>%s</span>\n", prefix, classes, value)

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
		renderSelect(buf, c, prefix)

	case "modal", "sheet":
		renderOverlayComponent(buf, c, prefix, depth)

	case "separator":
		fmt.Fprintf(buf, "%s<hr class=\"sigil-separator\"/>\n", prefix)

	case "progress":
		value := getPropString(c.Props, "value", "0")
		max := getPropString(c.Props, "max", "100")
		fmt.Fprintf(buf, "%s<progress class=\"sigil-progress\" value=%q max=%q></progress>\n", prefix, value, max)

	case "avatar":
		src := getPropString(c.Props, "src", "")
		alt := getPropString(c.Props, "alt", "")
		size := getPropString(c.Props, "size", "md")
		classes := fmt.Sprintf("sigil-avatar sigil-avatar-%s", size)
		fmt.Fprintf(buf, "%s<img class=%q src=%q alt=%q/>\n", prefix, classes, src, alt)

	case "alert":
		message := getPropString(c.Props, "message", "")
		variant := getPropString(c.Props, "variant", "info")
		classes := fmt.Sprintf("sigil-alert sigil-alert-%s", variant)
		fmt.Fprintf(buf, "%s<div class=%q>%s</div>\n", prefix, classes, message)

	case "icon":
		name := getPropString(c.Props, "name", "")
		fmt.Fprintf(buf, "%s<span class=\"icon icon-%s\"></span>\n", prefix, name)

	case "label":
		text := getPropString(c.Props, "text", "")
		fmt.Fprintf(buf, "%s<label class=\"sigil-label\">%s</label>\n", prefix, text)

	case "spacer":
		size := getPropString(c.Props, "size", "4")
		fmt.Fprintf(buf, "%s<div class=\"h-%s\"></div>\n", prefix, size)

	default:
		// Generic component — render as a div with data attributes
		idAttr := ""
		if c.ID != "" {
			idAttr = fmt.Sprintf(` id=%q`, c.ID)
		}
		htmx := renderHTMXAttrs(c.Actions)
		fmt.Fprintf(buf, "%s<div data-component=%q%s%s>\n", prefix, c.Type, idAttr, htmx)
		for i := range c.Children {
			renderComponent(buf, &c.Children[i], depth+1)
		}
		fmt.Fprintf(buf, "%s</div>\n", prefix)
	}
}

func renderDataTable(buf *bytes.Buffer, c *config.Component, prefix string) {
	datasource := getPropString(c.Props, "datasource", "")
	dsLower := strings.ToLower(datasource) + "s"
	tableID := "table-" + dsLower

	htmx := renderHTMXAttrs(c.Actions)
	fmt.Fprintf(buf, "%s<div id=%q class=\"sigil-data-table\"%s>\n", prefix, tableID, htmx)
	fmt.Fprintf(buf, "%s    <table class=\"w-full\">\n", prefix)
	fmt.Fprintf(buf, "%s        <thead>\n", prefix)
	fmt.Fprintf(buf, "%s            <tr>\n", prefix)

	// Extract columns from props
	if cols, ok := c.Props["columns"]; ok {
		if colSlice, ok := cols.([]interface{}); ok {
			for _, col := range colSlice {
				if colMap, ok := col.(map[string]interface{}); ok {
					label := fmt.Sprintf("%v", colMap["label"])
					field := fmt.Sprintf("%v", colMap["field"])
					sortable := false
					if s, ok := colMap["sortable"]; ok {
						sortable, _ = s.(bool)
					}
					sortAttr := ""
					if sortable {
						sortAttr = fmt.Sprintf(` hx-get="/api/%s?sort=%s" hx-target="#%s" class="cursor-pointer"`, dsLower, field, tableID)
					}
					fmt.Fprintf(buf, "%s                <th%s>%s</th>\n", prefix, sortAttr, label)
				}
			}
		}
	}

	fmt.Fprintf(buf, "%s            </tr>\n", prefix)
	fmt.Fprintf(buf, "%s        </thead>\n", prefix)
	fmt.Fprintf(buf, "%s        <tbody hx-get=\"/api/%s\" hx-trigger=\"load\" hx-target=\"this\">\n", prefix, dsLower)
	fmt.Fprintf(buf, "%s            <!-- Rows loaded via HTMX -->\n", prefix)
	fmt.Fprintf(buf, "%s        </tbody>\n", prefix)
	fmt.Fprintf(buf, "%s    </table>\n", prefix)
	fmt.Fprintf(buf, "%s</div>\n", prefix)
}

func renderForm(buf *bytes.Buffer, c *config.Component, prefix string, depth int) {
	htmx := renderHTMXAttrs(c.Actions)
	fmt.Fprintf(buf, "%s<form class=\"sigil-form flex flex-col gap-4\"%s>\n", prefix, htmx)

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
					required := false
					if r, ok := fm["required"]; ok {
						required, _ = r.(bool)
					}

					fmt.Fprintf(buf, "%s    <div class=\"flex flex-col gap-1\">\n", prefix)
					fmt.Fprintf(buf, "%s        <label class=\"sigil-label\">%s</label>\n", prefix, label)

					reqAttr := ""
					if required {
						reqAttr = " required"
					}

					switch fieldType {
					case "textarea":
						rows := "3"
						if r, ok := fm["rows"]; ok {
							rows = fmt.Sprintf("%v", r)
						}
						fmt.Fprintf(buf, "%s        <textarea name=%q class=\"sigil-input\" placeholder=%q rows=%q%s></textarea>\n",
							prefix, name, placeholder, rows, reqAttr)
					case "select":
						fmt.Fprintf(buf, "%s        <select name=%q class=\"sigil-input\"%s>\n", prefix, name, reqAttr)
						if opts, ok := fm["options"]; ok {
							if optSlice, ok := opts.([]interface{}); ok {
								for _, opt := range optSlice {
									if om, ok := opt.(map[string]interface{}); ok {
										val := fmt.Sprintf("%v", om["value"])
										lab := fmt.Sprintf("%v", om["label"])
										fmt.Fprintf(buf, "%s            <option value=%q>%s</option>\n", prefix, val, lab)
									}
								}
							}
						}
						fmt.Fprintf(buf, "%s        </select>\n", prefix)
					case "checkbox":
						fmt.Fprintf(buf, "%s        <input type=\"checkbox\" name=%q class=\"sigil-checkbox\"/>\n", prefix, name)
					default:
						fmt.Fprintf(buf, "%s        <input type=%q name=%q class=\"sigil-input\" placeholder=%q%s/>\n",
							prefix, fieldType, name, placeholder, reqAttr)
					}

					fmt.Fprintf(buf, "%s    </div>\n", prefix)
				}
			}
		}
	}

	// Submit button
	submitLabel := "Submit"
	if submit, ok := c.Props["submit"]; ok {
		if sm, ok := submit.(map[string]interface{}); ok {
			if l, ok := sm["label"]; ok {
				submitLabel = fmt.Sprintf("%v", l)
			}
		}
	}
	fmt.Fprintf(buf, "%s    <button type=\"submit\" class=\"sigil-btn sigil-btn-primary\">%s</button>\n", prefix, submitLabel)
	fmt.Fprintf(buf, "%s</form>\n", prefix)
}

func renderSelect(buf *bytes.Buffer, c *config.Component, prefix string) {
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
					fmt.Fprintf(buf, "%s    <option value=%q>%s</option>\n", prefix, val, lab)
				}
			}
		}
	}
	fmt.Fprintf(buf, "%s</select>\n", prefix)
}

func renderOverlayComponent(buf *bytes.Buffer, c *config.Component, prefix string, depth int) {
	title := getPropString(c.Props, "title", "")
	size := getPropString(c.Props, "size", "md")
	classes := fmt.Sprintf("sigil-%s sigil-%s-%s", c.Type, c.Type, size)

	fmt.Fprintf(buf, "%s<div class=%q>\n", prefix, classes)
	if title != "" {
		fmt.Fprintf(buf, "%s    <div class=\"sigil-%s-header\">\n", prefix, c.Type)
		fmt.Fprintf(buf, "%s        <h3>%s</h3>\n", prefix, title)
		fmt.Fprintf(buf, "%s    </div>\n", prefix)
	}
	fmt.Fprintf(buf, "%s    <div class=\"sigil-%s-content\">\n", prefix, c.Type)
	for i := range c.Children {
		renderComponent(buf, &c.Children[i], depth+2)
	}
	fmt.Fprintf(buf, "%s    </div>\n", prefix)
	fmt.Fprintf(buf, "%s</div>\n", prefix)
}

func renderHTMXAttrs(actions map[string]config.Action) string {
	if len(actions) == 0 {
		return ""
	}

	var attrs []string
	for _, action := range actions {
		switch action.Type {
		case "navigate":
			if action.URL != "" {
				attrs = append(attrs, fmt.Sprintf(` hx-get=%q hx-push-url="true"`, action.URL))
			} else if action.Page != "" {
				attrs = append(attrs, fmt.Sprintf(` hx-get="/ui/pages/%s" hx-push-url="true"`, action.Page))
			}
		case "modal":
			target := "#modal-container"
			if action.Target != "" {
				target = action.Target
			}
			attrs = append(attrs, fmt.Sprintf(` hx-get="/ui/modals/%s" hx-target=%q hx-swap="innerHTML"`,
				strings.ReplaceAll(strings.ToLower(action.Title), " ", "-"), target))
		case "sheet":
			page := action.Page
			if page == "" {
				page = "sheet"
			}
			attrs = append(attrs, fmt.Sprintf(` hx-get="/ui/sheets/%s" hx-target="#sheet-container" hx-swap="innerHTML"`, page))
		case "http":
			method := strings.ToLower(action.Method)
			switch method {
			case "post":
				attrs = append(attrs, fmt.Sprintf(` hx-post=%q`, action.URL))
			case "put":
				attrs = append(attrs, fmt.Sprintf(` hx-put=%q`, action.URL))
			case "delete":
				attrs = append(attrs, fmt.Sprintf(` hx-delete=%q`, action.URL))
			default:
				attrs = append(attrs, fmt.Sprintf(` hx-get=%q`, action.URL))
			}
			if action.Target != "" {
				attrs = append(attrs, fmt.Sprintf(` hx-target=%q`, action.Target))
			}
		case "emit":
			if action.Event != "" {
				attrs = append(attrs, fmt.Sprintf(` hx-trigger=%q`, action.Event))
			}
		case "confirm":
			if action.Message != "" {
				attrs = append(attrs, fmt.Sprintf(` hx-confirm=%q`, action.Message))
			}
			if action.OnConfirm != nil {
				inner := renderHTMXAttrs(map[string]config.Action{"_": *action.OnConfirm})
				attrs = append(attrs, inner)
			}
		case "close":
			attrs = append(attrs, ` onclick="this.closest('[data-modal]')?.remove()"`)
		}
	}
	return strings.Join(attrs, "")
}

// Helper functions

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

func indentString(s string, n int) string {
	prefix := strings.Repeat("    ", n)
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

// toSnakeCase converts "sprint-dashboard" to "sprint_dashboard".
func toSnakeCase(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), "-", "_")
}

// toPascalCase converts "sprint-dashboard" to "SprintDashboard".
func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

const pageTemplate = `// {{ .Page.ID }}.templ
// Generated by Sigil — do not edit manually
package pages

templ {{ .FuncName }}() {
    // Page: {{ .Page.Title }}
    // Overlay: {{ .Page.Overlay }}
{{ renderComponent .Page.Layout 0 }}}
`

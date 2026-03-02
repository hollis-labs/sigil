package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrispian/sigil/internal/components"
	"github.com/chrispian/sigil/internal/config"
)

// RegisterAllTools registers all Sigil MCP tools.
func RegisterAllTools(s *Server) {
	s.RegisterTool(toolListPages(s))
	s.RegisterTool(toolGetPage(s))
	s.RegisterTool(toolCreatePage(s))
	s.RegisterTool(toolUpdatePage(s))
	s.RegisterTool(toolValidate(s))
	s.RegisterTool(toolListComponents(s))
	s.RegisterTool(toolGetComponentSchema(s))
	s.RegisterTool(toolListDataSources(s))
	s.RegisterTool(toolCreateDataSource(s))
}

func toolListPages(s *Server) Tool {
	return Tool{
		Name:        "sigil_list_pages",
		Description: "List all page configs in the project",
		InputSchema: json.RawMessage(`{"type":"object","properties":{}}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			pagesDir := filepath.Join(s.SigilDir(), "pages")
			matches, err := filepath.Glob(filepath.Join(pagesDir, "*.yaml"))
			if err != nil {
				return nil, err
			}

			var pages []map[string]interface{}
			for _, path := range matches {
				page, err := config.ParseFile(path)
				if err != nil {
					pages = append(pages, map[string]interface{}{
						"path":  path,
						"error": err.Error(),
					})
					continue
				}
				pages = append(pages, map[string]interface{}{
					"id":      page.ID,
					"title":   page.Title,
					"overlay": page.Overlay,
					"module":  page.Module,
					"path":    path,
				})
			}

			data, _ := json.MarshalIndent(pages, "", "  ")
			return textResult(string(data)), nil
		},
	}
}

func toolGetPage(s *Server) Tool {
	return Tool{
		Name:        "sigil_get_page",
		Description: "Get the full config for a page",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"id":{"type":"string","description":"Page ID"}},"required":["id"]}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			var args struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(params, &args); err != nil {
				return nil, err
			}

			path := filepath.Join(s.SigilDir(), "pages", args.ID+".yaml")
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("page %q not found: %w", args.ID, err)
			}

			return textResult(string(data)), nil
		},
	}
}

func toolCreatePage(s *Server) Tool {
	return Tool{
		Name:        "sigil_create_page",
		Description: "Create a new page config (validates before saving)",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"id":{"type":"string","description":"Page ID"},"config":{"type":"string","description":"Page YAML content"}},"required":["id","config"]}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			var args struct {
				ID     string `json:"id"`
				Config string `json:"config"`
			}
			if err := json.Unmarshal(params, &args); err != nil {
				return nil, err
			}

			// Parse and validate
			page, err := config.Parse([]byte(args.Config))
			if err != nil {
				return errorResult(fmt.Sprintf("parse error: %s", err)), nil
			}

			registry := components.NewDefaultRegistry()
			result := config.Validate(page, registry)

			if !result.Valid {
				var errs []string
				for _, e := range result.Errors {
					errs = append(errs, e.String())
				}
				resp := map[string]interface{}{
					"valid":  false,
					"errors": errs,
				}
				data, _ := json.MarshalIndent(resp, "", "  ")
				return textResult(string(data)), nil
			}

			// Write file
			path := filepath.Join(s.SigilDir(), "pages", args.ID+".yaml")
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return nil, err
			}
			if err := os.WriteFile(path, []byte(args.Config), 0644); err != nil {
				return nil, err
			}

			var warnings []string
			for _, w := range result.Warnings {
				warnings = append(warnings, w.String())
			}

			resp := map[string]interface{}{
				"path":     path,
				"valid":    true,
				"warnings": warnings,
			}
			data, _ := json.MarshalIndent(resp, "", "  ")
			return textResult(string(data)), nil
		},
	}
}

func toolUpdatePage(s *Server) Tool {
	return Tool{
		Name:        "sigil_update_page",
		Description: "Update an existing page config",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"id":{"type":"string","description":"Page ID"},"config":{"type":"string","description":"Page YAML content"}},"required":["id","config"]}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			var args struct {
				ID     string `json:"id"`
				Config string `json:"config"`
			}
			if err := json.Unmarshal(params, &args); err != nil {
				return nil, err
			}

			path := filepath.Join(s.SigilDir(), "pages", args.ID+".yaml")
			if _, err := os.Stat(path); err != nil {
				return nil, fmt.Errorf("page %q not found", args.ID)
			}

			// Parse and validate
			page, err := config.Parse([]byte(args.Config))
			if err != nil {
				return errorResult(fmt.Sprintf("parse error: %s", err)), nil
			}

			registry := components.NewDefaultRegistry()
			result := config.Validate(page, registry)

			if !result.Valid {
				var errs []string
				for _, e := range result.Errors {
					errs = append(errs, e.String())
				}
				resp := map[string]interface{}{
					"valid":  false,
					"errors": errs,
				}
				data, _ := json.MarshalIndent(resp, "", "  ")
				return textResult(string(data)), nil
			}

			if err := os.WriteFile(path, []byte(args.Config), 0644); err != nil {
				return nil, err
			}

			resp := map[string]interface{}{
				"path":  path,
				"valid": true,
			}
			data, _ := json.MarshalIndent(resp, "", "  ")
			return textResult(string(data)), nil
		},
	}
}

func toolValidate(s *Server) Tool {
	return Tool{
		Name:        "sigil_validate",
		Description: "Validate a Sigil config (YAML string or file path)",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"config":{"type":"string","description":"YAML content to validate"},"path":{"type":"string","description":"Path to YAML file to validate"}}}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			var args struct {
				Config string `json:"config"`
				Path   string `json:"path"`
			}
			if err := json.Unmarshal(params, &args); err != nil {
				return nil, err
			}

			var yamlData []byte
			if args.Config != "" {
				yamlData = []byte(args.Config)
			} else if args.Path != "" {
				data, err := os.ReadFile(args.Path)
				if err != nil {
					return nil, err
				}
				yamlData = data
			} else {
				return nil, fmt.Errorf("either 'config' or 'path' is required")
			}

			page, err := config.Parse(yamlData)
			if err != nil {
				resp := map[string]interface{}{
					"valid":  false,
					"errors": []string{fmt.Sprintf("parse error: %s", err)},
				}
				data, _ := json.MarshalIndent(resp, "", "  ")
				return textResult(string(data)), nil
			}

			registry := components.NewDefaultRegistry()
			result := config.Validate(page, registry)

			var errors []string
			for _, e := range result.Errors {
				errors = append(errors, e.String())
			}
			var warnings []string
			for _, w := range result.Warnings {
				warnings = append(warnings, w.String())
			}

			resp := map[string]interface{}{
				"valid":    result.Valid,
				"errors":   errors,
				"warnings": warnings,
			}
			data, _ := json.MarshalIndent(resp, "", "  ")
			return textResult(string(data)), nil
		},
	}
}

func toolListComponents(s *Server) Tool {
	return Tool{
		Name:        "sigil_list_components",
		Description: "List available component types with their schemas",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"category":{"type":"string","description":"Filter by category (primitives, layouts, navigation, composites, data, forms)"}}}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			var args struct {
				Category string `json:"category"`
			}
			json.Unmarshal(params, &args)

			registry := components.NewDefaultRegistry()
			types := registry.Types()

			var result []map[string]interface{}
			for _, t := range types {
				schema, ok := registry.Get(t)
				if !ok {
					continue
				}
				if args.Category != "" && schema.Category != args.Category {
					continue
				}

				comp := map[string]interface{}{
					"type":        schema.Type,
					"category":    schema.Category,
					"description": schema.Description,
				}

				if len(schema.Props) > 0 {
					var props []map[string]interface{}
					for name, p := range schema.Props {
						prop := map[string]interface{}{
							"name": name,
							"type": p.Type,
						}
						if p.Required {
							prop["required"] = true
						}
						if p.Description != "" {
							prop["description"] = p.Description
						}
						if len(p.Enum) > 0 {
							prop["enum"] = p.Enum
						}
						if p.Default != nil {
							prop["default"] = p.Default
						}
						props = append(props, prop)
					}
					comp["props"] = props
				}

				result = append(result, comp)
			}

			data, _ := json.MarshalIndent(result, "", "  ")
			return textResult(string(data)), nil
		},
	}
}

func toolGetComponentSchema(s *Server) Tool {
	return Tool{
		Name:        "sigil_get_component_schema",
		Description: "Get the full schema for a component type",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"type":{"type":"string","description":"Component type name"}},"required":["type"]}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			var args struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal(params, &args); err != nil {
				return nil, err
			}

			registry := components.NewDefaultRegistry()
			schema, ok := registry.Get(args.Type)
			if !ok {
				return nil, fmt.Errorf("unknown component type %q", args.Type)
			}

			result := map[string]interface{}{
				"type":        schema.Type,
				"category":    schema.Category,
				"description": schema.Description,
			}
			if len(schema.Props) > 0 {
				result["props"] = schema.Props
			}
			if len(schema.Actions) > 0 {
				result["actions"] = schema.Actions
			}
			if len(schema.Slots) > 0 {
				result["slots"] = schema.Slots
			}
			if len(schema.Shortcuts) > 0 {
				result["shortcuts"] = schema.Shortcuts
			}

			data, _ := json.MarshalIndent(result, "", "  ")
			return textResult(string(data)), nil
		},
	}
}

func toolListDataSources(s *Server) Tool {
	return Tool{
		Name:        "sigil_list_datasources",
		Description: "List registered datasource manifests",
		InputSchema: json.RawMessage(`{"type":"object","properties":{}}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			dsDir := filepath.Join(s.SigilDir(), "datasources")
			matches, err := filepath.Glob(filepath.Join(dsDir, "*.yaml"))
			if err != nil {
				return textResult("[]"), nil
			}

			var result []map[string]interface{}
			for _, path := range matches {
				name := strings.TrimSuffix(filepath.Base(path), ".yaml")
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				// Quick parse for summary
				var ds struct {
					Alias        string   `yaml:"alias"`
					Description  string   `yaml:"description"`
					Capabilities []string `yaml:"capabilities"`
					Fields       []struct {
						Name string `yaml:"name"`
					} `yaml:"fields"`
				}
				if err := parseYAML(data, &ds); err != nil {
					continue
				}
				alias := ds.Alias
				if alias == "" {
					alias = name
				}
				result = append(result, map[string]interface{}{
					"alias":        alias,
					"description":  ds.Description,
					"capabilities": ds.Capabilities,
					"fieldCount":   len(ds.Fields),
				})
			}

			out, _ := json.MarshalIndent(result, "", "  ")
			return textResult(string(out)), nil
		},
	}
}

func toolCreateDataSource(s *Server) Tool {
	return Tool{
		Name:        "sigil_create_datasource",
		Description: "Create a new datasource manifest",
		InputSchema: json.RawMessage(`{"type":"object","properties":{"alias":{"type":"string","description":"DataSource alias"},"config":{"type":"string","description":"DataSource YAML content"}},"required":["alias","config"]}`),
		Handler: func(params json.RawMessage) (*ToolResult, error) {
			var args struct {
				Alias  string `json:"alias"`
				Config string `json:"config"`
			}
			if err := json.Unmarshal(params, &args); err != nil {
				return nil, err
			}

			dsDir := filepath.Join(s.SigilDir(), "datasources")
			if err := os.MkdirAll(dsDir, 0755); err != nil {
				return nil, err
			}

			path := filepath.Join(dsDir, strings.ToLower(args.Alias)+".yaml")
			if err := os.WriteFile(path, []byte(args.Config), 0644); err != nil {
				return nil, err
			}

			resp := map[string]interface{}{"path": path}
			data, _ := json.MarshalIndent(resp, "", "  ")
			return textResult(string(data)), nil
		},
	}
}

// Helpers

func textResult(text string) *ToolResult {
	return &ToolResult{
		Content: []ContentBlock{{Type: "text", Text: text}},
	}
}

func errorResult(text string) *ToolResult {
	return &ToolResult{
		Content: []ContentBlock{{Type: "text", Text: text}},
		IsError: true,
	}
}

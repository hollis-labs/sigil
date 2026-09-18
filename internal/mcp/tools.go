package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gomcpserver "github.com/hollis-labs/go-mcp/server"

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
		Name:           "sigil_list_pages",
		Description:    "List all page configs in the project",
		InputSchema:    gomcpserver.EmptyObjectSchema(),
		ReadOnlyHint:   true,
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
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

			return pages, nil
		},
	}
}

func toolGetPage(s *Server) Tool {
	return Tool{
		Name:        "sigil_get_page",
		Description: "Get the full config for a page",
		InputSchema: gomcpserver.ObjectSchema(map[string]interface{}{
			"id": map[string]interface{}{"type": "string", "description": "Page ID"},
		}, "id"),
		ReadOnlyHint:   true,
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			id, _ := args["id"].(string)

			path := filepath.Join(s.SigilDir(), "pages", id+".yaml")
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("page %q not found: %w", id, err)
			}

			return string(data), nil
		},
	}
}

func toolCreatePage(s *Server) Tool {
	return Tool{
		Name:        "sigil_create_page",
		Description: "Create a new page config (validates before saving)",
		InputSchema: gomcpserver.ObjectSchema(map[string]interface{}{
			"id":     map[string]interface{}{"type": "string", "description": "Page ID"},
			"config": map[string]interface{}{"type": "string", "description": "Page YAML content"},
		}, "id", "config"),
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			id, _ := args["id"].(string)
			cfg, _ := args["config"].(string)

			// Parse and validate
			page, err := config.Parse([]byte(cfg))
			if err != nil {
				return nil, fmt.Errorf("parse error: %w", err)
			}

			registry := components.NewDefaultRegistry()
			result := config.Validate(page, registry)

			if !result.Valid {
				var errs []string
				for _, e := range result.Errors {
					errs = append(errs, e.String())
				}
				return map[string]interface{}{
					"valid":  false,
					"errors": errs,
				}, nil
			}

			// Write file
			path := filepath.Join(s.SigilDir(), "pages", id+".yaml")
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return nil, err
			}
			if err := os.WriteFile(path, []byte(cfg), 0600); err != nil {
				return nil, err
			}

			var warnings []string
			for _, w := range result.Warnings {
				warnings = append(warnings, w.String())
			}

			return map[string]interface{}{
				"path":     path,
				"valid":    true,
				"warnings": warnings,
			}, nil
		},
	}
}

func toolUpdatePage(s *Server) Tool {
	return Tool{
		Name:        "sigil_update_page",
		Description: "Update an existing page config",
		InputSchema: gomcpserver.ObjectSchema(map[string]interface{}{
			"id":     map[string]interface{}{"type": "string", "description": "Page ID"},
			"config": map[string]interface{}{"type": "string", "description": "Page YAML content"},
		}, "id", "config"),
		DestructiveHint: true,
		IdempotentHint:  true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			id, _ := args["id"].(string)
			cfg, _ := args["config"].(string)

			path := filepath.Join(s.SigilDir(), "pages", id+".yaml")
			if _, err := os.Stat(path); err != nil {
				return nil, fmt.Errorf("page %q not found", id)
			}

			// Parse and validate
			page, err := config.Parse([]byte(cfg))
			if err != nil {
				return nil, fmt.Errorf("parse error: %w", err)
			}

			registry := components.NewDefaultRegistry()
			result := config.Validate(page, registry)

			if !result.Valid {
				var errs []string
				for _, e := range result.Errors {
					errs = append(errs, e.String())
				}
				return map[string]interface{}{
					"valid":  false,
					"errors": errs,
				}, nil
			}

			if err := os.WriteFile(path, []byte(cfg), 0600); err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"path":  path,
				"valid": true,
			}, nil
		},
	}
}

func toolValidate(s *Server) Tool {
	return Tool{
		Name:        "sigil_validate",
		Description: "Validate a Sigil config (YAML string or file path)",
		InputSchema: gomcpserver.ObjectSchema(map[string]interface{}{
			"config": map[string]interface{}{"type": "string", "description": "YAML content to validate"},
			"path":   map[string]interface{}{"type": "string", "description": "Path to YAML file to validate"},
		}),
		ReadOnlyHint:   true,
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			cfg, _ := args["config"].(string)
			path, _ := args["path"].(string)

			var yamlData []byte
			switch {
			case cfg != "":
				yamlData = []byte(cfg)
			case path != "":
				data, err := os.ReadFile(path) //nolint:gosec // sigil_validate's whole purpose is validating a caller-specified file
				if err != nil {
					return nil, err
				}
				yamlData = data
			default:
				return nil, fmt.Errorf("either 'config' or 'path' is required")
			}

			page, err := config.Parse(yamlData)
			if err != nil {
				return map[string]interface{}{
					"valid":  false,
					"errors": []string{fmt.Sprintf("parse error: %s", err)},
				}, nil
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

			return map[string]interface{}{
				"valid":    result.Valid,
				"errors":   errors,
				"warnings": warnings,
			}, nil
		},
	}
}

func toolListComponents(s *Server) Tool {
	return Tool{
		Name:        "sigil_list_components",
		Description: "List available component types with their schemas",
		InputSchema: gomcpserver.ObjectSchema(map[string]interface{}{
			"category": map[string]interface{}{"type": "string", "description": "Filter by category (primitives, layouts, navigation, composites, data, forms)"},
		}),
		ReadOnlyHint:   true,
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			category, _ := args["category"].(string)

			registry := components.NewDefaultRegistry()
			types := registry.Types()

			var result []map[string]interface{}
			for _, t := range types {
				schema, ok := registry.Get(t)
				if !ok {
					continue
				}
				if category != "" && schema.Category != category {
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

			return result, nil
		},
	}
}

func toolGetComponentSchema(s *Server) Tool {
	return Tool{
		Name:        "sigil_get_component_schema",
		Description: "Get the full schema for a component type",
		InputSchema: gomcpserver.ObjectSchema(map[string]interface{}{
			"type": map[string]interface{}{"type": "string", "description": "Component type name"},
		}, "type"),
		ReadOnlyHint:   true,
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			typeName, _ := args["type"].(string)

			registry := components.NewDefaultRegistry()
			schema, ok := registry.Get(typeName)
			if !ok {
				return nil, fmt.Errorf("unknown component type %q", typeName)
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

			return result, nil
		},
	}
}

func toolListDataSources(s *Server) Tool {
	return Tool{
		Name:           "sigil_list_datasources",
		Description:    "List registered datasource manifests",
		InputSchema:    gomcpserver.EmptyObjectSchema(),
		ReadOnlyHint:   true,
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			dsDir := filepath.Join(s.SigilDir(), "datasources")
			matches, err := filepath.Glob(filepath.Join(dsDir, "*.yaml"))
			if err != nil {
				return nil, err
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

			return result, nil
		},
	}
}

func toolCreateDataSource(s *Server) Tool {
	return Tool{
		Name:        "sigil_create_datasource",
		Description: "Create a new datasource manifest",
		InputSchema: gomcpserver.ObjectSchema(map[string]interface{}{
			"alias":  map[string]interface{}{"type": "string", "description": "DataSource alias"},
			"config": map[string]interface{}{"type": "string", "description": "DataSource YAML content"},
		}, "alias", "config"),
		IdempotentHint: true,
		Handler: func(ctx context.Context, args map[string]any) (any, error) {
			alias, _ := args["alias"].(string)
			cfg, _ := args["config"].(string)

			dsDir := filepath.Join(s.SigilDir(), "datasources")
			if err := os.MkdirAll(dsDir, 0755); err != nil {
				return nil, err
			}

			path := filepath.Join(dsDir, strings.ToLower(alias)+".yaml")
			if err := os.WriteFile(path, []byte(cfg), 0600); err != nil {
				return nil, err
			}

			return map[string]interface{}{"path": path}, nil
		},
	}
}

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

// RegisterAllResources registers MCP resources and the resource read handler.
func RegisterAllResources(s *Server) {
	s.SetResources([]Resource{
		{URI: "sigil://pages", Name: "Page listing", Description: "List of all page configs", MimeType: "application/json"},
		{URI: "sigil://pages/{id}", Name: "Page config", Description: "Full content of a specific page config", MimeType: "text/yaml"},
		{URI: "sigil://components", Name: "Component schemas", Description: "All available component types and schemas", MimeType: "application/json"},
		{URI: "sigil://components/{type}", Name: "Component schema", Description: "Schema for a specific component type", MimeType: "application/json"},
		{URI: "sigil://datasources", Name: "DataSource manifests", Description: "All datasource manifests", MimeType: "application/json"},
		{URI: "sigil://themes", Name: "Themes", Description: "All theme definitions", MimeType: "application/json"},
		{URI: "sigil://project", Name: "Project config", Description: "Project configuration", MimeType: "text/yaml"},
	})

	s.SetResourceHandler(func(uri string) (*ResourceContent, error) {
		return handleResourceRead(s, uri)
	})
}

func handleResourceRead(s *Server, uri string) (*ResourceContent, error) {
	switch {
	case uri == "sigil://pages":
		return readPagesResource(s)
	case strings.HasPrefix(uri, "sigil://pages/"):
		id := strings.TrimPrefix(uri, "sigil://pages/")
		return readPageResource(s, id)
	case uri == "sigil://components":
		return readComponentsResource()
	case strings.HasPrefix(uri, "sigil://components/"):
		typeName := strings.TrimPrefix(uri, "sigil://components/")
		return readComponentResource(typeName)
	case uri == "sigil://datasources":
		return readDataSourcesResource(s)
	case uri == "sigil://themes":
		return readThemesResource(s)
	case uri == "sigil://project":
		return readProjectResource(s)
	default:
		return nil, fmt.Errorf("unknown resource URI: %s", uri)
	}
}

func readPagesResource(s *Server) (*ResourceContent, error) {
	pagesDir := filepath.Join(s.SigilDir(), "pages")
	matches, _ := filepath.Glob(filepath.Join(pagesDir, "*.yaml"))

	var pages []map[string]interface{}
	for _, path := range matches {
		page, err := config.ParseFile(path)
		if err != nil {
			continue
		}
		pages = append(pages, map[string]interface{}{
			"id":      page.ID,
			"title":   page.Title,
			"overlay": page.Overlay,
			"module":  page.Module,
		})
	}

	data, _ := json.MarshalIndent(pages, "", "  ")
	return &ResourceContent{
		URI:      "sigil://pages",
		MimeType: "application/json",
		Text:     string(data),
	}, nil
}

func readPageResource(s *Server, id string) (*ResourceContent, error) {
	path := filepath.Join(s.SigilDir(), "pages", id+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("page %q not found", id)
	}
	return &ResourceContent{
		URI:      "sigil://pages/" + id,
		MimeType: "text/yaml",
		Text:     string(data),
	}, nil
}

func readComponentsResource() (*ResourceContent, error) {
	registry := components.NewDefaultRegistry()
	types := registry.Types()

	var result []map[string]interface{}
	for _, t := range types {
		schema, ok := registry.Get(t)
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"type":        schema.Type,
			"category":    schema.Category,
			"description": schema.Description,
		})
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return &ResourceContent{
		URI:      "sigil://components",
		MimeType: "application/json",
		Text:     string(data),
	}, nil
}

func readComponentResource(typeName string) (*ResourceContent, error) {
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

	data, _ := json.MarshalIndent(result, "", "  ")
	return &ResourceContent{
		URI:      "sigil://components/" + typeName,
		MimeType: "application/json",
		Text:     string(data),
	}, nil
}

func readDataSourcesResource(s *Server) (*ResourceContent, error) {
	dsDir := filepath.Join(s.SigilDir(), "datasources")
	matches, _ := filepath.Glob(filepath.Join(dsDir, "*.yaml"))

	var result []map[string]interface{}
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"alias": strings.TrimSuffix(filepath.Base(path), ".yaml"),
			"content": string(data),
		})
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return &ResourceContent{
		URI:      "sigil://datasources",
		MimeType: "application/json",
		Text:     string(data),
	}, nil
}

func readThemesResource(s *Server) (*ResourceContent, error) {
	themesDir := filepath.Join(s.SigilDir(), "themes")
	matches, _ := filepath.Glob(filepath.Join(themesDir, "*.yaml"))

	var result []map[string]interface{}
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var theme struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Extends     string `yaml:"extends"`
		}
		if err := parseYAML(data, &theme); err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"name":        theme.Name,
			"description": theme.Description,
			"extends":     theme.Extends,
		})
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return &ResourceContent{
		URI:      "sigil://themes",
		MimeType: "application/json",
		Text:     string(data),
	}, nil
}

func readProjectResource(s *Server) (*ResourceContent, error) {
	path := filepath.Join(s.SigilDir(), "sigil.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("project config not found: %w", err)
	}
	return &ResourceContent{
		URI:      "sigil://project",
		MimeType: "text/yaml",
		Text:     string(data),
	}, nil
}

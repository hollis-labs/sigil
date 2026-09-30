package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/hollis-labs/sigil/internal/components"
	"github.com/hollis-labs/sigil/internal/config"
)

// RegisterAllResources registers Sigil's MCP resources and resource
// templates. Resources are not wrapped by go-mcp, so registration goes
// directly against the underlying official-SDK server.
func RegisterAllResources(s *Server) {
	sdk := s.SDKServer()

	sdk.AddResource(&mcpsdk.Resource{
		URI:         "sigil://pages",
		Name:        "Page listing",
		Description: "List of all page configs",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
		return readPagesResource(s)
	})

	sdk.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		URITemplate: "sigil://pages/{id}",
		Name:        "Page config",
		Description: "Full content of a specific page config",
		MIMEType:    "text/yaml",
	}, func(ctx context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
		id := strings.TrimPrefix(req.Params.URI, "sigil://pages/")
		return readPageResource(s, id)
	})

	sdk.AddResource(&mcpsdk.Resource{
		URI:         "sigil://components",
		Name:        "Component schemas",
		Description: "All available component types and schemas",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
		return readComponentsResource()
	})

	sdk.AddResourceTemplate(&mcpsdk.ResourceTemplate{
		URITemplate: "sigil://components/{type}",
		Name:        "Component schema",
		Description: "Schema for a specific component type",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
		typeName := strings.TrimPrefix(req.Params.URI, "sigil://components/")
		return readComponentResource(typeName)
	})

	sdk.AddResource(&mcpsdk.Resource{
		URI:         "sigil://datasources",
		Name:        "DataSource manifests",
		Description: "All datasource manifests",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
		return readDataSourcesResource(s)
	})

	sdk.AddResource(&mcpsdk.Resource{
		URI:         "sigil://themes",
		Name:        "Themes",
		Description: "All theme definitions",
		MIMEType:    "application/json",
	}, func(ctx context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
		return readThemesResource(s)
	})

	sdk.AddResource(&mcpsdk.Resource{
		URI:         "sigil://project",
		Name:        "Project config",
		Description: "Project configuration",
		MIMEType:    "text/yaml",
	}, func(ctx context.Context, req *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
		return readProjectResource(s)
	})
}

func textResourceResult(uri, mimeType, text string) *mcpsdk.ReadResourceResult {
	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{{URI: uri, MIMEType: mimeType, Text: text}},
	}
}

func readPagesResource(s *Server) (*mcpsdk.ReadResourceResult, error) {
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
	return textResourceResult("sigil://pages", "application/json", string(data)), nil
}

func readPageResource(s *Server, id string) (*mcpsdk.ReadResourceResult, error) {
	path := filepath.Join(s.SigilDir(), "pages", id+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("page %q not found", id)
	}
	return textResourceResult("sigil://pages/"+id, "text/yaml", string(data)), nil
}

func readComponentsResource() (*mcpsdk.ReadResourceResult, error) {
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
	return textResourceResult("sigil://components", "application/json", string(data)), nil
}

func readComponentResource(typeName string) (*mcpsdk.ReadResourceResult, error) {
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
	return textResourceResult("sigil://components/"+typeName, "application/json", string(data)), nil
}

func readDataSourcesResource(s *Server) (*mcpsdk.ReadResourceResult, error) {
	dsDir := filepath.Join(s.SigilDir(), "datasources")
	matches, _ := filepath.Glob(filepath.Join(dsDir, "*.yaml"))

	var result []map[string]interface{}
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"alias":   strings.TrimSuffix(filepath.Base(path), ".yaml"),
			"content": string(data),
		})
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return textResourceResult("sigil://datasources", "application/json", string(data)), nil
}

func readThemesResource(s *Server) (*mcpsdk.ReadResourceResult, error) {
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
	return textResourceResult("sigil://themes", "application/json", string(data)), nil
}

func readProjectResource(s *Server) (*mcpsdk.ReadResourceResult, error) {
	path := filepath.Join(s.SigilDir(), "sigil.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("project config not found: %w", err)
	}
	return textResourceResult("sigil://project", "text/yaml", string(data)), nil
}

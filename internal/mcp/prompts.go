package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/hollis-labs/sigil/internal/components"
	"github.com/hollis-labs/sigil/internal/config"
)

// RegisterAllPrompts registers Sigil's MCP prompts. Prompts are not wrapped
// by go-mcp, so registration goes directly against the underlying
// official-SDK server.
func RegisterAllPrompts(s *Server) {
	sdk := s.SDKServer()

	sdk.AddPrompt(&mcpsdk.Prompt{
		Name:        "sigil_design_page",
		Description: "Guide for designing a Sigil page config",
		Arguments: []*mcpsdk.PromptArgument{
			{Name: "description", Description: "Description of the page to design", Required: true},
		},
	}, handleDesignPage)

	sdk.AddPrompt(&mcpsdk.Prompt{
		Name:        "sigil_review_config",
		Description: "Review a Sigil config for issues and suggest improvements",
		Arguments: []*mcpsdk.PromptArgument{
			{Name: "config", Description: "YAML config to review", Required: true},
		},
	}, handleReviewConfig(s))
}

func handleDesignPage(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	description := req.Params.Arguments["description"]
	if description == "" {
		return nil, fmt.Errorf("'description' argument is required")
	}

	// Build component catalog
	registry := components.NewDefaultRegistry()
	types := registry.Types()

	var catalog strings.Builder
	categories := map[string][]string{}
	for _, t := range types {
		schema, ok := registry.Get(t)
		if !ok {
			continue
		}
		categories[schema.Category] = append(categories[schema.Category], t)
	}

	catOrder := []string{"primitives", "layouts", "navigation", "composites", "data", "forms"}
	for _, cat := range catOrder {
		names, ok := categories[cat]
		if !ok {
			continue
		}
		label := strings.ToUpper(cat[:1]) + cat[1:]
		catalog.WriteString(fmt.Sprintf("## %s\n", label))
		for _, name := range names {
			schema, _ := registry.Get(name)
			catalog.WriteString(fmt.Sprintf("- **%s**: %s\n", name, schema.Description))
			if len(schema.Props) > 0 {
				catalog.WriteString("  Props: ")
				var props []string
				for pname, p := range schema.Props {
					req := ""
					if p.Required {
						req = " (required)"
					}
					props = append(props, fmt.Sprintf("%s (%s%s)", pname, p.Type, req))
				}
				catalog.WriteString(strings.Join(props, ", "))
				catalog.WriteString("\n")
			}
		}
		catalog.WriteString("\n")
	}

	systemPrompt := fmt.Sprintf(`You are designing a Sigil UI page config. The user wants: "%s"

Generate a valid Sigil YAML config following this structure:

%sPage config format:
%s
%sAvailable components:

%s

Rules:
- Every component needs a 'type' field
- Give each component a unique 'id'
- Layout types (rows, columns, grid) contain 'children'
- Actions map to HTMX interactions (navigate, modal, http, emit, confirm, close)
- Use datasource references for data-driven components
- Props must match the component schema
- Valid overlays: page, modal, sheet, drawer, fullscreen

Generate the complete YAML config now.`, description, "```yaml\n", configTemplate, "```\n", catalog.String())

	return &mcpsdk.GetPromptResult{
		Messages: []*mcpsdk.PromptMessage{
			{Role: "user", Content: &mcpsdk.TextContent{Text: systemPrompt}},
		},
	}, nil
}

func handleReviewConfig(s *Server) mcpsdk.PromptHandler {
	return func(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
		yamlConfig := req.Params.Arguments["config"]
		if yamlConfig == "" {
			return nil, fmt.Errorf("'config' argument is required")
		}

		// Parse and validate
		page, err := config.Parse([]byte(yamlConfig))
		if err != nil {
			return &mcpsdk.GetPromptResult{
				Messages: []*mcpsdk.PromptMessage{
					{
						Role: "user",
						Content: &mcpsdk.TextContent{
							Text: fmt.Sprintf("The following Sigil config has a parse error:\n\n```yaml\n%s\n```\n\nError: %s\n\nPlease fix the YAML syntax.", yamlConfig, err),
						},
					},
				},
			}, nil
		}

		registry := components.NewDefaultRegistry()
		result := config.Validate(page, registry)

		var review strings.Builder
		review.WriteString(fmt.Sprintf("Review the following Sigil config:\n\n```yaml\n%s\n```\n\n", yamlConfig))

		if result.Valid && len(result.Warnings) == 0 {
			review.WriteString("**Validation: PASS** — No errors or warnings.\n\n")
			review.WriteString("Suggest any improvements for:\n")
			review.WriteString("- Better component organization\n")
			review.WriteString("- Missing accessibility attributes\n")
			review.WriteString("- UX improvements\n")
			review.WriteString("- Performance considerations\n")
		} else {
			if !result.Valid {
				review.WriteString("**Validation: FAIL**\n\nErrors:\n")
				for _, e := range result.Errors {
					review.WriteString(fmt.Sprintf("- %s\n", e.String()))
				}
				review.WriteString("\n")
			}
			if len(result.Warnings) > 0 {
				review.WriteString("Warnings:\n")
				for _, w := range result.Warnings {
					review.WriteString(fmt.Sprintf("- %s\n", w.String()))
				}
				review.WriteString("\n")
			}
			review.WriteString("Please fix these issues and suggest improvements.\n")
		}

		data, _ := json.MarshalIndent(map[string]interface{}{
			"valid":    result.Valid,
			"errors":   len(result.Errors),
			"warnings": len(result.Warnings),
		}, "", "  ")
		review.WriteString(fmt.Sprintf("\nValidation summary:\n```json\n%s\n```\n", string(data)))

		return &mcpsdk.GetPromptResult{
			Messages: []*mcpsdk.PromptMessage{
				{Role: "user", Content: &mcpsdk.TextContent{Text: review.String()}},
			},
		}, nil
	}
}

const configTemplate = `sigil: "1.0"
kind: page
id: my-page
title: My Page
overlay: page
module: my-module

datasources:
  - alias: MyData
    capabilities: [search, filter, create]

layout:
  id: root
  type: rows
  props:
    gap: 4
    padding: 6
  children:
    - id: header
      type: heading
      props:
        level: 2
        text: Page Title
    - id: content
      type: data-table
      props:
        datasource: MyData

shortcuts:
  - key: Escape
    action: { type: close }
`

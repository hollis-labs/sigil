package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hollis-labs/sigil/internal/components"
	"github.com/spf13/cobra"
)

// NewSchemaCmd creates the schema command.
func NewSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Schema management commands",
	}
	cmd.AddCommand(newSchemaExportCmd())
	return cmd
}

func newSchemaExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export JSON Schema for page configs",
		Long: `Generate a JSON Schema file from the component registry.

The schema can be used with VS Code's YAML extension for autocomplete and validation.

Examples:
  sigil schema export                      Export to schemas/
  sigil schema export --output my-schema/  Export to custom directory`,
		RunE: runSchemaExport,
	}
	cmd.Flags().StringP("output", "o", "schemas", "Output directory")
	return cmd
}

func runSchemaExport(cmd *cobra.Command, args []string) error {
	outputDir, _ := cmd.Flags().GetString("output")

	registry := components.NewDefaultRegistry()
	schema := generatePageSchema(registry)

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling schema: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	outPath := filepath.Join(outputDir, "page.schema.json")
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return fmt.Errorf("writing schema: %w", err)
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "%s Generated %s (%d bytes)\n", green("✓"), outPath, len(data))

	return nil
}

// JSON Schema types
type jsonSchema map[string]interface{}

func generatePageSchema(registry *components.Registry) jsonSchema {
	componentTypes := registry.Types()
	sort.Strings(componentTypes)

	// Build component type enum
	typeEnum := make([]interface{}, len(componentTypes))
	for i, t := range componentTypes {
		typeEnum[i] = t
	}

	// Build per-component prop schemas
	componentSchemas := []jsonSchema{}
	for _, compType := range componentTypes {
		schema, ok := registry.GetSchema(compType)
		if !ok {
			continue
		}
		propSchema := jsonSchema{}
		required := []interface{}{}

		for name, prop := range schema.Props {
			ps := jsonSchema{}
			switch prop.Type {
			case "string":
				ps["type"] = "string"
			case "integer":
				ps["type"] = "integer"
			case "boolean":
				ps["type"] = "boolean"
			case "array":
				ps["type"] = "array"
			case "object":
				ps["type"] = "object"
			default:
				ps["type"] = "string"
			}
			if len(prop.Enum) > 0 {
				enumVals := make([]interface{}, len(prop.Enum))
				for i, v := range prop.Enum {
					enumVals[i] = v
				}
				ps["enum"] = enumVals
			}
			propSchema[name] = ps
			if prop.Required {
				required = append(required, name)
			}
		}

		ifClause := jsonSchema{
			"properties": jsonSchema{
				"type": jsonSchema{"const": compType},
			},
		}
		thenClause := jsonSchema{
			"properties": jsonSchema{
				"props": jsonSchema{
					"type":       "object",
					"properties": propSchema,
				},
			},
		}
		if len(required) > 0 {
			thenProps := thenClause["properties"].(jsonSchema)
			props := thenProps["props"].(jsonSchema)
			props["required"] = required
		}

		componentSchemas = append(componentSchemas, jsonSchema{
			"if":   ifClause,
			"then": thenClause,
		})
	}

	// Component definition (recursive)
	componentDef := jsonSchema{
		"type": "object",
		"properties": jsonSchema{
			"id":   jsonSchema{"type": "string"},
			"type": jsonSchema{"type": "string", "enum": typeEnum},
			"props": jsonSchema{
				"type": "object",
			},
			"actions": jsonSchema{
				"type": "object",
				"additionalProperties": jsonSchema{
					"$ref": "#/$defs/action",
				},
			},
			"children": jsonSchema{
				"type":  "array",
				"items": jsonSchema{"$ref": "#/$defs/component"},
			},
			"shortcuts": jsonSchema{
				"type":  "array",
				"items": jsonSchema{"$ref": "#/$defs/shortcut"},
			},
		},
		"required": []interface{}{"type"},
		"allOf":    componentSchemas,
	}

	// Action definition
	actionDef := jsonSchema{
		"type": "object",
		"properties": jsonSchema{
			"type":    jsonSchema{"type": "string", "enum": []interface{}{"navigate", "modal", "sheet", "http", "emit", "confirm", "close"}},
			"page":    jsonSchema{"type": "string"},
			"url":     jsonSchema{"type": "string"},
			"method":  jsonSchema{"type": "string", "enum": []interface{}{"GET", "POST", "PUT", "DELETE"}},
			"event":   jsonSchema{"type": "string"},
			"message": jsonSchema{"type": "string"},
			"target":  jsonSchema{"type": "string"},
		},
		"required": []interface{}{"type"},
	}

	// Shortcut definition
	shortcutDef := jsonSchema{
		"type": "object",
		"properties": jsonSchema{
			"key":         jsonSchema{"type": "string"},
			"action":      jsonSchema{"$ref": "#/$defs/action"},
			"description": jsonSchema{"type": "string"},
			"global":      jsonSchema{"type": "boolean"},
		},
		"required": []interface{}{"key", "action"},
	}

	return jsonSchema{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://sigil.dev/schemas/page.schema.json",
		"title":       "Sigil Page Config",
		"description": "Schema for Sigil page configuration YAML files",
		"type":        "object",
		"properties": jsonSchema{
			"sigil":       jsonSchema{"type": "string", "const": "1.0"},
			"kind":        jsonSchema{"type": "string", "enum": []interface{}{"page"}},
			"id":          jsonSchema{"type": "string", "pattern": "^[a-z0-9][a-z0-9-]*$"},
			"title":       jsonSchema{"type": "string"},
			"description": jsonSchema{"type": "string"},
			"overlay":     jsonSchema{"type": "string", "enum": []interface{}{"page", "modal", "sheet", "drawer", "fullscreen"}},
			"module":      jsonSchema{"type": "string"},
			"datasources": jsonSchema{
				"type": "array",
				"items": jsonSchema{
					"type": "object",
					"properties": jsonSchema{
						"alias":        jsonSchema{"type": "string"},
						"capabilities": jsonSchema{"type": "array", "items": jsonSchema{"type": "string"}},
					},
					"required": []interface{}{"alias"},
				},
			},
			"layout": jsonSchema{"$ref": "#/$defs/component"},
			"shortcuts": jsonSchema{
				"type":  "array",
				"items": jsonSchema{"$ref": "#/$defs/shortcut"},
			},
		},
		"required": []interface{}{"sigil", "kind", "id", "title", "overlay", "layout"},
		"$defs": jsonSchema{
			"component": componentDef,
			"action":    actionDef,
			"shortcut":  shortcutDef,
		},
	}
}

package components

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed builtin/*.yaml
var builtinSchemas embed.FS

// LoadBuiltinSchemas loads all embedded schema YAML files into the registry.
func LoadBuiltinSchemas(r *Registry) error {
	return fs.WalkDir(builtinSchemas, "builtin", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}
		data, err := builtinSchemas.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded %s: %w", path, err)
		}
		schema, err := parseSchema(data)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		r.Register(schema)
		return nil
	})
}

// LoadCustomSchemas loads schema YAML files from a directory into the registry.
// Custom schemas can override built-in ones (same type name).
func LoadCustomSchemas(r *Registry, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No custom schemas directory is fine
		}
		return fmt.Errorf("reading %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		schema, err := parseSchema(data)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		r.Register(schema)
	}
	return nil
}

func parseSchema(data []byte) (*Schema, error) {
	var schema Schema
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	if schema.Type == "" {
		return nil, fmt.Errorf("schema missing required 'type' field")
	}
	if schema.Category == "" {
		return nil, fmt.Errorf("schema %q missing required 'category' field", schema.Type)
	}
	return &schema, nil
}

// NewDefaultRegistry creates a registry loaded with all built-in schemas.
// This replaces the old hardcoded registry.
func NewDefaultRegistry() *Registry {
	r := NewRegistry()
	if err := LoadBuiltinSchemas(r); err != nil {
		// Built-in schemas are embedded, so this should never fail in production.
		// Fall back to empty registry rather than panic.
		return r
	}
	return r
}

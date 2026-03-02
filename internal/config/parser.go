package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseFile reads a YAML file and returns a Page config.
func ParseFile(path string) (*Page, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	return Parse(data)
}

// Parse parses YAML bytes into a Page config.
func Parse(data []byte) (*Page, error) {
	var page Page
	if err := yaml.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	ApplyDefaults(&page)
	return &page, nil
}

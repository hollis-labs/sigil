package mcp

import "gopkg.in/yaml.v3"

// parseYAML is a helper to parse YAML data into a struct.
func parseYAML(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

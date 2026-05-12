package config

import "fmt"

// ValidationResult contains all validation findings.
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationWarning
}

// ValidationError describes a structural error in the config.
type ValidationError struct {
	Path    string
	Message string
}

// ValidationWarning describes a non-fatal issue in the config.
type ValidationWarning struct {
	Path    string
	Message string
}

func (e ValidationError) String() string {
	return fmt.Sprintf("ERROR [%s]: %s", e.Path, e.Message)
}

func (w ValidationWarning) String() string {
	return fmt.Sprintf("WARN  [%s]: %s", w.Path, w.Message)
}

// Valid overlay types.
var validOverlays = map[string]bool{
	"page":       true,
	"modal":      true,
	"sheet":      true,
	"drawer":     true,
	"fullscreen": true,
}

// Valid action types.
var validActionTypes = map[string]bool{
	"navigate": true,
	"modal":    true,
	"sheet":    true,
	"http":     true,
	"emit":     true,
	"confirm":  true,
	"close":    true,
	"focus":    true,
	"toast":    true,
}

// Validate checks a Page config for structural correctness.
// If registry is nil, component type checks are skipped.
func Validate(page *Page, registry ComponentRegistry) *ValidationResult {
	result := &ValidationResult{Valid: true}

	validateRequired(page, result)
	validateOverlay(page.Overlay, result)

	seen := make(map[string]bool)
	validateComponent(&page.Layout, "layout", seen, registry, result)

	for i, s := range page.Shortcuts {
		validateShortcut(s, fmt.Sprintf("shortcuts[%d]", i), result)
	}

	for i, ds := range page.DataSources {
		validateDataSourceRef(ds, fmt.Sprintf("datasources[%d]", i), result)
	}

	result.Valid = len(result.Errors) == 0
	return result
}

func validateRequired(page *Page, result *ValidationResult) {
	if page.Sigil == "" {
		result.Errors = append(result.Errors, ValidationError{Path: "sigil", Message: "required field 'sigil' is missing"})
	}
	if page.ID == "" {
		result.Errors = append(result.Errors, ValidationError{Path: "id", Message: "required field 'id' is missing"})
	}
	if page.Title == "" {
		result.Errors = append(result.Errors, ValidationError{Path: "title", Message: "required field 'title' is missing"})
	}
	if page.Overlay == "" {
		result.Errors = append(result.Errors, ValidationError{Path: "overlay", Message: "required field 'overlay' is missing"})
	}
	if page.Layout.Type == "" {
		result.Errors = append(result.Errors, ValidationError{Path: "layout.type", Message: "required field 'layout.type' is missing"})
	}
}

func validateOverlay(overlay string, result *ValidationResult) {
	if overlay != "" && !validOverlays[overlay] {
		result.Errors = append(result.Errors, ValidationError{
			Path:    "overlay",
			Message: fmt.Sprintf("invalid overlay type %q (valid: page, modal, sheet, drawer, fullscreen)", overlay),
		})
	}
}

func validateComponent(c *Component, path string, seen map[string]bool, registry ComponentRegistry, result *ValidationResult) {
	if c.Type == "" {
		result.Errors = append(result.Errors, ValidationError{Path: path + ".type", Message: "component type is required"})
	} else if registry != nil && !registry.Has(c.Type) {
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("unknown component type %q", c.Type),
		})
	}

	// Deep prop validation if schema provider is available
	if c.Type != "" && registry != nil {
		if provider, ok := registry.(ComponentSchemaProvider); ok {
			if schema, found := provider.GetSchema(c.Type); found {
				validateProps(c.Props, schema, path, result)
			}
		}
	}

	if c.ID != "" {
		if seen[c.ID] {
			result.Errors = append(result.Errors, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("duplicate component id %q", c.ID),
			})
		}
		seen[c.ID] = true
	} else if c.Type != "" {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Path:    path,
			Message: "component has no id",
		})
	}

	for name, action := range c.Actions {
		validateAction(action, fmt.Sprintf("%s.actions.%s", path, name), result)
	}

	for i, s := range c.Shortcuts {
		validateShortcut(s, fmt.Sprintf("%s.shortcuts[%d]", path, i), result)
	}

	for i, child := range c.Children {
		validateComponent(&child, fmt.Sprintf("%s.children[%d]", path, i), seen, registry, result)
	}
}

func validateAction(action Action, path string, result *ValidationResult) {
	if action.Type == "" {
		result.Errors = append(result.Errors, ValidationError{Path: path + ".type", Message: "action type is required"})
		return
	}
	if !validActionTypes[action.Type] {
		result.Errors = append(result.Errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("invalid action type %q", action.Type),
		})
		return
	}

	switch action.Type {
	case "modal":
		if action.Title == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path, Message: "modal action requires 'title'"})
		}
		if len(action.Fields) == 0 {
			result.Errors = append(result.Errors, ValidationError{Path: path, Message: "modal action requires 'fields'"})
		}
		for i, f := range action.Fields {
			validateFormField(f, fmt.Sprintf("%s.fields[%d]", path, i), result)
		}
	case "navigate":
		if action.Page == "" && action.URL == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path, Message: "navigate action requires 'page' or 'url'"})
		}
	case "http":
		if action.URL == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path, Message: "http action requires 'url'"})
		}
		if action.Method == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path, Message: "http action requires 'method'"})
		}
	case "toast":
		if action.Message == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path, Message: "toast action requires 'message'"})
		}
	}

	if action.OnConfirm != nil {
		validateAction(*action.OnConfirm, path+".onConfirm", result)
	}
	if action.OnSuccess != nil {
		validateAction(*action.OnSuccess, path+".onSuccess", result)
	}
}

func validateFormField(f FormField, path string, result *ValidationResult) {
	if f.Name == "" {
		result.Errors = append(result.Errors, ValidationError{Path: path, Message: "form field requires 'name'"})
	}
	if f.Label == "" {
		result.Errors = append(result.Errors, ValidationError{Path: path, Message: "form field requires 'label'"})
	}
	if f.Type == "" {
		result.Errors = append(result.Errors, ValidationError{Path: path, Message: "form field requires 'type'"})
	}
	if (f.Type == "select" || f.Type == "multiselect") && len(f.Options) == 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{Path: path, Message: "select field has no options"})
	}
}

func validateShortcut(s Shortcut, path string, result *ValidationResult) {
	if s.Key == "" {
		result.Errors = append(result.Errors, ValidationError{Path: path, Message: "shortcut requires 'key'"})
	}
	validateAction(s.Action, path+".action", result)
}

func validateDataSourceRef(ds DataSourceRef, path string, result *ValidationResult) {
	if ds.Alias == "" {
		result.Errors = append(result.Errors, ValidationError{Path: path, Message: "datasource requires 'alias'"})
	}
}

func validateProps(props map[string]interface{}, schema *ComponentSchemaInfo, path string, result *ValidationResult) {
	// Check required props are present
	for name, propDef := range schema.Props {
		if propDef.Required {
			if _, ok := props[name]; !ok {
				result.Errors = append(result.Errors, ValidationError{
					Path:    path + ".props",
					Message: fmt.Sprintf("required prop %q is missing", name),
				})
			}
		}
	}

	// Check provided props
	for name, value := range props {
		propDef, known := schema.Props[name]
		if !known {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Path:    path + ".props." + name,
				Message: fmt.Sprintf("unknown prop %q", name),
			})
			continue
		}

		// Type checking
		if !checkPropType(value, propDef.Type) {
			result.Errors = append(result.Errors, ValidationError{
				Path:    path + ".props." + name,
				Message: fmt.Sprintf("expected type %s, got %T", propDef.Type, value),
			})
			continue
		}

		// Enum checking
		if len(propDef.Enum) > 0 {
			strVal := fmt.Sprintf("%v", value)
			valid := false
			for _, e := range propDef.Enum {
				if strVal == e {
					valid = true
					break
				}
			}
			if !valid {
				result.Errors = append(result.Errors, ValidationError{
					Path:    path + ".props." + name,
					Message: fmt.Sprintf("invalid value %q (valid: %v)", strVal, propDef.Enum),
				})
			}
		}
	}
}

// ValidateAppConfig checks an AppConfig for structural correctness.
func ValidateAppConfig(app *AppConfig) *ValidationResult {
	result := &ValidationResult{Valid: true}

	if app.Sigil == "" {
		result.Errors = append(result.Errors, ValidationError{Path: "sigil", Message: "required field 'sigil' is missing"})
	}
	if app.Name == "" {
		result.Warnings = append(result.Warnings, ValidationWarning{Path: "name", Message: "app name is empty"})
	}
	if len(app.Modules) == 0 {
		result.Errors = append(result.Errors, ValidationError{Path: "modules", Message: "at least one module is required"})
	}

	seenIDs := map[string]bool{}
	for i, mod := range app.Modules {
		path := fmt.Sprintf("modules[%d]", i)
		if mod.ID == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path + ".id", Message: "module id is required"})
		} else if seenIDs[mod.ID] {
			result.Errors = append(result.Errors, ValidationError{Path: path + ".id", Message: fmt.Sprintf("duplicate module id %q", mod.ID)})
		}
		seenIDs[mod.ID] = true
		if mod.Shell == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path + ".shell", Message: "module shell page is required"})
		}
		if len(mod.Pages) == 0 {
			result.Warnings = append(result.Warnings, ValidationWarning{Path: path + ".pages", Message: "module has no pages"})
		}
		// Per-module provider validation (sprint 10 phase 3.5). The same
		// rules that previously applied to app.Providers now apply per module.
		validateProviders(mod.Providers, path+".providers", result)
	}

	if app.API != nil && app.API.BaseURLEnv == "" {
		result.Errors = append(result.Errors, ValidationError{Path: "api.base_url_env", Message: "API base_url_env is required"})
	}

	switch app.TargetMode {
	case "", "spa", "app-router":
		// valid
	default:
		result.Errors = append(result.Errors, ValidationError{
			Path:    "target_mode",
			Message: fmt.Sprintf("invalid target_mode %q (valid: spa, app-router)", app.TargetMode),
		})
	}

	result.Valid = len(result.Errors) == 0
	return result
}

// validateProviders applies provider-config rules to a slice of providers
// (since sprint 10 phase 3.5, providers live per-module and this is reused
// from ValidateAppConfig).
//
// Rules:
//   - id required + unique within the slice
//   - built-in: datasource required, mounts not allowed (warning if present)
//   - custom (source.component set): source.component required,
//     datasource ignored (warning if present)
func validateProviders(providers []ProviderConfig, pathPrefix string, result *ValidationResult) {
	seenProviderIDs := map[string]bool{}
	for i := range providers {
		p := &providers[i]
		path := fmt.Sprintf("%s[%d]", pathPrefix, i)
		if p.ID == "" {
			result.Errors = append(result.Errors, ValidationError{Path: path + ".id", Message: "provider id is required"})
		} else if seenProviderIDs[p.ID] {
			result.Errors = append(result.Errors, ValidationError{Path: path + ".id", Message: fmt.Sprintf("duplicate provider id %q", p.ID)})
		}
		seenProviderIDs[p.ID] = true

		if p.IsCustom() {
			// Custom (source-based) provider: source.component required;
			// mounts must reference includes.
			if p.Source.Component == "" {
				result.Errors = append(result.Errors, ValidationError{Path: path + ".source.component", Message: "custom provider requires source.component"})
			}
			// datasource fields are ignored on custom providers — flag as warning
			// when set, so authors don't silently expect data-driven behavior.
			if p.Datasource != "" {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Path:    path + ".datasource",
					Message: "datasource is ignored on custom providers (source.component is set)",
				})
			}
		} else {
			// Built-in (data-driven) provider: datasource required.
			if p.Datasource == "" {
				result.Errors = append(result.Errors, ValidationError{
					Path:    path + ".datasource",
					Message: "built-in provider requires datasource (or set source.component for a custom provider)",
				})
			}
			if len(p.Mounts) > 0 {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Path:    path + ".mounts",
					Message: "mounts are ignored on built-in providers",
				})
			}
		}
	}
}

func checkPropType(value interface{}, expectedType string) bool {
	// Allow {{varName}} template references for any prop type — resolved at render time
	if s, ok := value.(string); ok && len(s) > 4 && s[:2] == "{{" && s[len(s)-2:] == "}}" {
		return true
	}

	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "integer":
		switch value.(type) {
		case int, int64, float64:
			return true
		}
		return false
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "array":
		switch value.(type) {
		case []interface{}, []map[string]interface{}:
			return true
		}
		return false
	case "object":
		switch value.(type) {
		case map[string]interface{}, map[interface{}]interface{}:
			return true
		}
		return false
	case "enum":
		// Enum values are strings
		_, ok := value.(string)
		return ok
	default:
		return true // Unknown type, don't block
	}
}

package config

import (
	"fmt"
	"strings"
)

// MigrateChange describes a single migration change.
type MigrateChange struct {
	Path    string
	Message string
}

// MigrateResult holds the results of a migration run.
type MigrateResult struct {
	FromVersion string
	ToVersion   string
	Changes     []MigrateChange
	Page        *Page
}

// HasChanges returns true if the migration made any changes.
func (r *MigrateResult) HasChanges() bool {
	return len(r.Changes) > 0
}

// Migrate applies all necessary migrations to bring a page up to the latest version.
func Migrate(page *Page) *MigrateResult {
	result := &MigrateResult{
		FromVersion: page.Sigil,
		ToVersion:   "1.0",
		Page:        page,
	}

	// Migration: ensure all components have IDs
	ensureComponentIDs(&page.Layout, "root", 0, result)

	// Migration: ensure sigil version is set
	if page.Sigil == "" {
		page.Sigil = "1.0"
		result.Changes = append(result.Changes, MigrateChange{
			Path:    "sigil",
			Message: "set version to 1.0",
		})
	}

	// Migration: ensure kind is set
	if page.Kind == "" {
		page.Kind = "page"
		result.Changes = append(result.Changes, MigrateChange{
			Path:    "kind",
			Message: "set kind to page",
		})
	}

	// Migration: ensure overlay is set
	if page.Overlay == "" {
		page.Overlay = "page"
		result.Changes = append(result.Changes, MigrateChange{
			Path:    "overlay",
			Message: "set overlay to page (default)",
		})
	}

	return result
}

func ensureComponentIDs(c *Component, parentID string, index int, result *MigrateResult) {
	if c.ID == "" {
		generated := generateComponentID(c.Type, parentID, index)
		c.ID = generated
		result.Changes = append(result.Changes, MigrateChange{
			Path:    fmt.Sprintf("component[%s]", generated),
			Message: fmt.Sprintf("auto-generated ID %q for %s component", generated, c.Type),
		})
	}

	for i := range c.Children {
		ensureComponentIDs(&c.Children[i], c.ID, i, result)
	}
}

func generateComponentID(compType, parentID string, index int) string {
	base := compType
	if base == "" {
		base = "component"
	}
	// Remove hyphens for cleaner IDs
	base = strings.ReplaceAll(base, "-", "_")
	return fmt.Sprintf("%s_%s_%d", parentID, base, index)
}

// FormatMigrateResult returns a human-readable summary of migration changes.
func FormatMigrateResult(result *MigrateResult, color bool) string {
	if !result.HasChanges() {
		return "No migrations needed"
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Migrated %s → %s (%d changes):\n\n",
		result.FromVersion, result.ToVersion, len(result.Changes)))

	for _, c := range result.Changes {
		if color {
			buf.WriteString(fmt.Sprintf("  \033[36m%s\033[0m: %s\n", c.Path, c.Message))
		} else {
			buf.WriteString(fmt.Sprintf("  %s: %s\n", c.Path, c.Message))
		}
	}

	return buf.String()
}

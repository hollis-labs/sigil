package config

import (
	"fmt"
	"sort"
	"strings"
)

// DiffChange represents a single semantic change between two configs.
type DiffChange struct {
	Type    string // "added", "removed", "modified"
	Path    string // e.g., "layout.children[0].props.text"
	OldVal  string // For modified/removed
	NewVal  string // For modified/added
	Context string // Component type or field name for context
}

// DiffResult holds all changes between two page configs.
type DiffResult struct {
	Changes []DiffChange
}

// HasChanges returns true if there are any differences.
func (r *DiffResult) HasChanges() bool {
	return len(r.Changes) > 0
}

// Diff compares two page configs and returns semantic differences.
func Diff(a, b *Page) *DiffResult {
	result := &DiffResult{}

	// Top-level fields
	if a.Title != b.Title {
		result.Changes = append(result.Changes, DiffChange{
			Type: "modified", Path: "title", OldVal: a.Title, NewVal: b.Title,
		})
	}
	if a.Overlay != b.Overlay {
		result.Changes = append(result.Changes, DiffChange{
			Type: "modified", Path: "overlay", OldVal: a.Overlay, NewVal: b.Overlay,
		})
	}
	if a.Description != b.Description {
		result.Changes = append(result.Changes, DiffChange{
			Type: "modified", Path: "description", OldVal: a.Description, NewVal: b.Description,
		})
	}
	if a.Module != b.Module {
		result.Changes = append(result.Changes, DiffChange{
			Type: "modified", Path: "module", OldVal: a.Module, NewVal: b.Module,
		})
	}

	// DataSources
	diffDataSources(a.DataSources, b.DataSources, result)

	// Layout component tree
	diffComponent(&a.Layout, &b.Layout, "layout", result)

	// Shortcuts
	if len(a.Shortcuts) != len(b.Shortcuts) {
		result.Changes = append(result.Changes, DiffChange{
			Type:   "modified",
			Path:   "shortcuts",
			OldVal: fmt.Sprintf("%d shortcuts", len(a.Shortcuts)),
			NewVal: fmt.Sprintf("%d shortcuts", len(b.Shortcuts)),
		})
	}

	return result
}

func diffDataSources(a, b []DataSourceRef, result *DiffResult) {
	aMap := map[string]DataSourceRef{}
	for _, ds := range a {
		aMap[ds.Alias] = ds
	}
	bMap := map[string]DataSourceRef{}
	for _, ds := range b {
		bMap[ds.Alias] = ds
	}

	for alias := range aMap {
		if _, ok := bMap[alias]; !ok {
			result.Changes = append(result.Changes, DiffChange{
				Type: "removed", Path: "datasources", OldVal: alias, Context: "datasource",
			})
		}
	}
	for alias := range bMap {
		if _, ok := aMap[alias]; !ok {
			result.Changes = append(result.Changes, DiffChange{
				Type: "added", Path: "datasources", NewVal: alias, Context: "datasource",
			})
		}
	}
}

func diffComponent(a, b *Component, path string, result *DiffResult) {
	// Type change
	if a.Type != b.Type {
		result.Changes = append(result.Changes, DiffChange{
			Type: "modified", Path: path + ".type", OldVal: a.Type, NewVal: b.Type,
			Context: a.Type + " → " + b.Type,
		})
		return // Type changed, don't compare props/children
	}

	// ID change
	if a.ID != b.ID {
		result.Changes = append(result.Changes, DiffChange{
			Type: "modified", Path: path + ".id", OldVal: a.ID, NewVal: b.ID, Context: a.Type,
		})
	}

	// Props diff
	diffProps(a.Props, b.Props, path+".props", a.Type, result)

	// Actions diff
	diffActions(a.Actions, b.Actions, path+".actions", a.Type, result)

	// Children diff
	diffChildren(a.Children, b.Children, path+".children", result)
}

func diffProps(a, b map[string]interface{}, path, compType string, result *DiffResult) {
	allKeys := map[string]bool{}
	for k := range a {
		allKeys[k] = true
	}
	for k := range b {
		allKeys[k] = true
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		aVal, aOk := a[k]
		bVal, bOk := b[k]

		if aOk && !bOk {
			result.Changes = append(result.Changes, DiffChange{
				Type: "removed", Path: path + "." + k, OldVal: fmt.Sprintf("%v", aVal), Context: compType,
			})
		} else if !aOk && bOk {
			result.Changes = append(result.Changes, DiffChange{
				Type: "added", Path: path + "." + k, NewVal: fmt.Sprintf("%v", bVal), Context: compType,
			})
		} else {
			aStr := fmt.Sprintf("%v", aVal)
			bStr := fmt.Sprintf("%v", bVal)
			if aStr != bStr {
				result.Changes = append(result.Changes, DiffChange{
					Type: "modified", Path: path + "." + k, OldVal: aStr, NewVal: bStr, Context: compType,
				})
			}
		}
	}
}

func diffActions(a, b map[string]Action, path, compType string, result *DiffResult) {
	allKeys := map[string]bool{}
	for k := range a {
		allKeys[k] = true
	}
	for k := range b {
		allKeys[k] = true
	}

	for k := range allKeys {
		_, aOk := a[k]
		_, bOk := b[k]

		if aOk && !bOk {
			result.Changes = append(result.Changes, DiffChange{
				Type: "removed", Path: path + "." + k, Context: compType,
			})
		} else if !aOk && bOk {
			result.Changes = append(result.Changes, DiffChange{
				Type: "added", Path: path + "." + k, Context: compType,
			})
		} else {
			aAction := a[k]
			bAction := b[k]
			if aAction.Type != bAction.Type {
				result.Changes = append(result.Changes, DiffChange{
					Type:    "modified",
					Path:    path + "." + k + ".type",
					OldVal:  aAction.Type,
					NewVal:  bAction.Type,
					Context: compType,
				})
			}
		}
	}
}

func diffChildren(a, b []Component, path string, result *DiffResult) {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	for i := 0; i < maxLen; i++ {
		childPath := fmt.Sprintf("%s[%d]", path, i)

		if i >= len(a) {
			// Added
			result.Changes = append(result.Changes, DiffChange{
				Type: "added", Path: childPath, NewVal: b[i].Type,
				Context: fmt.Sprintf("%s component", b[i].Type),
			})
			continue
		}
		if i >= len(b) {
			// Removed
			result.Changes = append(result.Changes, DiffChange{
				Type: "removed", Path: childPath, OldVal: a[i].Type,
				Context: fmt.Sprintf("%s component", a[i].Type),
			})
			continue
		}

		diffComponent(&a[i], &b[i], childPath, result)
	}
}

// FormatDiff returns a human-readable diff string with ANSI colors.
func FormatDiff(result *DiffResult, color bool) string {
	if !result.HasChanges() {
		return "No changes"
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%d change(s):\n\n", len(result.Changes)))

	for _, c := range result.Changes {
		var prefix, label string
		switch c.Type {
		case "added":
			if color {
				prefix = "\033[32m+ "
				label = fmt.Sprintf("%s%s\033[0m", prefix, c.Path)
			} else {
				label = fmt.Sprintf("+ %s", c.Path)
			}
			if c.NewVal != "" {
				buf.WriteString(fmt.Sprintf("  %s = %s\n", label, c.NewVal))
			} else {
				buf.WriteString(fmt.Sprintf("  %s\n", label))
			}
		case "removed":
			if color {
				prefix = "\033[31m- "
				label = fmt.Sprintf("%s%s\033[0m", prefix, c.Path)
			} else {
				label = fmt.Sprintf("- %s", c.Path)
			}
			if c.OldVal != "" {
				buf.WriteString(fmt.Sprintf("  %s = %s\n", label, c.OldVal))
			} else {
				buf.WriteString(fmt.Sprintf("  %s\n", label))
			}
		case "modified":
			if color {
				prefix = "\033[33m~ "
				label = fmt.Sprintf("%s%s\033[0m", prefix, c.Path)
			} else {
				label = fmt.Sprintf("~ %s", c.Path)
			}
			buf.WriteString(fmt.Sprintf("  %s: %s → %s\n", label, c.OldVal, c.NewVal))
		}
	}

	return buf.String()
}

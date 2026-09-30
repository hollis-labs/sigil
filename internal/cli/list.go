package cli

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/hollis-labs/sigil/internal/components"
	"github.com/hollis-labs/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command with subcommands.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <type>",
		Short: "List Sigil configs and components",
		Long:  "List pages, datasources, themes, or registered component types.",
	}
	cmd.AddCommand(newListPagesCmd())
	cmd.AddCommand(newListDataSourcesCmd())
	cmd.AddCommand(newListThemesCmd())
	cmd.AddCommand(newListComponentsCmd())
	return cmd
}

func newListPagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pages",
		Short: "List all page configs",
		RunE:  runListPages,
	}
}

func newListDataSourcesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "datasources",
		Short: "List all datasource manifests",
		RunE:  runListDataSources,
	}
}

func newListThemesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "themes",
		Short: "List all themes",
		RunE:  runListThemes,
	}
}

func newListComponentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "components",
		Short: "List registered component types",
		RunE:  runListComponents,
	}
	cmd.Flags().String("type", "", "Show detailed schema for a specific component type")
	return cmd
}

func runListPages(cmd *cobra.Command, args []string) error {
	pattern := filepath.Join(".sigil", "pages", "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "PAGES (%d found)\n", len(matches))
	if len(matches) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  (none)")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	for _, path := range matches {
		page, err := config.ParseFile(path)
		if err != nil {
			fmt.Fprintf(w, "  %s\t(parse error)\t\t\n", filepath.Base(path))
			continue
		}
		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", page.ID, page.Title, page.Overlay, page.Module)
	}
	w.Flush()
	return nil
}

func runListDataSources(cmd *cobra.Command, args []string) error {
	pattern := filepath.Join(".sigil", "datasources", "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "DATASOURCES (%d found)\n", len(matches))
	if len(matches) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  (none)")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	for _, path := range matches {
		name := strings.TrimSuffix(filepath.Base(path), ".yaml")
		fmt.Fprintf(w, "  %s\n", name)
	}
	w.Flush()
	return nil
}

func runListThemes(cmd *cobra.Command, args []string) error {
	pattern := filepath.Join(".sigil", "themes", "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "THEMES (%d found)\n", len(matches))
	if len(matches) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  (none)")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	for _, path := range matches {
		name := strings.TrimSuffix(filepath.Base(path), ".yaml")
		fmt.Fprintf(w, "  %s\n", name)
	}
	w.Flush()
	return nil
}

func runListComponents(cmd *cobra.Command, args []string) error {
	registry := components.NewDefaultRegistry()

	typeName, _ := cmd.Flags().GetString("type")
	if typeName != "" {
		return showComponentDetail(cmd, registry, typeName)
	}

	types := registry.Types()
	fmt.Fprintf(cmd.OutOrStdout(), "COMPONENTS (%d registered)\n", len(types))

	// Group by category
	categories := map[string][]string{}
	for _, t := range types {
		s, _ := registry.Get(t)
		categories[s.Category] = append(categories[s.Category], t)
	}

	// Print in a stable order
	catOrder := []string{"primitives", "layouts", "navigation", "composites", "data", "forms"}
	for _, cat := range catOrder {
		names := categories[cat]
		if len(names) == 0 {
			continue
		}
		sort.Strings(names)
		label := strings.ToUpper(cat[:1]) + cat[1:]
		fmt.Fprintf(cmd.OutOrStdout(), "  %-12s %s\n", label+":", strings.Join(names, ", "))
	}

	return nil
}

func showComponentDetail(cmd *cobra.Command, registry *components.Registry, typeName string) error {
	schema, ok := registry.Get(typeName)
	if !ok {
		return fmt.Errorf("unknown component type %q", typeName)
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "%s (%s)\n", schema.Type, schema.Category)
	fmt.Fprintf(out, "  %s\n\n", schema.Description)

	if len(schema.Props) > 0 {
		fmt.Fprintln(out, "PROPS:")
		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)

		// Sort prop names for stable output
		propNames := make([]string, 0, len(schema.Props))
		for name := range schema.Props {
			propNames = append(propNames, name)
		}
		sort.Strings(propNames)

		for _, name := range propNames {
			p := schema.Props[name]
			req := ""
			if p.Required {
				req = " (required)"
			}
			def := ""
			if p.Default != nil {
				def = fmt.Sprintf("  default: %v", p.Default)
			}
			enumStr := ""
			if len(p.Enum) > 0 {
				enumStr = fmt.Sprintf("  enum: [%s]", strings.Join(p.Enum, ", "))
			}
			desc := ""
			if p.Description != "" {
				desc = " \u2014 " + p.Description
			}
			fmt.Fprintf(w, "  %s\t%s%s%s%s%s\n", name, p.Type, req, def, enumStr, desc)
		}
		w.Flush()
		fmt.Fprintln(out)
	}

	if len(schema.Actions) > 0 {
		fmt.Fprintln(out, "ACTIONS:")
		for name, a := range schema.Actions {
			fmt.Fprintf(out, "  %s \u2014 %s\n", name, a.Description)
		}
		fmt.Fprintln(out)
	}

	if len(schema.Slots) > 0 {
		fmt.Fprintln(out, "SLOTS:")
		for name, s := range schema.Slots {
			accepts := ""
			if len(s.Accepts) > 0 {
				accepts = fmt.Sprintf(" accepts: [%s]", strings.Join(s.Accepts, ", "))
			}
			fmt.Fprintf(out, "  %s \u2014 %s%s\n", name, s.Description, accepts)
		}
		fmt.Fprintln(out)
	}

	if len(schema.Shortcuts) > 0 {
		fmt.Fprintln(out, "SHORTCUTS:")
		for _, s := range schema.Shortcuts {
			when := ""
			if s.When != "" {
				when = fmt.Sprintf(" (when: %s)", s.When)
			}
			fmt.Fprintf(out, "  %-12s %s%s\n", s.Key, s.Description, when)
		}
	}

	return nil
}

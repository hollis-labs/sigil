package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chrispian/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewInitCmd creates the init subcommand.
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a Sigil project",
		Long:  "Create a .sigil/ directory with project config, default theme, and empty directories.",
		RunE:  runInit,
	}
	cmd.Flags().String("name", "", "Project name (defaults to directory name)")
	cmd.Flags().String("theme", "dark", "Initial theme preset (dark or light)")
	cmd.Flags().Bool("force", false, "Overwrite existing .sigil/ directory")
	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	name, _ := cmd.Flags().GetString("name")
	themeName, _ := cmd.Flags().GetString("theme")

	if name == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		name = filepath.Base(cwd)
	}

	sigilDir := ".sigil"
	if _, err := os.Stat(sigilDir); err == nil && !force {
		return fmt.Errorf(".sigil/ already exists (use --force to overwrite)")
	}

	// Create directory structure
	dirs := []string{
		filepath.Join(sigilDir, "pages"),
		filepath.Join(sigilDir, "components"),
		filepath.Join(sigilDir, "datasources"),
		filepath.Join(sigilDir, "themes"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}

	// Create .gitkeep files
	for _, dir := range dirs[:3] { // pages, components, datasources
		keepFile := filepath.Join(dir, ".gitkeep")
		if err := os.WriteFile(keepFile, []byte{}, 0644); err != nil {
			return fmt.Errorf("creating %s: %w", keepFile, err)
		}
	}

	// Write project config
	projConfig := config.ProjectConfig{
		Version: "1.0",
		Name:    name,
	}
	projConfig.Defaults.Theme = "default"
	projConfig.Defaults.Renderer = "go-templ"
	projConfig.Defaults.Output = "internal/ui/"
	projConfig.Components.Builtin = true
	projConfig.Components.CustomDir = ".sigil/components"
	projConfig.DataSources.Dir = ".sigil/datasources"
	projConfig.Themes.Dir = ".sigil/themes"
	projConfig.Generation.Clean = false

	configPath := filepath.Join(sigilDir, "sigil.yaml")
	if err := config.WriteFile(configPath, &projConfig); err != nil {
		return fmt.Errorf("writing sigil.yaml: %w", err)
	}

	// Write default theme
	theme := defaultTheme(themeName)
	themePath := filepath.Join(sigilDir, "themes", "default.yaml")
	if err := config.WriteFile(themePath, theme); err != nil {
		return fmt.Errorf("writing default theme: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Initialized Sigil project %q\n", name)
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Created:")
	fmt.Fprintln(cmd.OutOrStdout(), "  .sigil/sigil.yaml          Project config")
	fmt.Fprintln(cmd.OutOrStdout(), "  .sigil/themes/default.yaml Default theme")
	fmt.Fprintln(cmd.OutOrStdout(), "  .sigil/pages/              Page definitions")
	fmt.Fprintln(cmd.OutOrStdout(), "  .sigil/components/         Custom component schemas")
	fmt.Fprintln(cmd.OutOrStdout(), "  .sigil/datasources/        DataSource manifests")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Next steps:")
	fmt.Fprintln(cmd.OutOrStdout(), "  sigil new page my-page     Create a page config")
	fmt.Fprintln(cmd.OutOrStdout(), "  sigil new datasource User  Create a datasource manifest")
	fmt.Fprintln(cmd.OutOrStdout(), "  sigil validate             Validate all configs")

	return nil
}

// themeData is a generic map for writing theme YAML.
type themeData struct {
	Name        string                            `yaml:"name"`
	Description string                            `yaml:"description"`
	Extends     *string                           `yaml:"extends"`
	Tokens      map[string]map[string]interface{} `yaml:"tokens"`
}

func defaultTheme(preset string) *themeData {
	if preset == "light" {
		return &themeData{
			Name:        "default",
			Description: "Default light theme",
			Extends:     nil,
			Tokens: map[string]map[string]interface{}{
				"colors": {
					"background":   "255 255 255",
					"surface":      "249 250 251",
					"surface-2":    "243 244 246",
					"border":       "209 213 219",
					"text":         "17 24 39",
					"text-muted":   "107 114 128",
					"accent":       "16 185 129",
					"accent-hover": "5 150 105",
					"danger":       "239 68 68",
					"danger-hover": "220 38 38",
					"warning":      "245 158 11",
					"info":         "59 130 246",
					"success":      "16 185 129",
				},
				"typography": {
					"font-sans":       "ui-sans-serif, system-ui, -apple-system, sans-serif",
					"font-mono":       "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace",
					"text-xs":         "0.75rem",
					"text-sm":         "0.875rem",
					"text-base":       "1rem",
					"text-lg":         "1.125rem",
					"text-xl":         "1.25rem",
					"text-2xl":        "1.5rem",
					"leading-tight":   "1.25",
					"leading-normal":  "1.5",
					"leading-relaxed": "1.625",
				},
				"spacing": {
					"unit": "0.25rem",
				},
				"radius": {
					"none": "0",
					"sm":   "0.125rem",
					"md":   "0.375rem",
					"lg":   "0.5rem",
					"xl":   "0.75rem",
					"full": "9999px",
				},
				"shadows": {
					"sm": "0 1px 2px 0 rgb(0 0 0 / 0.05)",
					"md": "0 4px 6px -1px rgb(0 0 0 / 0.1)",
					"lg": "0 10px 15px -3px rgb(0 0 0 / 0.1)",
				},
			},
		}
	}

	// Dark theme (default)
	return &themeData{
		Name:        "default",
		Description: "Default dark theme",
		Extends:     nil,
		Tokens: map[string]map[string]interface{}{
			"colors": {
				"background":   "9 9 11",
				"surface":      "24 24 27",
				"surface-2":    "39 39 42",
				"border":       "63 63 70",
				"text":         "244 244 245",
				"text-muted":   "161 161 170",
				"accent":       "16 185 129",
				"accent-hover": "5 150 105",
				"danger":       "239 68 68",
				"danger-hover": "220 38 38",
				"warning":      "245 158 11",
				"info":         "59 130 246",
				"success":      "16 185 129",
			},
			"typography": {
				"font-sans":       "ui-sans-serif, system-ui, -apple-system, sans-serif",
				"font-mono":       "ui-monospace, SFMono-Regular, Menlo, Consolas, monospace",
				"text-xs":         "0.75rem",
				"text-sm":         "0.875rem",
				"text-base":       "1rem",
				"text-lg":         "1.125rem",
				"text-xl":         "1.25rem",
				"text-2xl":        "1.5rem",
				"leading-tight":   "1.25",
				"leading-normal":  "1.5",
				"leading-relaxed": "1.625",
			},
			"spacing": {
				"unit": "0.25rem",
			},
			"radius": {
				"none": "0",
				"sm":   "0.125rem",
				"md":   "0.375rem",
				"lg":   "0.5rem",
				"xl":   "0.75rem",
				"full": "9999px",
			},
			"shadows": {
				"sm": "0 1px 2px 0 rgb(0 0 0 / 0.05)",
				"md": "0 4px 6px -1px rgb(0 0 0 / 0.1)",
				"lg": "0 10px 15px -3px rgb(0 0 0 / 0.1)",
			},
		},
	}
}

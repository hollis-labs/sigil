package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hollis-labs/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewNewThemeCmd creates the "new theme" subcommand.
func NewNewThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme <id>",
		Short: "Create a new theme",
		Long:  "Create a new theme from a preset or base theme.",
		Args:  cobra.ExactArgs(1),
		RunE:  runNewTheme,
	}
	cmd.Flags().String("extends", "", "Base theme to inherit from")
	cmd.Flags().String("preset", "dark", "Start from a preset (dark or light)")
	cmd.Flags().Bool("force", false, "Overwrite existing file")
	return cmd
}

// newThemeData is the YAML structure for a new theme file.
type newThemeData struct {
	Name        string                            `yaml:"name"`
	Description string                            `yaml:"description"`
	Extends     string                            `yaml:"extends,omitempty"`
	Tokens      map[string]map[string]interface{} `yaml:"tokens"`
}

func runNewTheme(cmd *cobra.Command, args []string) error {
	id := args[0]
	extends, _ := cmd.Flags().GetString("extends")
	preset, _ := cmd.Flags().GetString("preset")
	force, _ := cmd.Flags().GetBool("force")

	outDir := filepath.Join(".sigil", "themes")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating themes directory: %w", err)
	}

	outPath := filepath.Join(outDir, id+".yaml")
	if _, err := os.Stat(outPath); err == nil && !force {
		return fmt.Errorf("%s already exists (use --force to overwrite)", outPath)
	}

	theme := &newThemeData{
		Name:    id,
		Extends: extends,
	}

	if extends != "" {
		// When extending, just provide empty token overrides
		theme.Description = fmt.Sprintf("Theme extending %s", extends)
		theme.Tokens = map[string]map[string]interface{}{
			"colors": {},
		}
	} else {
		// Use preset tokens
		theme.Tokens = presetTokens(preset)
		if preset == "light" {
			theme.Description = "Light theme"
		} else {
			theme.Description = "Dark theme"
		}
	}

	if err := config.WriteFile(outPath, theme); err != nil {
		return fmt.Errorf("writing theme: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", outPath)
	return nil
}

func presetTokens(preset string) map[string]map[string]interface{} {
	if preset == "light" {
		return map[string]map[string]interface{}{
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
		}
	}

	// Dark preset (default)
	return map[string]map[string]interface{}{
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
	}
}

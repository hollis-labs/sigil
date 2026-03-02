package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrispian/sigil/internal/components"
	"github.com/chrispian/sigil/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewImportCmd creates the import command.
func NewImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import page configs from JSON",
		Long:  "Import page configs from JSON files, validating them before writing.",
		RunE:  runImport,
	}
	cmd.Flags().String("from", "json", "Import format (json)")
	cmd.Flags().StringP("input", "i", "", "Input file or directory (required)")
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	cmd.Flags().Bool("dry-run", false, "Validate without writing")
	cmd.MarkFlagRequired("input")
	return cmd
}

func runImport(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("from")
	input, _ := cmd.Flags().GetString("input")
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if format != "json" {
		return fmt.Errorf("unsupported import format %q (supported: json)", format)
	}

	// Determine input files
	var jsonFiles []string
	info, err := os.Stat(input)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	if info.IsDir() {
		matches, err := filepath.Glob(filepath.Join(input, "*.json"))
		if err != nil {
			return fmt.Errorf("scanning input directory: %w", err)
		}
		jsonFiles = matches
	} else {
		jsonFiles = []string{input}
	}

	if len(jsonFiles) == 0 {
		return fmt.Errorf("no JSON files found in %q", input)
	}

	pagesDir := filepath.Join(sigilDir, "pages")
	if !dryRun {
		if err := os.MkdirAll(pagesDir, 0755); err != nil {
			return fmt.Errorf("creating pages directory: %w", err)
		}
	}

	registry := components.NewDefaultRegistry()
	out := cmd.OutOrStdout()
	imported := 0
	var errors []string

	for _, path := range jsonFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", filepath.Base(path), err))
			continue
		}

		var page config.Page
		if err := json.Unmarshal(data, &page); err != nil {
			errors = append(errors, fmt.Sprintf("%s: invalid JSON: %v", filepath.Base(path), err))
			continue
		}

		// Apply defaults
		config.ApplyDefaults(&page)

		// Validate
		result := config.Validate(&page, registry)
		if !result.Valid {
			var msgs []string
			for _, e := range result.Errors {
				msgs = append(msgs, e.String())
			}
			errors = append(errors, fmt.Sprintf("%s: validation failed:\n    %s", filepath.Base(path), strings.Join(msgs, "\n    ")))
			continue
		}

		if dryRun {
			fmt.Fprintf(out, "  VALID: %s (id: %s)\n", filepath.Base(path), page.ID)
			imported++
			continue
		}

		// Marshal to YAML
		yamlData, err := yaml.Marshal(&page)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: YAML marshal error: %v", filepath.Base(path), err))
			continue
		}

		outPath := filepath.Join(pagesDir, page.ID+".yaml")
		if err := os.WriteFile(outPath, yamlData, 0644); err != nil {
			errors = append(errors, fmt.Sprintf("%s: write error: %v", filepath.Base(path), err))
			continue
		}

		imported++
		fmt.Fprintf(out, "  %s -> %s\n", filepath.Base(path), outPath)
	}

	if len(errors) > 0 {
		fmt.Fprintf(out, "\nErrors:\n")
		for _, e := range errors {
			fmt.Fprintf(out, "  %s\n", e)
		}
	}

	action := "Imported"
	if dryRun {
		action = "Validated"
	}
	fmt.Fprintf(out, "\n%s %d page(s)", action, imported)
	if len(errors) > 0 {
		fmt.Fprintf(out, ", %d error(s)", len(errors))
	}
	fmt.Fprintf(out, "\n")

	if len(errors) > 0 && imported == 0 {
		return fmt.Errorf("all imports failed")
	}
	return nil
}

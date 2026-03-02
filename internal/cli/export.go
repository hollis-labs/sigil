package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chrispian/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewExportCmd creates the export command.
func NewExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export page configs as JSON",
		Long:  "Export all page configs (or specific pages) to JSON files for sharing or backup.",
		RunE:  runExport,
	}
	cmd.Flags().String("format", "json", "Export format (json)")
	cmd.Flags().StringP("output", "o", ".", "Output directory for exported files")
	cmd.Flags().StringSlice("pages", nil, "Specific page IDs to export (default: all)")
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	return cmd
}

func runExport(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")
	pageIDs, _ := cmd.Flags().GetStringSlice("pages")
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")

	if format != "json" {
		return fmt.Errorf("unsupported export format %q (supported: json)", format)
	}

	pagesDir := filepath.Join(sigilDir, "pages")
	matches, err := filepath.Glob(filepath.Join(pagesDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("scanning pages: %w", err)
	}

	filterSet := map[string]bool{}
	for _, id := range pageIDs {
		filterSet[id] = true
	}

	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	out := cmd.OutOrStdout()
	exported := 0

	for _, path := range matches {
		page, err := config.ParseFile(path)
		if err != nil {
			fmt.Fprintf(out, "WARNING: skipping %s: %v\n", filepath.Base(path), err)
			continue
		}

		if len(filterSet) > 0 && !filterSet[page.ID] {
			continue
		}

		data, err := json.MarshalIndent(page, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling page %q: %w", page.ID, err)
		}

		outPath := filepath.Join(output, page.ID+".json")
		if err := os.WriteFile(outPath, append(data, '\n'), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}

		exported++
		fmt.Fprintf(out, "  %s -> %s\n", filepath.Base(path), outPath)
	}

	fmt.Fprintf(out, "\nExported %d page(s)\n", exported)
	return nil
}

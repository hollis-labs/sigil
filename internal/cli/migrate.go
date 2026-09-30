package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hollis-labs/sigil/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewMigrateCmd creates the migrate command.
func NewMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate page configs to latest schema version",
		Long: `Detect config schema version and apply migrations.

Current migrations:
  - Auto-generate missing component IDs
  - Set default sigil version, kind, and overlay

Examples:
  sigil migrate                    Migrate all pages in .sigil/
  sigil migrate --dry-run          Preview changes without writing
  sigil migrate --pages dashboard  Migrate specific pages`,
		RunE: runMigrate,
	}
	cmd.Flags().Bool("dry-run", false, "Preview changes without writing files")
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	cmd.Flags().StringSlice("pages", nil, "Specific page IDs to migrate (default: all)")
	return cmd
}

func runMigrate(cmd *cobra.Command, args []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")
	pageFilter, _ := cmd.Flags().GetStringSlice("pages")
	noColor, _ := cmd.Root().Flags().GetBool("no-color")

	pagesDir := filepath.Join(sigilDir, "pages")
	pattern := filepath.Join(pagesDir, "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("globbing pages: %w", err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("no page configs found in %s", pagesDir)
	}

	filterSet := map[string]bool{}
	for _, f := range pageFilter {
		filterSet[f] = true
	}

	out := cmd.OutOrStdout()
	totalChanges := 0
	pagesChanged := 0

	for _, path := range matches {
		page, err := config.ParseFile(path)
		if err != nil {
			fmt.Fprintf(out, "%s %s: %v\n", red("✗"), filepath.Base(path), err)
			continue
		}

		if len(filterSet) > 0 && !filterSet[page.ID] {
			continue
		}

		result := config.Migrate(page)

		if result.HasChanges() {
			pagesChanged++
			totalChanges += len(result.Changes)

			fmt.Fprintf(out, "\n%s %s\n", bold(page.ID), dim(path))
			fmt.Fprint(out, config.FormatMigrateResult(result, !noColor))

			if !dryRun {
				data, err := yaml.Marshal(result.Page)
				if err != nil {
					return fmt.Errorf("marshaling %s: %w", page.ID, err)
				}
				if err := os.WriteFile(path, data, 0644); err != nil {
					return fmt.Errorf("writing %s: %w", path, err)
				}
			}
		}
	}

	if totalChanges == 0 {
		fmt.Fprintf(out, "%s All pages up to date\n", green("✓"))
	} else if dryRun {
		fmt.Fprintf(out, "\n%s DRY RUN: %d changes across %d pages (no files written)\n",
			yellow("!"), totalChanges, pagesChanged)
	} else {
		fmt.Fprintf(out, "\n%s Migrated %d pages (%d changes)\n",
			green("✓"), pagesChanged, totalChanges)
	}

	return nil
}

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chrispian/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewNewCmd creates the "new" parent command.
func NewNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create new Sigil configs",
		Long:  "Create new page, datasource, or theme configs from templates.",
	}
	cmd.AddCommand(NewNewPageCmd())
	cmd.AddCommand(NewNewDataSourceCmd())
	cmd.AddCommand(NewNewThemeCmd())
	return cmd
}

// NewNewPageCmd creates the "new page" subcommand.
func NewNewPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page <id>",
		Short: "Create a new page config",
		Long:  "Create a new page config from a starter template.",
		Args:  cobra.ExactArgs(1),
		RunE:  runNewPage,
	}
	cmd.Flags().String("title", "", "Page title (inferred from ID if omitted)")
	cmd.Flags().String("overlay", "page", "Display mode (page, modal, sheet, drawer, fullscreen)")
	cmd.Flags().String("module", "", "Module grouping")
	cmd.Flags().String("layout", "rows", "Root layout type (rows, columns, grid)")
	cmd.Flags().String("datasource", "", "Add a datasource declaration")
	cmd.Flags().Bool("force", false, "Overwrite existing file")
	return cmd
}

func runNewPage(cmd *cobra.Command, args []string) error {
	id := args[0]
	title, _ := cmd.Flags().GetString("title")
	overlay, _ := cmd.Flags().GetString("overlay")
	module, _ := cmd.Flags().GetString("module")
	layout, _ := cmd.Flags().GetString("layout")
	datasource, _ := cmd.Flags().GetString("datasource")
	force, _ := cmd.Flags().GetBool("force")

	outDir := filepath.Join(".sigil", "pages")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating pages directory: %w", err)
	}

	outPath := filepath.Join(outDir, id+".yaml")
	if _, err := os.Stat(outPath); err == nil && !force {
		return fmt.Errorf("%s already exists (use --force to overwrite)", outPath)
	}

	page := config.NewPageTemplate(id, title, overlay, module, layout, datasource)
	if err := config.WriteFile(outPath, page); err != nil {
		return fmt.Errorf("writing page config: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", outPath)
	return nil
}

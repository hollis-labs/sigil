package cli

import (
	"fmt"
	"strings"

	"github.com/chrispian/sigil/internal/renderer"
	_ "github.com/chrispian/sigil/internal/renderer/gotempl" // register go-templ renderer
	"github.com/spf13/cobra"
)

// NewGenerateCmd creates the generate command.
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate framework-specific code from Sigil configs",
		Long:  "Parse, validate, and render all page configs into the target framework's source files.",
		RunE:  runGenerate,
	}
	cmd.Flags().StringP("target", "t", "go-templ", "Renderer target (go-templ)")
	cmd.Flags().StringP("output", "o", "", "Output directory (default: from sigil.yaml)")
	cmd.Flags().StringSlice("pages", nil, "Specific page IDs to generate (default: all)")
	cmd.Flags().String("theme", "", "Theme name override (default: from sigil.yaml)")
	cmd.Flags().Bool("clean", false, "Remove output directory contents before generating")
	cmd.Flags().Bool("dry-run", false, "List files that would be generated without writing")
	cmd.Flags().String("go-module", "", "Go module path for imports")
	return cmd
}

func runGenerate(cmd *cobra.Command, args []string) error {
	target, _ := cmd.Flags().GetString("target")
	output, _ := cmd.Flags().GetString("output")
	pages, _ := cmd.Flags().GetStringSlice("pages")
	theme, _ := cmd.Flags().GetString("theme")
	clean, _ := cmd.Flags().GetBool("clean")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	goModule, _ := cmd.Flags().GetString("go-module")

	if output == "" {
		output = "internal/ui"
	}

	cfg := renderer.GenerateConfig{
		Target:   target,
		OutputDir: output,
		SigilDir:  ".sigil",
		Pages:    pages,
		Theme:    theme,
		Clean:    clean,
		DryRun:   dryRun,
		GoModule: goModule,
	}

	result, err := renderer.Generate(cfg)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()

	if dryRun {
		fmt.Fprintf(out, "DRY RUN — %d pages, %d files would be generated:\n\n", result.PageCount, len(result.Files))
		for _, f := range result.Files {
			fmt.Fprintf(out, "  %s (%d bytes)\n", f.Path, len(f.Content))
		}
		return nil
	}

	fmt.Fprintf(out, "Generated %d files from %d pages\n", len(result.Files), result.PageCount)
	fmt.Fprintf(out, "Target: %s\n", target)
	fmt.Fprintf(out, "Output: %s/\n\n", output)

	// Group by directory for display
	dirs := map[string][]string{}
	for _, f := range result.Files {
		dir := f.Path
		if idx := strings.LastIndex(dir, "/"); idx >= 0 {
			dir = dir[:idx]
		}
		dirs[dir] = append(dirs[dir], f.Path)
	}
	for dir, files := range dirs {
		fmt.Fprintf(out, "  %s/ (%d files)\n", dir, len(files))
	}

	return nil
}

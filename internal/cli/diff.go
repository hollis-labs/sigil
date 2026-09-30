package cli

import (
	"fmt"

	"github.com/hollis-labs/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewDiffCmd creates the diff command.
func NewDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <file-a> <file-b>",
		Short: "Compare two page configs semantically",
		Long: `Parse two Sigil page configs and show semantic differences.

Unlike a text diff, this compares at the config level: added/removed components,
changed props, new datasources, etc.

Examples:
  sigil diff old.yaml new.yaml
  sigil diff .sigil/pages/dashboard.yaml backup/dashboard.yaml`,
		Args: cobra.ExactArgs(2),
		RunE: runDiff,
	}
	return cmd
}

func runDiff(cmd *cobra.Command, args []string) error {
	noColor, _ := cmd.Root().Flags().GetBool("no-color")

	pageA, err := config.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[0], err)
	}
	pageB, err := config.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[1], err)
	}

	result := config.Diff(pageA, pageB)
	out := cmd.OutOrStdout()

	if !result.HasChanges() {
		fmt.Fprintf(out, "%s No differences\n", green("✓"))
		return nil
	}

	fmt.Fprintf(out, "Comparing %s ↔ %s\n\n", args[0], args[1])
	fmt.Fprint(out, config.FormatDiff(result, !noColor))

	return nil
}

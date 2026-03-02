package cli

import (
	"fmt"
	"path/filepath"

	"github.com/chrispian/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewValidateCmd creates the validate subcommand.
func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate Sigil configs",
		Long:  "Validate Sigil YAML configs against component schemas.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runValidate,
	}
	cmd.Flags().Bool("strict", false, "Fail on warnings too")
	return cmd
}

func runValidate(cmd *cobra.Command, args []string) error {
	strict, _ := cmd.Flags().GetBool("strict")

	var paths []string
	if len(args) == 1 {
		paths = []string{args[0]}
	} else {
		// Scan .sigil/pages/ for YAML files
		pattern := filepath.Join(".sigil", "pages", "*.yaml")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("scanning .sigil/pages/: %w", err)
		}
		if len(matches) == 0 {
			fmt.Fprintln(cmd.ErrOrStderr(), "No config files found. Provide a path or create .sigil/pages/*.yaml")
			return nil
		}
		paths = matches
	}

	hasErrors := false
	for _, path := range paths {
		page, err := config.ParseFile(path)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "\u2717 %s \u2014 parse error: %v\n", path, err)
			hasErrors = true
			continue
		}

		result := config.Validate(page, nil)

		if result.Valid && len(result.Warnings) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "\u2713 %s \u2014 valid\n", path)
			continue
		}

		if result.Valid && !strict {
			fmt.Fprintf(cmd.OutOrStdout(), "\u2713 %s \u2014 valid (%d warning(s))\n", path, len(result.Warnings))
		} else {
			warnCount := len(result.Warnings)
			errCount := len(result.Errors)
			fmt.Fprintf(cmd.OutOrStdout(), "\u2717 %s \u2014 %d error(s), %d warning(s)\n", path, errCount, warnCount)
			hasErrors = true
		}

		for _, e := range result.Errors {
			fmt.Fprintf(cmd.OutOrStdout(), "  ERROR [%s]: %s\n", e.Path, e.Message)
		}
		for _, w := range result.Warnings {
			fmt.Fprintf(cmd.OutOrStdout(), "  WARN  [%s]: %s\n", w.Path, w.Message)
		}

		if strict && len(result.Warnings) > 0 {
			hasErrors = true
		}
	}

	if hasErrors {
		return fmt.Errorf("validation failed")
	}
	return nil
}

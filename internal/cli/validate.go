package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hollis-labs/sigil/internal/components"
	"github.com/hollis-labs/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewValidateCmd creates the validate subcommand.
func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate Sigil configs",
		Long: `Validate Sigil YAML configs against component schemas.

Examples:
  sigil validate                     Validate all pages in .sigil/pages/
  sigil validate my-page.yaml        Validate a specific file
  sigil validate --strict            Fail on warnings too`,
		Args: cobra.MaximumNArgs(1),
		RunE: runValidate,
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

	registry := components.NewDefaultRegistry()
	allTypes := registry.Types()

	hasErrors := false
	for _, path := range paths {
		page, err := config.ParseFile(path)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%s %s — parse error: %v\n", red("\u2717"), path, err)
			hasErrors = true
			continue
		}

		result := config.Validate(page, registry)

		if result.Valid && len(result.Warnings) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s — %s\n", green("\u2713"), path, green("valid"))
			continue
		}

		if result.Valid && !strict {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s — %s (%d warning(s))\n", green("\u2713"), path, green("valid"), len(result.Warnings))
		} else {
			warnCount := len(result.Warnings)
			errCount := len(result.Errors)
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s — %s\n", red("\u2717"), path, red(fmt.Sprintf("%d error(s), %d warning(s)", errCount, warnCount)))
			hasErrors = true
		}

		for _, e := range result.Errors {
			msg := fmt.Sprintf("  %s [%s]: %s", red("ERROR"), e.Path, e.Message)
			// Add "did you mean?" suggestions for unknown component types
			if strings.Contains(e.Message, "unknown component type") {
				typeName := extractQuoted(e.Message)
				if typeName != "" {
					if suggestion := SuggestComponentType(typeName, allTypes); suggestion != "" {
						msg += yellow(fmt.Sprintf(" (did you mean %q?)", suggestion))
					}
				}
			}
			fmt.Fprintln(cmd.OutOrStdout(), msg)
		}
		for _, w := range result.Warnings {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s [%s]: %s\n", yellow("WARN "), w.Path, w.Message)
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

// extractQuoted extracts the first quoted string from a message.
func extractQuoted(s string) string {
	start := strings.Index(s, `"`)
	if start < 0 {
		return ""
	}
	end := strings.Index(s[start+1:], `"`)
	if end < 0 {
		return ""
	}
	return s[start+1 : start+1+end]
}

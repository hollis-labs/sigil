package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewNewThemeCmd creates the "new theme" subcommand.
func NewNewThemeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "theme <id>",
		Short: "Create a new theme",
		Long:  "Create a new theme from a preset or base theme.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not yet implemented")
		},
	}
	cmd.Flags().String("extends", "", "Base theme to inherit from")
	cmd.Flags().String("preset", "dark", "Start from a preset (dark or light)")
	cmd.Flags().Bool("force", false, "Overwrite existing file")
	return cmd
}

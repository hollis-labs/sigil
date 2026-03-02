package cli

import "github.com/spf13/cobra"

// NewRootCmd creates the root sigil command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sigil",
		Short: "UI configuration and code generation tool",
		Long:  "Sigil turns declarative YAML configs into framework-specific UI code.",
	}
	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewNewCmd())
	cmd.AddCommand(NewValidateCmd())
	return cmd
}

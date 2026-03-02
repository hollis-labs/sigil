package cli

import "github.com/spf13/cobra"

// NewRootCmd creates the root sigil command.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sigil",
		Short: "UI configuration and code generation tool",
		Long: `Sigil turns declarative YAML configs into framework-specific UI code.

Examples:
  sigil init                        Initialize a new Sigil project
  sigil new page dashboard          Create a new page config
  sigil validate                    Validate all page configs
  sigil generate                    Generate Go/Templ code from configs
  sigil preview dashboard           Preview a page in the browser
  sigil list components             List all registered component types
  sigil export --format json        Export pages as JSON
  sigil import --from json -i dir/  Import pages from JSON
  sigil mcp serve                   Start MCP server for AI agents`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			noColor, _ := cmd.Root().Flags().GetBool("no-color")
			if noColor {
				DisableColor()
			}
		},
	}

	// Global flags
	cmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
	cmd.PersistentFlags().Bool("verbose", false, "Enable verbose/debug output")

	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(NewListCmd())
	cmd.AddCommand(NewNewCmd())
	cmd.AddCommand(NewValidateCmd())
	cmd.AddCommand(NewGenerateCmd())
	cmd.AddCommand(NewPreviewCmd())
	cmd.AddCommand(NewExportCmd())
	cmd.AddCommand(NewImportCmd())
	cmd.AddCommand(NewServeCmd())
	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(NewMCPCmd())
	return cmd
}

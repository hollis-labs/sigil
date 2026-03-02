package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewNewDataSourceCmd creates the "new datasource" subcommand.
func NewNewDataSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datasource <alias>",
		Short: "Create a new datasource manifest",
		Long:  "Create a new datasource manifest from a starter template.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not yet implemented")
		},
	}
	cmd.Flags().String("capabilities", "search,filter,sort,paginate", "Comma-separated capabilities")
	cmd.Flags().Bool("force", false, "Overwrite existing file")
	return cmd
}

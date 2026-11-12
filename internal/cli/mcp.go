package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/hollis-labs/substrate/module"
	"github.com/hollis-labs/substrate/modules/auth"
	substratemcp "github.com/hollis-labs/substrate/modules/mcp"

	"github.com/hollis-labs/sigil/internal/mcp"
)

var errModulesNotStarted = errors.New("the auth and mcp modules must be started")

// NewMCPCmd creates the mcp command with subcommands.
func NewMCPCmd(app *module.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server commands",
		Long:  "Model Context Protocol server for AI agent integration.",
	}
	cmd.AddCommand(newMCPServeCmd(app))
	return cmd
}

func newMCPServeCmd(app *module.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start MCP server (stdio)",
		Long: `Start the MCP JSON-RPC 2.0 server on stdin/stdout for AI agent communication.

The server comes from the Substrate mcp module. With auth off (the shipped setting)
no token is needed; if auth.mode is changed, the module requires the token named by
mcp.stdio_token_file and the scope sigil.mcp.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sigilDir, _ := cmd.Flags().GetString("sigil-dir")
			base := substratemcp.Server(app)
			table := auth.ScopeTable(app)
			if base == nil || table == nil {
				return errModulesNotStarted
			}
			mcp.RegisterScopes(table)
			server := mcp.NewServer(base, sigilDir)
			server.UseAuthorizer(func(ctx context.Context, tool string) error {
				return substratemcp.Authorize(ctx, app, tool)
			})
			mcp.RegisterAllTools(server)
			mcp.RegisterAllResources(server)
			mcp.RegisterAllPrompts(server)
			return substratemcp.ServeStdio(cmd.Context(), app, substratemcp.StdioOptions{})
		},
	}
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	return cmd
}

package cli

import (
	"github.com/chrispian/sigil/internal/mcp"
	"github.com/spf13/cobra"
)

// NewMCPCmd creates the mcp command with subcommands.
func NewMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server commands",
		Long:  "Model Context Protocol server for AI agent integration.",
	}
	cmd.AddCommand(newMCPServeCmd())
	return cmd
}

func newMCPServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start MCP server (stdio)",
		Long:  "Start the MCP JSON-RPC 2.0 server on stdin/stdout for AI agent communication.",
		RunE:  runMCPServe,
	}
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	return cmd
}

func runMCPServe(cmd *cobra.Command, args []string) error {
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")
	server := mcp.NewServer(sigilDir)
	mcp.RegisterAllTools(server)
	mcp.RegisterAllResources(server)
	mcp.RegisterAllPrompts(server)
	return server.Run(cmd.Context())
}

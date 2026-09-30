package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hollis-labs/sigil/internal/server"
	"github.com/spf13/cobra"
)

// NewServeCmd creates the serve command.
func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the live development server",
		Long: `Start an HTTP server that renders page previews in real-time.

The server watches .sigil/ for changes and reloads the browser automatically.
It also exposes a JSON/REST API under /api/* (pages, components, datasources,
themes, project config; CRUD for pages and datasources) for browser SPAs.

Examples:
  sigil serve                    Start on default port 3210
  sigil serve --port 8080        Start on custom port
  sigil serve --sigil-dir path   Use custom .sigil directory`,
		RunE: runServe,
	}
	cmd.Flags().Int("port", 3210, "HTTP server port")
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	return cmd
}

func runServe(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")

	// Verify .sigil directory exists
	if _, err := os.Stat(sigilDir); os.IsNotExist(err) {
		return fmt.Errorf("%s not found — run 'sigil init' first", sigilDir)
	}

	srv := server.New(sigilDir, port)

	// Start file watcher
	watcher := server.NewWatcher(sigilDir, 500*time.Millisecond, func() {
		srv.NotifyReload()
	})
	watcher.Start()

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "\n  %s Sigil dev server running\n\n", green("●"))
	fmt.Fprintf(out, "  %s  http://localhost:%d\n", bold("Local:"), port)
	fmt.Fprintf(out, "  %s  %s\n\n", bold("Watching:"), sigilDir)
	fmt.Fprintf(out, "  %s\n\n", dim("Press Ctrl+C to stop"))

	// Handle shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-done:
		fmt.Fprintf(out, "\n  Shutting down...\n")
		watcher.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-errCh:
		watcher.Stop()
		return err
	}
}

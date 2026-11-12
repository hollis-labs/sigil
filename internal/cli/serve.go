package cli

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/hollis-labs/substrate/module"
	"github.com/hollis-labs/substrate/modules/auth"
	"github.com/hollis-labs/substrate/modules/httphost"

	"github.com/hollis-labs/sigil/internal/hostguard"
	"github.com/hollis-labs/sigil/internal/server"
)

// NewServeCmd creates the serve command. The listener is the Substrate httphost module's:
// its address is httphost.addr (127.0.0.1:3210 unless the SUBSTRATE_CONFIG file says
// otherwise), or that host with --port. Sigil refuses any address that is not loopback
// (see loopbackAddr), which is stricter than the framework's check in mode off.
func NewServeCmd(app *module.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the live development server",
		Long: `Start an HTTP server that renders page previews in real-time.

The server watches .sigil/ for changes and reloads the browser automatically.
It also exposes a JSON/REST API under /api/* (pages, components, datasources,
themes, project config; CRUD for pages and datasources) for browser SPAs.

It listens on httphost.addr (default 127.0.0.1:3210), set in the YAML file named
by SUBSTRATE_CONFIG, or on that host with --port. Authentication is off, so only a
loopback address (127.0.0.0/8, ::1 or localhost) is accepted. Health is GET /healthz.

Examples:
  sigil serve                    Start on the configured address (port 3210)
  sigil serve --port 8080        Start on a custom port
  sigil serve --sigil-dir path   Use custom .sigil directory`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error { return runServe(cmd, app) },
	}
	cmd.Flags().Int("port", 0, "HTTP server port (default: the port of httphost.addr, 3210)")
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	return cmd
}

func runServe(cmd *cobra.Command, app *module.App) error {
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")
	port, _ := cmd.Flags().GetInt("port")
	addr, err := loopbackAddr(app, port)
	if err != nil {
		return err
	}

	// Verify .sigil directory exists
	if _, statErr := os.Stat(sigilDir); os.IsNotExist(statErr) {
		return fmt.Errorf("%s not found — run 'sigil init' first", sigilDir)
	}

	srv := server.New(sigilDir)
	api := server.NewAPI(sigilDir)
	table := auth.ScopeTable(app)
	if table == nil {
		return fmt.Errorf("the auth module must be started")
	}
	server.RegisterScopes(table, srv.Previews(), api.Routes())

	// The API writes project files, so its writes need X-Requested-With (see
	// server.RequireRequestedWith). The preview pages and the reload stream are
	// mounted outside /api/, each with the scope its route declares.
	routes := http.NewServeMux()
	routes.Handle("/api/", server.RequireRequestedWith(api.Handler()))
	var mounts []httphost.Mount
	for _, r := range srv.Previews() {
		mounts = append(mounts, httphost.Mount{Pattern: r.Method + " " + r.Pattern, Handler: r.Handler})
	}
	handler, err := httphost.Handler(app, httphost.Options{Routes: routes, Mounts: mounts})
	if err != nil {
		return err
	}
	// Outermost: refuse a Host the server is not meant to answer (DNS rebinding).
	handler = hostguard.Handler(nil, handler)

	ctx := cmd.Context()
	watcher := server.NewWatcher(sigilDir, 500*time.Millisecond, srv.NotifyReload)
	watcher.Start()
	defer watcher.Stop()
	go func() {
		<-ctx.Done()
		srv.StopStreams()
	}()

	out := cmd.OutOrStdout()
	err = httphost.Serve(ctx, app, handler, httphost.ServeOptions{
		Addr: addr,
		Ready: func(addr net.Addr) {
			fmt.Fprintf(out, "\n  %s Sigil dev server running\n\n", green("●"))
			fmt.Fprintf(out, "  %s  http://%s\n", bold("Local:"), addr)
			fmt.Fprintf(out, "  %s  %s\n\n", bold("Watching:"), sigilDir)
			fmt.Fprintf(out, "  %s\n\n", dim("Press Ctrl+C to stop"))
		},
	})
	if ctx.Err() != nil && err == nil {
		fmt.Fprintf(out, "\n  Shut down.\n")
	}
	return err
}

// loopbackAddr returns the address serve listens on: httphost.addr, with its port replaced by
// port when that is non-zero. It refuses anything but a loopback host. The framework's check
// in mode off also accepts tailnet ranges and *.ts.net names; this server has no
// authentication and writes project files, so only 127.0.0.0/8, ::1 and localhost are allowed.
func loopbackAddr(app *module.App, port int) (string, error) {
	cfg, _ := app.Config["httphost"].(map[string]any)
	addr, _ := cfg["addr"].(string)
	if addr == "" {
		addr = "127.0.0.1:3210"
	}
	host, p, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("sigil serve: bad address %q: %w", addr, err)
	}
	if port != 0 {
		if port < 0 || port > 65535 {
			return "", fmt.Errorf("sigil serve: --port %d is out of range", port)
		}
		p = strconv.Itoa(port)
	}
	if ip := net.ParseIP(host); (ip == nil || !ip.IsLoopback()) && host != "localhost" {
		return "", fmt.Errorf("sigil serve: %q is not a loopback address; the server has no authentication, so it listens on 127.0.0.1, ::1 or localhost only", host)
	}
	return net.JoinHostPort(host, p), nil
}

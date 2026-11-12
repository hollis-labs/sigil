// Command sigil is the Sigil UI compiler, built on Substrate. The shell (config, logging,
// health, HTTP host, auth guard, MCP server, CLI skeleton with version, config and health)
// comes from the manifest; the commands in internal/cli are the app. There is no database
// or identity: serve runs with auth off, on loopback only.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hollis-labs/substrate/boot"
	substratecli "github.com/hollis-labs/substrate/modules/cli"

	"github.com/hollis-labs/sigil/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "sigil:", err)
		stop()
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	started, err := boot.Start(ctx, substrateApp, substrateModules, substrateConfigJSON)
	if err != nil {
		return err
	}
	defer func() { _ = started.Stop(context.Background()) }()

	root := substratecli.Root(started.App)
	cli.Install(root, started.App)
	return root.ExecuteContext(ctx)
}

package main

import (
	"os"

	"github.com/hollis-labs/sigil/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

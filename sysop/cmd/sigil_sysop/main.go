// Command sigil_sysop serves the Sigil Sysop UI: a single binary that
// exposes Sigil's JSON API and the embedded React admin SPA, same-origin.
//
// The API surface is Sigil's own internal/server.API — the same handler
// `sigil serve` mounts — so the Sysop UI reads and mutates a real Sigil
// project (pages, components, datasources, themes) with no API duplication.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/hollis-labs/sigil/internal/server"
	"github.com/hollis-labs/sigil/sysop/internal/webui"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "TCP address the server listens on (use :8080 to listen on all interfaces)")
	sigilDir := flag.String("sigil-dir", ".sigil",
		"path to the Sigil project's .sigil directory to administer")
	flag.Parse()

	if info, err := os.Stat(*sigilDir); err != nil || !info.IsDir() {
		log.Printf("warning: sigil-dir %q is not a readable directory — "+
			"the API will report an empty project until it exists", *sigilDir)
	}

	mux := http.NewServeMux()

	// Sigil's JSON API — the same handler `sigil serve` mounts. Serving it
	// here lets the SPA read/mutate the project same-origin.
	server.NewAPI(*sigilDir).Register(mux)

	// The Sysop UI SPA — served from the embedded frontend build by go-webui.
	webui.Mount(mux)

	// Send the bare root to the UI so visiting the host is enough.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, webui.BasePath+"/", http.StatusFound)
	})

	log.Printf("Sigil Sysop UI: administering %q", *sigilDir)
	log.Printf("listening on %s — UI at %s/", *addr, webui.BasePath)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

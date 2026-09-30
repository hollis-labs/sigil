module github.com/hollis-labs/sigil/sysop

go 1.26.1

require (
	github.com/hollis-labs/go-webui v0.1.0
	github.com/hollis-labs/sigil v0.0.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

// The Sysop UI binary serves Sigil's own JSON API (internal/server). It
// lives inside the Sigil repo as a nested module, so the parent module is
// wired through a filesystem replace rather than a published version.
replace github.com/hollis-labs/sigil => ../

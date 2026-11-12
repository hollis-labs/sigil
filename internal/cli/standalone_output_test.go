package cli

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGeneratedOutputIsFreeOfSubstrate holds the rule that Substrate applies to the
// sigil program and never to what it emits: generated code must stand alone.
func TestGeneratedOutputIsFreeOfSubstrate(t *testing.T) {
	for _, target := range []string{"go-templ", "react-shadcn"} {
		t.Run(target, func(t *testing.T) {
			t.Chdir(t.TempDir())
			run := func(args ...string) {
				t.Helper()
				cmd := NewRootCmd()
				cmd.SetOut(&bytes.Buffer{})
				cmd.SetErr(&bytes.Buffer{})
				cmd.SetArgs(args)
				if err := cmd.Execute(); err != nil {
					t.Fatalf("sigil %s: %v", strings.Join(args, " "), err)
				}
			}
			run("init")
			run("new", "page", "dashboard")
			out := filepath.Join(t.TempDir(), "out")
			run("generate", "--target", target, "--output", out)

			files := 0
			err := filepath.WalkDir(out, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				files++
				b, err := os.ReadFile(path) //nolint:gosec // path comes from walking the test's own output directory
				if err != nil {
					return err
				}
				if bytes.Contains(bytes.ToLower(b), []byte("hollis-labs/substrate")) {
					t.Errorf("%s mentions hollis-labs/substrate", path)
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if files == 0 {
				t.Fatalf("target %s generated no files", target)
			}
		})
	}
}

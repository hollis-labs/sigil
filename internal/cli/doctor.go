package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// NewDoctorCmd creates the doctor command.
func NewDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and project health",
		Long: `Verify that your environment is set up correctly for Sigil development.

Checks:
  - Go version
  - templ installed (for go-templ renderer)
  - .sigil directory exists and is valid
  - Page configs parse correctly

Examples:
  sigil doctor
  sigil doctor --sigil-dir path/to/.sigil`,
		RunE: runDoctor,
	}
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	return cmd
}

type checkResult struct {
	Name   string
	Status string // "pass", "warn", "fail", "info"
	Detail string
}

func runDoctor(cmd *cobra.Command, args []string) error {
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")
	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "\n  Sigil Doctor\n\n")

	var checks []checkResult

	// Go version
	goVer := runtime.Version()
	checks = append(checks, checkResult{
		Name: "Go", Status: "pass", Detail: goVer,
	})

	// OS/Arch
	checks = append(checks, checkResult{
		Name: "Platform", Status: "info", Detail: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	})

	// templ installed
	templPath, err := exec.LookPath("templ")
	if err != nil {
		checks = append(checks, checkResult{
			Name: "templ", Status: "warn", Detail: "not found (needed for go-templ renderer)",
		})
	} else {
		checks = append(checks, checkResult{
			Name: "templ", Status: "pass", Detail: templPath,
		})
	}

	// Node/npm (for react-shadcn)
	nodePath, err := exec.LookPath("node")
	if err != nil {
		checks = append(checks, checkResult{
			Name: "Node.js", Status: "warn", Detail: "not found (needed for react-shadcn renderer)",
		})
	} else {
		checks = append(checks, checkResult{
			Name: "Node.js", Status: "pass", Detail: nodePath,
		})
	}

	// .sigil directory
	if info, err := os.Stat(sigilDir); err != nil {
		checks = append(checks, checkResult{
			Name: ".sigil dir", Status: "fail", Detail: fmt.Sprintf("%s not found", sigilDir),
		})
	} else if !info.IsDir() {
		checks = append(checks, checkResult{
			Name: ".sigil dir", Status: "fail", Detail: fmt.Sprintf("%s is not a directory", sigilDir),
		})
	} else {
		checks = append(checks, checkResult{
			Name: ".sigil dir", Status: "pass", Detail: sigilDir,
		})

		// Check subdirectories
		for _, sub := range []string{"pages", "themes"} {
			subPath := filepath.Join(sigilDir, sub)
			if _, err := os.Stat(subPath); err != nil {
				checks = append(checks, checkResult{
					Name: sub + "/", Status: "warn", Detail: "directory missing",
				})
			} else {
				checks = append(checks, checkResult{
					Name: sub + "/", Status: "pass", Detail: "exists",
				})
			}
		}

		// sigil.yaml
		configPath := filepath.Join(sigilDir, "sigil.yaml")
		if _, err := os.Stat(configPath); err != nil {
			checks = append(checks, checkResult{
				Name: "sigil.yaml", Status: "fail", Detail: "not found",
			})
		} else {
			checks = append(checks, checkResult{
				Name: "sigil.yaml", Status: "pass", Detail: "found",
			})
		}

		// Count pages
		pagePattern := filepath.Join(sigilDir, "pages", "*.yaml")
		pageMatches, _ := filepath.Glob(pagePattern)
		checks = append(checks, checkResult{
			Name: "Pages", Status: "info", Detail: fmt.Sprintf("%d page config(s)", len(pageMatches)),
		})
	}

	// Display results
	for _, c := range checks {
		var icon string
		switch c.Status {
		case "pass":
			icon = green("✓")
		case "warn":
			icon = yellow("!")
		case "fail":
			icon = red("✗")
		case "info":
			icon = cyan("i")
		}
		fmt.Fprintf(out, "  %s %-14s %s\n", icon, c.Name, dim(c.Detail))
	}
	fmt.Fprintln(out)

	return nil
}

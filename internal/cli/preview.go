package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hollis-labs/sigil/internal/config"
	"github.com/hollis-labs/sigil/internal/renderer"
	"github.com/hollis-labs/sigil/internal/renderer/preview"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewPreviewCmd creates the preview command.
func NewPreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview <page-id>",
		Short: "Preview a page config in the browser",
		Long: `Generate a standalone HTML file from a page config and open it in the default browser.

Examples:
  sigil preview dashboard             Preview dashboard page
  sigil preview users --theme light   Preview with specific theme
  sigil preview modal-form --no-open  Generate HTML without opening browser
  sigil preview page -o preview.html  Write to specific path`,
		Args: cobra.ExactArgs(1),
		RunE: runPreview,
	}
	cmd.Flags().String("sigil-dir", ".sigil", "Path to .sigil directory")
	cmd.Flags().String("theme", "", "Theme name override (default: from sigil.yaml)")
	cmd.Flags().StringP("output", "o", "", "Write HTML to this path instead of a temp file")
	cmd.Flags().Bool("no-open", false, "Don't open the browser, just write the file")
	return cmd
}

func runPreview(cmd *cobra.Command, args []string) error {
	pageID := args[0]
	sigilDir, _ := cmd.Flags().GetString("sigil-dir")
	themeName, _ := cmd.Flags().GetString("theme")
	outputPath, _ := cmd.Flags().GetString("output")
	noOpen, _ := cmd.Flags().GetBool("no-open")

	// Load page
	pagePath := filepath.Join(sigilDir, "pages", pageID+".yaml")
	page, err := config.ParseFile(pagePath)
	if err != nil {
		return fmt.Errorf("loading page %q: %w", pageID, err)
	}

	// Load theme (optional)
	theme := loadPreviewTheme(sigilDir, themeName)

	// Load datasources (optional)
	dataSources := loadPreviewDataSources(sigilDir)

	cfg := &preview.PreviewConfig{
		Page:        page,
		Theme:       theme,
		DataSources: dataSources,
	}

	html, err := preview.GenerateHTML(cfg)
	if err != nil {
		return fmt.Errorf("generating preview: %w", err)
	}

	// Determine output path
	if outputPath == "" {
		tmpDir := os.TempDir()
		outputPath = filepath.Join(tmpDir, fmt.Sprintf("sigil-preview-%s.html", pageID))
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.WriteFile(outputPath, html, 0644); err != nil {
		return fmt.Errorf("writing preview: %w", err)
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Preview written to %s\n", outputPath)

	if !noOpen {
		if err := openBrowser(outputPath); err != nil {
			fmt.Fprintf(out, "Could not open browser: %v\n", err)
			fmt.Fprintf(out, "Open the file manually: %s\n", outputPath)
		}
	}

	return nil
}

func loadPreviewTheme(sigilDir, themeName string) *renderer.ThemeConfig {
	// Try project config for default theme
	if themeName == "" {
		projPath := filepath.Join(sigilDir, "sigil.yaml")
		data, err := os.ReadFile(projPath)
		if err == nil {
			var proj config.ProjectConfig
			if yaml.Unmarshal(data, &proj) == nil && proj.Defaults.Theme != "" {
				themeName = proj.Defaults.Theme
			}
		}
	}
	if themeName == "" {
		themeName = "default"
	}

	themesDir := filepath.Join(sigilDir, "themes")
	theme, err := renderer.LoadTheme(themesDir, themeName)
	if err != nil {
		return nil
	}
	return theme
}

func loadPreviewDataSources(sigilDir string) map[string]*renderer.DataSourceManifest {
	dsDir := filepath.Join(sigilDir, "datasources")
	matches, err := filepath.Glob(filepath.Join(dsDir, "*.yaml"))
	if err != nil || len(matches) == 0 {
		return nil
	}

	result := map[string]*renderer.DataSourceManifest{}
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var ds renderer.DataSourceManifest
		if yaml.Unmarshal(data, &ds) != nil {
			continue
		}
		alias := ds.Alias
		if alias == "" {
			alias = filepath.Base(path)
			alias = alias[:len(alias)-len(filepath.Ext(alias))]
		}
		result[alias] = &ds
	}
	return result
}

func openBrowser(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	url := "file://" + absPath

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return fmt.Errorf("unsupported platform %q", runtime.GOOS)
	}
}

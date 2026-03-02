// Package gotempl implements the Go/Templ+HTMX renderer for Sigil.
package gotempl

import (
	"github.com/chrispian/sigil/internal/renderer"
)

const rendererName = "go-templ"

// GoTemplRenderer generates Go Templ files, HTMX interactions, and CSS.
type GoTemplRenderer struct{}

func init() {
	renderer.RegisterRenderer(&GoTemplRenderer{})
}

// Name returns the renderer identifier.
func (r *GoTemplRenderer) Name() string {
	return rendererName
}

// Render generates .templ files for a single page config.
func (r *GoTemplRenderer) Render(ctx *renderer.RenderContext) ([]renderer.OutputFile, error) {
	return renderPage(ctx)
}

// RenderTheme generates theme CSS files.
func (r *GoTemplRenderer) RenderTheme(theme *renderer.ThemeConfig) ([]renderer.OutputFile, error) {
	return renderThemeCSS(theme)
}

// RenderDataSourceStubs generates Go handler stubs.
func (r *GoTemplRenderer) RenderDataSourceStubs(ds *renderer.DataSourceManifest) ([]renderer.OutputFile, error) {
	return renderDataSourceStubs(ds)
}

// SharedComponents returns shared component .templ files.
func (r *GoTemplRenderer) SharedComponents(usedTypes []string) ([]renderer.OutputFile, error) {
	return renderSharedComponents(usedTypes)
}

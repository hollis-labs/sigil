// Package reactshadcn implements the React/shadcn+Tailwind renderer for Sigil.
package reactshadcn

import (
	"github.com/chrispian/sigil/internal/renderer"
)

const rendererName = "react-shadcn"

// ReactShadcnRenderer generates React TSX files with shadcn/ui components.
type ReactShadcnRenderer struct{}

func init() {
	renderer.RegisterRenderer(&ReactShadcnRenderer{})
}

// Name returns the renderer identifier.
func (r *ReactShadcnRenderer) Name() string {
	return rendererName
}

// Render generates .tsx files for a single page config.
func (r *ReactShadcnRenderer) Render(ctx *renderer.RenderContext) ([]renderer.OutputFile, error) {
	return renderPage(ctx)
}

// RenderTheme generates Tailwind config and CSS files.
func (r *ReactShadcnRenderer) RenderTheme(theme *renderer.ThemeConfig) ([]renderer.OutputFile, error) {
	return renderTheme(theme)
}

// RenderDataSourceStubs generates TypeScript API hooks for datasources.
func (r *ReactShadcnRenderer) RenderDataSourceStubs(ds *renderer.DataSourceManifest) ([]renderer.OutputFile, error) {
	return renderDataSourceHooks(ds)
}

// SharedComponents returns shared component and utility files.
func (r *ReactShadcnRenderer) SharedComponents(usedTypes []string) ([]renderer.OutputFile, error) {
	return renderSharedComponents(usedTypes)
}

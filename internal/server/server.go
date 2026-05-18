// Package server implements the Sigil live development server.
package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chrispian/sigil/internal/config"
	"github.com/chrispian/sigil/internal/renderer"
	"github.com/chrispian/sigil/internal/renderer/preview"
	"gopkg.in/yaml.v3"
)

// Server is the Sigil live development server.
type Server struct {
	SigilDir string
	Port     int

	mu      sync.RWMutex
	clients map[chan struct{}]struct{}
	mux     *http.ServeMux
	srv     *http.Server
}

// New creates a new dev server.
func New(sigilDir string, port int) *Server {
	s := &Server{
		SigilDir: sigilDir,
		Port:     port,
		clients:  make(map[chan struct{}]struct{}),
	}
	s.mux = s.buildMux()
	return s
}

func (s *Server) buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleIndex)
	mux.HandleFunc("GET /pages/{id}", s.handlePage)
	mux.HandleFunc("GET /assets/theme.css", s.handleThemeCSS)
	mux.HandleFunc("GET /events", s.handleSSE)

	// Mount the JSON/REST API. The /api/* routes coexist with the HTML
	// preview routes above; a browser SPA consumes the API while the
	// preview pages keep working unchanged.
	NewAPI(s.SigilDir).Register(mux)

	return mux
}

// Handler returns the HTTP handler (for testing).
func (s *Server) Handler() http.Handler {
	return s.mux
}

// ListenAndServe starts the server.
func (s *Server) ListenAndServe() error {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: s.mux,
	}
	return s.srv.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv != nil {
		return s.srv.Shutdown(ctx)
	}
	return nil
}

// NotifyReload sends a reload event to all connected SSE clients.
func (s *Server) NotifyReload() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for ch := range s.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// handlePage renders a page preview on-the-fly.
func (s *Server) handlePage(w http.ResponseWriter, r *http.Request) {
	pageID := r.PathValue("id")
	if pageID == "" {
		http.Error(w, "page ID required", http.StatusBadRequest)
		return
	}

	page, err := s.loadPage(pageID)
	if err != nil {
		http.Error(w, fmt.Sprintf("page %q not found: %v", pageID, err), http.StatusNotFound)
		return
	}

	theme := s.loadTheme()
	dataSources := s.loadDataSources()

	cfg := &preview.PreviewConfig{
		Page:        page,
		Theme:       theme,
		DataSources: dataSources,
	}

	html, err := preview.GenerateHTML(cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("render error: %v", err), http.StatusInternalServerError)
		return
	}

	// Inject reload script and navigation header before </body>
	html = injectDevExtras(html, pageID, s.SigilDir)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

// handleThemeCSS serves the theme as a CSS file.
func (s *Server) handleThemeCSS(w http.ResponseWriter, r *http.Request) {
	theme := s.loadTheme()
	if theme == nil {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write([]byte("/* no theme found */\n"))
		return
	}

	var buf bytes.Buffer
	buf.WriteString(":root {\n")
	categories := []string{"colors", "typography", "spacing", "radius", "shadows"}
	for _, cat := range categories {
		tokens, ok := theme.Tokens[cat]
		if !ok {
			continue
		}
		for key, val := range tokens {
			buf.WriteString(fmt.Sprintf("  --sigil-%s: %v;\n", key, val))
		}
	}
	buf.WriteString("}\n")

	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(buf.Bytes())
}

// handleSSE sends server-sent events for live reload.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan struct{}, 1)
	s.mu.Lock()
	s.clients[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, ch)
		s.mu.Unlock()
	}()

	// Send initial connected event
	fmt.Fprintf(w, "event: connected\ndata: ok\n\n")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ch:
			fmt.Fprintf(w, "event: reload\ndata: %d\n\n", time.Now().Unix())
			flusher.Flush()
		}
	}
}

// Data loading helpers

func (s *Server) loadPage(id string) (*config.Page, error) {
	path := filepath.Join(s.SigilDir, "pages", id+".yaml")
	return config.ParseFile(path)
}

func (s *Server) loadTheme() *renderer.ThemeConfig {
	// Try project config for default theme name
	themeName := "default"
	projPath := filepath.Join(s.SigilDir, "sigil.yaml")
	data, err := os.ReadFile(projPath)
	if err == nil {
		var proj config.ProjectConfig
		if yaml.Unmarshal(data, &proj) == nil && proj.Defaults.Theme != "" {
			themeName = proj.Defaults.Theme
		}
	}

	path := filepath.Join(s.SigilDir, "themes", themeName+".yaml")
	data, err = os.ReadFile(path)
	if err != nil {
		return nil
	}
	var theme renderer.ThemeConfig
	if yaml.Unmarshal(data, &theme) != nil {
		return nil
	}
	return &theme
}

func (s *Server) loadDataSources() map[string]*renderer.DataSourceManifest {
	dsDir := filepath.Join(s.SigilDir, "datasources")
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
			alias = strings.TrimSuffix(filepath.Base(path), ".yaml")
		}
		result[alias] = &ds
	}
	return result
}

// ListPages returns metadata for all pages in the sigil dir.
func (s *Server) ListPages() []PageInfo {
	pagesDir := filepath.Join(s.SigilDir, "pages")
	matches, _ := filepath.Glob(filepath.Join(pagesDir, "*.yaml"))

	var pages []PageInfo
	for _, path := range matches {
		page, err := config.ParseFile(path)
		if err != nil {
			pages = append(pages, PageInfo{
				ID:    strings.TrimSuffix(filepath.Base(path), ".yaml"),
				Error: err.Error(),
			})
			continue
		}

		info := PageInfo{
			ID:             page.ID,
			Title:          page.Title,
			Overlay:        page.Overlay,
			Module:         page.Module,
			ComponentCount: countComponents(&page.Layout),
		}

		pages = append(pages, info)
	}
	return pages
}

// PageInfo holds metadata about a page for the index.
type PageInfo struct {
	ID             string
	Title          string
	Overlay        string
	Module         string
	ComponentCount int
	Error          string
}

func countComponents(c *config.Component) int {
	count := 1
	for i := range c.Children {
		count += countComponents(&c.Children[i])
	}
	return count
}

// injectDevExtras adds the reload script and a nav bar to the preview HTML.
func injectDevExtras(html []byte, currentPage, sigilDir string) []byte {
	reloadScript := `<script>
const es = new EventSource('/events');
es.addEventListener('reload', () => location.reload());
es.addEventListener('error', () => setTimeout(() => location.reload(), 2000));
</script>`

	navBar := fmt.Sprintf(`<div style="position:fixed;top:0;left:0;right:0;z-index:9999;background:rgba(0,0,0,0.85);backdrop-filter:blur(8px);padding:8px 16px;display:flex;align-items:center;gap:16px;font-family:system-ui;font-size:13px;color:#a1a1aa;border-bottom:1px solid rgba(255,255,255,0.1)">
  <a href="/" style="color:#818cf8;text-decoration:none;font-weight:600">Sigil</a>
  <span style="color:#3f3f46">|</span>
  <span style="color:#e4e4e7">%s</span>
  <span style="flex:1"></span>
  <span style="font-size:11px;color:#52525b">Live</span>
  <span style="width:6px;height:6px;border-radius:50%%;background:#22c55e;display:inline-block"></span>
</div>
<div style="height:40px"></div>`, currentPage)

	content := string(html)

	// Inject nav bar after <body...>
	bodyIdx := strings.Index(content, "<body")
	if bodyIdx >= 0 {
		closeTag := strings.Index(content[bodyIdx:], ">")
		if closeTag >= 0 {
			insertAt := bodyIdx + closeTag + 1
			content = content[:insertAt] + "\n" + navBar + "\n" + content[insertAt:]
		}
	}

	// Inject script before </body>
	content = strings.Replace(content, "</body>", reloadScript+"\n</body>", 1)

	return []byte(content)
}

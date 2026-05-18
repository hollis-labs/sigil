package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chrispian/sigil/internal/components"
	"github.com/chrispian/sigil/internal/config"
	"gopkg.in/yaml.v3"
)

// maxBodyBytes caps request bodies for the JSON API write endpoints.
const maxBodyBytes = 4 << 20 // 4 MiB

// API is a JSON/REST handler over a Sigil project's data layer. It exposes the
// same pages / components / datasources / themes / project data that the MCP
// server serves, but over plain HTTP so a browser SPA can read and mutate it.
//
// The routes are registered under /api/. They reuse the existing data layer
// (internal/config, internal/components) rather than duplicating any business
// logic — page create/update run the same parse + validate path as the MCP
// tools sigil_create_page / sigil_update_page.
//
// An API value can be mounted on the dev server's mux (see Server.buildMux) or
// served standalone via Handler() so a Sysop UI binary can serve a static SPA
// and this API same-origin.
type API struct {
	// SigilDir is the path to the .sigil project directory.
	SigilDir string
}

// NewAPI creates a JSON API handler rooted at the given .sigil directory.
func NewAPI(sigilDir string) *API {
	return &API{SigilDir: sigilDir}
}

// Register adds the /api/* routes to mux. Use this to mount the API on an
// existing server (the dev server) or alongside a static SPA handler.
func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", a.handleHealth)

	mux.HandleFunc("GET /api/pages", a.handleListPages)
	mux.HandleFunc("POST /api/pages", a.handleCreatePage)
	mux.HandleFunc("GET /api/pages/{id}", a.handleGetPage)
	mux.HandleFunc("PUT /api/pages/{id}", a.handleUpdatePage)

	mux.HandleFunc("GET /api/components", a.handleListComponents)
	mux.HandleFunc("GET /api/components/{type}", a.handleGetComponent)

	mux.HandleFunc("GET /api/datasources", a.handleListDataSources)
	mux.HandleFunc("POST /api/datasources", a.handleCreateDataSource)
	mux.HandleFunc("GET /api/datasources/{alias}", a.handleGetDataSource)

	mux.HandleFunc("GET /api/themes", a.handleListThemes)
	mux.HandleFunc("GET /api/themes/{name}", a.handleGetTheme)

	mux.HandleFunc("GET /api/project", a.handleProject)
}

// Handler returns a standalone http.Handler serving only the /api/* routes.
// A Sysop UI binary can mount this next to a static SPA handler to serve both
// same-origin.
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	a.Register(mux)
	return mux
}

// --- JSON helpers ---

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

func writeAPIError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// safeName rejects path segments that could escape the project directory.
// Route patterns already constrain {id}/{alias}/{name} to a single path
// segment; this is defense in depth against "." and ".." style values.
func safeName(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	return !strings.ContainsAny(name, "/\\")
}

// readProjectFile reads a file from the project directory. Callers pass
// either a fixed filename or a path built from a safeName-validated segment
// or a filepath.Glob match within the project dir.
func readProjectFile(path string) ([]byte, error) {
	return os.ReadFile(path) //nolint:gosec // path is a fixed name, safeName-validated segment, or glob match within the project dir
}

// writeProjectFile writes a project file using the same permissions
// (0644 file / 0755 dir) as `sigil init` and the MCP write tools, so files
// created through the API stay consistent with the rest of the project.
func writeProjectFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { //nolint:gosec // project dirs are user-editable, matches sigil init
		return err
	}
	return os.WriteFile(path, data, 0o644) //nolint:gosec // project config files are user-editable, matches sigil init
}

// isYAMLContentType reports whether ct denotes a YAML request body.
func isYAMLContentType(ct string) bool {
	ct = strings.ToLower(strings.TrimSpace(strings.SplitN(ct, ";", 2)[0]))
	switch ct {
	case "text/yaml", "application/yaml", "application/x-yaml", "text/x-yaml":
		return true
	default:
		return false
	}
}

// readYAMLAsJSON loads a YAML file and decodes it into a generic value that
// marshals cleanly to JSON (yaml.v3 produces map[string]interface{}).
func readYAMLAsJSON(path string) (interface{}, error) {
	data, err := readProjectFile(path)
	if err != nil {
		return nil, err
	}
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// --- health ---

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "sigil",
	})
}

// --- pages ---

type pageSummary struct {
	ID      string `json:"id"`
	Title   string `json:"title,omitempty"`
	Overlay string `json:"overlay,omitempty"`
	Module  string `json:"module,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (a *API) pagePath(id string) string {
	return filepath.Join(a.SigilDir, "pages", id+".yaml")
}

func (a *API) handleListPages(w http.ResponseWriter, r *http.Request) {
	matches, _ := filepath.Glob(filepath.Join(a.SigilDir, "pages", "*.yaml"))
	sort.Strings(matches)

	pages := []pageSummary{}
	for _, path := range matches {
		page, err := config.ParseFile(path)
		if err != nil {
			pages = append(pages, pageSummary{
				ID:    strings.TrimSuffix(filepath.Base(path), ".yaml"),
				Error: err.Error(),
			})
			continue
		}
		pages = append(pages, pageSummary{
			ID:      page.ID,
			Title:   page.Title,
			Overlay: page.Overlay,
			Module:  page.Module,
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"pages": pages})
}

func (a *API) handleGetPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !safeName(id) {
		writeAPIError(w, http.StatusBadRequest, "invalid page id")
		return
	}
	page, err := config.ParseFile(a.pagePath(id))
	if err != nil {
		writeAPIError(w, http.StatusNotFound, fmt.Sprintf("page %q not found", id))
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (a *API) handleCreatePage(w http.ResponseWriter, r *http.Request) {
	page, raw, err := decodePageBody(r)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !safeName(page.ID) {
		writeAPIError(w, http.StatusUnprocessableEntity, "page id is required")
		return
	}
	path := a.pagePath(page.ID)
	if _, err := os.Stat(path); err == nil {
		writeAPIError(w, http.StatusConflict, fmt.Sprintf("page %q already exists", page.ID))
		return
	}
	a.writePage(w, page, raw, path, http.StatusCreated)
}

func (a *API) handleUpdatePage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !safeName(id) {
		writeAPIError(w, http.StatusBadRequest, "invalid page id")
		return
	}
	path := a.pagePath(id)
	if _, err := os.Stat(path); err != nil {
		writeAPIError(w, http.StatusNotFound, fmt.Sprintf("page %q not found", id))
		return
	}

	page, raw, err := decodePageBody(r)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if page.ID != "" && page.ID != id {
		writeAPIError(w, http.StatusUnprocessableEntity,
			fmt.Sprintf("body id %q does not match URL id %q", page.ID, id))
		return
	}
	// The URL is authoritative for the id; re-marshal so the persisted YAML
	// matches when the body omitted it.
	if page.ID != id {
		page.ID = id
		if raw, err = config.MarshalYAML(page); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	a.writePage(w, page, raw, path, http.StatusOK)
}

// writePage validates page and, if valid, persists raw to path. It mirrors the
// parse + validate path of the MCP sigil_create_page / sigil_update_page tools.
func (a *API) writePage(w http.ResponseWriter, page *config.Page, raw []byte, path string, okStatus int) {
	registry := components.NewDefaultRegistry()
	result := config.Validate(page, registry)
	if !result.Valid {
		errs := make([]string, 0, len(result.Errors))
		for _, e := range result.Errors {
			errs = append(errs, e.String())
		}
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error":   "validation failed",
			"details": errs,
		})
		return
	}

	if err := writeProjectFile(path, raw); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}

	warnings := make([]string, 0, len(result.Warnings))
	for _, wn := range result.Warnings {
		warnings = append(warnings, wn.String())
	}
	writeJSON(w, okStatus, map[string]interface{}{
		"id":       page.ID,
		"valid":    true,
		"warnings": warnings,
	})
}

// decodePageBody reads a page from the request. The body may be a JSON
// config.Page object (the default) or a raw YAML document when the request
// carries a YAML Content-Type. It returns the parsed page and the YAML bytes
// to persist.
func decodePageBody(r *http.Request) (*config.Page, []byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("reading request body: %w", err)
	}
	if len(body) == 0 {
		return nil, nil, fmt.Errorf("empty request body")
	}

	if isYAMLContentType(r.Header.Get("Content-Type")) {
		page, parseErr := config.Parse(body)
		if parseErr != nil {
			return nil, nil, parseErr
		}
		return page, body, nil
	}

	var page config.Page
	if jsonErr := json.Unmarshal(body, &page); jsonErr != nil {
		return nil, nil, fmt.Errorf("invalid JSON body: %w", jsonErr)
	}
	config.ApplyDefaults(&page)
	raw, marshalErr := config.MarshalYAML(&page)
	if marshalErr != nil {
		return nil, nil, marshalErr
	}
	return &page, raw, nil
}

// --- components ---

func componentSummary(s *components.Schema) map[string]interface{} {
	return map[string]interface{}{
		"type":        s.Type,
		"category":    s.Category,
		"description": s.Description,
	}
}

func (a *API) handleListComponents(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	registry := components.NewDefaultRegistry()

	result := []map[string]interface{}{}
	for _, t := range registry.Types() {
		schema, ok := registry.Get(t)
		if !ok {
			continue
		}
		if category != "" && schema.Category != category {
			continue
		}
		result = append(result, componentSummary(schema))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i]["type"].(string) < result[j]["type"].(string)
	})
	writeJSON(w, http.StatusOK, map[string]interface{}{"components": result})
}

func (a *API) handleGetComponent(w http.ResponseWriter, r *http.Request) {
	typeName := r.PathValue("type")
	registry := components.NewDefaultRegistry()
	schema, ok := registry.Get(typeName)
	if !ok {
		writeAPIError(w, http.StatusNotFound, fmt.Sprintf("unknown component type %q", typeName))
		return
	}

	result := componentSummary(schema)
	if len(schema.Props) > 0 {
		result["props"] = schema.Props
	}
	if len(schema.Actions) > 0 {
		result["actions"] = schema.Actions
	}
	if len(schema.Slots) > 0 {
		result["slots"] = schema.Slots
	}
	if len(schema.Shortcuts) > 0 {
		result["shortcuts"] = schema.Shortcuts
	}
	writeJSON(w, http.StatusOK, result)
}

// --- datasources ---

func (a *API) datasourcePath(alias string) string {
	return filepath.Join(a.SigilDir, "datasources", strings.ToLower(alias)+".yaml")
}

func (a *API) handleListDataSources(w http.ResponseWriter, r *http.Request) {
	matches, _ := filepath.Glob(filepath.Join(a.SigilDir, "datasources", "*.yaml"))
	sort.Strings(matches)

	result := []map[string]interface{}{}
	for _, path := range matches {
		data, err := readProjectFile(path)
		if err != nil {
			continue
		}
		var ds struct {
			Alias        string   `yaml:"alias"`
			Description  string   `yaml:"description"`
			Capabilities []string `yaml:"capabilities"`
			Fields       []struct {
				Name string `yaml:"name"`
			} `yaml:"fields"`
		}
		if err := yaml.Unmarshal(data, &ds); err != nil {
			continue
		}
		alias := ds.Alias
		if alias == "" {
			alias = strings.TrimSuffix(filepath.Base(path), ".yaml")
		}
		result = append(result, map[string]interface{}{
			"alias":        alias,
			"description":  ds.Description,
			"capabilities": ds.Capabilities,
			"fieldCount":   len(ds.Fields),
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"datasources": result})
}

func (a *API) handleGetDataSource(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")
	if !safeName(alias) {
		writeAPIError(w, http.StatusBadRequest, "invalid datasource alias")
		return
	}
	manifest, err := readYAMLAsJSON(a.datasourcePath(alias))
	if err != nil {
		writeAPIError(w, http.StatusNotFound, fmt.Sprintf("datasource %q not found", alias))
		return
	}
	writeJSON(w, http.StatusOK, manifest)
}

func (a *API) handleCreateDataSource(w http.ResponseWriter, r *http.Request) {
	alias, raw, err := decodeDataSourceBody(r)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !safeName(alias) {
		writeAPIError(w, http.StatusUnprocessableEntity, "datasource alias is required")
		return
	}
	path := a.datasourcePath(alias)
	if _, err := os.Stat(path); err == nil {
		writeAPIError(w, http.StatusConflict, fmt.Sprintf("datasource %q already exists", alias))
		return
	}
	if err := writeProjectFile(path, raw); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"alias": strings.ToLower(alias),
	})
}

// decodeDataSourceBody reads a datasource manifest from the request. The body
// may be a JSON object (the default) or a raw YAML document. It returns the
// alias (from the manifest) and the YAML bytes to persist.
func decodeDataSourceBody(r *http.Request) (string, []byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		return "", nil, fmt.Errorf("reading request body: %w", err)
	}
	if len(body) == 0 {
		return "", nil, fmt.Errorf("empty request body")
	}

	if isYAMLContentType(r.Header.Get("Content-Type")) {
		var meta struct {
			Alias string `yaml:"alias"`
		}
		if yamlErr := yaml.Unmarshal(body, &meta); yamlErr != nil {
			return "", nil, fmt.Errorf("invalid YAML body: %w", yamlErr)
		}
		return meta.Alias, body, nil
	}

	var obj map[string]interface{}
	if jsonErr := json.Unmarshal(body, &obj); jsonErr != nil {
		return "", nil, fmt.Errorf("invalid JSON body: %w", jsonErr)
	}
	alias, _ := obj["alias"].(string)
	raw, marshalErr := config.MarshalYAML(obj)
	if marshalErr != nil {
		return "", nil, marshalErr
	}
	return alias, raw, nil
}

// --- themes ---

func (a *API) handleListThemes(w http.ResponseWriter, r *http.Request) {
	matches, _ := filepath.Glob(filepath.Join(a.SigilDir, "themes", "*.yaml"))
	sort.Strings(matches)

	result := []map[string]interface{}{}
	for _, path := range matches {
		data, err := readProjectFile(path)
		if err != nil {
			continue
		}
		var theme struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Extends     string `yaml:"extends"`
		}
		if err := yaml.Unmarshal(data, &theme); err != nil {
			continue
		}
		name := theme.Name
		if name == "" {
			name = strings.TrimSuffix(filepath.Base(path), ".yaml")
		}
		result = append(result, map[string]interface{}{
			"name":        name,
			"description": theme.Description,
			"extends":     theme.Extends,
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"themes": result})
}

func (a *API) handleGetTheme(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if !safeName(name) {
		writeAPIError(w, http.StatusBadRequest, "invalid theme name")
		return
	}
	theme, err := readYAMLAsJSON(filepath.Join(a.SigilDir, "themes", name+".yaml"))
	if err != nil {
		writeAPIError(w, http.StatusNotFound, fmt.Sprintf("theme %q not found", name))
		return
	}
	writeJSON(w, http.StatusOK, theme)
}

// --- project ---

func (a *API) handleProject(w http.ResponseWriter, r *http.Request) {
	project, err := readYAMLAsJSON(filepath.Join(a.SigilDir, "sigil.yaml"))
	if err != nil {
		writeAPIError(w, http.StatusNotFound, "project config not found")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

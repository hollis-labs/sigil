package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chrispian/sigil/internal/config"
)

// validPageJSON returns the JSON body for a minimal valid Sigil page.
func validPageJSON(t *testing.T, id, title string) []byte {
	t.Helper()
	page := config.Page{
		Sigil:   "1.0",
		Kind:    "page",
		ID:      id,
		Title:   title,
		Overlay: "page",
		Layout: config.Component{
			ID:   "root",
			Type: "rows",
			Children: []config.Component{
				{ID: "h", Type: "heading", Props: map[string]interface{}{"level": 2, "text": "Hi"}},
			},
		},
	}
	data, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("marshal page: %v", err)
	}
	return data
}

func apiRequest(t *testing.T, h http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body != nil {
		r = httptest.NewRequest(method, path, bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestAPIHealth(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/health", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %q", resp["status"])
	}
}

func TestAPIListPages(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/pages", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Pages []pageSummary `json:"pages"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Pages) != 1 || resp.Pages[0].ID != "test-page" {
		t.Fatalf("expected one page test-page, got %+v", resp.Pages)
	}
}

func TestAPIGetPage(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/pages/test-page", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var page config.Page
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if page.ID != "test-page" || page.Title != "Test Page" {
		t.Errorf("unexpected page: %+v", page)
	}
}

func TestAPIGetPageNotFound(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/pages/nope", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAPICreatePage(t *testing.T) {
	sigilDir := setupTestProject(t)
	h := NewAPI(sigilDir).Handler()

	w := apiRequest(t, h, "POST", "/api/pages", validPageJSON(t, "new-page", "New Page"))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	path := filepath.Join(sigilDir, "pages", "new-page.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected page file written: %v", err)
	}
	parsed, err := config.ParseFile(path)
	if err != nil || parsed.Title != "New Page" {
		t.Fatalf("written page not valid: %v %+v", err, parsed)
	}
}

func TestAPICreatePageInvalid(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()

	// Missing required fields (no title / overlay / layout type).
	bad, _ := json.Marshal(config.Page{Sigil: "1.0", ID: "broken"})
	w := apiRequest(t, h, "POST", "/api/pages", bad)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Error   string   `json:"error"`
		Details []string `json:"details"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Details) == 0 {
		t.Error("expected validation error details")
	}
}

func TestAPICreatePageConflict(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "POST", "/api/pages", validPageJSON(t, "test-page", "Dup"))
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestAPIUpdatePage(t *testing.T) {
	sigilDir := setupTestProject(t)
	h := NewAPI(sigilDir).Handler()

	w := apiRequest(t, h, "PUT", "/api/pages/test-page", validPageJSON(t, "test-page", "Renamed"))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	parsed, err := config.ParseFile(filepath.Join(sigilDir, "pages", "test-page.yaml"))
	if err != nil || parsed.Title != "Renamed" {
		t.Fatalf("page not updated: %v %+v", err, parsed)
	}
}

func TestAPIUpdatePageNotFound(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "PUT", "/api/pages/ghost", validPageJSON(t, "ghost", "Ghost"))
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAPIUpdatePageIDMismatch(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "PUT", "/api/pages/test-page", validPageJSON(t, "other-id", "X"))
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestAPIListComponents(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/components", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp struct {
		Components []map[string]interface{} `json:"components"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Components) == 0 {
		t.Error("expected at least one component")
	}
}

func TestAPIGetComponent(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/components/heading", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["type"] != "heading" {
		t.Errorf("expected type heading, got %v", resp["type"])
	}
}

func TestAPIGetComponentNotFound(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/components/does-not-exist", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAPIDataSources(t *testing.T) {
	sigilDir := setupTestProject(t)
	h := NewAPI(sigilDir).Handler()

	manifest := map[string]interface{}{
		"alias":       "users",
		"description": "User records",
		"fields": []map[string]interface{}{
			{"name": "id", "type": "string"},
			{"name": "email", "type": "string"},
		},
	}
	body, _ := json.Marshal(manifest)

	w := apiRequest(t, h, "POST", "/api/datasources", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := os.Stat(filepath.Join(sigilDir, "datasources", "users.yaml")); err != nil {
		t.Fatalf("datasource file not written: %v", err)
	}

	w = apiRequest(t, h, "GET", "/api/datasources", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", w.Code)
	}
	var list struct {
		DataSources []map[string]interface{} `json:"datasources"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list.DataSources) != 1 {
		t.Fatalf("expected one datasource, got %+v", list.DataSources)
	}

	w = apiRequest(t, h, "GET", "/api/datasources/users", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", w.Code)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if got["alias"] != "users" {
		t.Errorf("expected alias users, got %v", got["alias"])
	}
}

func TestAPICreateDataSourceConflict(t *testing.T) {
	sigilDir := setupTestProject(t)
	h := NewAPI(sigilDir).Handler()
	os.WriteFile(filepath.Join(sigilDir, "datasources", "dup.yaml"), []byte("alias: dup\n"), 0o600)

	body, _ := json.Marshal(map[string]interface{}{"alias": "dup"})
	w := apiRequest(t, h, "POST", "/api/datasources", body)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestAPIThemes(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()

	w := apiRequest(t, h, "GET", "/api/themes", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", w.Code)
	}
	var list struct {
		Themes []map[string]interface{} `json:"themes"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(list.Themes) != 1 || list.Themes[0]["name"] != "default" {
		t.Fatalf("expected default theme, got %+v", list.Themes)
	}

	w = apiRequest(t, h, "GET", "/api/themes/default", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "tokens") {
		t.Error("expected theme tokens in response")
	}
}

func TestAPIProject(t *testing.T) {
	h := NewAPI(setupTestProject(t)).Handler()
	w := apiRequest(t, h, "GET", "/api/project", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["name"] != "test-project" {
		t.Errorf("expected project name test-project, got %v", resp["name"])
	}
}

// TestAPIMountedOnDevServer verifies the API coexists with the HTML preview
// routes on the same mux: /api/* serves JSON and the existing routes still work.
func TestAPIMountedOnDevServer(t *testing.T) {
	srv := New(setupTestProject(t), 0)
	h := srv.Handler()

	w := apiRequest(t, h, "GET", "/api/health", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "ok") {
		t.Fatalf("api health failed: %d %s", w.Code, w.Body.String())
	}

	// Existing HTML preview routes unaffected.
	w = apiRequest(t, h, "GET", "/", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Sigil Dev Server") {
		t.Errorf("index route broken: %d", w.Code)
	}
	w = apiRequest(t, h, "GET", "/pages/test-page", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Hello World") {
		t.Errorf("page preview route broken: %d", w.Code)
	}
}

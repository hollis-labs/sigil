package server

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func setupTestProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	sigilDir := filepath.Join(dir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)
	os.MkdirAll(filepath.Join(sigilDir, "themes"), 0755)
	os.MkdirAll(filepath.Join(sigilDir, "datasources"), 0755)

	// Project config
	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte(`
version: "1.0"
name: test-project
defaults:
  theme: default
`), 0644)

	// Theme
	os.WriteFile(filepath.Join(sigilDir, "themes", "default.yaml"), []byte(`
name: default
tokens:
  colors:
    background: "9 9 11"
    surface: "24 24 27"
    text: "244 244 245"
    text-muted: "161 161 170"
    accent: "79 70 229"
    border: "63 63 70"
  radius:
    md: "0.375rem"
`), 0644)

	// Test page
	os.WriteFile(filepath.Join(sigilDir, "pages", "test-page.yaml"), []byte(`
sigil: "1.0"
kind: page
id: test-page
title: Test Page
overlay: page
layout:
  id: root
  type: rows
  props:
    gap: 4
  children:
    - id: header
      type: heading
      props:
        level: 2
        text: Hello World
`), 0644)

	return sigilDir
}

func TestHandlePage(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	req := httptest.NewRequest("GET", "/pages/test-page", nil)
	req.SetPathValue("id", "test-page")
	w := httptest.NewRecorder()

	srv.handlePage(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()
	checks := []string{
		"Hello World",
		"cdn.tailwindcss.com",
		"EventSource('/events')",
		"Sigil",
	}
	for _, check := range checks {
		if !strings.Contains(body, check) {
			t.Errorf("expected %q in page response", check)
		}
	}
}

func TestHandlePageNotFound(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	req := httptest.NewRequest("GET", "/pages/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	srv.handlePage(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleIndex(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	srv.handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	checks := []string{
		"test-project",
		"test-page",
		"Test Page",
		"/pages/test-page",
		"1 pages",
		"Sigil Dev Server",
	}
	for _, check := range checks {
		if !strings.Contains(body, check) {
			t.Errorf("expected %q in index response", check)
		}
	}
}

func TestHandleIndexNoPages(t *testing.T) {
	dir := t.TempDir()
	sigilDir := filepath.Join(dir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)
	os.WriteFile(filepath.Join(sigilDir, "sigil.yaml"), []byte("version: \"1.0\"\nname: empty"), 0644)

	srv := New(sigilDir, 0)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	srv.handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "No pages found") {
		t.Error("expected empty state message")
	}
}

func TestHandleThemeCSS(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	req := httptest.NewRequest("GET", "/assets/theme.css", nil)
	w := httptest.NewRecorder()

	srv.handleThemeCSS(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/css") {
		t.Errorf("expected text/css, got %s", ct)
	}

	body := w.Body.String()
	if !strings.Contains(body, "--sigil-accent") {
		t.Error("expected theme CSS variable")
	}
	if !strings.Contains(body, ":root") {
		t.Error("expected :root selector")
	}
}

func TestHandleThemeCSSNoTheme(t *testing.T) {
	dir := t.TempDir()
	sigilDir := filepath.Join(dir, ".sigil")
	os.MkdirAll(filepath.Join(sigilDir, "pages"), 0755)

	srv := New(sigilDir, 0)
	req := httptest.NewRequest("GET", "/assets/theme.css", nil)
	w := httptest.NewRecorder()

	srv.handleThemeCSS(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "no theme found") {
		t.Error("expected fallback comment")
	}
}

func TestHandleSSE(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	// Connect to SSE endpoint
	resp, err := http.Get(ts.URL + "/events")
	if err != nil {
		t.Fatalf("connecting to SSE: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected text/event-stream, got %s", resp.Header.Get("Content-Type"))
	}

	// Read the initial connected event
	scanner := bufio.NewScanner(resp.Body)
	var lines []string
	timeout := time.After(2 * time.Second)
	done := make(chan bool)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
			if strings.Contains(line, "data: ok") {
				done <- true
				return
			}
		}
		done <- false
	}()

	select {
	case ok := <-done:
		if !ok {
			t.Error("expected connected event")
		}
	case <-timeout:
		t.Error("timeout waiting for SSE connected event")
	}

	// Now trigger a reload
	srv.NotifyReload()

	// Read the reload event
	reloadDone := make(chan bool)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "event: reload") {
				reloadDone <- true
				return
			}
		}
		reloadDone <- false
	}()

	select {
	case ok := <-reloadDone:
		if !ok {
			t.Error("expected reload event")
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout waiting for SSE reload event")
	}
}

func TestListPages(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	pages := srv.ListPages()
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	if pages[0].ID != "test-page" {
		t.Errorf("expected 'test-page', got %q", pages[0].ID)
	}
	if pages[0].Title != "Test Page" {
		t.Errorf("expected 'Test Page', got %q", pages[0].Title)
	}
	if pages[0].ComponentCount < 2 {
		t.Errorf("expected at least 2 components, got %d", pages[0].ComponentCount)
	}
}

func TestIntegrationServer(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	// Test index
	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("GET / status: %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Test page
	resp, err = http.Get(ts.URL + "/pages/test-page")
	if err != nil {
		t.Fatalf("GET /pages/test-page: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("GET /pages/test-page status: %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Test theme CSS
	resp, err = http.Get(ts.URL + "/assets/theme.css")
	if err != nil {
		t.Fatalf("GET /assets/theme.css: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("GET /assets/theme.css status: %d", resp.StatusCode)
	}
	resp.Body.Close()

	// Test 404
	resp, err = http.Get(ts.URL + "/pages/nonexistent")
	if err != nil {
		t.Fatalf("GET /pages/nonexistent: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestGracefulShutdown(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	// Use httptest for this test
	ts := httptest.NewServer(srv.Handler())
	ts.Close() // immediate close should not panic

	// Test Shutdown on unstarted server
	err := srv.Shutdown(context.Background())
	if err != nil {
		t.Errorf("shutdown on unstarted server: %v", err)
	}
}

func TestNotifyReloadNoClients(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	// Should not panic with no connected clients
	srv.NotifyReload()
}

func TestCountComponents(t *testing.T) {
	sigilDir := setupTestProject(t)
	srv := New(sigilDir, 0)

	pages := srv.ListPages()
	// Root (rows) + heading = 2
	if pages[0].ComponentCount != 2 {
		t.Errorf("expected 2 components, got %d", pages[0].ComponentCount)
	}
}

func TestInjectDevExtras(t *testing.T) {
	html := []byte(`<!DOCTYPE html>
<html>
<body class="test">
<h1>Hello</h1>
</body>
</html>`)

	result := injectDevExtras(html, "test-page", ".sigil")
	content := string(result)

	if !strings.Contains(content, "EventSource('/events')") {
		t.Error("expected reload script injection")
	}
	if !strings.Contains(content, "test-page") {
		t.Error("expected page name in nav bar")
	}
	if !strings.Contains(content, "Sigil") {
		t.Error("expected Sigil link in nav bar")
	}
}

// Suppress unused import warning
var _ = fmt.Sprintf

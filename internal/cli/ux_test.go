package cli

import (
	"bytes"
	"testing"
)

func TestSuggestComponentType(t *testing.T) {
	known := []string{"data-table", "heading", "button", "search-bar", "columns", "rows", "grid", "form", "badge", "text"}

	tests := []struct {
		input    string
		expected string
	}{
		{"data-tabl", "data-table"},
		{"datatable", "data-table"},
		{"header", "heading"},
		{"search-ba", "search-bar"},
		{"colums", "columns"},
		{"row", "rows"},
		{"badeg", "badge"},
		{"zzzzzzzzz", ""}, // nothing close
	}

	for _, tt := range tests {
		got := SuggestComponentType(tt.input, known)
		if got != tt.expected {
			t.Errorf("SuggestComponentType(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"abc", "abc", 0},
		{"abc", "abd", 1},
		{"kitten", "sitting", 3},
	}

	for _, tt := range tests {
		got := levenshtein(tt.a, tt.b)
		if got != tt.expected {
			t.Errorf("levenshtein(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestVersionCmd(t *testing.T) {
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"version"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}
	if out.String() != "sigil dev\n" {
		t.Errorf("expected 'sigil dev\\n', got %q", out.String())
	}
}

func TestNoColorFlag(t *testing.T) {
	// Save original state
	origColor := colorEnabled
	defer func() { colorEnabled = origColor }()

	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"version", "--no-color"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version --no-color failed: %v", err)
	}
	if colorEnabled {
		t.Error("expected color to be disabled after --no-color")
	}
}

func TestColorHelpers(t *testing.T) {
	// With color enabled
	colorEnabled = true
	defer func() { colorEnabled = true }()

	if g := green("ok"); g == "ok" {
		t.Error("expected ANSI codes in green output when color is enabled")
	}
	if r := red("err"); r == "err" {
		t.Error("expected ANSI codes in red output when color is enabled")
	}
	if y := yellow("warn"); y == "warn" {
		t.Error("expected ANSI codes in yellow output when color is enabled")
	}

	// Without color
	colorEnabled = false
	if g := green("ok"); g != "ok" {
		t.Errorf("expected plain 'ok' when color disabled, got %q", g)
	}
	if r := red("err"); r != "err" {
		t.Errorf("expected plain 'err' when color disabled, got %q", r)
	}
	if y := yellow("warn"); y != "warn" {
		t.Errorf("expected plain 'warn' when color disabled, got %q", y)
	}
}

func TestExtractQuoted(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`unknown component type "data-grid"`, "data-grid"},
		{`no quotes here`, ""},
		{`"only-start`, ""},
		{`"hello" world "second"`, "hello"},
	}
	for _, tt := range tests {
		got := extractQuoted(tt.input)
		if got != tt.expected {
			t.Errorf("extractQuoted(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestRootCmdHasExamples(t *testing.T) {
	rootCmd := NewRootCmd()
	if rootCmd.Long == "" {
		t.Error("expected root command to have Long description with examples")
	}
}

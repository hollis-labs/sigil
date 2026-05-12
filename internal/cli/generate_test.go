package cli

import (
	"bytes"
	"strings"
	"testing"
)

// TestGenerateTargetModeFlagValidation exercises the CLI-level --target-mode
// validation. Only the flag parsing + early reject path is tested here;
// end-to-end engine behavior (where the flag overrides app.yaml's
// target_mode) is covered by the renderer engine tests.
func TestGenerateTargetModeFlagValidation(t *testing.T) {
	cases := []struct {
		name        string
		args        []string
		wantErr     bool
		wantErrSubs string
	}{
		{
			name:    "unset is valid",
			args:    []string{"--dry-run"},
			wantErr: false,
		},
		{
			name:    "spa is valid",
			args:    []string{"--target-mode", "spa", "--dry-run"},
			wantErr: false,
		},
		{
			name:    "app-router is valid",
			args:    []string{"--target-mode", "app-router", "--dry-run"},
			wantErr: false,
		},
		{
			name:        "garbage value is rejected up-front",
			args:        []string{"--target-mode", "wails", "--dry-run"},
			wantErr:     true,
			wantErrSubs: `invalid --target-mode "wails"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGenerateCmd()
			// Quiet — we only check the early validation path, not the full
			// generate pipeline.
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrSubs)
				}
				if tc.wantErrSubs != "" && !strings.Contains(err.Error(), tc.wantErrSubs) {
					t.Fatalf("error %q did not contain %q", err.Error(), tc.wantErrSubs)
				}
				return
			}
			// Without a .sigil/ in this test's working dir, Execute() will
			// fail in the renderer engine. We only care that the failure
			// is NOT the up-front --target-mode rejection.
			if err != nil && strings.Contains(err.Error(), "invalid --target-mode") {
				t.Fatalf("unexpected --target-mode rejection: %v", err)
			}
		})
	}
}

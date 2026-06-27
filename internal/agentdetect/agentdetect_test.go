package agentdetect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectKnownAgentFiles(t *testing.T) {
	for _, tc := range []struct {
		name   string
		setup  func(string) error
		agent  string
		reason string
	}{
		{
			name: "codex",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# rules\n"), 0644)
			},
			agent:  "codex",
			reason: "AGENTS.md",
		},
		{
			name: "claude",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# rules\n"), 0644)
			},
			agent:  "claude",
			reason: "CLAUDE.md",
		},
		{
			name: "cursor",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".cursor", "rules"), 0755)
			},
			agent:  "cursor",
			reason: ".cursor/rules",
		},
		{
			name: "antigravity",
			setup: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".agents", "rules"), 0755)
			},
			agent:  "antigravity",
			reason: ".agents/rules",
		},
		{
			name: "gemini",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "GEMINI.md"), []byte("# rules\n"), 0644)
			},
			agent:  "gemini",
			reason: "GEMINI.md",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := tc.setup(dir); err != nil {
				t.Fatal(err)
			}
			got := Detect(dir, nil)
			if got.Agent != tc.agent || got.Reason != tc.reason {
				t.Fatalf("Detect()=%#v want agent=%q reason=%q", got, tc.agent, tc.reason)
			}
		})
	}
}

func TestDetectFallsBackToGeneric(t *testing.T) {
	got := Detect(t.TempDir(), nil)
	if got.Agent != "generic" || got.Reason != "no signal" {
		t.Fatalf("Detect()=%#v", got)
	}
}

func TestDetectAmbiguousSignalsUseGeneric(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# rules\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("# rules\n"), 0644); err != nil {
		t.Fatal(err)
	}
	got := Detect(dir, nil)
	if got.Agent != "generic" || got.Reason != "ambiguous" {
		t.Fatalf("Detect()=%#v", got)
	}
}

package instructions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLintDetectsAndFixesSafeNoise(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "AGENTS.md")
	text := "Generated at 2026-06-27 12:00:00\nKeep this\nKeep this\n```sh\nC:\\Users\\yoons\\secret\n```\nC:\\Users\\yoons\\repo\n"
	if err := os.WriteFile(file, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	findings, err := Lint(dir)
	if err != nil {
		t.Fatalf("Lint: %v", err)
	}
	if len(findings) < 3 {
		t.Fatalf("expected findings, got %#v", findings)
	}
	report, err := Fix(dir)
	if err != nil {
		t.Fatalf("Fix: %v", err)
	}
	if !strings.Contains(report, "AGENTS.md") {
		t.Fatalf("expected report to mention file: %s", report)
	}
	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	s := string(got)
	if strings.Contains(s, "Generated at") {
		t.Fatalf("timestamp not removed:\n%s", s)
	}
	if !strings.Contains(s, "```sh\nC:\\Users\\yoons\\secret\n```") {
		t.Fatalf("code fence changed:\n%s", s)
	}
	if _, err := os.Stat(filepath.Join(dir, ".logdiet", "backup")); err != nil {
		t.Fatalf("backup dir missing: %v", err)
	}
}

func TestInstallRulesIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	if _, err := InstallRules(dir, "codex", false); err != nil {
		t.Fatal(err)
	}
	if _, err := InstallRules(dir, "codex", false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(b), BeginMarker) != 1 {
		t.Fatalf("managed section not idempotent:\n%s", string(b))
	}
	if _, err := os.Stat(filepath.Join(dir, ".logdiet", "backup")); err != nil {
		t.Fatalf("backup dir missing: %v", err)
	}
}

package instructions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestRulesInstallReplaceAndRemoveManagedSection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	original := "before user content\n\n" + BeginMarker + "\nold managed text\n" + EndMarker + "\n\nafter user content\n"
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := InstallRules(dir, "codex", false); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if strings.Count(s, BeginMarker) != 1 || strings.Count(s, EndMarker) != 1 {
		t.Fatalf("install should leave exactly one managed section:\n%s", s)
	}
	if strings.Contains(s, "old managed text") {
		t.Fatalf("install did not replace old managed section:\n%s", s)
	}
	if !strings.Contains(s, "before user content") || !strings.Contains(s, "after user content") {
		t.Fatalf("install did not preserve user content:\n%s", s)
	}

	if _, err := RemoveRules(dir, "codex"); err != nil {
		t.Fatal(err)
	}
	b, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s = string(b)
	if strings.Contains(s, BeginMarker) || strings.Contains(s, EndMarker) || strings.Contains(s, RulesText) {
		t.Fatalf("uninstall did not remove managed section:\n%s", s)
	}
	if !strings.Contains(s, "before user content") || !strings.Contains(s, "after user content") {
		t.Fatalf("uninstall did not preserve user content:\n%s", s)
	}
}

func TestReplaceManagedSectionsHandlesMalformedMarkers(t *testing.T) {
	cases := []string{
		"no markers\n",
		BeginMarker + "\nunterminated managed section\n",
		EndMarker + "\nend before begin\n",
		"prefix\n" + EndMarker + "\n" + BeginMarker + "\nsuffix\n",
	}
	for _, tc := range cases {
		done := make(chan string, 1)
		go func(in string) {
			done <- replaceManagedSections(in, ManagedRulesBlock())
		}(tc)
		select {
		case got := <-done:
			if got == "" {
				t.Fatalf("malformed marker case returned empty output for %q", tc)
			}
		case <-time.After(250 * time.Millisecond):
			t.Fatalf("replaceManagedSections did not return for malformed input %q", tc)
		}
	}
}

func TestFixTextReplacesAbsolutePathsWithPlaceholders(t *testing.T) {
	input := strings.Join([]string{
		"linux /home/alice/project/file.go",
		"mac /Users/bob/project/file.go",
		"windows C:\\Users\\carol\\project\\file.go",
		"repo /repo/logdiet/internal/cli.go",
		"```",
		"fenced /home/alice/project/file.go",
		"fenced C:\\Users\\carol\\project\\file.go",
		"```",
		"",
	}, "\n")
	got := fixText(input, "", "/repo/logdiet")
	for _, want := range []string{
		"linux <home>/project/file.go",
		"mac <home>/project/file.go",
		"windows <home>\\project\\file.go",
		"repo <repo>/internal/cli.go",
		"fenced /home/alice/project/file.go",
		"fenced C:\\Users\\carol\\project\\file.go",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("fixed text missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "/<home>/") || strings.Contains(got, "$1") || strings.Contains(got, "$2") {
		t.Fatalf("fixed text contains invalid replacement artifact:\n%s", got)
	}
}

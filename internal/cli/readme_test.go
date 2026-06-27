package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestREADMEIncludesAgentAndVerificationSections(t *testing.T) {
	path := filepath.Join("..", "..", "README.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	readme := string(data)

	for _, want := range []string{
		"Put your coding agent on a token diet.",
		"LogDiet keeps full command logs locally",
		"## TL;DR",
		"## For AI agents",
		"## Try It In 60 Seconds",
		"## Agent Quickstarts",
		"## Trust but verify",
		"go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest",
		"logdiet setup codex",
		"logdiet setup claude",
		"logdiet setup cursor",
		"logdiet setup antigravity",
		"logdiet show latest:F1 --around 40",
		".agents/rules/logdiet.md",
	} {
		if !strings.Contains(readme, want) {
			t.Fatalf("README.md missing %q", want)
		}
	}
}

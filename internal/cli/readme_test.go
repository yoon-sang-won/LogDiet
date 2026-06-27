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
		"## The problem",
		"## The LogDiet loop",
		"## Before and after",
		"## TL;DR",
		"## For AI agents",
		"## Try It In 60 Seconds",
		"## Who this is for",
		"## Agent Quickstarts",
		"## What LogDiet is not",
		"## FAQ",
		"## Trust but verify",
		"go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest",
		"logdiet setup codex",
		"logdiet setup claude",
		"logdiet setup cursor",
		"logdiet setup antigravity",
		"logdiet show latest:F1 --around 40",
		"logdiet grep latest \"pattern\"",
		".agents/rules/logdiet.md",
		"No network. No telemetry. No model/API calls.",
	} {
		if !strings.Contains(readme, want) {
			t.Fatalf("README.md missing %q", want)
		}
	}
}

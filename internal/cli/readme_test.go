package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestREADMEIncludesLandingPageSections(t *testing.T) {
	readme := readProjectFile(t, "README.md")

	for _, want := range []string{
		"English",
		"한국어",
		"Agent-native token diet for coding agents.",
		"Agent-first. CLI-powered. No network. No telemetry.",
		"## Before / After",
		"## How LogDiet works",
		"mermaid",
		"## Quickstart: install for your agent",
		"## Hook rewrite bridge",
		"## Works with",
		"## Core commands",
		"## Setup modes",
		"## What LogDiet is not",
		"No network. No telemetry. No model/API calls.",
		"Automatic command rewriting is available where the agent supports command hooks.",
		"go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest",
		"logdiet setup codex --mode all",
		"logdiet show latest:F1 --around 40",
		"logdiet grep latest \"pattern\"",
		"logdiet hook rewrite --command \"go test ./...\"",
		"docs/agent-native.md",
		"README.ko.md",
	} {
		if !strings.Contains(readme, want) {
			t.Fatalf("README.md missing %q", want)
		}
	}
}

func TestKoreanREADMEIncludesLandingPageSections(t *testing.T) {
	readme := readProjectFile(t, "README.ko.md")

	for _, want := range []string{
		"English",
		"한국어",
		"agent-native",
		"Agent-first. CLI-powered. No network. No telemetry.",
		"## 왜 필요한가",
		"## Before / After",
		"## LogDiet의 레이어",
		"mermaid",
		"## Quickstart: 에이전트에 설치",
		"## Hook rewrite bridge",
		"## 지원 패키지",
		"## 주요 명령",
		"## setup mode",
		"## LogDiet이 아닌 것",
		"No network. No telemetry. No model/API calls.",
		"command hook",
		"go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest",
		"logdiet setup codex --mode all",
		"logdiet show latest:F1 --around 40",
		"logdiet hook rewrite --command \"go test ./...\"",
		"README.md",
	} {
		if !strings.Contains(readme, want) {
			t.Fatalf("README.ko.md missing %q", want)
		}
	}
}

func readProjectFile(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join("..", "..", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

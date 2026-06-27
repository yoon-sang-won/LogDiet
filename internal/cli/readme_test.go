package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestREADMEIncludesAgentAndVerificationSections(t *testing.T) {
	readme := readProjectFile(t, "README.md")

	for _, want := range []string{
		"English",
		"한국어",
		"Put your coding agent on a token diet.",
		"LogDiet keeps full command logs locally",
		"## The problem",
		"## The LogDiet loop",
		"mermaid",
		"## Before and after",
		"## Works with",
		"## Core commands",
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

func TestKoreanREADMEIncludesCoreSections(t *testing.T) {
	readme := readProjectFile(t, "README.ko.md")

	for _, want := range []string{
		"English",
		"한국어",
		"코딩 에이전트에게 토큰 다이어트를 시키세요.",
		"## 한눈에 보기",
		"## 왜 필요할까요",
		"## LogDiet의 동작 방식",
		"mermaid",
		"## Before / After",
		"## 60초 안에 써보기",
		"## AI 에이전트를 위한 사용법",
		"## 함께 쓰기 좋은 환경",
		"## 주요 명령어",
		"## FAQ",
		"No network. No telemetry. No model/API calls.",
		"go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest",
		"logdiet show latest:F1 --around 40",
		".agents/rules/logdiet.md",
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

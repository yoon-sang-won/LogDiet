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
		"Put your coding agent on a token diet.",
		"Keep the logs. Cut the noise.",
		"Stop feeding log walls to your coding agent",
		"## Before / After",
		"## How LogDiet works",
		"mermaid",
		"## Try It In 60 Seconds",
		"## For AI agents",
		"## Works with",
		"## Core commands",
		"## Agent quickstarts",
		"## What LogDiet is not",
		"## FAQ",
		"No network. No telemetry. No model/API calls.",
		"go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest",
		"logdiet show latest:F1 --around 40",
		"logdiet grep latest \"pattern\"",
		".agents/rules/logdiet.md",
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
		"코딩 에이전트에게 토큰 다이어트를 시키세요.",
		"전체 로그는 보관하고, 노이즈는 줄입니다.",
		"AI 에이전트에게 로그 벽을 먹이지 마세요",
		"## Before / After",
		"## LogDiet의 동작 방식",
		"mermaid",
		"## 60초 안에 써보기",
		"## AI 에이전트를 위한 사용법",
		"## 함께 쓰기 좋은 환경",
		"## 주요 명령어",
		"## LogDiet이 아닌 것",
		"## FAQ",
		"No network. No telemetry. No model/API calls.",
		"go install github.com/yoon-sang-won/LogDiet/cmd/logdiet@latest",
		"logdiet show latest:F1 --around 40",
		".agents/rules/logdiet.md",
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

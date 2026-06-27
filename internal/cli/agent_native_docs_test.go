package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentNativeDocsAndIntegrationFiles(t *testing.T) {
	files := map[string][]string{
		"docs/agent-native.md": {
			"# Agent-Native LogDiet",
			"## Integration levels",
			"## Agent support matrix",
			"Automatic command rewriting is available where the agent supports command hooks.",
			"## Security",
		},
		"docs/v0.2-verification.md": {
			"# v0.2 Verification",
			"logdiet hook rewrite --command \"go test ./...\"",
			"logdiet setup codex --mode all",
			"Doctor should report rules, shims, and native template status",
		},
		"docs/codex-integration-verification.md": {
			"# Codex Integration Verification",
			"## What is verified automatically",
			"## What requires manual verification",
			"/hooks",
			"Automatic command rewriting is available only where Codex command hooks are supported and trusted.",
		},
		"integrations/codex/README.md": {
			"LogDiet for Codex",
			"not magically built into Codex",
			"rules fallback",
			"hook rewrite template",
			"/hooks",
			"full raw logs stay local",
		},
		"integrations/codex/AGENTS.md": {
			"# LogDiet for Codex",
			"use LogDiet for noisy commands",
			"prefer `logdiet wrap -- <command>`",
			"logdiet show latest:F1 --around 40",
			"logdiet grep latest \"pattern\"",
			"logdiet raw latest",
			"do not ask the user to paste full terminal logs",
		},
		"integrations/codex/logdiet-instructions.md": {
			"Automatic command rewriting is available where the agent supports command hooks.",
		},
		"integrations/codex/hook-rewrite-template.sh": {
			"logdiet hook rewrite --command \"$COMMAND\"",
			"template",
		},
		"integrations/claude-code/README.md": {
			"LogDiet for Claude Code",
			"not an official Claude plugin",
		},
		"integrations/claude-code/skill.md": {
			"do not paste log walls",
			"logdiet wrap -- <cmd>",
		},
		"integrations/claude-code/hook-rewrite-template.sh": {
			"logdiet hook rewrite --command \"$COMMAND\"",
		},
		"integrations/cursor/README.md": {
			"LogDiet for Cursor",
			"template",
		},
		"integrations/cursor/logdiet.mdc": {
			"prefer LogDiet for noisy commands",
			"do not paste full logs",
		},
		"integrations/cursor/hook-rewrite-template.sh": {
			"logdiet hook rewrite --command \"$COMMAND\"",
		},
		"integrations/gemini/README.md": {
			"LogDiet for Gemini",
			"template",
		},
		"integrations/gemini/GEMINI.md": {
			"use compact evidence",
			"avoid full log pastes",
		},
		"integrations/gemini/hook-rewrite-template.sh": {
			"logdiet hook rewrite --command \"$COMMAND\"",
		},
		"integrations/antigravity/README.md": {
			"LogDiet for Antigravity",
			".agents/rules/logdiet.md",
		},
		"integrations/antigravity/logdiet.md": {
			"rules fallback first",
			"automatic command rewrite is not guaranteed",
		},
		"integrations/generic/README.md": {
			"LogDiet for generic terminal agents",
			"manual wrapper mode",
			"PATH shim mode",
			"hook adapter",
		},
		"integrations/generic/logdiet-rules.md": {
			"Prefer compact LogDiet evidence",
			"Do not paste full logs",
		},
		"integrations/generic/logdiet-hook-rewrite.sh": {
			"#!/bin/sh",
			"logdiet hook rewrite --command \"$COMMAND\"",
			"must not execute the command",
		},
		"scripts/verify-codex-integration.sh": {
			"#!/bin/sh",
			"Codex native hook template verified. Runtime trust must be verified manually in Codex with /hooks.",
			"logdiet setup codex --mode rules",
			"logdiet setup codex --mode all",
			"logdiet hook rewrite --command \"go test ./...\"",
			"\"wrap\":true",
			"\"wrap\":false",
		},
	}

	for name, wants := range files {
		t.Run(name, func(t *testing.T) {
			content := readProjectFile(t, filepath.FromSlash(name))
			for _, want := range wants {
				if !strings.Contains(content, want) {
					t.Fatalf("%s missing %q", name, want)
				}
			}
		})
	}
}

func TestReadmesAndChangelogIncludeV02Positioning(t *testing.T) {
	readme := readProjectFile(t, "README.md")
	for _, want := range []string{
		"Agent-native token diet for coding agents.",
		"Agent-first. CLI-powered. No network. No telemetry.",
		"LogDiet has two layers:",
		"Automatic command rewriting is available where the agent supports command hooks.",
		"logdiet setup codex --mode all",
		"logdiet hook rewrite --command \"go test ./...\"",
		"docs/agent-native.md",
		"### Codex verification",
		"./scripts/verify-codex-integration.sh",
		"/hooks",
	} {
		if !strings.Contains(readme, want) {
			t.Fatalf("README.md missing %q", want)
		}
	}

	korean := readProjectFile(t, "README.ko.md")
	for _, want := range []string{
		"agent-native",
		"Agent-first. CLI-powered. No network. No telemetry.",
		"command hook",
		"logdiet setup codex --mode all",
		"logdiet hook rewrite --command \"go test ./...\"",
		"Codex 검증",
		"./scripts/verify-codex-integration.sh",
		"/hooks",
	} {
		if !strings.Contains(korean, want) {
			t.Fatalf("README.ko.md missing %q", want)
		}
	}

	changelogBytes, err := os.ReadFile(filepath.Join("..", "..", "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	changelog := string(changelogBytes)
	for _, want := range []string{
		"## v0.2.0 - Unreleased",
		"agent-native token diet layer",
		"Agent integration packages under `integrations/`.",
		"`logdiet hook rewrite`",
		"Automatic command rewriting is available where an agent supports command hooks.",
	} {
		if !strings.Contains(changelog, want) {
			t.Fatalf("CHANGELOG.md missing %q", want)
		}
	}
}

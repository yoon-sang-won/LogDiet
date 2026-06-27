package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentNativeDocsAndIntegrationFiles(t *testing.T) {
	files := map[string][]string{
		"AI_INSTALL.md": {
			"# AI Install Guide for LogDiet",
			"You are an AI coding agent.",
			"The user gave you the LogDiet repository URL and asked you to install it.",
			"## What to do now",
			"logdiet bootstrap --agent auto",
			"logdiet agent-instructions --agent auto",
			"## Choose the best available integration",
			"logdiet init --agent auto",
			"## For the rest of this session",
			"## Do not require hooks",
			"Native hooks are optional advanced mode.",
		},
		"docs/agent-self-install.md": {
			"# Agent Self-Install",
			"Install https://github.com/yoon-sang-won/LogDiet",
			"logdiet bootstrap --agent auto",
			"Native adapters and fallback",
			"Why hooks are optional",
			"logdiet wrap -- <command>",
		},
		"docs/native-adapters.md": {
			"# Native Agent Adapters",
			"Native where possible. Fallback everywhere. Raw logs always local.",
			"## Thin adapter rule",
			"rules plus explicit",
			"## Supported agents matrix",
			"## Scope of transparent rewrite",
			"## Verification",
		},
		"docs/first-agent-prompt.md": {
			"# First Agent Prompt",
			"Install https://github.com/yoon-sang-won/LogDiet",
			"logdiet wrap -- <command>",
			"logdiet show latest:F1 --around 40",
			"logdiet grep latest \"pattern\"",
			"logdiet raw latest",
			"Hooks are optional. Do not block on hook setup.",
		},
		"integrations/ADAPTER_CONTRACT.md": {
			"# LogDiet Adapter Contract",
			"thin delegates",
			"logdiet hook rewrite --command",
			"logdiet wrap -- COMMAND",
			"Adapters must not",
		},
		"integrations/fixtures/README.md": {
			"# Adapter Fixtures",
			"scripts/verify-adapter-fixtures.sh",
		},
		"integrations/fixtures/commands.txt": {
			"go test ./...",
			"pytest -q",
			"npm test",
			`rg "TODO"`,
			"vim file.go",
			"logdiet wrap -- go test ./...",
		},
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
			"template / not verified",
			"logdiet wrap -- <command>",
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
			"template / not verified",
			"Runtime Cursor hook behavior is not claimed.",
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
			"template / not verified",
			"Runtime Gemini hook behavior is not verified here.",
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
			"Automatic command rewrite is not guaranteed.",
		},
		"integrations/antigravity/native-template.md": {
			"Antigravity Native Adapter Notes",
			"does not currently claim a verified Antigravity native hook adapter",
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
		"scripts/verify-agent-self-install.sh": {
			"#!/bin/sh",
			"Native hooks are not required for this verification.",
			"go install ./cmd/logdiet",
			"logdiet bootstrap --agent auto",
			"logdiet doctor",
			"logdiet agent-instructions --agent auto",
			"logdiet wrap -- sh -c",
			"logdiet raw latest",
			"logdiet grep latest \"line 2\"",
			"Agent self-install verification passed.",
		},
		"scripts/verify-adapter-fixtures.sh": {
			"#!/bin/sh",
			"go install ./cmd/logdiet",
			"integrations/fixtures/commands.txt",
			"logdiet hook rewrite --command",
			`check_wrap 'go test ./...' true`,
			`check_wrap 'vim file.go' false`,
			"Adapter fixture verification passed.",
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

func TestShellVerificationScriptsAreExecutable(t *testing.T) {
	for _, script := range []string{
		"scripts/verify-agent-self-install.sh",
		"scripts/verify-adapter-fixtures.sh",
	} {
		t.Run(script, func(t *testing.T) {
			cmd := exec.Command("git", "ls-files", "-s", script)
			cmd.Dir = filepath.Join("..", "..")
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("git ls-files: %v", err)
			}
			if !strings.HasPrefix(string(out), "100755 ") {
				t.Fatalf("%s should be executable in git index, got %q", script, string(out))
			}
		})
	}
}

func TestReadmesAndChangelogIncludeV02Positioning(t *testing.T) {
	readme := readProjectFile(t, "README.md")
	for _, want := range []string{
		"Agent-native token diet for coding agents.",
		"Agent-first. CLI-powered. No network. No telemetry.",
		"A token diet kit your coding agent can install and use by itself.",
		"## Easiest path: tell your agent",
		"Install https://github.com/yoon-sang-won/LogDiet and use it for noisy test/build/git/search output.",
		"logdiet bootstrap --agent auto",
		"logdiet agent-instructions --agent auto",
		"logdiet wrap -- pytest -q",
		"logdiet wrap -- npm test",
		"logdiet wrap -- git diff",
		"logdiet wrap -- rg \"pattern\"",
		"## What happens after bootstrap?",
		"Hooks are optional advanced mode.",
		"docs/agent-self-install.md",
		"docs/first-agent-prompt.md",
		"## Agent support matrix",
		"Transparent rewrite applies only to shell/terminal command paths exposed to hooks.",
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
		"A token diet kit your coding agent can install and use by itself.",
		"## 가장 쉬운 사용법: 에이전트에게 맡기기",
		"logdiet bootstrap --agent auto",
		"logdiet agent-instructions --agent auto",
		"logdiet wrap -- pytest -q",
		"logdiet wrap -- npm test",
		"logdiet wrap -- git diff",
		"logdiet wrap -- rg \"pattern\"",
		"## bootstrap 이후에는 무엇이 달라지나요?",
		"docs/first-agent-prompt.md",
		"## 에이전트 지원 매트릭스",
		"shell/terminal command",
		"### Codex 검증",
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
		"`AI_INSTALL.md` for agents installing LogDiet from a GitHub link.",
		"`logdiet bootstrap` for agent self-install flows.",
		"`logdiet agent-instructions` for current-session operating rules.",
		"Agent self-install documentation.",
		"Tests for bootstrap and agent instruction flows.",
		"`scripts/verify-agent-self-install.sh` for hookless self-install verification.",
		"`docs/first-agent-prompt.md` with a copy-paste prompt for coding agents.",
		"Native adapter architecture documentation.",
		"`logdiet init` entrypoint for agent integration setup/status.",
		"Adapter contract documentation.",
		"Adapter fixture verification script.",
		"Agent support matrix.",
		"README now leads with the agent self-install path.",
		"Native hooks are documented as optional advanced mode, not the default requirement.",
		"README and README.ko.md now surface the agent self-install flow earlier.",
		"`AI_INSTALL.md`, `bootstrap`, and `agent-instructions` now more clearly tell agents to continue with `logdiet wrap` without requiring hooks.",
		"Integration READMEs now distinguish native adapters, rules fallback, and manual CLI usage.",
		"`doctor` now reports native adapter status more clearly.",
		"Documentation now frames hooks as best-effort native automation with wrapper fallback everywhere.",
	} {
		if !strings.Contains(changelog, want) {
			t.Fatalf("CHANGELOG.md missing %q", want)
		}
	}
}

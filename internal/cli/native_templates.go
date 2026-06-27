package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type nativeTemplate struct {
	Name    string
	Content string
}

func installNativeTemplates(root, target string) (string, error) {
	dir, files, err := nativeTemplates(target)
	if err != nil {
		return "", err
	}
	base := filepath.Join(root, ".logdiet", "integrations", dir)
	if err := os.MkdirAll(base, 0755); err != nil {
		return "", err
	}
	for _, file := range files {
		mode := os.FileMode(0644)
		if filepath.Ext(file.Name) == ".sh" {
			mode = 0755
		}
		if err := os.WriteFile(filepath.Join(base, file.Name), []byte(file.Content), mode); err != nil {
			return "", err
		}
	}
	return filepath.ToSlash(filepath.Join(".logdiet", "integrations", dir)), nil
}

func nativeTemplateStatus(root, target string) string {
	dir, files, err := nativeTemplates(target)
	if err != nil || len(files) == 0 {
		return "unknown"
	}
	path := filepath.Join(root, ".logdiet", "integrations", dir, files[0].Name)
	switch target {
	case "codex":
		path = filepath.Join(root, ".logdiet", "integrations", dir, "hook-rewrite-template.sh")
	case "claude":
		path = filepath.Join(root, ".logdiet", "integrations", dir, "skill.md")
	}
	if fileExists(path) {
		return "template installed"
	}
	return "template missing"
}

func nativeTemplates(target string) (string, []nativeTemplate, error) {
	switch target {
	case "codex":
		return "codex", []nativeTemplate{
			{"README.md", nativeReadme("Codex")},
			{"AGENTS.md", codexInstructions()},
			{"logdiet-instructions.md", commonAgentInstructions("Codex")},
			{"hook-rewrite-template.sh", hookRewriteTemplate("Codex")},
		}, nil
	case "claude":
		return "claude-code", []nativeTemplate{
			{"README.md", nativeReadme("Claude Code")},
			{"skill.md", commonAgentInstructions("Claude Code")},
			{"hook-rewrite-template.sh", hookRewriteTemplate("Claude Code")},
		}, nil
	case "cursor":
		return "cursor", []nativeTemplate{
			{"README.md", nativeReadme("Cursor")},
			{"logdiet.mdc", commonAgentInstructions("Cursor")},
			{"hook-rewrite-template.sh", hookRewriteTemplate("Cursor")},
		}, nil
	case "gemini":
		return "gemini", []nativeTemplate{
			{"README.md", nativeReadme("Gemini")},
			{"GEMINI.md", commonAgentInstructions("Gemini")},
			{"hook-rewrite-template.sh", hookRewriteTemplate("Gemini")},
		}, nil
	case "antigravity":
		return "antigravity", []nativeTemplate{
			{"README.md", nativeReadme("Antigravity")},
			{"logdiet.md", commonAgentInstructions("Antigravity")},
		}, nil
	case "generic":
		return "generic", []nativeTemplate{
			{"README.md", nativeReadme("generic terminal agents")},
			{"logdiet-rules.md", commonAgentInstructions("generic terminal agents")},
			{"logdiet-hook-rewrite.sh", hookRewriteTemplate("generic terminal agents")},
		}, nil
	default:
		return "", nil, fmt.Errorf("unknown native template target %q", target)
	}
}

func nativeReadme(agent string) string {
	return fmt.Sprintf(`# LogDiet for %s

These files are local templates for installing LogDiet as an agent-native token diet layer.

Review the rules and hook template before enabling them in your agent. Template installation does not enable hooks automatically.
`, agent)
}

func codexInstructions() string {
	return `# LogDiet instructions for Codex

When working in this repository:

- prefer ` + "`logdiet wrap -- <command>`" + ` for noisy test/build/search/git commands if automatic hooks are unavailable;
- read compact evidence first;
- use ` + "`logdiet show latest:F1 --around 40`" + ` to expand a handle;
- use ` + "`logdiet grep latest \"pattern\"`" + ` to search exact raw output;
- use ` + "`logdiet raw latest`" + ` only when compact evidence is insufficient;
- do not ask the user to paste full terminal logs unless LogDiet expansion commands are insufficient;
- keep final responses concise:
  - success: changed files, tests run, risks;
  - failure: exact error, likely cause, next command.
`
}

func commonAgentInstructions(agent string) string {
	return fmt.Sprintf(`# LogDiet instructions for %s

Use LogDiet compact evidence for noisy terminal output.

- Do not paste log walls.
- Prefer `+"`logdiet wrap -- <cmd>`"+` when hook rewrite is unavailable.
- Use `+"`logdiet show latest:F1 --around 40`"+` to expand a handle.
- Use `+"`logdiet grep latest \"pattern\"`"+` to search exact raw output.
- Use `+"`logdiet raw latest`"+` only when compact evidence is insufficient.
- Keep success and failure reports short.

Automatic command rewriting is available where the agent supports command hooks. Other agents use rules/instructions fallback or manual `+"`logdiet wrap`"+`.
`, agent)
}

func hookRewriteTemplate(agent string) string {
	return fmt.Sprintf(`#!/bin/sh
set -eu

# LogDiet hook rewrite template for %s.
# Agent hook protocols differ. Adapt this script to the exact input/output
# contract of your agent before enabling it.
# This script must not execute the command itself.

: "${COMMAND:?COMMAND is required}"
logdiet hook rewrite --command "$COMMAND"
`, agent)
}

func oneOf(s string, vals ...string) bool {
	for _, v := range vals {
		if s == v {
			return true
		}
	}
	return false
}

package instructions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	BeginMarker = "<!-- BEGIN LOGDIET MANAGED RESPONSE CONTRACT -->"
	EndMarker   = "<!-- END LOGDIET MANAGED RESPONSE CONTRACT -->"
)

const RulesText = `LogDiet response contract:
- No prefaces, apologies, filler, or progress narration.
- Do not say "I'll inspect", "Now I'll test", or "Summary:" unless it adds new information.
- On success, report only:
  1. changed files
  2. tests/commands run
  3. remaining risks
- On failure, report only:
  1. exact error
  2. likely cause
  3. next command or smallest next edit
- Preserve exact commands, paths, line numbers, error messages, and code identifiers.
- Prefer the smallest working change.
- Reuse existing code before adding code.
- Prefer stdlib/native platform features before new dependencies.
- Do not remove validation, security checks, accessibility, error handling, or data-loss protection just to be short.
- Ask a question only when blocked.
`

func ManagedRulesBlock() string {
	return BeginMarker + "\n" + RulesText + EndMarker + "\n"
}

func InstallRules(root, target string, dryRun bool) (string, error) {
	path, prefix, err := rulesPath(root, target)
	if err != nil {
		return "", err
	}
	block := ManagedRulesBlock()
	content := prefix + block
	if b, err := os.ReadFile(path); err == nil {
		existing := string(b)
		if hasManagedSection(existing) {
			content = replaceManagedSections(existing, block)
		} else {
			sep := ""
			if existing != "" && !strings.HasSuffix(existing, "\n") {
				sep = "\n"
			}
			content = existing + sep + "\n" + block
		}
	}
	if dryRun {
		return fmt.Sprintf("would write %s\n", rel(root, path)), nil
	}
	if err := backupExisting(root, path); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	return fmt.Sprintf("installed rules: %s\n", rel(root, path)), nil
}

func RemoveRules(root, target string) (string, error) {
	path, _, err := rulesPath(root, target)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("no rules found: %s\n", rel(root, path)), nil
	}
	next := replaceManagedSections(string(b), "")
	if next == string(b) {
		return fmt.Sprintf("no managed rules found: %s\n", rel(root, path)), nil
	}
	if err := backupExisting(root, path); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(next), 0644); err != nil {
		return "", err
	}
	return fmt.Sprintf("removed rules: %s\n", rel(root, path)), nil
}

func rulesPath(root, target string) (string, string, error) {
	switch target {
	case "", "generic":
		return filepath.Join(root, ".logdiet", "LOGDIET_RULES.md"), "", nil
	case "codex":
		return filepath.Join(root, "AGENTS.md"), "", nil
	case "claude":
		return filepath.Join(root, "CLAUDE.md"), "", nil
	case "gemini":
		return filepath.Join(root, "GEMINI.md"), "", nil
	case "cursor":
		return filepath.Join(root, ".cursor", "rules", "logdiet.mdc"), "---\nalwaysApply: true\n---\n\n", nil
	default:
		return "", "", fmt.Errorf("unknown rules target %q", target)
	}
}

func hasManagedSection(s string) bool {
	return strings.Contains(s, BeginMarker) && strings.Contains(s, EndMarker)
}

func replaceManagedSections(s, block string) string {
	var out strings.Builder
	pos := 0
	wrote := false
	for {
		relStart := strings.Index(s[pos:], BeginMarker)
		if relStart < 0 {
			out.WriteString(s[pos:])
			return out.String()
		}
		start := pos + relStart
		out.WriteString(s[pos:start])
		relEnd := strings.Index(s[start:], EndMarker)
		if relEnd < 0 {
			out.WriteString(s[start:])
			return out.String()
		}
		end := start + relEnd + len(EndMarker)
		if end < len(s) && s[end] == '\n' {
			end++
		}
		if !wrote && block != "" {
			out.WriteString(block)
			wrote = true
		}
		pos = end
	}
}

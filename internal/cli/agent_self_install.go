package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/agentdetect"
	"github.com/yoon-sang-won/LogDiet/internal/instructions"
	"github.com/yoon-sang-won/LogDiet/internal/shim"
)

var supportedAgents = map[string]bool{
	"auto":        true,
	"codex":       true,
	"claude":      true,
	"cursor":      true,
	"antigravity": true,
	"gemini":      true,
	"generic":     true,
}

func bootstrapCommand(root string, args []string, stdout, stderr io.Writer) int {
	agent, err := parseAgentFlag(args)
	if err != nil {
		fmt.Fprintln(stderr, "usage: logdiet bootstrap [--agent auto|codex|claude|cursor|antigravity|gemini|generic]")
		return 2
	}
	agent = resolveAgent(root, agent)
	if _, err := shim.Install(root, "", shim.InstallOptions{}); err != nil {
		fmt.Fprintf(stderr, "error: installing shims: %v\n", err)
		return 1
	}
	if _, err := instructions.InstallRules(root, agent, false); err != nil {
		fmt.Fprintf(stderr, "error: installing %s rules: %v\n", agent, err)
		return 1
	}
	fmt.Fprint(stdout, bootstrapText(agent))
	return 0
}

func agentInstructionsCommand(root string, args []string, stdout, stderr io.Writer) int {
	agent, err := parseAgentFlag(args)
	if err != nil {
		fmt.Fprintln(stderr, "usage: logdiet agent-instructions [--agent auto|codex|claude|cursor|antigravity|gemini|generic]")
		return 2
	}
	agent = resolveAgent(root, agent)
	fmt.Fprint(stdout, agentInstructionsText(agent))
	return 0
}

func parseAgentFlag(args []string) (string, error) {
	agent := "auto"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--agent":
			if i+1 >= len(args) {
				return "", fmt.Errorf("--agent requires value")
			}
			agent = args[i+1]
			i++
		default:
			return "", fmt.Errorf("unknown argument %q", args[i])
		}
	}
	if !supportedAgents[agent] {
		return "", fmt.Errorf("unknown agent %q", agent)
	}
	return agent, nil
}

func resolveAgent(root, agent string) string {
	if agent != "auto" {
		return agent
	}
	return agentdetect.Detect(root, os.Environ()).Agent
}

func bootstrapText(agent string) string {
	var b strings.Builder
	fmt.Fprintln(&b, "LogDiet bootstrap complete")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "agent: %s\n", agent)
	fmt.Fprintln(&b, "engine: OK")
	fmt.Fprintln(&b, "state: .logdiet OK")
	fmt.Fprintln(&b, "runs: .logdiet/runs OK")
	fmt.Fprintf(&b, "rules: %s installed\n", ruleDisplayPath(agent))
	fmt.Fprintln(&b, "shims: .logdiet/bin OK")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Continue this session using LogDiet:")
	fmt.Fprintln(&b, "  logdiet wrap -- <command>")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Use it for:")
	fmt.Fprintln(&b, "  tests, builds, type checks, git diffs/status/logs, search, long logs")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Expand evidence:")
	fmt.Fprintln(&b, "  logdiet show latest:F1 --around 40")
	fmt.Fprintln(&b, "  logdiet grep latest \"pattern\"")
	fmt.Fprintln(&b, "  logdiet raw latest")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Do not ask the user to paste full logs unless these commands are insufficient.")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Native hooks are optional advanced mode.")
	fmt.Fprintln(&b, "This bootstrap did not require hook trust.")
	fmt.Fprintln(&b)
	if agent != "generic" {
		fmt.Fprintln(&b, "Optional native hooks:")
		fmt.Fprintf(&b, "  logdiet setup %s --mode native\n", agent)
		fmt.Fprintln(&b, "  Hooks may require review/trust inside the agent.")
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b, "Verify:")
	fmt.Fprintln(&b, "  logdiet doctor")
	return b.String()
}

func agentInstructionsText(agent string) string {
	var b strings.Builder
	fmt.Fprintln(&b, "LogDiet session instructions")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "For the rest of this session, use LogDiet for noisy commands.")
	fmt.Fprintln(&b)
	if note := agentRulesNote(agent); note != "" {
		fmt.Fprintln(&b, note)
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b, "Run noisy commands as:")
	fmt.Fprintln(&b, "  logdiet wrap -- <command>")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "- Prefer it for tests, builds, type checks, git diffs, search, and long logs.")
	fmt.Fprintln(&b, "- Read compact evidence first.")
	fmt.Fprintln(&b, "- Expand one handle with `logdiet show latest:F1 --around 40`.")
	fmt.Fprintln(&b, "- Search raw output with `logdiet grep latest \"pattern\"`.")
	fmt.Fprintln(&b, "- Print full raw logs with `logdiet raw latest` only if compact evidence is insufficient.")
	fmt.Fprintln(&b, "- Do not ask the user to paste full terminal logs unless LogDiet expansion commands are insufficient.")
	fmt.Fprintln(&b, "- Keep final responses concise:")
	fmt.Fprintln(&b, "  - success: changed files, tests run, risks;")
	fmt.Fprintln(&b, "  - failure: exact error, likely cause, next command.")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Manual examples:")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "  logdiet wrap -- go test ./...")
	fmt.Fprintln(&b, "  logdiet wrap -- pytest -q")
	fmt.Fprintln(&b, "  logdiet wrap -- npm test")
	fmt.Fprintln(&b, "  logdiet show latest:F1 --around 40")
	fmt.Fprintln(&b, "  logdiet grep latest \"panic\"")
	return b.String()
}

func agentRulesNote(agent string) string {
	switch agent {
	case "codex":
		return "Codex rules are usually installed in AGENTS.md."
	case "claude":
		return "Claude Code rules are usually installed in CLAUDE.md."
	case "cursor":
		return "Cursor rules are usually installed in .cursor/rules/logdiet.mdc."
	case "antigravity":
		return "Antigravity rules are usually installed in .agents/rules/logdiet.md."
	case "gemini":
		return "Gemini rules are usually installed in GEMINI.md."
	default:
		return ""
	}
}

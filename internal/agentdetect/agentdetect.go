package agentdetect

import (
	"os"
	"path/filepath"
)

type Result struct {
	Agent  string
	Reason string
}

func Detect(cwd string, env []string) Result {
	signals := detectFileSignals(cwd)
	if len(signals) == 0 {
		return Result{Agent: "generic", Reason: "no signal"}
	}
	if len(signals) > 1 {
		return Result{Agent: "generic", Reason: "ambiguous"}
	}
	return signals[0]
}

func detectFileSignals(cwd string) []Result {
	checks := []struct {
		agent  string
		reason string
		path   string
		isDir  bool
	}{
		{"codex", "AGENTS.md", "AGENTS.md", false},
		{"claude", "CLAUDE.md", "CLAUDE.md", false},
		{"cursor", ".cursor/rules", filepath.Join(".cursor", "rules"), true},
		{"antigravity", ".agents/rules", filepath.Join(".agents", "rules"), true},
		{"gemini", "GEMINI.md", "GEMINI.md", false},
	}
	var signals []Result
	for _, check := range checks {
		info, err := os.Stat(filepath.Join(cwd, check.path))
		if err != nil {
			continue
		}
		if check.isDir && !info.IsDir() {
			continue
		}
		if !check.isDir && info.IsDir() {
			continue
		}
		signals = append(signals, Result{Agent: check.agent, Reason: check.reason})
	}
	return signals
}

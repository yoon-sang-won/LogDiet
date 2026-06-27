package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLIWrapRawShowAndGrep(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "fake")
	body := "#!/bin/sh\necho 'alpha ok'\necho 'beta failed' >&2\nexit 5\n"
	if runtime.GOOS == "windows" {
		script = filepath.Join(dir, "fake.cmd")
		body = "@echo off\necho alpha ok\necho beta failed 1>&2\nexit /b 5\n"
	}
	if err := os.WriteFile(script, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
	oldwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	var out, errb bytes.Buffer
	code := Run([]string{"wrap", "--", script}, &out, &errb)
	if code != 5 {
		t.Fatalf("wrap exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if !strings.Contains(out.String(), "logdiet run") || !strings.Contains(out.String(), "raw:") {
		t.Fatalf("bad wrap output:\n%s", out.String())
	}
	for _, want := range []string{
		"cmd:",
		"summary:",
		"show: logdiet show latest:F1 --around 40",
		"raw: logdiet raw latest",
		"stats:",
		"approx_saved=",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("wrap output missing %q:\n%s", want, out.String())
		}
	}
	runIDBytes, err := os.ReadFile(filepath.Join(dir, ".logdiet", "latest"))
	if err != nil {
		t.Fatalf("latest pointer missing: %v", err)
	}
	runID := strings.TrimSpace(string(runIDBytes))
	if runID == "" {
		t.Fatal("latest pointer is empty")
	}
	for _, name := range []string{"stdout.txt", "stderr.txt", "combined.txt", "meta.json", "index.json"} {
		if _, err := os.Stat(filepath.Join(dir, ".logdiet", "runs", runID, name)); err != nil {
			t.Fatalf("missing run file %s: %v", name, err)
		}
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"raw", "latest"}, &out, &errb); code != 0 {
		t.Fatalf("raw exit=%d err=%s", code, errb.String())
	}
	if out.String() != "alpha ok\nbeta failed\n" &&
		out.String() != "alpha ok\nbeta failed \n" &&
		out.String() != "alpha ok\r\nbeta failed \r\n" {
		t.Fatalf("raw did not print exact combined output:\n%q", out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"grep", "latest", "beta"}, &out, &errb); code != 0 {
		t.Fatalf("grep exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "beta failed") {
		t.Fatalf("grep missing match:\n%s", out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"grep", "latest", "["}, &out, &errb); code != 2 {
		t.Fatalf("invalid regexp exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"show", "F1", "--around", "2"}, &out, &errb); code != 0 {
		t.Fatalf("show exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "     2 | beta failed") {
		t.Fatalf("show missing raw line:\n%s", out.String())
	}
}

func TestCLICommonCommands(t *testing.T) {
	dir := t.TempDir()
	oldwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	var out, errb bytes.Buffer
	if code := Run([]string{"--version"}, &out, &errb); code != 0 {
		t.Fatalf("version exit=%d err=%s", code, errb.String())
	}
	if strings.TrimSpace(out.String()) != "logdiet 0.2.0-dev" {
		t.Fatalf("bad version output: %q", out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"rules", "--print"}, &out, &errb); code != 0 {
		t.Fatalf("rules print exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "LogDiet response contract:") {
		t.Fatalf("rules print missing contract:\n%s", out.String())
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"help"}, &out, &errb); code != 0 {
		t.Fatalf("help exit=%d err=%s", code, errb.String())
	}
	help := out.String()
	for _, want := range []string{
		"LogDiet keeps full command logs locally and feeds AI coding agents compact, expandable evidence.",
		"logdiet wrap -- pytest -q",
		"logdiet show latest:F1 --around 40",
		"logdiet raw latest",
		"logdiet grep latest \"panic\"",
		"logdiet hook rewrite --command \"go test ./...\"",
		"bootstrap              install LogDiet rules/shims for an AI agent",
		"agent-instructions     print session instructions for an AI agent",
		"init                    install or inspect agent integrations",
		"logdiet bootstrap --agent auto",
		"logdiet agent-instructions --agent auto",
		"logdiet init --agent auto",
		"logdiet init --agent claude --mode native",
		"logdiet init --agent cursor --mode all",
		"logdiet init --show",
		"logdiet init --uninstall --agent codex",
	} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q:\n%s", want, help)
		}
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"env", "--shell", "powershell"}, &out, &errb); code != 0 {
		t.Fatalf("env exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), ".logdiet\\bin") {
		t.Fatalf("env missing shim path: %s", out.String())
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"install"}, &out, &errb); code != 0 {
		t.Fatalf("install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	for _, want := range []string{
		"LogDiet installed",
		"state: .logdiet OK",
		"shims: .logdiet/bin OK",
		"runs: .logdiet/runs OK",
		"activate:",
		`eval "$(logdiet env)"`,
		"PowerShell:",
		"Invoke-Expression (logdiet env --shell powershell)",
		"verify:",
		"logdiet doctor",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("install output missing %q:\n%s", want, out.String())
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".logdiet", "bin")); err != nil {
		t.Fatalf("install did not create shim dir: %v", err)
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"rules", "--install", "codex"}, &out, &errb); code != 0 {
		t.Fatalf("rules install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(agents), "BEGIN LOGDIET MANAGED RESPONSE CONTRACT") != 1 {
		t.Fatalf("rules install did not create one managed section:\n%s", string(agents))
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"rules", "--install", "all", "--dry-run"}, &out, &errb); code != 0 {
		t.Fatalf("rules dry-run all exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if !strings.Contains(out.String(), ".agents/rules/logdiet.md") {
		t.Fatalf("rules all dry-run missing antigravity target:\n%s", out.String())
	}

	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("Generated at 2026-06-27 12:00:00\nPath "+dir+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"lint-instructions", "--fix"}, &out, &errb); code != 0 {
		t.Fatalf("lint fix exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	fixed, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(fixed), "Generated at") || !strings.Contains(string(fixed), "<repo>") {
		t.Fatalf("lint fix did not sanitize as expected:\n%s", string(fixed))
	}

	out.Reset()
	errb.Reset()
}

func TestInitShowAndInstallModes(t *testing.T) {
	t.Run("show", func(t *testing.T) {
		dir := t.TempDir()
		oldwd, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldwd)

		var out, errb bytes.Buffer
		if code := Run([]string{"init", "--show"}, &out, &errb); code != 0 {
			t.Fatalf("init show exit=%d out=%s err=%s", code, out.String(), errb.String())
		}
		for _, want := range []string{
			"LogDiet init status",
			"auto-detected agent: generic",
			"rules fallback: available",
			"explicit wrapper: available",
			"Native where possible. Fallback everywhere. Raw logs always local.",
			"Codex:",
			"Generic:",
		} {
			if !strings.Contains(out.String(), want) {
				t.Fatalf("init --show output missing %q:\n%s", want, out.String())
			}
		}
		if _, err := os.Stat(filepath.Join(dir, ".logdiet")); !os.IsNotExist(err) {
			t.Fatalf("init --show should not modify files, err=%v", err)
		}
	})

	for _, tc := range []struct {
		name      string
		args      []string
		wantFile  string
		wantOut   []string
		absentDir string
	}{
		{
			name:     "generic default rules",
			args:     []string{"init", "--agent", "generic"},
			wantFile: filepath.Join(".logdiet", "LOGDIET_RULES.md"),
			wantOut:  []string{"LogDiet init: generic", "mode: rules", "rules: .logdiet/LOGDIET_RULES.md installed", "native: skipped"},
		},
		{
			name:     "codex rules",
			args:     []string{"init", "--agent", "codex", "--mode", "rules"},
			wantFile: "AGENTS.md",
			wantOut:  []string{"LogDiet init: codex", "mode: rules", "rules: AGENTS.md installed", "native: skipped"},
		},
		{
			name:     "claude native",
			args:     []string{"init", "--agent", "claude", "--mode", "native"},
			wantFile: "CLAUDE.md",
			wantOut:  []string{"LogDiet init: claude", "mode: native", "rules: CLAUDE.md installed", "native: template installed .logdiet/integrations/claude-code"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			oldwd, _ := os.Getwd()
			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldwd)

			var out, errb bytes.Buffer
			if code := Run(tc.args, &out, &errb); code != 0 {
				t.Fatalf("init exit=%d out=%s err=%s", code, out.String(), errb.String())
			}
			if _, err := os.Stat(filepath.Join(dir, tc.wantFile)); err != nil {
				t.Fatalf("init missing %s: %v\n%s", tc.wantFile, err, out.String())
			}
			for _, want := range tc.wantOut {
				if !strings.Contains(out.String(), want) {
					t.Fatalf("init output missing %q:\n%s", want, out.String())
				}
			}
		})
	}
}

func TestInitUninstallRemovesManagedRulesOnly(t *testing.T) {
	dir := t.TempDir()
	oldwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	var out, errb bytes.Buffer
	if code := Run([]string{"init", "--agent", "codex", "--mode", "rules"}, &out, &errb); code != 0 {
		t.Fatalf("init install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"init", "--uninstall", "--agent", "codex"}, &out, &errb); code != 0 {
		t.Fatalf("init uninstall exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if !strings.Contains(out.String(), "removed rules: AGENTS.md") ||
		!strings.Contains(out.String(), "raw logs preserved") {
		t.Fatalf("init uninstall output missing conservative removal details:\n%s", out.String())
	}
	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(agents), "BEGIN LOGDIET MANAGED RESPONSE CONTRACT") {
		t.Fatalf("init uninstall left managed block:\n%s", string(agents))
	}
}

func TestBootstrapGenericCodexAndAuto(t *testing.T) {
	for _, tc := range []struct {
		name      string
		args      []string
		rulesFile string
		agentLine string
	}{
		{
			name:      "generic",
			args:      []string{"bootstrap", "--agent", "generic"},
			rulesFile: filepath.Join(".logdiet", "LOGDIET_RULES.md"),
			agentLine: "agent: generic",
		},
		{
			name:      "codex",
			args:      []string{"bootstrap", "--agent", "codex"},
			rulesFile: "AGENTS.md",
			agentLine: "agent: codex",
		},
		{
			name:      "auto no signals",
			args:      []string{"bootstrap", "--agent", "auto"},
			rulesFile: filepath.Join(".logdiet", "LOGDIET_RULES.md"),
			agentLine: "agent: generic",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			oldwd, _ := os.Getwd()
			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldwd)

			var out, errb bytes.Buffer
			if code := Run(tc.args, &out, &errb); code != 0 {
				t.Fatalf("bootstrap exit=%d out=%s err=%s", code, out.String(), errb.String())
			}
			for _, rel := range []string{filepath.Join(".logdiet", "runs"), filepath.Join(".logdiet", "bin"), tc.rulesFile} {
				if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
					t.Fatalf("bootstrap missing %s: %v\n%s", rel, err, out.String())
				}
			}
			for _, want := range []string{
				"LogDiet bootstrap complete",
				tc.agentLine,
				"engine: OK",
				"state: .logdiet OK",
				"runs: .logdiet/runs OK",
				"shims: .logdiet/bin OK",
				"Continue this session using LogDiet:",
				"logdiet wrap -- <command>",
				"Use it for:",
				"tests, builds, type checks, git diffs/status/logs, search, long logs",
				"Expand evidence:",
				"logdiet show latest:F1 --around 40",
				"logdiet grep latest \"pattern\"",
				"logdiet raw latest",
				"Do not ask the user to paste full logs unless these commands are insufficient.",
				"Native hooks are optional advanced mode.",
				"This bootstrap did not require hook trust.",
				"Verify:",
				"logdiet doctor",
			} {
				if !strings.Contains(out.String(), want) {
					t.Fatalf("bootstrap output missing %q:\n%s", want, out.String())
				}
			}
			if strings.Contains(out.String(), "hook enabled") {
				t.Fatalf("bootstrap should not claim hooks are enabled:\n%s", out.String())
			}
		})
	}
}

func TestBootstrapUsageErrors(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run([]string{"bootstrap", "--agent", "unknown"}, &out, &errb); code != 2 {
		t.Fatalf("bootstrap invalid agent exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
}

func TestAgentInstructionsGenericAndCodex(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "generic",
			args: []string{"agent-instructions", "--agent", "generic"},
			want: []string{
				"LogDiet session instructions",
				"For the rest of this session, use LogDiet for noisy commands.",
				"Run noisy commands as:",
				"logdiet wrap -- <command>",
				"logdiet show latest:F1 --around 40",
				"logdiet grep latest \"pattern\"",
				"logdiet raw latest",
				"Do not ask the user to paste full terminal logs unless LogDiet expansion commands are insufficient.",
				"logdiet wrap -- go test ./...",
			},
		},
		{
			name: "codex",
			args: []string{"agent-instructions", "--agent", "codex"},
			want: []string{
				"LogDiet session instructions",
				"Codex rules are usually installed in AGENTS.md.",
				"logdiet wrap -- pytest -q",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			oldwd, _ := os.Getwd()
			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldwd)

			var out, errb bytes.Buffer
			if code := Run(tc.args, &out, &errb); code != 0 {
				t.Fatalf("agent-instructions exit=%d out=%s err=%s", code, out.String(), errb.String())
			}
			for _, want := range tc.want {
				if !strings.Contains(out.String(), want) {
					t.Fatalf("agent-instructions output missing %q:\n%s", want, out.String())
				}
			}
			if _, err := os.Stat(filepath.Join(dir, ".logdiet")); !os.IsNotExist(err) {
				t.Fatalf("agent-instructions should not modify files, err=%v", err)
			}
		})
	}
}

func TestAgentInstructionsUsageErrors(t *testing.T) {
	var out, errb bytes.Buffer
	if code := Run([]string{"agent-instructions", "--agent", "unknown"}, &out, &errb); code != 2 {
		t.Fatalf("agent-instructions invalid agent exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
}

func TestCLIHookRewriteJSON(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want string
	}{
		{
			name: "wrap",
			in:   "go test ./...",
			want: `{"wrap":true,"command":"logdiet wrap -- go test ./...","reason":"known noisy developer command"}` + "\n",
		},
		{
			name: "no wrap",
			in:   "echo hello",
			want: `{"wrap":false,"command":"echo hello","reason":"not selected"}` + "\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out, errb bytes.Buffer
			code := Run([]string{"hook", "rewrite", "--command", tc.in}, &out, &errb)
			if code != 0 {
				t.Fatalf("hook rewrite exit=%d out=%s err=%s", code, out.String(), errb.String())
			}
			if out.String() != tc.want {
				t.Fatalf("bad JSON:\n%s\nwant:\n%s", out.String(), tc.want)
			}
			if errb.String() != "" {
				t.Fatalf("stderr should be empty: %s", errb.String())
			}
		})
	}
}

func TestCLIHookRewriteUsageErrors(t *testing.T) {
	var out, errb bytes.Buffer
	code := Run([]string{"hook", "rewrite"}, &out, &errb)
	if code != 2 {
		t.Fatalf("hook rewrite usage exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if out.String() != "" {
		t.Fatalf("usage error should not write JSON stdout: %s", out.String())
	}
	if !strings.Contains(errb.String(), "usage: logdiet hook rewrite --command <command>") {
		t.Fatalf("usage error missing help: %s", errb.String())
	}
}

func TestCLIBenchFixturesFromRepo(t *testing.T) {
	oldwd, _ := os.Getwd()
	repo := filepath.Join("..", "..")
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	var out, errb bytes.Buffer
	if code := Run([]string{"bench-fixtures"}, &out, &errb); code != 0 {
		t.Fatalf("bench exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	for _, want := range []string{
		"fixture",
		"raw_bytes",
		"compact_bytes",
		"approx_raw_tokens",
		"approx_compact_tokens",
		"reduction",
		"handles",
		"pytest_failure.txt",
		"rg_matches.txt",
		"approx token estimates use ceil(bytes / 4)",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("bench output missing %q:\n%s", want, out.String())
		}
	}
}

func TestSetupCodexAndAntigravity(t *testing.T) {
	for _, tc := range []struct {
		agent string
		file  string
	}{
		{"codex", "AGENTS.md"},
		{"antigravity", filepath.Join(".agents", "rules", "logdiet.md")},
	} {
		t.Run(tc.agent, func(t *testing.T) {
			dir := t.TempDir()
			oldwd, _ := os.Getwd()
			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldwd)

			var out, errb bytes.Buffer
			if code := Run([]string{"setup", tc.agent}, &out, &errb); code != 0 {
				t.Fatalf("setup exit=%d out=%s err=%s", code, out.String(), errb.String())
			}
			if _, err := os.Stat(filepath.Join(dir, ".logdiet", "bin")); err != nil {
				t.Fatalf("setup did not create bin: %v", err)
			}
			b, err := os.ReadFile(filepath.Join(dir, tc.file))
			if err != nil {
				t.Fatalf("setup did not create %s: %v", tc.file, err)
			}
			if strings.Count(string(b), "BEGIN LOGDIET MANAGED RESPONSE CONTRACT") != 1 {
				t.Fatalf("setup wrote duplicate or missing managed section:\n%s", string(b))
			}
			if !strings.Contains(out.String(), "logdiet doctor") ||
				!strings.Contains(out.String(), "PowerShell:") ||
				!strings.Contains(out.String(), "Invoke-Expression") {
				t.Fatalf("setup output missing activation/doctor hints:\n%s", out.String())
			}
			out.Reset()
			errb.Reset()
			if code := Run([]string{"setup", tc.agent}, &out, &errb); code != 0 {
				t.Fatalf("setup idempotent exit=%d out=%s err=%s", code, out.String(), errb.String())
			}
			b, err = os.ReadFile(filepath.Join(dir, tc.file))
			if err != nil {
				t.Fatal(err)
			}
			if strings.Count(string(b), "BEGIN LOGDIET MANAGED RESPONSE CONTRACT") != 1 {
				t.Fatalf("second setup duplicated managed section:\n%s", string(b))
			}
		})
	}
}

func TestSetupModes(t *testing.T) {
	t.Run("rules only", func(t *testing.T) {
		dir := t.TempDir()
		oldwd, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldwd)

		var out, errb bytes.Buffer
		if code := Run([]string{"setup", "codex", "--mode", "rules"}, &out, &errb); code != 0 {
			t.Fatalf("setup rules exit=%d out=%s err=%s", code, out.String(), errb.String())
		}
		if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err != nil {
			t.Fatalf("rules mode did not install AGENTS.md: %v", err)
		}
		if _, err := os.Stat(filepath.Join(dir, ".logdiet", "bin")); !os.IsNotExist(err) {
			t.Fatalf("rules mode should not install shims, err=%v", err)
		}
		if !strings.Contains(out.String(), "mode: rules") ||
			!strings.Contains(out.String(), "rules: AGENTS.md") ||
			strings.Contains(out.String(), "hook enabled") {
			t.Fatalf("setup rules output is not explicit:\n%s", out.String())
		}
	})

	t.Run("all installs shims and native templates", func(t *testing.T) {
		dir := t.TempDir()
		oldwd, _ := os.Getwd()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(oldwd)

		var out, errb bytes.Buffer
		if code := Run([]string{"setup", "codex", "--mode", "all"}, &out, &errb); code != 0 {
			t.Fatalf("setup all exit=%d out=%s err=%s", code, out.String(), errb.String())
		}
		if _, err := os.Stat(filepath.Join(dir, ".logdiet", "bin")); err != nil {
			t.Fatalf("all mode did not install shims: %v", err)
		}
		if _, err := os.Stat(filepath.Join(dir, ".logdiet", "integrations", "codex", "hook-rewrite-template.sh")); err != nil {
			t.Fatalf("all mode did not install native template: %v", err)
		}
		for _, want := range []string{
			"mode: all",
			"rules: AGENTS.md",
			"shims: .logdiet/bin",
			"native: template installed",
			"review generated hook/plugin files",
			"logdiet doctor",
		} {
			if !strings.Contains(out.String(), want) {
				t.Fatalf("setup all output missing %q:\n%s", want, out.String())
			}
		}
	})
}

func TestCodexSetupOutputsContainOperationalInstructions(t *testing.T) {
	dir := t.TempDir()
	oldwd, _ := os.Getwd()
	oldPath := os.Getenv("PATH")
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	defer os.Setenv("PATH", oldPath)

	var out, errb bytes.Buffer
	if code := Run([]string{"setup", "codex", "--mode", "rules"}, &out, &errb); code != 0 {
		t.Fatalf("setup codex rules exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"logdiet wrap",
		"logdiet show latest:F1 --around 40",
		"logdiet grep latest",
		"logdiet raw latest",
		"do not ask the user to paste full logs",
	} {
		if !strings.Contains(string(agents), want) {
			t.Fatalf("AGENTS.md missing %q:\n%s", want, string(agents))
		}
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"setup", "codex", "--mode", "all"}, &out, &errb); code != 0 {
		t.Fatalf("setup codex all exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if _, err := os.Stat(filepath.Join(dir, ".logdiet", "integrations", "codex", "hook-rewrite-template.sh")); err != nil {
		t.Fatalf("codex native hook template missing: %v", err)
	}
	agents, err = os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(agents), "BEGIN LOGDIET MANAGED RESPONSE CONTRACT") != 1 {
		t.Fatalf("setup all should preserve one managed rules block:\n%s", string(agents))
	}

	bin := filepath.Join(dir, ".logdiet", "bin")
	if err := os.Setenv("PATH", bin+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"doctor"}, &out, &errb); code != 0 {
		t.Fatalf("doctor exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	for _, want := range []string{
		"Codex AGENTS.md: installed",
		"Codex rules: AGENTS.md installed",
		"Codex native: template installed",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("doctor output missing %q:\n%s", want, out.String())
		}
	}
}

func TestDoctorBeforeAndAfterInstall(t *testing.T) {
	dir := t.TempDir()
	oldwd, _ := os.Getwd()
	oldPath := os.Getenv("PATH")
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	defer os.Setenv("PATH", oldPath)

	var out, errb bytes.Buffer
	if code := Run([]string{"doctor"}, &out, &errb); code != 1 {
		t.Fatalf("doctor before install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if !strings.Contains(out.String(), "state: .logdiet missing") {
		t.Fatalf("doctor before install missing state report:\n%s", out.String())
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"install"}, &out, &errb); code != 0 {
		t.Fatalf("install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	bin := filepath.Join(dir, ".logdiet", "bin")
	if err := os.Setenv("PATH", bin+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"doctor"}, &out, &errb); code != 0 {
		t.Fatalf("doctor after install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	for _, want := range []string{"PATH: .logdiet/bin is first OK", "agent rules:", "Codex AGENTS.md:", "latest run: none"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("doctor output missing %q:\n%s", want, out.String())
		}
	}
}

func TestDoctorShowsAgentNativeStatus(t *testing.T) {
	dir := t.TempDir()
	oldwd, _ := os.Getwd()
	oldPath := os.Getenv("PATH")
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	defer os.Setenv("PATH", oldPath)

	var out, errb bytes.Buffer
	if code := Run([]string{"setup", "codex", "--mode", "all"}, &out, &errb); code != 0 {
		t.Fatalf("setup all exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	bin := filepath.Join(dir, ".logdiet", "bin")
	if err := os.Setenv("PATH", bin+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"doctor"}, &out, &errb); code != 0 {
		t.Fatalf("doctor exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	for _, want := range []string{
		"Agent integrations",
		"auto-detected agent: codex",
		"rules fallback: available",
		"explicit wrapper: available",
		"Codex:",
		"  rules: installed",
		"  native adapter: installed",
		"  transparent rewrite: partial",
		"  trust required: yes",
		"Claude Code:",
		"  native adapter: template",
		"Generic:",
		"  native adapter: not applicable",
		"  transparent rewrite: no",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("doctor output missing %q:\n%s", want, out.String())
		}
	}
}

func TestShimBypassRunsRealCommandWithoutWrapping(t *testing.T) {
	dir := t.TempDir()
	shimDir := filepath.Join(dir, ".logdiet", "bin")
	realDir := filepath.Join(dir, "real")
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(realDir, 0755); err != nil {
		t.Fatal(err)
	}
	name := "fake"
	script := filepath.Join(realDir, name)
	body := "#!/bin/sh\necho direct ok\n"
	if runtime.GOOS == "windows" {
		script += ".cmd"
		body = "@echo off\necho direct ok\n"
	}
	if err := os.WriteFile(script, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
	oldwd, _ := os.Getwd()
	oldPath := os.Getenv("PATH")
	oldBypass := os.Getenv("LOGDIET_BYPASS")
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	defer os.Setenv("PATH", oldPath)
	defer os.Setenv("LOGDIET_BYPASS", oldBypass)
	if err := os.Setenv("PATH", shimDir+string(os.PathListSeparator)+realDir); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("LOGDIET_BYPASS", "1"); err != nil {
		t.Fatal(err)
	}
	var out, errb bytes.Buffer
	code := Run([]string{"shim", "--shim-dir", shimDir, "--", name}, &out, &errb)
	if code != 0 {
		t.Fatalf("shim bypass exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if strings.Contains(out.String(), "logdiet run") || !strings.Contains(out.String(), "direct ok") {
		t.Fatalf("shim bypass did not run real command directly:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, ".logdiet", "latest")); !os.IsNotExist(err) {
		t.Fatalf("bypass should not create latest pointer, err=%v", err)
	}
}

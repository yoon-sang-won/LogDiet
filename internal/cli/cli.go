package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"logdiet/internal/bench"
	"logdiet/internal/compact"
	"logdiet/internal/instructions"
	"logdiet/internal/run"
	"logdiet/internal/shim"
	"logdiet/internal/store"
	"logdiet/internal/textutil"
	"logdiet/internal/version"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(stdout, helpText())
		return 0
	}
	if args[0] == "--version" || args[0] == "version" {
		fmt.Fprintf(stdout, "logdiet %s\n", version.Version)
		return 0
	}
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	switch args[0] {
	case "wrap":
		return wrapCommand(root, args[1:], nil, stdout, stderr)
	case "show":
		return showCommand(root, args[1:], stdout, stderr)
	case "raw":
		return rawCommand(root, args[1:], stdout, stderr)
	case "grep":
		return grepCommand(root, args[1:], stdout, stderr)
	case "install":
		return installCommand(root, args[1:], stdout, stderr)
	case "uninstall":
		return uninstallCommand(root, args[1:], stdout, stderr)
	case "shim":
		return shimCommand(root, args[1:], stdout, stderr)
	case "env":
		return envCommand(args[1:], stdout)
	case "rules":
		return rulesCommand(root, args[1:], stdout, stderr)
	case "lint-instructions":
		return lintCommand(root, args[1:], stdout, stderr)
	case "bench-fixtures":
		return benchCommand(root, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "usage error: unknown command %q\n", args[0])
		return 2
	}
}

func helpText() string {
	return `LogDiet keeps full command logs locally and feeds compact, expandable evidence.

common commands:
  logdiet install
  logdiet env
  logdiet wrap -- <cmd>
  logdiet show <run-id>:<handle> --around 40
  logdiet raw
  logdiet grep latest "pattern"
  logdiet lint-instructions
  logdiet rules --print
  logdiet bench-fixtures
`
}

func wrapCommand(root string, args []string, display []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "--" || len(args) == 1 {
		fmt.Fprintln(stderr, "usage: logdiet wrap -- <command> [args...]")
		return 2
	}
	execArgs := args[1:]
	if display == nil {
		display = execArgs
	}
	capres, capErr := run.Capture(execArgs)
	code := capres.ExitCode
	if capErr != nil && code == 127 && len(capres.Stdout) == 0 && len(capres.Stderr) == 0 {
		fmt.Fprintf(stderr, "error: executable not found: %s\n", execArgs[0])
		return 127
	}
	runID := store.GenerateRunID()
	res := compact.Compact(display, capres.Stdout, capres.Stderr, capres.Combined, code)
	res.RunID = runID
	rendered := compact.Render(res)
	compact.SetRenderedStats(&res, rendered)
	rendered = compact.Render(res)
	compact.SetRenderedStats(&res, rendered)
	if err := store.SaveRun(root, store.RunData{
		RunID: runID, CWD: root, Cmd: display, StartedAt: capres.StartedAt, EndedAt: capres.EndedAt,
		ExitCode: code, Stdout: capres.Stdout, Stderr: capres.Stderr, Combined: capres.Combined, Result: res,
	}); err != nil {
		fmt.Fprintf(stderr, "error: storing run: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, rendered)
	return code
}

func showCommand(root string, args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "usage: logdiet show <run-id>:<handle> --around <N>")
		return 2
	}
	target := args[0]
	around := 20
	for i := 1; i < len(args); i++ {
		if args[i] == "--around" && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err != nil || n < 0 {
				fmt.Fprintln(stderr, "usage error: --around requires a non-negative integer")
				return 2
			}
			around = n
			i++
		} else {
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", args[i])
			return 2
		}
	}
	runID, handle := splitTarget(target)
	if handle == "" {
		handle = runID
		runID = "latest"
	}
	idx, err := store.LoadIndex(root, runID)
	if err != nil {
		fmt.Fprintf(stderr, "error: run not found: %v\n", err)
		return 1
	}
	var found *compact.EvidenceItem
	for i := range idx.Items {
		if idx.Items[i].ID == handle {
			found = &idx.Items[i]
			break
		}
	}
	if found == nil {
		fmt.Fprintf(stderr, "error: handle %s not found in run %s\n", handle, idx.RunID)
		return 1
	}
	b, err := store.ReadRaw(root, idx.RunID, found.Stream)
	if err != nil {
		fmt.Fprintf(stderr, "error: reading raw output: %v\n", err)
		return 1
	}
	lines := textutil.SplitLines(b)
	start, end := textutil.ClampRange(found.StartLine-around, found.EndLine+around, len(lines))
	fmt.Fprintf(stdout, "run %s handle %s lines %d-%d\n", idx.RunID, found.ID, start, end)
	for i := start; i <= end; i++ {
		fmt.Fprintf(stdout, "%6d | %s\n", i, lines[i-1])
	}
	return 0
}

func rawCommand(root string, args []string, stdout, stderr io.Writer) int {
	runID := "latest"
	stream := "combined"
	head, tail := -1, -1
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--stdout":
			stream = "stdout"
		case "--stderr":
			stream = "stderr"
		case "--combined":
			stream = "combined"
		case "--head", "--tail":
			if i+1 >= len(args) {
				fmt.Fprintf(stderr, "usage error: %s requires N\n", args[i])
				return 2
			}
			n, err := strconv.Atoi(args[i+1])
			if err != nil || n < 0 {
				fmt.Fprintf(stderr, "usage error: %s requires non-negative N\n", args[i])
				return 2
			}
			if args[i] == "--head" {
				head = n
			} else {
				tail = n
			}
			i++
		default:
			if strings.HasPrefix(args[i], "--") {
				fmt.Fprintf(stderr, "usage error: unknown argument %q\n", args[i])
				return 2
			}
			runID = args[i]
		}
	}
	b, err := store.ReadRaw(root, runID, stream)
	if err != nil {
		fmt.Fprintf(stderr, "error: reading raw output: %v\n", err)
		return 1
	}
	lines := textutil.SplitLines(b)
	if head >= 0 || tail >= 0 {
		if head >= 0 && head < len(lines) {
			lines = lines[:head]
		}
		if tail >= 0 && tail < len(lines) {
			lines = lines[len(lines)-tail:]
		}
		for _, line := range lines {
			fmt.Fprintln(stdout, line)
		}
		return 0
	}
	_, _ = stdout.Write(b)
	return 0
}

func grepCommand(root string, args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "usage: logdiet grep <run-id> <pattern> [--ignore-case] [--around N]")
		return 2
	}
	runID, pattern := args[0], args[1]
	ignoreCase := false
	around := 0
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--ignore-case":
			ignoreCase = true
		case "--around":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "usage error: --around requires N")
				return 2
			}
			n, err := strconv.Atoi(args[i+1])
			if err != nil || n < 0 {
				fmt.Fprintln(stderr, "usage error: --around requires non-negative N")
				return 2
			}
			around = n
			i++
		default:
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", args[i])
			return 2
		}
	}
	if ignoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(stderr, "regexp error: %v\n", err)
		return 2
	}
	b, err := store.ReadRaw(root, runID, "combined")
	if err != nil {
		fmt.Fprintf(stderr, "error: reading raw output: %v\n", err)
		return 1
	}
	lines := textutil.SplitLines(b)
	matches := 0
	printed := map[int]bool{}
	for i, line := range lines {
		if !re.MatchString(line) {
			continue
		}
		matches++
		start, end := textutil.ClampRange(i+1-around, i+1+around, len(lines))
		for n := start; n <= end; n++ {
			if printed[n] {
				continue
			}
			printed[n] = true
			fmt.Fprintf(stdout, "%d:%s\n", n, lines[n-1])
		}
	}
	if matches == 0 {
		return 1
	}
	return 0
}

func installCommand(root string, args []string, stdout, stderr io.Writer) int {
	opts := shim.InstallOptions{}
	rules := false
	for _, arg := range args {
		switch arg {
		case "--local":
		case "--force":
			opts.Force = true
		case "--rules":
			rules = true
		default:
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", arg)
			return 2
		}
	}
	msg, err := shim.Install(root, "", opts)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, msg)
	if rules {
		rmsg, err := instructions.InstallRules(root, "generic", false)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		fmt.Fprint(stdout, rmsg)
	}
	return 0
}

func uninstallCommand(root string, args []string, stdout, stderr io.Writer) int {
	rules := false
	for _, arg := range args {
		if arg == "--rules" {
			rules = true
		} else {
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", arg)
			return 2
		}
	}
	msg, err := shim.Uninstall(root)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, msg)
	if rules {
		for _, target := range []string{"generic", "codex", "claude", "gemini", "cursor"} {
			msg, err := instructions.RemoveRules(root, target)
			if err != nil {
				fmt.Fprintf(stderr, "error: %v\n", err)
				return 1
			}
			fmt.Fprint(stdout, msg)
		}
	}
	return 0
}

func shimCommand(root string, args []string, stdout, stderr io.Writer) int {
	shimDir := ""
	var cmd []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--shim-dir":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "usage error: --shim-dir requires dir")
				return 2
			}
			shimDir = args[i+1]
			i++
		case "--":
			cmd = args[i+1:]
			i = len(args)
		default:
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", args[i])
			return 2
		}
	}
	if shimDir == "" || len(cmd) == 0 {
		fmt.Fprintln(stderr, "usage: logdiet shim --shim-dir <dir> -- <command-name> [args...]")
		return 2
	}
	exe, _ := os.Executable()
	pathValue := os.Getenv("PATH")
	real, err := shim.ResolveRealCommand(cmd[0], shimDir, pathValue, exe)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 127
	}
	cleanPath := shim.SanitizePATH(pathValue, shimDir)
	childEnv := os.Environ()
	childEnv = setEnv(childEnv, "PATH", cleanPath)
	childEnv = setEnv(childEnv, "LOGDIET_ACTIVE", "1")
	realArgs := append([]string{real}, cmd[1:]...)
	if os.Getenv("LOGDIET_BYPASS") == "1" || os.Getenv("LOGDIET_ACTIVE") == "1" || os.Getenv("LOGDIET_MODE") == "off" {
		return run.Execute(realArgs, childEnv, stdout, stderr)
	}
	mode := os.Getenv("LOGDIET_MODE")
	if mode == "" {
		mode = "auto"
	}
	if !shim.ShouldWrap(cmd, mode) {
		return run.Execute(realArgs, childEnv, stdout, stderr)
	}
	oldPath := os.Getenv("PATH")
	_ = os.Setenv("PATH", cleanPath)
	defer os.Setenv("PATH", oldPath)
	return wrapCommand(root, append([]string{"--"}, realArgs...), cmd, stdout, stderr)
}

func envCommand(args []string, stdout io.Writer) int {
	shell := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--shell" && i+1 < len(args) {
			shell = args[i+1]
			i++
		} else {
			return 2
		}
	}
	fmt.Fprint(stdout, shim.Env(shell))
	return 0
}

func rulesCommand(root string, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "--print" {
		fmt.Fprint(stdout, instructions.RulesText)
		return 0
	}
	install := ""
	dryRun := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--install":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "usage error: --install requires target")
				return 2
			}
			install = args[i+1]
			i++
		case "--dry-run":
			dryRun = true
		default:
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", args[i])
			return 2
		}
	}
	if install == "all" {
		for _, target := range []string{"generic", "codex", "claude", "gemini", "cursor"} {
			msg, err := instructions.InstallRules(root, target, dryRun)
			if err != nil {
				fmt.Fprintf(stderr, "error: %v\n", err)
				return 1
			}
			fmt.Fprint(stdout, msg)
		}
		return 0
	}
	msg, err := instructions.InstallRules(root, install, dryRun)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprint(stdout, msg)
	return 0
}

func lintCommand(root string, args []string, stdout, stderr io.Writer) int {
	fix, jsonOut := false, false
	for _, arg := range args {
		switch arg {
		case "--fix":
			fix = true
		case "--json":
			jsonOut = true
		default:
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", arg)
			return 2
		}
	}
	if fix {
		msg, err := instructions.Fix(root)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		fmt.Fprint(stdout, msg)
		return 0
	}
	findings, err := instructions.Lint(root)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if jsonOut {
		s, err := instructions.FindingsJSON(findings)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		fmt.Fprint(stdout, s)
		return 0
	}
	fmt.Fprint(stdout, instructions.FormatFindings(findings))
	return 0
}

func benchCommand(root string, args []string, stdout, stderr io.Writer) int {
	jsonOut := false
	for _, arg := range args {
		if arg == "--json" {
			jsonOut = true
		} else {
			fmt.Fprintf(stderr, "usage error: unknown argument %q\n", arg)
			return 2
		}
	}
	results, err := bench.Run(root)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if jsonOut {
		s, err := bench.JSON(results)
		if err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		fmt.Fprint(stdout, s)
		return 0
	}
	fmt.Fprint(stdout, bench.Format(results))
	return 0
}

func splitTarget(s string) (string, string) {
	if i := strings.LastIndex(s, ":"); i >= 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

func setEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, v := range env {
		if strings.HasPrefix(v, prefix) {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}

func executableName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/agentdetect"
	"github.com/yoon-sang-won/LogDiet/internal/instructions"
	"github.com/yoon-sang-won/LogDiet/internal/shim"
)

type initOptions struct {
	agent     string
	mode      string
	show      bool
	uninstall bool
}

func initCommand(root string, args []string, stdout, stderr io.Writer) int {
	opts, err := parseInitArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, "usage: logdiet init [--agent auto|codex|claude|cursor|antigravity|gemini|generic] [--mode rules|shim|native|all] [--show] [--uninstall]")
		return 2
	}
	if opts.show {
		fmt.Fprint(stdout, initStatusText(root))
		return 0
	}
	agent := resolveAgent(root, opts.agent)
	if opts.uninstall {
		fmt.Fprintf(stdout, "LogDiet init uninstall: %s\n\n", agent)
		msg, err := instructions.RemoveRules(root, agent)
		if err != nil {
			fmt.Fprintf(stderr, "error: removing %s rules: %v\n", agent, err)
			return 1
		}
		fmt.Fprint(stdout, msg)
		fmt.Fprintln(stdout, "native templates: preserved for manual review/removal")
		fmt.Fprintln(stdout, "raw logs preserved")
		return 0
	}
	return runInitInstall(root, agent, opts.mode, stdout, stderr)
}

func parseInitArgs(args []string) (initOptions, error) {
	opts := initOptions{agent: "auto", mode: "rules"}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--agent":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--agent requires value")
			}
			opts.agent = args[i+1]
			i++
		case "--mode":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--mode requires value")
			}
			opts.mode = args[i+1]
			i++
		case "--show":
			opts.show = true
		case "--uninstall":
			opts.uninstall = true
		default:
			return opts, fmt.Errorf("unknown argument %q", args[i])
		}
	}
	if !supportedAgents[opts.agent] {
		return opts, fmt.Errorf("unknown agent %q", opts.agent)
	}
	if !oneOf(opts.mode, "rules", "shim", "native", "all") {
		return opts, fmt.Errorf("unknown mode %q", opts.mode)
	}
	return opts, nil
}

func runInitInstall(root, agent, mode string, stdout, stderr io.Writer) int {
	if _, err := setupTargets(agent); err != nil {
		fmt.Fprintf(stderr, "usage error: %v\n", err)
		return 2
	}
	installShims := mode == "shim" || mode == "all"
	installNative := mode == "native" || mode == "all"
	if installShims {
		if _, err := shim.Install(root, "", shim.InstallOptions{}); err != nil {
			fmt.Fprintf(stderr, "error: installing shims: %v\n", err)
			return 1
		}
	}
	fmt.Fprintf(stdout, "LogDiet init: %s\n\n", agent)
	fmt.Fprintf(stdout, "mode: %s\n\n", mode)
	fmt.Fprintln(stdout, "installed:")
	if _, err := instructions.InstallRules(root, agent, false); err != nil {
		fmt.Fprintf(stderr, "error: installing %s rules: %v\n", agent, err)
		return 1
	}
	fmt.Fprintf(stdout, "  rules: %s installed\n", ruleDisplayPath(agent))
	if installShims {
		fmt.Fprintln(stdout, "  shims: .logdiet/bin OK")
	} else {
		fmt.Fprintln(stdout, "  shims: skipped")
	}
	if installNative {
		path, err := installNativeTemplates(root, agent)
		if err != nil {
			fmt.Fprintf(stderr, "error: installing %s native templates: %v\n", agent, err)
			return 1
		}
		fmt.Fprintf(stdout, "  native: template installed %s\n", path)
	} else {
		fmt.Fprintln(stdout, "  native: skipped")
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "next:")
	fmt.Fprintln(stdout, "  1. run: logdiet doctor")
	fmt.Fprintln(stdout, "  2. use: logdiet wrap -- <command>")
	fmt.Fprintln(stdout, "  3. enable/trust native hooks only after review, if your agent supports them")
	return 0
}

func initStatusText(root string) string {
	var b strings.Builder
	fmt.Fprintln(&b, "LogDiet init status")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Native where possible. Fallback everywhere. Raw logs always local.")
	fmt.Fprintln(&b)
	b.WriteString(agentIntegrationStatusText(root))
	return b.String()
}

func agentIntegrationStatusText(root string) string {
	detected := agentdetect.Detect(root, os.Environ()).Agent
	var b strings.Builder
	fmt.Fprintln(&b, "Agent integrations")
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "auto-detected agent: %s\n", detected)
	fmt.Fprintln(&b, "rules fallback: available")
	fmt.Fprintln(&b, "explicit wrapper: available")
	for _, entry := range agentStatusEntries(root) {
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "%s:\n", entry.Label)
		fmt.Fprintf(&b, "  rules: %s\n", entry.Rules)
		fmt.Fprintf(&b, "  native adapter: %s\n", entry.Native)
		fmt.Fprintf(&b, "  transparent rewrite: %s\n", entry.TransparentRewrite)
		fmt.Fprintf(&b, "  trust required: %s\n", entry.TrustRequired)
	}
	return b.String()
}

type agentStatusEntry struct {
	Label              string
	Target             string
	Rules              string
	Native             string
	TransparentRewrite string
	TrustRequired      string
}

func agentStatusEntries(root string) []agentStatusEntry {
	specs := []struct {
		label       string
		target      string
		transparent string
		trust       string
	}{
		{"Codex", "codex", "partial", "yes"},
		{"Claude Code", "claude", "template / not verified", "yes"},
		{"Cursor", "cursor", "template / not verified", "yes"},
		{"Gemini", "gemini", "template / not verified", "yes"},
		{"Antigravity", "antigravity", "unknown", "unknown"},
		{"Generic", "generic", "no", "no"},
	}
	entries := make([]agentStatusEntry, 0, len(specs))
	for _, spec := range specs {
		entries = append(entries, agentStatusEntry{
			Label:              spec.label,
			Target:             spec.target,
			Rules:              ruleInstallStatus(root, spec.target),
			Native:             nativeAdapterStatus(root, spec.target),
			TransparentRewrite: spec.transparent,
			TrustRequired:      spec.trust,
		})
	}
	return entries
}

func nativeAdapterStatus(root, target string) string {
	if target == "generic" {
		return "not applicable"
	}
	dir, files, err := nativeTemplates(target)
	if err != nil || len(files) == 0 {
		return "not supported yet"
	}
	base := filepath.Join(root, ".logdiet", "integrations", dir)
	if dirExists(base) {
		return "installed"
	}
	return "template"
}

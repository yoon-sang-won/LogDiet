package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var ShimCommands = []string{"git", "rg", "grep", "pytest", "go", "cargo", "npm", "pnpm", "yarn", "bun", "make", "just"}

func ResolveRealCommand(name, shimDir, pathValue, logdietPath string) (string, error) {
	for _, dir := range filepath.SplitList(pathValue) {
		if dir == "" || shouldSkipDir(dir, shimDir) {
			continue
		}
		for _, candidate := range commandCandidates(name) {
			full := filepath.Join(dir, candidate)
			if samePath(full, logdietPath) {
				continue
			}
			info, err := os.Stat(full)
			if err != nil || info.IsDir() {
				continue
			}
			return full, nil
		}
	}
	return "", fmt.Errorf("command %q not found outside LogDiet shims", name)
}

func SanitizePATH(pathValue, shimDir string) string {
	var keep []string
	for _, dir := range filepath.SplitList(pathValue) {
		if dir == "" || shouldSkipDir(dir, shimDir) {
			continue
		}
		keep = append(keep, dir)
	}
	return strings.Join(keep, string(os.PathListSeparator))
}

func shouldSkipDir(dir, shimDir string) bool {
	clean := filepath.Clean(dir)
	if samePath(clean, shimDir) {
		return true
	}
	norm := filepath.ToSlash(clean)
	return strings.HasSuffix(norm, "/.logdiet/bin")
}

func samePath(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	ca := filepath.Clean(a)
	cb := filepath.Clean(b)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(ca, cb)
	}
	return ca == cb
}

func commandCandidates(name string) []string {
	if runtime.GOOS != "windows" {
		return []string{name}
	}
	return windowsCommandCandidates(name)
}

func windowsCommandCandidates(name string) []string {
	ext := strings.ToLower(filepath.Ext(name))
	if ext != "" {
		return []string{name}
	}
	return []string{name + ".exe", name + ".cmd", name + ".bat", name + ".com"}
}

func ShouldWrap(cmd []string, mode string) bool {
	if len(cmd) == 0 {
		return false
	}
	if mode == "" {
		mode = "auto"
	}
	if mode == "force" {
		return true
	}
	if mode == "off" {
		return false
	}
	name := filepath.Base(cmd[0])
	name = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, ".exe"), ".cmd"), ".bat"), ".com")
	switch name {
	case "rg", "grep", "pytest", "make", "just":
		return true
	case "git":
		return len(cmd) > 1 && oneOf(cmd[1], "status", "diff", "show", "log")
	case "go":
		return len(cmd) > 1 && cmd[1] == "test"
	case "cargo":
		return len(cmd) > 1 && oneOf(cmd[1], "test", "build", "check", "clippy")
	case "npm", "pnpm":
		return len(cmd) > 1 && (cmd[1] == "test" || (len(cmd) > 2 && cmd[1] == "run" && oneOf(cmd[2], "test", "build", "lint")))
	case "yarn":
		return len(cmd) > 1 && oneOf(cmd[1], "test", "build", "lint")
	case "bun":
		return len(cmd) > 1 && cmd[1] == "test"
	default:
		return false
	}
}

func oneOf(s string, vals ...string) bool {
	for _, v := range vals {
		if s == v {
			return true
		}
	}
	return false
}

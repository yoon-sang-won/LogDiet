package instructions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"logdiet/internal/store"
	"logdiet/internal/textutil"
)

type Finding struct {
	ID      string `json:"id"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
	Fixable bool   `json:"fixable"`
}

var (
	dateRE      = regexp.MustCompile(`\b20[0-9]{2}-[01][0-9]-[0-3][0-9]\b`)
	timeRE      = regexp.MustCompile(`\b[0-3]?[0-9]:[0-5][0-9](:[0-5][0-9])?\b`)
	stampWordRE = regexp.MustCompile(`(?i)(generated|created|updated|last updated|timestamp|current[- ]date)`)
	winUserRE   = regexp.MustCompile(`(?i)C:\\Users\\[^\\\s]+((?:\\[^\s\\]+)*)`)
	unixUserRE  = regexp.MustCompile(`/(?:Users|home)/[^/\s]+((?:/[^\s/]+)*)`)
)

func Lint(root string) ([]Finding, error) {
	files, err := instructionFiles(root)
	if err != nil {
		return nil, err
	}
	var findings []Finding
	blockOwners := map[string]Finding{}
	for _, path := range files {
		fs, err := lintFile(root, path)
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			findings = append(findings, f)
		}
		b, _ := os.ReadFile(path)
		for _, block := range normalizedParagraphBlocks(string(b)) {
			if block == "" || len(block) < 80 {
				continue
			}
			if first, ok := blockOwners[block]; ok {
				findings = append(findings, Finding{
					File: rel(root, path), Line: 1, Kind: "duplicate-block", Fixable: false,
					Message: fmt.Sprintf("instruction block duplicates %s:%d and can waste input tokens", first.File, first.Line),
				})
			} else {
				blockOwners[block] = Finding{File: rel(root, path), Line: 1}
			}
		}
	}
	for i := range findings {
		findings[i].ID = fmt.Sprintf("I%d", i+1)
	}
	return findings, nil
}

func FormatFindings(findings []Finding) string {
	var b strings.Builder
	fmt.Fprintf(&b, "instruction lint: %d findings\n", len(findings))
	for _, f := range findings {
		fmt.Fprintf(&b, "%s %s:%d %s\n", f.ID, f.File, f.Line, f.Message)
	}
	if len(findings) > 0 {
		b.WriteString("fix: logdiet lint-instructions --fix\n")
	}
	return b.String()
}

func FindingsJSON(findings []Finding) (string, error) {
	b, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return "", err
	}
	return string(append(b, '\n')), nil
}

func lintFile(root, path string) ([]Finding, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	home, _ := os.UserHomeDir()
	repoAbs, _ := filepath.Abs(root)
	lines := textutil.SplitLines(b)
	var findings []Finding
	seen := map[string]int{}
	inFence := false
	fenceStart := 0
	managedCount := 0
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~") {
			if inFence {
				if i+1-fenceStart > 80 {
					findings = append(findings, finding(root, path, fenceStart, "large-code-fence", "large code fence can dominate repeated prompts", false))
				}
				inFence = false
			} else {
				inFence = true
				fenceStart = i + 1
			}
			continue
		}
		if strings.Contains(line, BeginMarker) {
			managedCount++
		}
		if inFence {
			continue
		}
		if looksVolatileTimestamp(line) {
			findings = append(findings, finding(root, path, i+1, "timestamp", "volatile timestamp changes repeated prompt prefixes and can reduce cache friendliness", true))
		}
		if hasAbsPath(line, home, repoAbs) {
			findings = append(findings, finding(root, path, i+1, "absolute-path", "absolute local path can change across machines and waste prompt cache", true))
		}
		norm := strings.ToLower(textutil.NormalizeSpace(line))
		if norm != "" {
			if prev, ok := seen[norm]; ok {
				findings = append(findings, finding(root, path, i+1, "duplicate-line", fmt.Sprintf("duplicate line repeats %s:%d", rel(root, path), prev), true))
			} else {
				seen[norm] = i + 1
			}
		}
		if strings.Contains(strings.ToLower(line), "describe every step") || strings.Contains(strings.ToLower(line), "explain everything in detail") || strings.Contains(strings.ToLower(line), "always include long summaries") {
			findings = append(findings, finding(root, path, i+1, "verbose-rule", "rule conflicts with terse output and can increase token use", false))
		}
	}
	if managedCount > 1 {
		findings = append(findings, finding(root, path, 1, "duplicate-managed-section", "duplicate LogDiet managed sections repeat instruction text", true))
	}
	if inFence && len(lines)+1-fenceStart > 80 {
		findings = append(findings, finding(root, path, fenceStart, "large-code-fence", "large code fence can dominate repeated prompts", false))
	}
	return findings, nil
}

func Fix(root string) (string, error) {
	files, err := instructionFiles(root)
	if err != nil {
		return "", err
	}
	var changed []string
	home, _ := os.UserHomeDir()
	repoAbs, _ := filepath.Abs(root)
	for _, path := range files {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		next := fixText(string(b), home, repoAbs)
		if next == string(b) {
			continue
		}
		if err := backupExisting(root, path); err != nil {
			return "", err
		}
		if err := os.WriteFile(path, []byte(next), 0644); err != nil {
			return "", err
		}
		changed = append(changed, rel(root, path))
	}
	if len(changed) == 0 {
		return "fixed 0 files\n", nil
	}
	return "fixed files: " + strings.Join(changed, ", ") + "\n", nil
}

func fixText(s, home, repoAbs string) string {
	lines := textutil.SplitLines([]byte(s))
	var out []string
	seen := map[string]bool{}
	blank := 0
	inFence := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~") {
			inFence = !inFence
			out = append(out, line)
			continue
		}
		next := line
		if !inFence {
			if looksVolatileTimestamp(line) {
				continue
			}
			if repoAbs != "" {
				next = strings.ReplaceAll(next, repoAbs, "<repo>")
			}
			if home != "" {
				next = strings.ReplaceAll(next, home, "<home>")
			}
			next = winUserRE.ReplaceAllString(next, `<home>${1}`)
			next = unixUserRE.ReplaceAllString(next, `<home>${1}`)
			norm := strings.ToLower(textutil.NormalizeSpace(next))
			if norm != "" {
				if seen[norm] {
					continue
				}
				seen[norm] = true
			}
		}
		if strings.TrimSpace(next) == "" {
			blank++
			if blank > 2 {
				continue
			}
		} else {
			blank = 0
		}
		out = append(out, next)
	}
	fixed := strings.Join(out, "\n")
	fixed = replaceManagedSectionsKeepOne(fixed)
	if strings.HasSuffix(s, "\n") || fixed != "" {
		fixed += "\n"
	}
	return fixed
}

func replaceManagedSectionsKeepOne(s string) string {
	firstStart := strings.Index(s, BeginMarker)
	if firstStart < 0 {
		return s
	}
	firstEndRel := strings.Index(s[firstStart:], EndMarker)
	if firstEndRel < 0 {
		return s
	}
	firstEnd := firstStart + firstEndRel + len(EndMarker)
	keep := s[firstStart:firstEnd]
	without := replaceManagedSections(s, "")
	prefix := strings.TrimRight(s[:firstStart], "\n")
	suffix := strings.TrimLeft(without[firstStart:], "\n")
	parts := []string{}
	if prefix != "" {
		parts = append(parts, prefix)
	}
	parts = append(parts, keep)
	if suffix != "" {
		parts = append(parts, suffix)
	}
	return strings.Join(parts, "\n")
}

func instructionFiles(root string) ([]string, error) {
	var files []string
	add := func(path string) {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			files = append(files, path)
		}
	}
	for _, relp := range []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md", filepath.Join(".github", "copilot-instructions.md")} {
		add(filepath.Join(root, relp))
	}
	for _, pattern := range []string{filepath.Join(root, ".cursor", "rules", "*.mdc"), filepath.Join(root, ".codex", "*.md"), filepath.Join(root, ".opencode", "*.md"), filepath.Join(root, ".aider*")} {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, m := range matches {
			add(m)
		}
	}
	sort.Strings(files)
	return files, nil
}

func looksVolatileTimestamp(line string) bool {
	return stampWordRE.MatchString(line) || (dateRE.MatchString(line) && timeRE.MatchString(line))
}

func hasAbsPath(line, home, repoAbs string) bool {
	return (home != "" && strings.Contains(line, home)) ||
		(repoAbs != "" && strings.Contains(line, repoAbs)) ||
		winUserRE.MatchString(line) || unixUserRE.MatchString(line)
}

func finding(root, path string, line int, kind, msg string, fixable bool) Finding {
	return Finding{File: rel(root, path), Line: line, Kind: kind, Message: msg, Fixable: fixable}
}

func normalizedParagraphBlocks(s string) []string {
	parts := regexp.MustCompile(`\n\s*\n`).Split(s, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.ToLower(textutil.NormalizeSpace(p)))
	}
	return out
}

func backupExisting(root, path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	relp := rel(root, path)
	dst := filepath.Join(store.BackupDir(root), time.Now().UTC().Format("20060102T150405Z"), relp)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, b, 0644)
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(r)
}

package compact

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/textutil"
)

func Compact(cmd []string, stdout []byte, stderr []byte, combined []byte, exitCode int) Result {
	if len(combined) == 0 {
		combined = append(append([]byte{}, stdout...), stderr...)
	}
	name := ""
	if len(cmd) > 0 {
		name = filepath.Base(cmd[0])
		name = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, ".exe"), ".cmd"), ".bat"), ".com")
	}
	rawBytes := len(combined)
	res := Result{
		Cmd:      append([]string{}, cmd...),
		ExitCode: exitCode,
		Stats: Stats{
			RawBytes:        rawBytes,
			ApproxRawTokens: textutil.ApproxTokens(rawBytes),
		},
	}
	switch {
	case name == "pytest":
		res = compactPytest(res, combined)
	case name == "go" && len(cmd) > 1 && cmd[1] == "test":
		res = compactGoTest(res, combined)
	case name == "cargo" && hasAnyArg(cmd, "test", "build", "check", "clippy"):
		res = compactCargo(res, combined)
	case name == "npm" || name == "pnpm" || name == "yarn" || name == "bun":
		res = compactJSTestOrBuild(res, combined)
	case name == "git" && len(cmd) > 1 && cmd[1] == "status":
		res = compactGitStatus(res, combined)
	case name == "git" && len(cmd) > 1 && cmd[1] == "diff":
		res = compactGitDiff(res, combined)
	case name == "git" && len(cmd) > 1 && (cmd[1] == "show" || cmd[1] == "log"):
		res = compactGitLog(res, combined)
	case name == "rg" || name == "grep":
		res = compactSearch(res, combined)
	default:
		res = compactGeneric(res, combined)
	}
	if res.Summary == "" {
		if exitCode == 0 {
			res.Summary = "ok"
		} else {
			res.Summary = fmt.Sprintf("exit %d", exitCode)
		}
	}
	res.Stats.RawBytes = rawBytes
	res.Stats.ApproxRawTokens = textutil.ApproxTokens(rawBytes)
	return res
}

func Render(res Result) string {
	var b bytes.Buffer
	if res.RunID != "" {
		fmt.Fprintf(&b, "logdiet run %s exit=%d raw=.logdiet/runs/%s\n", res.RunID, res.ExitCode, res.RunID)
	}
	if len(res.Cmd) > 0 {
		fmt.Fprintf(&b, "cmd: %s\n", textutil.JoinCommand(res.Cmd))
	}
	fmt.Fprintf(&b, "summary: %s\n", res.Summary)
	for _, item := range res.Items {
		line := item.Title
		if line == "" {
			line = item.Kind
		}
		fmt.Fprintf(&b, "%s %s\n", item.ID, line)
	}
	for _, line := range res.Lines {
		fmt.Fprintf(&b, "%s\n", line)
	}
	if res.RunID != "" {
		for _, item := range res.Items {
			fmt.Fprintf(&b, "show: logdiet show %s:%s --around 40\n", res.RunID, item.ID)
		}
		fmt.Fprintf(&b, "raw: logdiet raw %s\n", res.RunID)
		fmt.Fprintf(&b, "grep: logdiet grep %s \"pattern\"\n", res.RunID)
	} else {
		for _, item := range res.Items {
			fmt.Fprintf(&b, "show: logdiet show <run-id>:%s --around 40\n", item.ID)
		}
		fmt.Fprintf(&b, "raw: logdiet raw <run-id>\n")
		fmt.Fprintf(&b, "grep: logdiet grep <run-id> \"pattern\"\n")
	}
	stats := res.Stats
	if stats.CompactBytes == 0 {
		stats.CompactBytes = b.Len()
	}
	stats.ReductionPct = textutil.ReductionPct(stats.RawBytes, stats.CompactBytes)
	stats.ApproxCompactTokens = textutil.ApproxTokens(stats.CompactBytes)
	fmt.Fprintf(&b, "stats: raw=%dB compact=%dB approx_saved=%.1f%% approx_tokens=%d->%d\n",
		stats.RawBytes, stats.CompactBytes, stats.ReductionPct, stats.ApproxRawTokens, stats.ApproxCompactTokens)
	return b.String()
}

func SetRenderedStats(res *Result, rendered string) {
	res.Stats.CompactBytes = len([]byte(rendered))
	res.Stats.ReductionPct = textutil.ReductionPct(res.Stats.RawBytes, res.Stats.CompactBytes)
	res.Stats.ApproxCompactTokens = textutil.ApproxTokens(res.Stats.CompactBytes)
}

func hasAnyArg(cmd []string, wants ...string) bool {
	set := map[string]bool{}
	for _, w := range wants {
		set[w] = true
	}
	for _, arg := range cmd[1:] {
		if set[arg] {
			return true
		}
	}
	return false
}

func item(id, kind, title string, start, end int) EvidenceItem {
	if start <= 0 {
		start = 1
	}
	if end < start {
		end = start
	}
	file, line, _ := textutil.FirstFileLine(title)
	return EvidenceItem{
		ID: id, Kind: kind, Title: strings.TrimSpace(title), File: file, Line: line,
		Stream: "combined", StartLine: start, EndLine: end, Confidence: "medium",
	}
}

func findLine(lines []string, needle string, from int) int {
	if from < 0 {
		from = 0
	}
	for i := from; i < len(lines); i++ {
		if needle == "" || strings.Contains(lines[i], needle) {
			return i + 1
		}
	}
	return 1
}

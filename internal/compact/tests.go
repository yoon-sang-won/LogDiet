package compact

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/textutil"
)

var (
	pytestFailedRE = regexp.MustCompile(`^FAILED\s+(.+?)(?:\s+-\s+(.+))?$`)
	goFailRE       = regexp.MustCompile(`^--- FAIL:\s+([A-Za-z0-9_/-]+)`)
	cargoErrorRE   = regexp.MustCompile(`^error(?:\[[A-Za-z0-9]+\])?:\s+(.+)`)
	tscParenRE     = regexp.MustCompile(`^(.+\.[A-Za-z0-9]+)\(([0-9]+),([0-9]+)\):\s+error\s+(TS[0-9]+):\s+(.+)`)
	tscColonRE     = regexp.MustCompile(`^(.+\.[A-Za-z0-9]+):([0-9]+):([0-9]+)\s+-\s+error\s+(TS[0-9]+):\s+(.+)`)
)

func compactPytest(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	var count int
	for i, line := range lines {
		if m := pytestFailedRE.FindStringSubmatch(line); len(m) > 0 {
			count++
			title := strings.TrimSpace(m[1])
			if len(m) > 2 && strings.TrimSpace(m[2]) != "" {
				title += " " + strings.TrimSpace(m[2])
			}
			start, end := textutil.ClampRange(i-20, i+20, len(lines))
			it := item(fmt.Sprintf("F%d", count), "failure", title, start, end)
			it.Confidence = "high"
			if it.File == "" {
				for j := i + 1; j < len(lines) && j < i+8; j++ {
					if f, n, ok := textutil.FirstFileLine(lines[j]); ok {
						it.File, it.Line = f, n
						break
					}
				}
			}
			res.Items = append(res.Items, it)
		}
		if strings.Contains(line, "Traceback (most recent call last):") {
			start, end := textutil.ClampRange(i+1, i+12, len(lines))
			res.Items = append(res.Items, item(fmt.Sprintf("S%d", len(res.Items)+1), "stack", "Traceback", start, end))
		}
		if strings.Contains(line, " failed") || strings.Contains(line, " passed") {
			res.Summary = strings.TrimSpace(line)
		}
	}
	if res.Summary == "" && count > 0 {
		res.Summary = fmt.Sprintf("%d pytest failures", count)
	}
	if count == 0 {
		return compactGeneric(res, combined)
	}
	return res
}

func compactGoTest(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	var failCount, okCount int
	var failedPkgs []string
	for i, line := range lines {
		if m := goFailRE.FindStringSubmatch(line); len(m) > 0 {
			failCount++
			title := m[1]
			start, end := textutil.ClampRange(i-8, i+16, len(lines))
			it := item(fmt.Sprintf("F%d", failCount), "failure", title, start, end)
			it.Confidence = "high"
			for j := i + 1; j < len(lines) && j < i+8; j++ {
				if f, n, ok := textutil.FirstFileLine(lines[j]); ok {
					it.File, it.Line = f, n
					it.Title = fmt.Sprintf("%s %s:%d", title, f, n)
					break
				}
			}
			res.Items = append(res.Items, it)
		}
		if strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "FAIL    ") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				failedPkgs = append(failedPkgs, fields[1])
			}
		}
		if strings.HasPrefix(line, "ok\t") || strings.HasPrefix(line, "ok      ") || strings.HasPrefix(line, "?   \t") {
			okCount++
		}
		if strings.Contains(line, "panic:") {
			res.Items = append(res.Items, item(fmt.Sprintf("S%d", len(res.Items)+1), "stack", strings.TrimSpace(line), i+1, minLine(i+12, len(lines))))
		}
	}
	if failCount > 0 {
		res.Summary = fmt.Sprintf("%d failed tests", failCount)
		if len(failedPkgs) > 0 {
			res.Summary += ", failed packages: " + strings.Join(failedPkgs, ", ")
		}
		if okCount > 0 {
			res.Summary += fmt.Sprintf(", %d packages ok", okCount)
		}
		return res
	}
	if okCount > 0 {
		res.Summary = fmt.Sprintf("ok, %d packages", okCount)
		return res
	}
	return compactGeneric(res, combined)
}

func compactCargo(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	var errors int
	for i, line := range lines {
		if m := cargoErrorRE.FindStringSubmatch(line); len(m) > 0 {
			errors++
			title := strings.TrimSpace(line)
			if i+1 < len(lines) && strings.Contains(lines[i+1], "-->") {
				title += " " + strings.TrimSpace(lines[i+1])
			}
			res.Items = append(res.Items, item(fmt.Sprintf("E%d", errors), "error", title, i+1, minLine(i+8, len(lines))))
		}
		if strings.Contains(line, "panicked at") {
			res.Items = append(res.Items, item(fmt.Sprintf("S%d", len(res.Items)+1), "stack", strings.TrimSpace(line), i+1, minLine(i+10, len(lines))))
		}
		if strings.Contains(line, "test result: FAILED") {
			res.Summary = strings.TrimSpace(line)
		}
	}
	if errors > 0 && res.Summary == "" {
		res.Summary = fmt.Sprintf("%d rust errors", errors)
	}
	if len(res.Items) == 0 {
		return compactGeneric(res, combined)
	}
	return res
}

func compactJSTestOrBuild(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	var failures, tsErrors int
	for i, line := range lines {
		if m := tscParenRE.FindStringSubmatch(line); len(m) > 0 {
			tsErrors++
			res.Items = append(res.Items, item(fmt.Sprintf("E%d", tsErrors), "error", strings.TrimSpace(line), i+1, i+1))
			continue
		}
		if m := tscColonRE.FindStringSubmatch(line); len(m) > 0 {
			tsErrors++
			res.Items = append(res.Items, item(fmt.Sprintf("E%d", tsErrors), "error", strings.TrimSpace(line), i+1, i+1))
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "FAIL ") {
			failures++
			start, end := textutil.ClampRange(i-4, i+12, len(lines))
			res.Items = append(res.Items, item(fmt.Sprintf("F%d", failures), "failure", strings.TrimSpace(line), start, end))
		}
		if strings.Contains(line, "Expected:") || strings.Contains(line, "Received:") {
			res.Lines = append(res.Lines, strings.TrimSpace(line))
		}
		if strings.Contains(line, "Test Suites:") || strings.Contains(line, "Tests:") {
			res.Summary = strings.TrimSpace(line)
		}
	}
	if tsErrors > 0 {
		res.Summary = fmt.Sprintf("%d TypeScript errors", tsErrors)
	}
	if failures > 0 && res.Summary == "" {
		res.Summary = fmt.Sprintf("%d JavaScript test failures", failures)
	}
	if len(res.Items) == 0 {
		return compactGeneric(res, combined)
	}
	return res
}

func minLine(a, b int) int {
	if a < b {
		return a
	}
	return b
}

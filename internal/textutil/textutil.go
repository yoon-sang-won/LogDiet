package textutil

import (
	"fmt"
	"regexp"
	"strings"
)

func SplitLines(b []byte) []string {
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func JoinCommand(cmd []string) string {
	if len(cmd) == 0 {
		return ""
	}
	parts := make([]string, len(cmd))
	for i, p := range cmd {
		if strings.ContainsAny(p, " \t\"'") {
			parts[i] = fmt.Sprintf("%q", p)
		} else {
			parts[i] = p
		}
	}
	return strings.Join(parts, " ")
}

func ApproxTokens(bytes int) int {
	if bytes <= 0 {
		return 0
	}
	return (bytes + 3) / 4
}

func ReductionPct(raw, compact int) float64 {
	if raw <= 0 {
		return 0
	}
	saved := raw - compact
	if saved < 0 {
		saved = 0
	}
	return float64(saved) * 100 / float64(raw)
}

func ClampRange(start, end, max int) (int, int) {
	if max <= 0 {
		return 0, 0
	}
	if start < 1 {
		start = 1
	}
	if end > max {
		end = max
	}
	if start > end {
		start = end
	}
	return start, end
}

var fileLineRE = regexp.MustCompile(`([A-Za-z0-9_./\\:-]+\.[A-Za-z0-9_]+)[:(]([0-9]+)`)

func FirstFileLine(s string) (string, int, bool) {
	m := fileLineRE.FindStringSubmatch(s)
	if len(m) < 3 {
		return "", 0, false
	}
	var line int
	for _, r := range m[2] {
		line = line*10 + int(r-'0')
	}
	return m[1], line, true
}

func NormalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func LastN[T any](in []T, n int) []T {
	if n <= 0 || len(in) <= n {
		return in
	}
	return in[len(in)-n:]
}

package compact

import (
	"fmt"
	"strings"

	"logdiet/internal/textutil"
)

var signalWords = []string{
	"error", "failed", "failure", "fatal", "panic", "exception", "traceback",
	"assertion", "warning", "denied", "timeout", "cannot", "undefined", "not found",
}

func compactGeneric(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	type seenLine struct {
		line  string
		count int
		num   int
	}
	seen := map[string]int{}
	var sig []seenLine
	for i, line := range lines {
		if !isSignal(line) {
			continue
		}
		key := strings.TrimSpace(line)
		if key == "" {
			continue
		}
		if pos, ok := seen[key]; ok {
			sig[pos].count++
			continue
		}
		seen[key] = len(sig)
		sig = append(sig, seenLine{line: key, count: 1, num: i + 1})
	}
	sig = textutil.LastN(sig, 20)
	if len(sig) == 0 {
		if res.ExitCode == 0 {
			res.Summary = "ok"
		} else {
			res.Summary = fmt.Sprintf("exit %d with no recognized high-signal lines", res.ExitCode)
		}
		return res
	}
	res.Summary = fmt.Sprintf("exit %d, %d high-signal lines", res.ExitCode, len(sig))
	for i, s := range sig {
		line := s.line
		if s.count > 1 {
			line = fmt.Sprintf("%s x%d", line, s.count)
		}
		res.Lines = append(res.Lines, line)
		if res.ExitCode != 0 && i == 0 {
			it := item("F1", "failure", s.line, s.num, s.num)
			it.Confidence = "medium"
			res.Items = append(res.Items, it)
		}
	}
	return res
}

func isSignal(line string) bool {
	lower := strings.ToLower(line)
	for _, word := range signalWords {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

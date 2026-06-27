package compact

import (
	"fmt"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/textutil"
)

func compactGitDiff(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	var currentFile string
	var handle int
	adds, dels := 0, 0
	for i, line := range lines {
		if strings.HasPrefix(line, "diff --git ") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				currentFile = strings.TrimPrefix(parts[3], "b/")
				res.Lines = append(res.Lines, "file: "+currentFile)
			}
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			adds++
		}
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			dels++
		}
		if strings.HasPrefix(line, "@@") {
			handle++
			title := strings.TrimSpace(line)
			if currentFile != "" {
				title = currentFile + " " + title
			}
			start, end := textutil.ClampRange(i-3, i+12, len(lines))
			res.Items = append(res.Items, EvidenceItem{
				ID: fmt.Sprintf("D%d", handle), Kind: "diff", Title: title, File: currentFile,
				Stream: "combined", StartLine: start, EndLine: end, Confidence: "high",
			})
			res.Lines = append(res.Lines, title)
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if strings.HasPrefix(lines[j], "+") || strings.HasPrefix(lines[j], "-") {
					res.Lines = append(res.Lines, lines[j])
				}
			}
		}
		if strings.Contains(strings.ToLower(line), "binary files") {
			res.Lines = append(res.Lines, strings.TrimSpace(line))
		}
	}
	if handle == 0 {
		return compactGeneric(res, combined)
	}
	res.Summary = fmt.Sprintf("%d diff hunks, +%d -%d", handle, adds, dels)
	return res
}

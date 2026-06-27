package compact

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"logdiet/internal/textutil"
)

func compactSearch(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	type match struct {
		lineNo int
		text   string
	}
	groups := map[string][]match{}
	order := []string{}
	for i, line := range lines {
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		if _, err := strconv.Atoi(parts[1]); err != nil {
			continue
		}
		file := parts[0]
		if _, ok := groups[file]; !ok {
			order = append(order, file)
		}
		groups[file] = append(groups[file], match{i + 1, line})
	}
	if len(groups) == 0 {
		return compactGeneric(res, combined)
	}
	sort.Strings(order)
	total := 0
	for _, file := range order {
		total += len(groups[file])
	}
	res.Summary = fmt.Sprintf("%d matches in %d files", total, len(groups))
	for i, file := range order {
		matches := groups[file]
		start, end := textutil.ClampRange(matches[0].lineNo, matches[len(matches)-1].lineNo, len(lines))
		res.Items = append(res.Items, EvidenceItem{
			ID: fmt.Sprintf("G%d", i+1), Kind: "search", Title: fmt.Sprintf("%s %d matches", file, len(matches)),
			File: file, Stream: "combined", StartLine: start, EndLine: end, Confidence: "high",
		})
		limit := len(matches)
		if limit > 3 {
			limit = 3
		}
		for j := 0; j < limit; j++ {
			res.Lines = append(res.Lines, matches[j].text)
		}
		if len(matches) > limit {
			res.Lines = append(res.Lines, fmt.Sprintf("%s +%d more", file, len(matches)-limit))
		}
	}
	return res
}

package compact

import (
	"fmt"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/textutil"
)

func compactGitStatus(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	branch := ""
	staged, unstaged, untracked := 0, 0, 0
	var files []string
	section := ""
	for _, line := range lines {
		t := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(t, "On branch "):
			branch = strings.TrimPrefix(t, "On branch ")
		case strings.Contains(t, "Changes to be committed"):
			section = "staged"
		case strings.Contains(t, "Changes not staged"):
			section = "unstaged"
		case strings.Contains(t, "Untracked files"):
			section = "untracked"
		case strings.HasPrefix(t, "modified:") || strings.HasPrefix(t, "new file:") || strings.HasPrefix(t, "deleted:"):
			if section == "staged" {
				staged++
			} else {
				unstaged++
			}
			files = append(files, t)
		case section == "untracked" && t != "" && !strings.HasPrefix(t, "(") && !strings.Contains(t, "use \"git"):
			untracked++
			files = append(files, "untracked: "+t)
		}
	}
	if branch == "" && len(files) == 0 {
		return compactGeneric(res, combined)
	}
	if branch == "" {
		branch = "unknown"
	}
	res.Summary = fmt.Sprintf("branch=%s staged=%d unstaged=%d untracked=%d", branch, staged, unstaged, untracked)
	limit := len(files)
	if limit > 20 {
		limit = 20
	}
	res.Lines = append(res.Lines, files[:limit]...)
	if len(files) > limit {
		res.Lines = append(res.Lines, fmt.Sprintf("+%d more", len(files)-limit))
	}
	return res
}

func compactGitLog(res Result, combined []byte) Result {
	lines := textutil.SplitLines(combined)
	var out []string
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		if strings.HasPrefix(t, "commit ") {
			fields := strings.Fields(t)
			if len(fields) >= 2 {
				out = append(out, fields[1])
			}
			continue
		}
		if strings.HasPrefix(t, "Author:") || strings.HasPrefix(t, "Date:") {
			continue
		}
		if len(out) > 0 && !strings.Contains(out[len(out)-1], " ") {
			out[len(out)-1] += " " + t
		} else {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return compactGeneric(res, combined)
	}
	out = textutil.LastN(out, 30)
	res.Summary = fmt.Sprintf("%d commit lines", len(out))
	res.Lines = out
	return res
}

package bench

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/compact"
	"github.com/yoon-sang-won/LogDiet/internal/textutil"
)

type FixtureResult struct {
	Name                string  `json:"name"`
	RawBytes            int     `json:"raw_bytes"`
	CompactBytes        int     `json:"compact_bytes"`
	ReductionPct        float64 `json:"reduction_pct"`
	ApproxRawTokens     int     `json:"approx_raw_tokens"`
	ApproxCompactTokens int     `json:"approx_compact_tokens"`
	HandleCount         int     `json:"handle_count"`
}

func Run(root string) ([]FixtureResult, error) {
	dir := filepath.Join(root, "testdata", "fixtures")
	matches, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	var results []FixtureResult
	for _, path := range matches {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		cmd := fixtureCmd(filepath.Base(path))
		res := compact.Compact(cmd, b, nil, b, fixtureExit(filepath.Base(path)))
		rendered := compact.Render(res)
		results = append(results, FixtureResult{
			Name: filepath.Base(path), RawBytes: len(b), CompactBytes: len([]byte(rendered)),
			ReductionPct:    textutil.ReductionPct(len(b), len([]byte(rendered))),
			ApproxRawTokens: textutil.ApproxTokens(len(b)), ApproxCompactTokens: textutil.ApproxTokens(len([]byte(rendered))),
			HandleCount: len(res.Items),
		})
	}
	return results, nil
}

func Format(results []FixtureResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%-24s %9s %13s %18s %21s %10s %7s\n",
		"fixture", "raw_bytes", "compact_bytes", "approx_raw_tokens", "approx_compact_tokens", "reduction", "handles")
	for _, r := range results {
		fmt.Fprintf(&b, "%-24s %9d %13d %18d %21d %9.1f%% %7d\n",
			r.Name, r.RawBytes, r.CompactBytes, r.ApproxRawTokens, r.ApproxCompactTokens, r.ReductionPct, r.HandleCount)
	}
	b.WriteString("\napprox token estimates use ceil(bytes / 4) and are not provider billing measurements\n")
	return b.String()
}

func JSON(results []FixtureResult) (string, error) {
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", err
	}
	return string(append(b, '\n')), nil
}

func fixtureCmd(name string) []string {
	switch {
	case strings.Contains(name, "cargo"):
		return []string{"cargo", "test"}
	case strings.Contains(name, "pytest"):
		return []string{"pytest", "-q"}
	case strings.Contains(name, "go_test"):
		return []string{"go", "test", "./..."}
	case strings.Contains(name, "jest"):
		return []string{"npm", "test"}
	case strings.Contains(name, "git_status"):
		return []string{"git", "status"}
	case strings.Contains(name, "git_diff"):
		return []string{"git", "diff"}
	case strings.Contains(name, "rg"):
		return []string{"rg", "TODO"}
	case strings.Contains(name, "tsc"):
		return []string{"npm", "run", "build"}
	default:
		return []string{"make"}
	}
}

func fixtureExit(name string) int {
	if strings.Contains(name, "success") || strings.Contains(name, "status") || strings.Contains(name, "diff") || strings.Contains(name, "rg") {
		return 0
	}
	return 1
}

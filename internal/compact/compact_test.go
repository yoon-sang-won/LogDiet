package compact

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenericCompactorKeepsSignalsAndDeduplicates(t *testing.T) {
	out := []byte("progress 1\nprogress 1\nERROR src/app.go:12 cannot open file\nwarning: retrying\nERROR src/app.go:12 cannot open file\n")
	res := Compact([]string{"unknown"}, out, nil, out, 1)
	if len(res.Items) == 0 {
		t.Fatal("expected evidence for failed command")
	}
	if res.Summary == "" {
		t.Fatal("expected summary")
	}
	rendered := Render(res)
	if !contains(rendered, "x2") {
		t.Fatalf("expected repeated count in rendered output:\n%s", rendered)
	}
}

func TestSpecificCompactorsFindHandles(t *testing.T) {
	cases := []struct {
		name string
		cmd  []string
		out  string
		want string
	}{
		{
			name: "pytest",
			cmd:  []string{"pytest", "-q"},
			out:  "FAILED tests/test_api.py::test_returns_200 - AssertionError: expected 200, got 500\nE   AssertionError: expected 200, got 500\ntests/test_api.py:42: AssertionError\n2 failed, 31 passed in 1.24s\n",
			want: "F1",
		},
		{
			name: "go test",
			cmd:  []string{"go", "test", "./..."},
			out:  "--- FAIL: TestLogin (0.00s)\n    auth_test.go:17: missing token\nFAIL\tlogdiet/internal/auth\t0.011s\n",
			want: "F1",
		},
		{
			name: "tsc",
			cmd:  []string{"npm", "run", "build"},
			out:  "src/app.ts(12,5): error TS2322: Type 'string' is not assignable to type 'number'.\n",
			want: "E1",
		},
		{
			name: "rg",
			cmd:  []string{"rg", "TODO"},
			out:  "internal/a.go:10:TODO: first\ninternal/a.go:21:TODO: second\n",
			want: "G1",
		},
		{
			name: "git diff",
			cmd:  []string{"git", "diff"},
			out:  "diff --git a/a.go b/a.go\n@@ -1,2 +1,2 @@\n-old\n+new\n",
			want: "D1",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := Compact(tc.cmd, []byte(tc.out), nil, []byte(tc.out), 1)
			if !hasItem(res.Items, tc.want) {
				t.Fatalf("expected handle %s, got %#v", tc.want, res.Items)
			}
		})
	}
}

func TestFixturesProduceCompactEvidence(t *testing.T) {
	root := filepath.Join("..", "..")
	matches, err := filepath.Glob(filepath.Join(root, "testdata", "fixtures", "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("expected fixtures")
	}
	for _, path := range matches {
		name := filepath.Base(path)
		t.Run(name, func(t *testing.T) {
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			res := Compact(cmdForFixture(name), b, nil, b, exitForFixture(name))
			rendered := Render(res)
			if strings.Contains(name, "failure") || strings.Contains(name, "cargo") || strings.Contains(name, "tsc") || strings.Contains(name, "jest") {
				if len(res.Items) == 0 {
					t.Fatalf("expected handles for %s: %#v", name, res)
				}
			}
			if len(rendered) >= len(b) && strings.Contains(name, "noisy") {
				t.Fatalf("expected noisy fixture to compact below raw size: raw=%d compact=%d", len(b), len(rendered))
			}
			for _, item := range res.Items {
				if item.StartLine < 1 || item.EndLine < item.StartLine {
					t.Fatalf("invalid range for %#v", item)
				}
			}
		})
	}
}

func cmdForFixture(name string) []string {
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

func exitForFixture(name string) int {
	if strings.Contains(name, "success") || strings.Contains(name, "status") || strings.Contains(name, "diff") || strings.Contains(name, "rg") {
		return 0
	}
	return 1
}

func hasItem(items []EvidenceItem, id string) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return sub == ""
}

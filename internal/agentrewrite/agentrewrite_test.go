package agentrewrite

import (
	"reflect"
	"testing"
)

func TestDecideWrapsKnownNoisyDeveloperCommands(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want []string
	}{
		{"go test", "go test ./...", []string{"go", "test", "./..."}},
		{"pytest", "pytest -q", []string{"pytest", "-q"}},
		{"npm test", "npm test", []string{"npm", "test"}},
		{"npm run", "npm run build -- --watch=false", []string{"npm", "run", "build", "--", "--watch=false"}},
		{"git diff", "git diff --stat", []string{"git", "diff", "--stat"}},
		{"rg", `rg "TODO"`, []string{"rg", "TODO"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Decide(tc.in)
			if !got.Wrap {
				t.Fatalf("Wrap=false, reason=%q", got.Reason)
			}
			if got.Reason != "known noisy developer command" {
				t.Fatalf("Reason=%q", got.Reason)
			}
			if !reflect.DeepEqual(got.Command, tc.want) {
				t.Fatalf("Command=%#v want %#v", got.Command, tc.want)
			}
		})
	}
}

func TestDecideDoesNotWrapUnsafeOrUnselectedCommands(t *testing.T) {
	for _, tc := range []struct {
		name   string
		in     string
		reason string
	}{
		{"already logdiet", "logdiet wrap -- go test ./...", "already logdiet command"},
		{"interactive editor", "vim file.go", "interactive command"},
		{"interactive ssh", "ssh example.com", "interactive command"},
		{"unselected", "echo hello", "not selected"},
		{"shell operator", "go test ./... && echo ok", "ambiguous shell control operator"},
		{"python repl", "python", "interactive command"},
		{"node repl", "node", "interactive command"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Decide(tc.in)
			if got.Wrap {
				t.Fatalf("Wrap=true, command=%#v", got.Command)
			}
			if got.Reason != tc.reason {
				t.Fatalf("Reason=%q want %q", got.Reason, tc.reason)
			}
		})
	}
}

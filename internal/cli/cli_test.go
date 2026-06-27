package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLIWrapRawShowAndGrep(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "fake")
	body := "#!/bin/sh\necho 'alpha ok'\necho 'beta failed' >&2\nexit 5\n"
	if runtime.GOOS == "windows" {
		script = filepath.Join(dir, "fake.cmd")
		body = "@echo off\necho alpha ok\necho beta failed 1>&2\nexit /b 5\n"
	}
	if err := os.WriteFile(script, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
	oldwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	var out, errb bytes.Buffer
	code := Run([]string{"wrap", "--", script}, &out, &errb)
	if code != 5 {
		t.Fatalf("wrap exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if !strings.Contains(out.String(), "logdiet run") || !strings.Contains(out.String(), "raw:") {
		t.Fatalf("bad wrap output:\n%s", out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"raw"}, &out, &errb); code != 0 {
		t.Fatalf("raw exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "beta failed") {
		t.Fatalf("raw missing output:\n%s", out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"grep", "latest", "beta"}, &out, &errb); code != 0 {
		t.Fatalf("grep exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "beta failed") {
		t.Fatalf("grep missing match:\n%s", out.String())
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"show", "F1", "--around", "2"}, &out, &errb); code != 0 {
		t.Fatalf("show exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), "beta failed") {
		t.Fatalf("show missing raw line:\n%s", out.String())
	}
}

func TestCLICommonCommands(t *testing.T) {
	dir := t.TempDir()
	oldwd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	var out, errb bytes.Buffer
	if code := Run([]string{"--version"}, &out, &errb); code != 0 {
		t.Fatalf("version exit=%d err=%s", code, errb.String())
	}
	if strings.TrimSpace(out.String()) != "logdiet 0.1.0" {
		t.Fatalf("bad version output: %q", out.String())
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"help"}, &out, &errb); code != 0 {
		t.Fatalf("help exit=%d err=%s", code, errb.String())
	}
	help := out.String()
	for _, want := range []string{
		"LogDiet keeps full command logs locally and feeds AI coding agents compact, expandable evidence.",
		"logdiet wrap -- pytest -q",
		"logdiet show latest:F1 --around 40",
		"logdiet raw latest",
		"logdiet grep latest \"panic\"",
	} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q:\n%s", want, help)
		}
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"env", "--shell", "powershell"}, &out, &errb); code != 0 {
		t.Fatalf("env exit=%d err=%s", code, errb.String())
	}
	if !strings.Contains(out.String(), ".logdiet\\bin") {
		t.Fatalf("env missing shim path: %s", out.String())
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"install"}, &out, &errb); code != 0 {
		t.Fatalf("install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	if _, err := os.Stat(filepath.Join(dir, ".logdiet", "bin")); err != nil {
		t.Fatalf("install did not create shim dir: %v", err)
	}

	out.Reset()
	errb.Reset()
	if code := Run([]string{"rules", "--install", "codex"}, &out, &errb); code != 0 {
		t.Fatalf("rules install exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(agents), "BEGIN LOGDIET MANAGED RESPONSE CONTRACT") != 1 {
		t.Fatalf("rules install did not create one managed section:\n%s", string(agents))
	}

	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte("Generated at 2026-06-27 12:00:00\nPath "+dir+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errb.Reset()
	if code := Run([]string{"lint-instructions", "--fix"}, &out, &errb); code != 0 {
		t.Fatalf("lint fix exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	fixed, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(fixed), "Generated at") || !strings.Contains(string(fixed), "<repo>") {
		t.Fatalf("lint fix did not sanitize as expected:\n%s", string(fixed))
	}

	out.Reset()
	errb.Reset()
}

func TestCLIBenchFixturesFromRepo(t *testing.T) {
	oldwd, _ := os.Getwd()
	repo := filepath.Join("..", "..")
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	var out, errb bytes.Buffer
	if code := Run([]string{"bench-fixtures"}, &out, &errb); code != 0 {
		t.Fatalf("bench exit=%d out=%s err=%s", code, out.String(), errb.String())
	}
	for _, want := range []string{"pytest_failure.txt", "rg_matches.txt", "approx token counts"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("bench output missing %q:\n%s", want, out.String())
		}
	}
}

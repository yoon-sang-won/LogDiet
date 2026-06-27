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

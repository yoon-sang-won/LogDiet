package run

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCapturePreservesOutputAndExitCode(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "fake")
	body := "#!/bin/sh\necho out\necho err >&2\nexit 7\n"
	if runtime.GOOS == "windows" {
		script = filepath.Join(dir, "fake.cmd")
		body = "@echo off\necho out\necho err 1>&2\nexit /b 7\n"
	}
	if err := os.WriteFile(script, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
	res, err := Capture([]string{script})
	if err == nil {
		t.Fatal("expected nonzero exit error")
	}
	if res.ExitCode != 7 {
		t.Fatalf("exit=%d want 7", res.ExitCode)
	}
	if string(res.Stdout) != "out\n" && string(res.Stdout) != "out\r\n" {
		t.Fatalf("stdout=%q", res.Stdout)
	}
	if len(res.Stderr) == 0 || len(res.Combined) == 0 {
		t.Fatalf("expected stderr and combined output")
	}
}

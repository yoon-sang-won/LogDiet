package shim

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveRealCommandSkipsShimDir(t *testing.T) {
	root := t.TempDir()
	shimDir := filepath.Join(root, ".logdiet", "bin")
	realDir := filepath.Join(root, "real")
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(realDir, 0755); err != nil {
		t.Fatal(err)
	}
	name := "tool"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if err := os.WriteFile(filepath.Join(shimDir, name), []byte("shim"), 0755); err != nil {
		t.Fatal(err)
	}
	realPath := filepath.Join(realDir, name)
	if err := os.WriteFile(realPath, []byte("real"), 0755); err != nil {
		t.Fatal(err)
	}
	path := shimDir + string(os.PathListSeparator) + realDir
	found, err := ResolveRealCommand("tool", shimDir, path, "")
	if err != nil {
		t.Fatalf("ResolveRealCommand: %v", err)
	}
	if found != realPath {
		t.Fatalf("found=%q want %q", found, realPath)
	}
	clean := SanitizePATH(path, shimDir)
	if strings.Contains(clean, shimDir) {
		t.Fatalf("sanitized PATH still has shim dir: %q", clean)
	}
}

func TestInstallCreatesManagedShimsIdempotently(t *testing.T) {
	root := t.TempDir()
	exe := filepath.Join(root, "logdiet-test.exe")
	if err := os.WriteFile(exe, []byte("binary"), 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(root, exe, InstallOptions{}); err != nil {
		t.Fatalf("Install first: %v", err)
	}
	if _, err := Install(root, exe, InstallOptions{}); err != nil {
		t.Fatalf("Install second: %v", err)
	}
	bin := filepath.Join(root, ".logdiet", "bin")
	for _, cmd := range ShimCommands {
		path := filepath.Join(bin, cmd)
		if runtime.GOOS == "windows" {
			path += ".cmd"
		}
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("missing shim %s: %v", cmd, err)
		}
		if !strings.Contains(string(b), marker) {
			t.Fatalf("shim %s missing marker:\n%s", cmd, string(b))
		}
		if runtime.GOOS != "windows" {
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm()&0111 == 0 {
				t.Fatalf("unix shim %s is not executable: %v", path, info.Mode())
			}
		}
	}
}

func TestCommandCandidatesIncludesWindowsExtensions(t *testing.T) {
	got := windowsCommandCandidates("tool")
	want := []string{"tool.exe", "tool.cmd", "tool.bat", "tool.com"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("candidates=%v want %v", got, want)
	}
}

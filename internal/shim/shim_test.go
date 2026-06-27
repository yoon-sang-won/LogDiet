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

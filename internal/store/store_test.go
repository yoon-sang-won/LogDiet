package store

import (
	"os"
	"path/filepath"
	"testing"

	"logdiet/internal/compact"
)

func TestSaveRunWritesFilesAndLatest(t *testing.T) {
	dir := t.TempDir()
	data := RunData{
		RunID:    "20260627T120000Z-1234-abcd",
		CWD:      dir,
		Cmd:      []string{"fake"},
		ExitCode: 7,
		Stdout:   []byte("out\n"),
		Stderr:   []byte("err\n"),
		Combined: []byte("out\nerr\n"),
		Result: compact.Result{
			RunID:    "20260627T120000Z-1234-abcd",
			Cmd:      []string{"fake"},
			ExitCode: 7,
			Summary:  "failed",
			Items: []compact.EvidenceItem{{
				ID: "F1", Kind: "failure", Title: "failed", Stream: "combined", StartLine: 1, EndLine: 2,
			}},
		},
	}
	if err := SaveRun(dir, data); err != nil {
		t.Fatalf("SaveRun: %v", err)
	}
	for _, name := range []string{"meta.json", "stdout.txt", "stderr.txt", "combined.txt", "index.json"} {
		if _, err := os.Stat(filepath.Join(dir, ".logdiet", "runs", data.RunID, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	latest, err := LatestRunID(dir)
	if err != nil {
		t.Fatalf("LatestRunID: %v", err)
	}
	if latest != data.RunID {
		t.Fatalf("latest=%q want %q", latest, data.RunID)
	}
	idx, err := LoadIndex(dir, data.RunID)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	if idx.Summary != "failed" || len(idx.Items) != 1 {
		t.Fatalf("bad index: %#v", idx)
	}
}

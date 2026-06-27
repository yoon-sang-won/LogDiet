package store

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"logdiet/internal/compact"
	"logdiet/internal/version"
)

type RunData struct {
	RunID     string
	CWD       string
	Cmd       []string
	StartedAt time.Time
	EndedAt   time.Time
	ExitCode  int
	Stdout    []byte
	Stderr    []byte
	Combined  []byte
	Result    compact.Result
}

func StateDir(root string) string {
	return filepath.Join(root, ".logdiet")
}

func RunsDir(root string) string {
	return filepath.Join(StateDir(root), "runs")
}

func BackupDir(root string) string {
	return filepath.Join(StateDir(root), "backup")
}

func GenerateRunID() string {
	var b [2]byte
	if _, err := rand.Read(b[:]); err != nil {
		now := time.Now().UnixNano()
		b[0] = byte(now)
		b[1] = byte(now >> 8)
	}
	return fmt.Sprintf("%s-%d-%s", time.Now().UTC().Format("20060102T150405Z"), os.Getpid(), hex.EncodeToString(b[:]))
}

func EnsureState(root string) error {
	if err := os.MkdirAll(filepath.Join(StateDir(root), "bin"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(RunsDir(root), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(BackupDir(root), 0755); err != nil {
		return err
	}
	gitignore := filepath.Join(StateDir(root), ".gitignore")
	if _, err := os.Stat(gitignore); errors.Is(err, os.ErrNotExist) {
		return os.WriteFile(gitignore, []byte("runs/\nbackup/\n*.log\n"), 0644)
	}
	return nil
}

func SaveRun(root string, data RunData) error {
	if data.RunID == "" {
		data.RunID = GenerateRunID()
	}
	if data.StartedAt.IsZero() {
		data.StartedAt = time.Now().UTC()
	}
	if data.EndedAt.IsZero() {
		data.EndedAt = data.StartedAt
	}
	if data.CWD == "" {
		data.CWD = root
	}
	if err := EnsureState(root); err != nil {
		return err
	}
	runDir := filepath.Join(RunsDir(root), data.RunID)
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return err
	}
	meta := Meta{
		RunID: data.RunID, Version: version.Version, CWD: data.CWD, Cmd: data.Cmd,
		StartedAt:  data.StartedAt.UTC().Format(time.RFC3339Nano),
		EndedAt:    data.EndedAt.UTC().Format(time.RFC3339Nano),
		DurationMS: data.EndedAt.Sub(data.StartedAt).Milliseconds(),
		ExitCode:   data.ExitCode, StdoutBytes: len(data.Stdout), StderrBytes: len(data.Stderr), CombinedBytes: len(data.Combined),
	}
	idx := Index{
		RunID: data.RunID, Cmd: data.Cmd, ExitCode: data.ExitCode, Summary: data.Result.Summary,
		RawPaths: RawPaths{Stdout: "stdout.txt", Stderr: "stderr.txt", Combined: "combined.txt"},
		Stats:    data.Result.Stats, Items: data.Result.Items,
	}
	if err := writeJSON(filepath.Join(runDir, "meta.json"), meta); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(runDir, "stdout.txt"), data.Stdout, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(runDir, "stderr.txt"), data.Stderr, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(runDir, "combined.txt"), data.Combined, 0644); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(runDir, "index.json"), idx); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(StateDir(root), "latest"), []byte(data.RunID+"\n"), 0644)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0644)
}

func LatestRunID(root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(StateDir(root), "latest"))
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(b))
	if id == "" {
		return "", fmt.Errorf("latest run pointer is empty")
	}
	return id, nil
}

func ResolveRunID(root, id string) (string, error) {
	if id == "" || id == "latest" {
		return LatestRunID(root)
	}
	return id, nil
}

func LoadIndex(root, runID string) (Index, error) {
	var idx Index
	runID, err := ResolveRunID(root, runID)
	if err != nil {
		return idx, err
	}
	b, err := os.ReadFile(filepath.Join(RunsDir(root), runID, "index.json"))
	if err != nil {
		return idx, err
	}
	if err := json.Unmarshal(b, &idx); err != nil {
		return idx, err
	}
	return idx, nil
}

func ReadRaw(root, runID, stream string) ([]byte, error) {
	runID, err := ResolveRunID(root, runID)
	if err != nil {
		return nil, err
	}
	name := "combined.txt"
	switch stream {
	case "stdout":
		name = "stdout.txt"
	case "stderr":
		name = "stderr.txt"
	case "combined", "":
		name = "combined.txt"
	default:
		return nil, fmt.Errorf("unknown stream %q", stream)
	}
	return os.ReadFile(filepath.Join(RunsDir(root), runID, name))
}

func RunDir(root, runID string) string {
	return filepath.Join(RunsDir(root), runID)
}

func BackupPath(root string) string {
	return filepath.Join(BackupDir(root), time.Now().UTC().Format("20060102T150405Z"))
}

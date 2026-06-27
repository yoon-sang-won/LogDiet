package store

import "github.com/yoon-sang-won/LogDiet/internal/compact"

type RawPaths struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Combined string `json:"combined"`
}

type Index struct {
	RunID    string                 `json:"run_id"`
	Cmd      []string               `json:"cmd"`
	ExitCode int                    `json:"exit_code"`
	Summary  string                 `json:"summary"`
	RawPaths RawPaths               `json:"raw_paths"`
	Stats    compact.Stats          `json:"stats"`
	Items    []compact.EvidenceItem `json:"items"`
}

type Meta struct {
	RunID         string   `json:"run_id"`
	Version       string   `json:"version"`
	CWD           string   `json:"cwd"`
	Cmd           []string `json:"cmd"`
	StartedAt     string   `json:"started_at"`
	EndedAt       string   `json:"ended_at"`
	DurationMS    int64    `json:"duration_ms"`
	ExitCode      int      `json:"exit_code"`
	StdoutBytes   int      `json:"stdout_bytes"`
	StderrBytes   int      `json:"stderr_bytes"`
	CombinedBytes int      `json:"combined_bytes"`
}

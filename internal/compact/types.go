package compact

type Result struct {
	RunID    string         `json:"run_id"`
	Cmd      []string       `json:"cmd"`
	ExitCode int            `json:"exit_code"`
	Summary  string         `json:"summary"`
	Items    []EvidenceItem `json:"items"`
	Lines    []string       `json:"lines"`
	Stats    Stats          `json:"stats"`
}

type EvidenceItem struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Title      string `json:"title"`
	File       string `json:"file,omitempty"`
	Line       int    `json:"line,omitempty"`
	Stream     string `json:"stream"`
	StartLine  int    `json:"start_line"`
	EndLine    int    `json:"end_line"`
	Confidence string `json:"confidence"`
}

type Stats struct {
	RawBytes            int     `json:"raw_bytes"`
	CompactBytes        int     `json:"compact_bytes"`
	ReductionPct        float64 `json:"reduction_pct"`
	ApproxRawTokens     int     `json:"approx_raw_tokens"`
	ApproxCompactTokens int     `json:"approx_compact_tokens"`
}

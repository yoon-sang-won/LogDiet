package run

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type CaptureResult struct {
	Stdout    []byte
	Stderr    []byte
	Combined  []byte
	ExitCode  int
	StartedAt time.Time
	EndedAt   time.Time
}

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code %d", e.Code)
	}
	return e.Err.Error()
}

func (e *ExitError) Unwrap() error { return e.Err }

func Capture(args []string) (CaptureResult, error) {
	var res CaptureResult
	if len(args) == 0 {
		res.ExitCode = 2
		return res, &ExitError{Code: 2, Err: fmt.Errorf("missing command")}
	}
	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	res.StartedAt = time.Now().UTC()
	waitErr := cmd.Run()
	res.EndedAt = time.Now().UTC()
	res.Stdout = stdout.Bytes()
	res.Stderr = stderr.Bytes()
	res.Combined = append(append([]byte{}, res.Stdout...), res.Stderr...)
	if waitErr != nil {
		var ee *exec.ExitError
		if errors.As(waitErr, &ee) {
			res.ExitCode = ee.ExitCode()
			return res, &ExitError{Code: res.ExitCode, Err: waitErr}
		}
		res.ExitCode = 127
		return res, &ExitError{Code: 127, Err: waitErr}
	}
	res.ExitCode = 0
	return res, nil
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var ee *ExitError
	if errors.As(err, &ee) {
		return ee.Code
	}
	return 1
}

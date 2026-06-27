package run

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		res.ExitCode = 1
		return res, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		res.ExitCode = 1
		return res, err
	}
	var stdout, stderr bytes.Buffer
	res.StartedAt = time.Now().UTC()
	if err := cmd.Start(); err != nil {
		res.EndedAt = time.Now().UTC()
		res.ExitCode = 127
		return res, &ExitError{Code: 127, Err: err}
	}
	outDone := make(chan error, 1)
	errDone := make(chan error, 1)
	go func() { _, e := io.Copy(&stdout, stdoutPipe); outDone <- e }()
	go func() { _, e := io.Copy(&stderr, stderrPipe); errDone <- e }()
	waitErr := cmd.Wait()
	copyErr1 := <-outDone
	copyErr2 := <-errDone
	res.EndedAt = time.Now().UTC()
	res.Stdout = stdout.Bytes()
	res.Stderr = stderr.Bytes()
	res.Combined = append(append([]byte{}, res.Stdout...), res.Stderr...)
	if copyErr1 != nil {
		res.ExitCode = 1
		return res, copyErr1
	}
	if copyErr2 != nil {
		res.ExitCode = 1
		return res, copyErr2
	}
	if waitErr != nil {
		var ee *exec.ExitError
		if errors.As(waitErr, &ee) {
			res.ExitCode = ee.ExitCode()
			return res, &ExitError{Code: res.ExitCode, Err: waitErr}
		}
		res.ExitCode = 1
		return res, waitErr
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

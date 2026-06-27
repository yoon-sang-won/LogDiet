package run

import (
	"io"
	"os"
	"os/exec"
)

func Execute(args []string, env []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return 2
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	if env != nil {
		cmd.Env = env
	}
	if err := cmd.Run(); err != nil {
		return ExitCode(err)
	}
	return 0
}

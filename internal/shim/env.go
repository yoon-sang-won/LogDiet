package shim

import (
	"fmt"
	"runtime"
)

func Env(shell string) string {
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "powershell"
		} else {
			shell = "sh"
		}
	}
	switch shell {
	case "fish":
		return "set -gx PATH \"$PWD/.logdiet/bin\" $PATH\n"
	case "powershell", "pwsh":
		return "$env:PATH = \"$PWD\\.logdiet\\bin;$env:PATH\"\n"
	case "cmd":
		return "set PATH=%CD%\\.logdiet\\bin;%PATH%\n"
	default:
		return "export PATH=\"$PWD/.logdiet/bin:$PATH\"\n"
	}
}

func ActivationInstructions() string {
	return fmt.Sprintf("activate POSIX:\n%spowershell:\n%scmd:\n%s", Env("sh"), Env("powershell"), Env("cmd"))
}

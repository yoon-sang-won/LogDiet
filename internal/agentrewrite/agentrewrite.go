package agentrewrite

import (
	"path"
	"strings"
	"unicode"
)

const (
	ReasonKnownNoisy     = "known noisy developer command"
	ReasonAlreadyLogDiet = "already logdiet command"
	ReasonShellOperator  = "ambiguous shell control operator"
	ReasonInteractive    = "interactive command"
	ReasonNotSelected    = "not selected"
)

type Decision struct {
	Wrap    bool
	Reason  string
	Command []string
}

func Decide(command string) Decision {
	command = strings.TrimSpace(command)
	if command == "" {
		return noWrap(ReasonNotSelected)
	}
	if hasShellControlOperator(command) {
		return noWrap(ReasonShellOperator)
	}
	args, ok := splitCommand(command)
	if !ok || len(args) == 0 {
		return noWrap(ReasonNotSelected)
	}
	name := commandName(args[0])
	if name == "logdiet" {
		return noWrap(ReasonAlreadyLogDiet)
	}
	if isInteractive(name, args) {
		return noWrap(ReasonInteractive)
	}
	if shouldWrap(name, args) {
		return Decision{Wrap: true, Reason: ReasonKnownNoisy, Command: args}
	}
	return noWrap(ReasonNotSelected)
}

func noWrap(reason string) Decision {
	return Decision{Reason: reason}
}

func shouldWrap(name string, args []string) bool {
	switch name {
	case "rg", "grep", "pytest":
		return true
	case "go":
		return len(args) > 1 && args[1] == "test"
	case "npm":
		return len(args) > 1 && (args[1] == "test" || args[1] == "run")
	case "pnpm", "yarn", "cargo", "mvn", "gradle":
		return len(args) > 1 && args[1] == "test"
	case "git":
		return len(args) > 1 && oneOf(args[1], "diff", "status", "log")
	default:
		return false
	}
}

func isInteractive(name string, args []string) bool {
	switch name {
	case "vim", "nano", "less", "top", "ssh":
		return true
	case "python", "python3", "node":
		return len(args) == 1
	default:
		return false
	}
}

func commandName(arg string) string {
	arg = strings.ReplaceAll(arg, "\\", "/")
	base := strings.ToLower(path.Base(arg))
	for _, ext := range []string{".exe", ".cmd", ".bat", ".com"} {
		base = strings.TrimSuffix(base, ext)
	}
	return base
}

func splitCommand(command string) ([]string, bool) {
	var args []string
	var b strings.Builder
	runes := []rune(command)
	var quote rune
	inArg := false
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case quote == 0 && unicode.IsSpace(r):
			if inArg {
				args = append(args, b.String())
				b.Reset()
				inArg = false
			}
		case quote == 0 && (r == '\'' || r == '"'):
			quote = r
			inArg = true
		case quote != 0 && r == quote:
			quote = 0
		case r == '\\' && quote != '\'' && i+1 < len(runes) && canEscape(quote, runes[i+1]):
			i++
			b.WriteRune(runes[i])
			inArg = true
		default:
			b.WriteRune(r)
			inArg = true
		}
	}
	if quote != 0 {
		return nil, false
	}
	if inArg {
		args = append(args, b.String())
	}
	return args, true
}

func canEscape(quote rune, next rune) bool {
	if quote == '"' {
		return next == '"' || next == '\\' || next == '$' || next == '`'
	}
	return unicode.IsSpace(next) || next == '\'' || next == '"' || next == '\\'
}

func hasShellControlOperator(command string) bool {
	runes := []rune(command)
	var quote rune
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if quote == 0 && (r == '\'' || r == '"') {
			quote = r
			continue
		}
		if quote != 0 {
			if r == quote {
				quote = 0
			}
			if r == '\\' && quote == '"' && i+1 < len(runes) {
				i++
			}
			continue
		}
		switch r {
		case ';', '|', '<', '>', '`', '\n', '\r':
			return true
		case '&':
			return true
		case '$':
			if i+1 < len(runes) && runes[i+1] == '(' {
				return true
			}
		}
	}
	return false
}

func oneOf(s string, vals ...string) bool {
	for _, v := range vals {
		if s == v {
			return true
		}
	}
	return false
}

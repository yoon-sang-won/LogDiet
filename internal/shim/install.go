package shim

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yoon-sang-won/LogDiet/internal/store"
)

const marker = "LOGDIET MANAGED SHIM"

type InstallOptions struct {
	Force bool
}

func Install(root, exePath string, opts InstallOptions) (string, error) {
	if err := store.EnsureState(root); err != nil {
		return "", err
	}
	if exePath == "" {
		var err error
		exePath, err = os.Executable()
		if err != nil {
			return "", err
		}
	}
	bin := filepath.Join(store.StateDir(root), "bin")
	var written []string
	for _, cmd := range ShimCommands {
		path := shimPath(bin, cmd)
		if err := writeShim(path, exePath, cmd, opts.Force); err != nil {
			return "", err
		}
		written = append(written, path)
	}
	return fmt.Sprintf("installed %d shims in %s\n%s", len(written), bin, ActivationInstructions()), nil
}

func Uninstall(root string) (string, error) {
	bin := filepath.Join(store.StateDir(root), "bin")
	var removed []string
	for _, cmd := range ShimCommands {
		path := shimPath(bin, cmd)
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if !strings.Contains(string(b), marker) {
			continue
		}
		if err := os.Remove(path); err != nil {
			return "", err
		}
		removed = append(removed, filepath.Base(path))
	}
	if len(removed) == 0 {
		return "removed 0 shims; .logdiet/runs left intact\n", nil
	}
	return fmt.Sprintf("removed shims: %s\nleft intact: .logdiet/runs\n", strings.Join(removed, ", ")), nil
}

func shimPath(bin, cmd string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(bin, cmd+".cmd")
	}
	return filepath.Join(bin, cmd)
}

func writeShim(path, exePath, cmd string, force bool) error {
	if b, err := os.ReadFile(path); err == nil {
		if !strings.Contains(string(b), marker) {
			if !force {
				return fmt.Errorf("%s exists and is not a LogDiet shim; use --force to overwrite managed shims only", path)
			}
			return fmt.Errorf("%s exists and is not a LogDiet shim", path)
		}
	}
	var content string
	if runtime.GOOS == "windows" {
		content = fmt.Sprintf("@echo off\r\nrem %s\r\n\"%s\" shim --shim-dir \"%%~dp0\" -- \"%s\" %%*\r\n", marker, exePath, cmd)
		return os.WriteFile(path, []byte(content), 0644)
	}
	content = fmt.Sprintf("#!/bin/sh\n# %s\nexec \"%s\" shim --shim-dir \"$(CDPATH= cd -- \"$(dirname -- \"$0\")\" && pwd)\" -- \"%s\" \"$@\"\n", marker, exePath, cmd)
	return os.WriteFile(path, []byte(content), 0755)
}

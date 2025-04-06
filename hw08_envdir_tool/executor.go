package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}
	newCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	newCmd.Env = updateEnv(env)

	newCmd.Stderr = os.Stderr
	newCmd.Stdin = os.Stdin
	newCmd.Stdout = os.Stdout

	err := newCmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}

func updateEnv(env Environment) []string {
	originEnv := os.Environ()

	finalEnv := make([]string, 0, len(originEnv)+len(env))

	customEnvMap := make(map[string]struct{}, len(env))

	for key, v := range env {
		customEnvMap[key] = struct{}{}
		if v.NeedRemove {
			continue
		}
		finalEnv = append(finalEnv, fmt.Sprintf("%v=%v", key, v.Value))
	}
	for _, pair := range originEnv {
		parts := strings.SplitN(pair, "=", 2)

		if _, exists := customEnvMap[parts[0]]; !exists {
			finalEnv = append(finalEnv, pair)
		}
	}
	return finalEnv
}

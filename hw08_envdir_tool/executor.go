package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}
	newCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	newCmd.Env = getEnvs(env)

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

func getEnvs(env Environment) []string {
	res := make([]string, 0, len(env))
	for key, v := range env {
		if v.NeedRemove {
			continue
		}
		res = append(res, fmt.Sprintf("%v=%v", key, v.Value))
	}
	return res
}

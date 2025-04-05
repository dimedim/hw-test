package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrTODO = errors.New("TODO")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntrys, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	envMap := make(Environment)

	for _, dirEntr := range dirEntrys {
		var envVal EnvValue
		name := dirEntr.Name()

		if strings.Contains(name, "=") {
			continue
		}
		fileData, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}

		if len(fileData) == 0 {
			envVal.NeedRemove = true
			envMap[name] = envVal
			continue
		}
		fileData = bytes.TrimRight(fileData, " \t")
		idx := bytes.Index(fileData, []byte{'\n'})

		if idx != -1 {
			fileData = fileData[:idx]
		}
		fileData = bytes.ReplaceAll(fileData, []byte("\000"), []byte{'\n'})

		// fileData = bytes.TrimSpace(fileData)
		envVal.Value = string(fileData)
		envMap[name] = envVal
	}
	return envMap, nil
}

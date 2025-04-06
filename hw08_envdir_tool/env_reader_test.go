package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("./", "readDirTest")
	if err != nil {
		t.Fatalf("cannot create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	testsFiles := []struct {
		name    string
		content []byte
	}{
		{
			name:    "VAR1",
			content: []byte("value1\nsecond-line\nthird-line"),
		},
		{
			name:    "VAR2",
			content: []byte("  value \t "),
		},
		{
			name:    "UNSET",
			content: []byte{},
		},
		{
			name:    "NULL_BYTE",
			content: []byte("some\000value\nanother"),
		},
		{
			name:    "EMPTY",
			content: []byte("\n"),
		},
		{
			name:    "EQUALS=IN_NAME",
			content: []byte("ignored"),
		},
	}

	for _, testFile := range testsFiles {
		filePath := filepath.Join(tempDir, testFile.name)
		if err := os.WriteFile(filePath, testFile.content, 0o644); err != nil {
			t.Fatalf("WriteFile: %s: %v", filePath, err)
		}
	}

	env, err := ReadDir(tempDir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	expected := Environment{
		"VAR1": EnvValue{
			Value:      "value1",
			NeedRemove: false,
		},
		"VAR2": EnvValue{
			Value:      "  value",
			NeedRemove: false,
		},
		"EMPTY": EnvValue{
			Value:      "",
			NeedRemove: false,
		},
		"UNSET": EnvValue{
			Value:      "",
			NeedRemove: true,
		},
		"NULL_BYTE": EnvValue{
			Value:      "some\nvalue",
			NeedRemove: false,
		},
	}

	if len(env) != len(expected) {
		t.Fatalf("expected %d env vars, got %d", len(expected), len(env))
	}

	for k, v := range expected {
		got, ok := env[k]
		if !ok {
			t.Errorf("expected key %q not found in result", k)
			continue
		}

		require.Equal(t, got, v)
	}
}

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	DIR     = "./testdata"
	TMPFILE = "testCopy"
)

type TestCase struct {
	name          string
	fromPath      string
	offset        int64
	limit         int64
	expectedError error
	expectedFile  string
}

func TestCopyCustom(t *testing.T) {
	TC := []TestCase{
		{
			name:         "common",
			fromPath:     filepath.Join(DIR, "input.txt"),
			offset:       0,
			limit:        0,
			expectedFile: filepath.Join(DIR, "out_offset0_limit0.txt"),
		},
		{
			name:         "big data",
			fromPath:     filepath.Join(DIR, "big_data.txt"),
			offset:       0,
			limit:        0,
			expectedFile: filepath.Join(DIR, "big_data.txt"),
		},
		{
			name:         "limit is more than file size",
			fromPath:     filepath.Join(DIR, "input.txt"),
			offset:       0,
			limit:        9999999,
			expectedFile: filepath.Join(DIR, "out_offset0_limit0.txt"),
		},
		{
			name:         "offset 20",
			fromPath:     filepath.Join(DIR, "mini_data.txt"),
			offset:       20,
			limit:        0,
			expectedFile: filepath.Join(DIR, "mini_of20.txt"),
		},
		{
			name:         "limit 200",
			fromPath:     filepath.Join(DIR, "input.txt"),
			offset:       0,
			limit:        200,
			expectedFile: filepath.Join(DIR, "limit20.txt"),
		},
		{
			name:         "limit 2000 offset 20",
			fromPath:     filepath.Join(DIR, "input.txt"),
			offset:       20,
			limit:        2000,
			expectedFile: filepath.Join(DIR, "limit2000_of20.txt"),
		},
		{
			name:          "empty file path",
			fromPath:      "",
			offset:        0,
			limit:         0,
			expectedError: ErrEmptyFilePath,
		},
		{
			name:          "no such file or directory",
			fromPath:      "NotExist.txt",
			offset:        0,
			limit:         0,
			expectedError: os.ErrNotExist,
		},
		{
			name:          "directory instead of file",
			fromPath:      DIR,
			offset:        0,
			limit:         0,
			expectedError: ErrUnsupportedFile,
		},
		{
			name:          "unsupported file",
			fromPath:      "/dev/urandom",
			offset:        0,
			limit:         0,
			expectedError: ErrUnsupportedFile,
		},
		{
			name:          "offset is more than the file size",
			fromPath:      filepath.Join(DIR, "mini_data.txt"),
			offset:        100,
			limit:         0,
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			name:          "unexpected error",
			offset:        0,
			limit:         0,
			expectedError: nil,
		},
	}

	for _, tc := range TC {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp(DIR, TMPFILE)
			require.NoError(t, err)
			dstName := tmpFile.Name()
			defer os.Remove(dstName)
			tmpFile.Close()

			err = Copy(tc.fromPath, dstName, tc.offset, tc.limit)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
				return
			}
			if tc.expectedFile == "" {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			expected, err := os.ReadFile(tc.expectedFile)
			require.NoError(t, err)
			actual, err := os.ReadFile(dstName)
			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	}

	t.Run("there is no new file", func(t *testing.T) {
		filename := filepath.Join(DIR, "mini_data.txt")
		err := Copy(filename, filename, 0, 0)
		require.ErrorIs(t, err, ErrNoNewFile)
	})
	t.Run("empty toPath file", func(t *testing.T) {
		toPath := "filename.txt"
		filename := filepath.Join(DIR, "mini_data.txt")
		err := Copy(filename, toPath, 0, 0)
		require.NoError(t, err)
		err = os.Remove(toPath)
		require.NoError(t, err)
	})
	t.Run("same file, but different filename", func(t *testing.T) {
		filename := filepath.Join(DIR, "mini_data.txt")
		err := Copy(filename, "./"+filename, 0, 0)
		require.ErrorIs(t, err, ErrNoNewFile)
	})
}

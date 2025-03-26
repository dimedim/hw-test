package main

// cd hw07_file_copying/

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyFilePath         = errors.New("empty file path")
	ErrNoNewFile             = errors.New("there is no new file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrEmptyFilePath
	}
	if fromPath == toPath {
		return ErrNoNewFile
	}

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return fmt.Errorf("get stat: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer src.Close()
	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek func: %w", err)
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(src)

	_, err = io.CopyN(dst, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		dst.Close()
		os.Remove(toPath)
		return fmt.Errorf("copy file: %w", err)
	}
	dst.Close()

	return nil
}

package logger

import (
	"io"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

func OpenLogFile(filepath string) (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	file, err := os.OpenFile(filepath, flags, 0o644)
	if err != nil {
		slog.Error("filepath ERROR", "err", err)
	}

	return file, err
}

func New(level string, output io.Writer) Logger {
	var slogLvl slog.Level

	switch level {
	case "info":
		slogLvl = slog.LevelInfo
	case "debug":
		slogLvl = slog.LevelDebug
	case "warn":
		slogLvl = slog.LevelWarn
	case "error":
		slogLvl = slog.LevelError
	default:
		slogLvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLvl,
	}
	logger := slog.New(slog.NewTextHandler(output, opts))
	slog.SetDefault(logger)
	return logger
}

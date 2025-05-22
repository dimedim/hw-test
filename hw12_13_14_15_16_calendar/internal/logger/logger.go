package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type LoggerTodo interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Logger struct { // TODO
	LoggerTodo
}

func New(level string) LoggerTodo {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// // TODO
// func NewSlogLogger(level slog.Level, logsType LogsType) *Logger {
// 	opts := &slog.HandlerOptions{
// 		Level: level,
// 		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
// 			if a.Key == slog.TimeKey {
// 				if t, ok := a.Value.Any().(time.Time); ok {
// 					return slog.String(slog.TimeKey, t.Format("2006-01-02 15:04:05"))
// 				}
// 			}
// 			return a
// 		},
// 	}

// 	var handler slog.Handler
// 	switch logsType {
// 	case JSONLogs:
// 		handler = slog.NewJSONHandler(os.Stdout, opts)
// 	default:
// 		handler = slog.NewTextHandler(os.Stdout, opts)
// 	}

// 	logger := slog.New(handler)
// 	slog.SetDefault(logger)
// 	logger.Info("logger started")
// 	return &Logger{logger}
// }

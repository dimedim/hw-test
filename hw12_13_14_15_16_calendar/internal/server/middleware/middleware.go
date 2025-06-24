package mware

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
	"github.com/gorilla/mux"
)

func LoggingMiddleware(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			statusRec := newStatusRecorder(w)

			next.ServeHTTP(statusRec, r)

			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if ip == "" {
				ip = r.RemoteAddr
			}

			userAgent := r.UserAgent()
			if userAgent == "" {
				userAgent = "-"
			}
			if statusRec.status >= 500 {
				logger.Error(
					"msg",
					slog.String("ip", ip),
					slog.String("method", r.Method),
					slog.String("uri", r.RequestURI),
					slog.String("protocol", r.Proto),
					slog.String("duration", time.Since(start).String()),
					slog.Int("status_code", statusRec.status),
					slog.String("user_agent", userAgent),
				)
			} else {
				logger.Info(
					"msg",
					slog.String("ip", ip),
					slog.String("method", r.Method),
					slog.String("uri", r.RequestURI),
					slog.String("protocol", r.Proto),
					slog.String("duration", time.Since(start).String()),
					slog.Int("status_code", statusRec.status),
					slog.String("user_agent", userAgent),
				)
			}
		})
	}
}

func PanicRecover(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered", slog.Any("ERROR", err))
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

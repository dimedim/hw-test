package interceptors

import (
	"context"
	"log/slog"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LogInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		clientAddr := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			clientAddr = p.Addr.String()
		}

		now := time.Now()
		resp, err = handler(ctx, req)
		duration := time.Since(now)

		code := status.Code(err)

		if err != nil {
			logger.Error("GRPC error",
				slog.String("method", info.FullMethod),
				slog.String("duration", duration.String()),
				slog.String("client", clientAddr),
				slog.String("code", code.String()),
				slog.String("ERR", err.Error()),
			)
		} else {
			logger.Info("GRPC response",
				slog.String("method", info.FullMethod),
				slog.String("duration", duration.String()),
				slog.String("client", clientAddr),
				slog.String("code", code.String()),
			)
		}
		return resp, err
	}
}

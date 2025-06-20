package internalgrpc

import (
	"context"
	"log"
	"log/slog"
	"net"
	"time"

	ge "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/api/grpcevents"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/server/interceptors"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
	"google.golang.org/grpc"
)

// type EventStorage interface {
// 	CreateEvent(ctx context.Context, e *models.Event) (*models.Event, error)
// 	UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error)
// 	DeleteEvent(ctx context.Context, eventID string) error

// 	ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error)
// 	ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error)
// 	ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error)

// 	Close() error
// }

type EventApp interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error)
	DeleteEvent(ctx context.Context, eventID string) error

	ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error)
	ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error)
	ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error)
}

//	type GRPCServer struct {
//		Cfg *config.Config
//		ge.UnimplementedCalendarServiceServer
//		Store storage.EventStorage
//	}
type GRPCServer struct {
	Cfg *config.Config
	ge.UnimplementedCalendarServiceServer
	App EventApp
}

func New(app EventApp, cfg *config.Config) *GRPCServer {
	return &GRPCServer{App: app, Cfg: cfg}
}

func (s *GRPCServer) Start(ctx context.Context, cfg *config.Config, logger logger.Logger) {

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.LogInterceptor(logger)),
	)
	ge.RegisterCalendarServiceServer(grpcSrv, s)

	listener, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("GRPC starts", slog.String("port", cfg.GRPC.Port))

	go func() {
		if err := grpcSrv.Serve(listener); err != nil {
			log.Fatal("GRPC server stoped", err)
		}
	}()
	go func() {
		<-ctx.Done()
		slog.Info("GRPC graceful shutting down")
		grpcSrv.GracefulStop()
	}()
}

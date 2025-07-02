package storage

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/database"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	memorystorage "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/storage/memory"
	sqlstorage "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/storage/sql"
	"github.com/pressly/goose/v3"
)

type EventStorage interface {
	CreateEvent(ctx context.Context, e *models.Event) (*models.Event, error)
	UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error)
	DeleteEvent(ctx context.Context, eventID string) error

	ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error)
	ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error)
	ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error)

	Close() error

	SchedulerStorrage
}

// TODO: Заглушка
type SchedulerStorrage interface {
	// notifyAt = event.StartsAt - event.NotifyBefore
	//  и если now >= notifyAt, то событие нужно вернуть.
	ListEventsToNotify(ctx context.Context, now time.Time) ([]*models.Event, error)
	DeleteOlderThan(ctx context.Context, expire time.Time) (int, error)
}

func NewStorage(ctx context.Context, dbType, dsn, migrationsPath string) EventStorage {
	switch dbType {
	case "memory":
		return memorystorage.New()
	case "postgres":
		pgxConn := database.MustConnectDatabase(ctx, dsn)
		psqlStorage := sqlstorage.New(pgxConn)
		if err := psqlStorage.Connect(ctx); err != nil {
			log.Fatal("cant connect to db: ", err)
		}
		if migrationsPath != "" {
			migrate(ctx, pgxConn.DB, migrationsPath)
		}
		return psqlStorage
	}
	slog.Warn("storage type not set", slog.String("type", dbType))
	return memorystorage.New()
}

func migrate(ctx context.Context, db *sql.DB, migrationsPath string) {
	err := goose.UpContext(ctx, db, migrationsPath)
	if err != nil {
		log.Fatal("migration error: %w", err)
	}
	// if err := goose.DownContext(ctx, db, migrationsPath); err != nil {
	// 	log.Fatal("down migration: %w", err)
	// }
}

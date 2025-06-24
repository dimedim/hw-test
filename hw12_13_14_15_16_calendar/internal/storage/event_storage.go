package storage

import (
	"context"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
)

type EventStorage interface {
	CreateEvent(ctx context.Context, e *models.Event) (*models.Event, error)
	UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error)
	DeleteEvent(ctx context.Context, eventID string) error

	ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error)
	ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error)
	ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error)

	Close() error
}

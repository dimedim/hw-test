package app

import (
	"context"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/google/uuid"
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

type App struct {
	Repo EventStorage
}

func New(storage EventStorage) *App {
	return &App{
		Repo: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	event.ID = uuid.NewString()
	return a.Repo.CreateEvent(ctx, event)
}

// TODO: not impl
func (a *App) UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error) {
	return nil, nil
}
func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return nil
}
func (a *App) ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error) {
	return nil, nil
}
func (a *App) ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error) {
	return nil, nil
}
func (a *App) ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error) {
	return nil, nil
}

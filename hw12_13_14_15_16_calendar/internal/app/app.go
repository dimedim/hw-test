package app

import (
	"context"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	Repo storage.EventStorage
}

func New(storage storage.EventStorage) *App {
	return &App{
		Repo: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	event.ID = uuid.NewString()
	return a.Repo.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error) {
	return a.Repo.UpdateEvent(ctx, eventID, e)
}

func (a *App) DeleteEvent(ctx context.Context, eventID string) error {
	return a.Repo.DeleteEvent(ctx, eventID)
}

func (a *App) ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error) {
	return a.Repo.ListEventsByDay(ctx, userID, day)
}

func (a *App) ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error) {
	return a.Repo.ListEventsByWeek(ctx, userID, week)
}

func (a *App) ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error) {
	return a.Repo.ListEventsByMonth(ctx, userID, month)
}

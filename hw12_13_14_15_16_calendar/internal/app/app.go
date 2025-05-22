package app

import (
	"context"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	l    Logger
	stor Storage
}

type Logger interface { // TODO
}

type Storage interface { // TODO
	CreateEvent(
		ctx context.Context,
		e *storage.Event,
	)
	UpdateEvent(
		ctx context.Context,
		eventID string,
		e *storage.Event,
	) error

	DeleteEvent(
		ctx context.Context,
		eventID string,
	) error
	GetEventByDay(
		ctx context.Context,
		eventID string,
		day time.Time,
	) ([]*storage.Event, error)

	GetEventByWeek(ctx context.Context,
		eventID string,
		week time.Time,
	) ([]*storage.Event, error)

	GetEventByMounth(ctx context.Context,
		eventID string,
		mounth time.Time,
	) ([]*storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		l:    logger,
		stor: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO

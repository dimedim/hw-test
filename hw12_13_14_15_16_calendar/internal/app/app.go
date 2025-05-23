package app

import (
	"context"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	l    Logger
	stor Storage
}

type Logger interface { // TODO
}

type Storage interface { // TODO
	storage.EventStorage
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

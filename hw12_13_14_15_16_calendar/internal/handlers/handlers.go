package handlers

import (
	"context"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
)

type Application interface {
	CreateEvent(ctx context.Context, event *models.Event) error
}

type Hadlers struct {
	Log logger.Logger
	App Application
}

func NewHadnlers(log logger.Logger, app Application) *Hadlers {
	return &Hadlers{
		Log: log,
		App: app,
	}
}

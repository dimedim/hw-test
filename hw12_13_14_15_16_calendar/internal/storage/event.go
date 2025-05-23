package storage

import (
	"context"
	"errors"
	"time"
)

type Event struct {
	ID          string
	UserID      string
	Title       string
	Description string `json:"description,omitempty"`

	StartsAt     time.Time     `json:"starts_at"`
	EndsAt       time.Time     `json:"ends_at,omitempty"`
	NotifyOffset time.Duration `json:"notify_offset,omitempty"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Notification struct {
	EventID  string
	Title    string
	StartsAt time.Time
	UserID   string
}

type EventStorage interface {
	CreateEvent(ctx context.Context, e *Event) error
	UpdateEvent(ctx context.Context, eventID string, e *Event) error
	DeleteEvent(ctx context.Context, eventID string) error

	ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*Event, error)
	ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*Event, error)
	ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*Event, error)
}

var ErrEventNotExists = errors.New("event not exists")

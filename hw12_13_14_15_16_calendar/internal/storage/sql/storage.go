package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.DB.Close()
}

func (s *Storage) CreateEvent(
	ctx context.Context,
	e *models.Event,
) error {
	const query = `INSERT INTO events(id, user_id, title, description, 
	starts_at, ends_at, notify_offset)
	VALUES (:id, :user_id, :title, :description, 
	:starts_at, :ends_at, :notify_offset);`

	_, err := s.DB.NamedExecContext(ctx, query, e)
	if err != nil {
		return fmt.Errorf("new event insert: %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(
	ctx context.Context,
	eventID string,
	e *models.Event,
) error {
	e.UpdatedAt = time.Now()
	e.ID = eventID

	const query = `UPDATE events SET
		title = :title,
		description   = :description,
		starts_at     = :starts_at,
		ends_at       = :ends_at,
		notify_offset = :notify_offset,
		updated_at    = :updated_at
	WHERE id = :id;
	`

	res, err := s.DB.NamedExecContext(ctx, query, e)
	if err != nil {
		return fmt.Errorf("update event exec: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update event RowsAffected: %w", err)
	}
	if count == 0 {
		return models.ErrEventNotExists
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	const query = `DELETE FROM events WHERE id = $1;`

	res, err := s.DB.ExecContext(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("DeleteEvent exec: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteEvent rows affected: %w", err)
	}
	if rows == 0 {
		return models.ErrEventNotExists
	}
	return nil
}

func (s *Storage) ListEventsByDay(
	ctx context.Context, userID string, day time.Time,
) ([]*models.Event, error) {
	loc := day.Location()
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)

	return s.ListByAnyTime(ctx, userID, start, end)
}

func (s *Storage) ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error) {
	loc := week.Location()
	start := time.Date(week.Year(), week.Month(), week.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 7)
	return s.ListByAnyTime(ctx, userID, start, end)
}

func (s *Storage) ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error) {
	loc := month.Location()
	start := time.Date(month.Year(), month.Month(), month.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)

	return s.ListByAnyTime(ctx, userID, start, end)
}

func (s *Storage) ListByAnyTime(
	ctx context.Context, userID string, start, end time.Time,
) ([]*models.Event, error) {
	var events []*models.Event
	const query = `SELECT
	id, user_id, title, description, starts_at, ends_at, notify_offset, created_at, updated_at
	FROM events
	WHERE user_id = $1 AND starts_at >= $2 AND starts_at < $3
    ORDER BY starts_at ASC;`

	err := s.DB.SelectContext(ctx, &events, query, userID, start, end)
	if err != nil {
		return nil, fmt.Errorf("list by day select: %w", err)
	}
	return events, nil
}

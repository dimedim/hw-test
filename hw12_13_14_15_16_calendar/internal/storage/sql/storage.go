package sqlstorage

import (
	"context"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

/*
Видел вот такой юз миграций в sqlx: https://github.com/jmoiron/sqlx


var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`
*/
// exec the schema or fail; multi-statement Exec behavior varies between
// database drivers;  pq will exec them all, sqlite3 won't, ymmv
//     db.MustExec(schema)

type Storage struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.DB.Ping()
}

func (s *Storage) Close(ctx context.Context) error {
	// return s.DB.Close(ctx)
	return s.DB.Close()
}

func (s *Storage) CreateEvent(
	ctx context.Context,
	e *storage.Event,
) error {

	// query := `INSERT INTO `
	return nil
}
func (s *Storage) UpdateEvent(
	ctx context.Context,
	eventID string,
	e *storage.Event,
) error {
	return nil
}
func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	return nil
}
func (s *Storage) ListEventsByDay(
	ctx context.Context, userID string, day time.Time,
) ([]*storage.Event, error) {
	return nil, nil
}
func (s *Storage) ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*storage.Event, error) {
	return nil, nil
}
func (s *Storage) ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*storage.Event, error) {
	return nil, nil
}

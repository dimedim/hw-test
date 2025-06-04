package memorystorage

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
)

type Storage struct {
	mu sync.RWMutex
	DB map[string]*models.Event
}

func New() *Storage {
	return &Storage{
		DB: make(map[string]*models.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, e *models.Event) (*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("create event: %w", ctx.Err())
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		e.CreatedAt = time.Now().UTC()
		s.DB[e.ID] = e
		return e, nil
	}
}

func (s *Storage) UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("update event: %w", ctx.Err())
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.CheckExistance(eventID) {
			return nil, models.ErrEventNotExists
		}
		e.UpdatedAt = time.Now().UTC()
		s.DB[eventID] = e
		return e, nil
	}
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("delete event: %w", ctx.Err())
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.DB, eventID)
		return nil
	}
}

func (s *Storage) ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("list event by day: %w", ctx.Err())
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()
		loc := day.Location()
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 1)
		return s.GetEventListByAnyDate(userID, start, end)
	}
}

func (s *Storage) ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("list event by week: %w", ctx.Err())
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()
		loc := week.Location()
		start := time.Date(week.Year(), week.Month(), week.Day(), 0, 0, 0, 0, loc)
		end := start.AddDate(0, 0, 7)
		return s.GetEventListByAnyDate(userID, start, end)
	}
}

func (s *Storage) ListEventsByMonth(
	ctx context.Context,
	userID string,
	month time.Time,
) ([]*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("list event by month: %w", ctx.Err())
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()
		loc := month.Location()
		start := time.Date(month.Year(), month.Month(), month.Day(), 0, 0, 0, 0, loc)
		end := start.AddDate(0, 1, 0)
		return s.GetEventListByAnyDate(userID, start, end)
	}
}

func (s *Storage) CheckExistance(eventID string) bool {
	_, exists := s.DB[eventID]
	return exists
}

func (s *Storage) GetEventListByAnyDate(userID string, start, end time.Time) ([]*models.Event, error) {
	res := make([]*models.Event, 0, len(s.DB))
	for _, ev := range s.DB {
		if userID == ev.UserID {
			if ev.StartsAt.Compare(start) >= 0 && ev.StartsAt.Compare(end) == -1 {
				res = append(res, ev)
			}
		}
	}
	slices.SortFunc(res, func(a, b *models.Event) int {
		return a.StartsAt.Compare(b.StartsAt)
	})
	return res, nil
}

func (s *Storage) Close() error {
	return nil
}

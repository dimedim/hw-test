package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
	mu sync.RWMutex //nolint:unused
	DB map[string]*storage.Event
}

func New() *Storage {
	return &Storage{
		DB: make(map[string]*storage.Event),
	}
}

func (s *Storage) CreateEvent(
	ctx context.Context,
	e *storage.Event,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DB[e.ID] = e
}

func (s *Storage) UpdateEvent(
	ctx context.Context,
	eventID string,
	e *storage.Event,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.CheckExistance(eventID) {
		return storage.ErrEventNotExists
	}
	s.DB[eventID] = e
	return nil
}

func (s *Storage) DeleteEvent(
	ctx context.Context,
	eventID string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.DB, eventID)
	return nil
}

func (s *Storage) GetEventByDay(
	ctx context.Context,
	eventID string,
	day time.Time,
) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*storage.Event, 50)
	for _, ev := range s.DB {
		if ev.Deadline.Compare(day) >= 0 && ev.Deadline.Compare(day.Add(time.Hour*24)) < -1 {
			res = append(res, ev)
		}
	}
	return res, nil
}

func (s *Storage) GetEventByWeek(
	ctx context.Context,
	eventID string,
	week time.Time,
) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*storage.Event, 50)
	for _, ev := range s.DB {
		if ev.Deadline.Compare(week) >= 0 && ev.Deadline.Compare(week.Add(time.Hour*24)) < -1 {
			res = append(res, ev)
		}
	}
	return res, nil
}

func (s *Storage) GetEventByMounth(
	ctx context.Context,
	eventID string,
	mounth time.Time,
) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*storage.Event, 50)
	for _, ev := range s.DB {
		if ev.Deadline.Compare(mounth) >= 0 && ev.Deadline.Compare(mounth.Add(time.Hour*24)) < -1 {
			res = append(res, ev)
		}
	}
	return res, nil
}

func (s *Storage) CheckExistance(eventID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.DB[eventID]
	return exists
}

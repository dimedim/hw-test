package service

import "github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/storage"

type EventService struct {
	DB storage.EventStorage
}

func NewEventService(db storage.EventStorage) *EventService {
	return &EventService{DB: db}
}

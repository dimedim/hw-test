package storage

import (
	"context"
	"errors"
	"time"
)

type Event struct {
	ID        string
	Title     string
	Desc      string `json:"desc,omitempty"`
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deadline  time.Time
	Notify    time.Time `json:"notify,omitempty"`
}

/*
   ID события;
   Заголовок события;
   Дата события;
   Пользователь, которому отправлять.
*/

type Notification struct {
	EventID    string
	EventTitle string
	EventDate  time.Time
	UserID     string
}

/*
? Описание методов эндпоинтов:?
    Создать (событие);
    Обновить (ID события, событие);
    Удалить (ID события);
    СписокСобытийНаДень (дата);
    СписокСобытийНаНеделю (дата начала недели);
    СписокСобытийНaМесяц (дата начала месяца).
*/

type EventStorage interface {
	CreateEvent(
		ctx context.Context,
		e *Event,
	)
	UpdateEvent(
		ctx context.Context,
		eventID string,
		e *Event,
	) error

	DeleteEvent(
		ctx context.Context,
		eventID string,
	) error
	GetEventByDay(
		ctx context.Context,
		eventID string,
		day time.Time,
	) ([]*Event, error)

	GetEventByWeek(ctx context.Context,
		eventID string,
		week time.Time,
	) ([]*Event, error)

	GetEventByMounth(ctx context.Context,
		eventID string,
		mounth time.Time,
	) ([]*Event, error)
}

var ErrEventNotExists = errors.New("event not exists")

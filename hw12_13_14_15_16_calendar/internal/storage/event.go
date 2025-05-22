package storage

import (
	"context"
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
	CreateEvent(ctx context.Context)
	UpdateEvent(ctx context.Context)
	DeleteEvent(ctx context.Context)
	GetEventByDay(ctx context.Context)
	GetEventByWeek(ctx context.Context)
	GetEventByMounth(ctx context.Context)
}

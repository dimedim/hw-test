package internalgrpc

import (
	"context"
	"fmt"
	"time"

	ge "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/api/grpcevents"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventStorage interface {
	CreateEvent(ctx context.Context, e *models.Event) (*models.Event, error)
	UpdateEvent(ctx context.Context, eventID string, e *models.Event) (*models.Event, error)
	DeleteEvent(ctx context.Context, eventID string) error

	ListEventsByDay(ctx context.Context, userID string, day time.Time) ([]*models.Event, error)
	ListEventsByWeek(ctx context.Context, userID string, week time.Time) ([]*models.Event, error)
	ListEventsByMonth(ctx context.Context, userID string, month time.Time) ([]*models.Event, error)

	Close() error
}

type GRPCServer struct {
	Cfg *config.Config
	ge.UnimplementedCalendarServiceServer
	Store EventStorage
}

func New(storage EventStorage, cfg *config.Config) *GRPCServer {
	return &GRPCServer{Store: storage, Cfg: cfg}
}

func (s *GRPCServer) CreateEvent(
	ctx context.Context,
	req *ge.CreateEventRequest,
) (*ge.CreateEventResponse, error) {
	event := req.GetEvent()
	if event == nil {
		return nil, models.ErrEventNotSet
	}

	e := gEventToEvent(event)

	created, err := s.Store.CreateEvent(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("store create event")
	}

	return &ge.CreateEventResponse{
		Event: eventToGEvent(created),
	}, nil
}

func gEventToEvent(event *ge.Event) *models.Event {
	start := event.GetStartsAt().AsTime()
	end := event.GetEndsAt().AsTime()
	return &models.Event{
		ID:          event.GetId(),
		Title:       event.GetTitle(),
		StartsAt:    start,
		EndsAt:      end,
		Description: event.GetDescription(),
		UserID:      event.GetUserId(),
		// TODO: need test
		// NotifyOffset: func() time.Duration {
		// 	if event.GetNotifyOffset() != nil {
		// 		return event.GetNotifyOffset().AsDuration()
		// 	}
		// 	return 0
		// }(),
		NotifyOffset: event.GetNotifyOffset().AsDuration(),
	}
}

func eventToGEvent(event *models.Event) *ge.Event {
	return &ge.Event{
		Id:          event.ID,
		Title:       event.Title,
		StartsAt:    timestamppb.New(event.StartsAt),
		EndsAt:      timestamppb.New(event.EndsAt),
		Description: event.Description,
		UserId:      event.UserID,
		NotifyOffset: func() *durationpb.Duration {
			if event.NotifyOffset > 0 {
				return durationpb.New(event.NotifyOffset)
			}
			return nil
		}(),
	}
}

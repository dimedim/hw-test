package internalgrpc

import (
	"context"
	"fmt"

	ge "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/api/grpcevents"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func (s *GRPCServer) UpdateEvent(ctx context.Context, req *ge.UpdateEventRequest) (*ge.UpdateEventResponse, error) {
	return nil, nil
}

func (s *GRPCServer) DeleteEvent(ctx context.Context, req *ge.DeleteEventRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *GRPCServer) ListDay(ctx context.Context, req *ge.ListDayRequest) (*ge.ListDayResponse, error) {
	return nil, nil
}

func (s *GRPCServer) ListWeek(ctx context.Context, req *ge.ListWeekRequest) (*ge.ListWeekResponse, error) {
	return nil, nil
}

func (s *GRPCServer) ListMonth(ctx context.Context, req *ge.ListMonthRequest) (*ge.ListMonthResponse, error) {
	return nil, nil
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

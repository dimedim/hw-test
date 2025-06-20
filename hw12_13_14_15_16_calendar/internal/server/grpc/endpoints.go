package internalgrpc

import (
	"context"

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
	requestEvent := req.GetEvent()
	if requestEvent == nil {
		return nil, models.ErrEventNotSet
	}

	ev := GRPCToEvent(requestEvent)

	res, err := s.App.CreateEvent(ctx, ev)
	if err != nil {
		return nil, err
	}

	return &ge.CreateEventResponse{
		Event: eventToGRPC(res),
	}, nil
}

func (s *GRPCServer) UpdateEvent(ctx context.Context, req *ge.UpdateEventRequest) (*ge.UpdateEventResponse, error) {

	requestEvent := req.GetEvent()
	eventID := req.GetEventId()
	if requestEvent == nil || eventID == "" {
		return nil, models.ErrEventNotSet
	}
	event := GRPCToEvent(requestEvent)
	res, err := s.App.UpdateEvent(ctx, eventID, event)
	if err != nil {
		return nil, err
	}
	return &ge.UpdateEventResponse{
		Event: eventToGRPC(res),
	}, nil
}

func (s *GRPCServer) DeleteEvent(ctx context.Context, req *ge.DeleteEventRequest) (*emptypb.Empty, error) {
	eventID := req.GetEventId()
	err := s.App.DeleteEvent(ctx, eventID)
	return nil, err
}

func (s *GRPCServer) ListDay(ctx context.Context, req *ge.ListDayRequest) (*ge.ListDayResponse, error) {

	userID := req.GetUserId()

	events, err := s.App.ListEventsByDay(ctx, userID, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}
	return &ge.ListDayResponse{Events: manyEventsToGRPC(events)}, nil
}

func (s *GRPCServer) ListWeek(ctx context.Context, req *ge.ListWeekRequest) (*ge.ListWeekResponse, error) {
	userID := req.GetUserId()

	events, err := s.App.ListEventsByWeek(ctx, userID, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}
	return &ge.ListWeekResponse{Events: manyEventsToGRPC(events)}, nil
}

func (s *GRPCServer) ListMonth(ctx context.Context, req *ge.ListMonthRequest) (*ge.ListMonthResponse, error) {
	userID := req.GetUserId()

	events, err := s.App.ListEventsByMonth(ctx, userID, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}
	return &ge.ListMonthResponse{Events: manyEventsToGRPC(events)}, nil
}

func GRPCToEvent(event *ge.Event) *models.Event {
	if event == nil {
		return nil
	}
	start := event.GetStartsAt().AsTime()
	end := event.GetEndsAt().AsTime()
	return &models.Event{
		ID:           event.GetId(),
		UserID:       event.GetUserId(),
		Title:        event.GetTitle(),
		Description:  event.GetDescription(),
		StartsAt:     start,
		EndsAt:       end,
		NotifyOffset: event.GetNotifyOffset().AsDuration(),
		CreatedAt:    event.GetCreatedAt().AsTime(),
		UpdatedAt:    event.GetUpdatedAt().AsTime(),
	}
	//  if t := pb.GetCreatedAt(); t != nil {
	//     ev.CreatedAt = t.AsTime()
	// }
	// if t := pb.GetUpdatedAt(); t != nil {
	//     ev.UpdatedAt = &t.AsTime()
	// }
}

func eventToGRPC(event *models.Event) *ge.Event {
	if event == nil {
		return nil
	}
	return &ge.Event{
		Id:           event.ID,
		UserId:       event.UserID,
		Title:        event.Title,
		Description:  event.Description,
		StartsAt:     timestamppb.New(event.StartsAt),
		EndsAt:       timestamppb.New(event.EndsAt),
		NotifyOffset: durationpb.New(event.NotifyOffset),
		CreatedAt:    timestamppb.New(event.CreatedAt),
		UpdatedAt:    timestamppb.New(event.UpdatedAt),
	}
	// на случай если буду делать указатель
	// if ev.UpdatedAt != nil {
	// 	pb.UpdatedAt = timestamppb.New(*ev.UpdatedAt)
	// }
}

func manyEventsToGRPC(events []*models.Event) []*ge.Event {
	res := make([]*ge.Event, 0, len(events))

	for _, v := range events {
		res = append(res, eventToGRPC(v))
	}
	return res
}

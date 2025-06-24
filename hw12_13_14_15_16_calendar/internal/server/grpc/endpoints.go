package internalgrpc

import (
	"context"

	ge "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/api/grpcevents"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GRPCServer) CreateEvent(ctx context.Context, req *ge.CreateEventRequest) (*ge.CreateEventResponse, error) {
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

func (s *GRPCServer) ListDay(ctx context.Context, req *ge.ListEventsRequest) (*ge.ListEventsResponse, error) {
	userID := req.GetUserId()

	events, err := s.App.ListEventsByDay(ctx, userID, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}
	return &ge.ListEventsResponse{Events: manyEventsToGRPC(events)}, nil
}

func (s *GRPCServer) ListWeek(ctx context.Context, req *ge.ListEventsRequest) (*ge.ListEventsResponse, error) {
	userID := req.GetUserId()

	events, err := s.App.ListEventsByWeek(ctx, userID, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}
	return &ge.ListEventsResponse{Events: manyEventsToGRPC(events)}, nil
}

func (s *GRPCServer) ListMonth(ctx context.Context, req *ge.ListEventsRequest) (*ge.ListEventsResponse, error) {
	userID := req.GetUserId()

	events, err := s.App.ListEventsByMonth(ctx, userID, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}
	return &ge.ListEventsResponse{Events: manyEventsToGRPC(events)}, nil
}

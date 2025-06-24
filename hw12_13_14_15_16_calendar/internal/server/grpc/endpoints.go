package internalgrpc

import (
	"context"
	"time"

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
	time.Unix()
	return &ge.ListEventsResponse{Events: manyEventsToGRPC(events)}, nil
}

// TODO: мб лучше через time.Unix()
/*
func (serv *server) GetListByPeriod(_ context.Context, req *eventproto.EventGetListByPeriodRequest) (*eventproto.EventGetListByPeriodResponse, error) {
	events, err := serv.Calendar.GetEventsByPeriod(
		time.Unix(req.DateFrom.GetSeconds(), int64(req.DateFrom.GetNanos())),
		time.Unix(req.DateTo.GetSeconds(), int64(req.DateTo.GetNanos())),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	respEvents := make([]*eventproto.Event, 0, len(events))
	for _, event := range events {
		respEvents = append(respEvents, convertToProtoEvent(event))
	}

	return &eventproto.EventGetListByPeriodResponse{Events: respEvents}, nil
}

func convertToProtoEvent(event *domain.Event) *eventproto.Event {
	return &eventproto.Event{
		Id:       event.ID.String(),
		Title:    event.Title,
		DateFrom: &timestamp.Timestamp{Seconds: event.DateFrom.Unix(), Nanos: int32(event.DateFrom.UnixNano())},
		DateTo:   &timestamp.Timestamp{Seconds: event.DateTo.Unix(), Nanos: int32(event.DateTo.UnixNano())},
	}
}

*/

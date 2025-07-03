package internalgrpc

import (
	ge "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/api/grpcevents"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
}

func manyEventsToGRPC(events []*models.Event) []*ge.Event {
	res := make([]*ge.Event, 0, len(events))

	for _, v := range events {
		res = append(res, eventToGRPC(v))
	}
	return res
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
	"github.com/gorilla/mux"
)

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		models.JSONError(h.Log, w, http.StatusBadRequest, "invalid JSON body")
		h.Log.Error("create event bad request", logger.Err(err))
		return
	}

	res, err := h.App.CreateEvent(r.Context(), &event)
	if err != nil {
		models.JSONError(h.Log, w, http.StatusInternalServerError, "internal")
		h.Log.Error("create event internal", logger.Err(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		models.JSONError(h.Log, w, http.StatusInternalServerError, "internal")
		h.Log.Error("encode JSON", logger.Err(err))
		return
	}
}

func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var event models.Event

	vars := mux.Vars(r)
	eventID := vars[EventIDParam]
	if eventID == "" {
		models.JSONError(h.Log, w, http.StatusBadRequest, "event_id is required")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		models.JSONError(h.Log, w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	res, err := h.App.UpdateEvent(r.Context(), eventID, &event)
	if errors.Is(err, models.ErrEventNotExists) {
		models.JSONError(h.Log, w, http.StatusBadRequest, models.ErrEventNotExists.Error())
		return
	}
	if err != nil {
		models.JSONError(h.Log, w, http.StatusInternalServerError, "internal")
		h.Log.Error("update event internal", logger.Err(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		models.JSONError(h.Log, w, http.StatusInternalServerError, "internal")
		h.Log.Error("encode JSON", logger.Err(err))
		return
	}
}

func (h *Handlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars[EventIDParam]
	if eventID == "" {
		models.JSONError(h.Log, w, http.StatusBadRequest, "event_id is required")
		return
	}
	err := h.App.DeleteEvent(r.Context(), eventID)
	if errors.Is(err, models.ErrEventNotExists) {
		models.JSONError(h.Log, w, http.StatusBadRequest, models.ErrEventNotExists.Error())
		return
	}
	if err != nil {
		models.JSONError(h.Log, w, http.StatusInternalServerError, "internal")
		h.Log.Error("delete event internal", logger.Err(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) ListEventsByDay(w http.ResponseWriter, r *http.Request) {
	h.listEvents(w, r, Day)
}

func (h *Handlers) ListEventsByWeek(w http.ResponseWriter, r *http.Request) {
	h.listEvents(w, r, Week)
}

func (h *Handlers) ListEventsByMonth(w http.ResponseWriter, r *http.Request) {
	h.listEvents(w, r, Month)
}

func (h *Handlers) listEvents(w http.ResponseWriter, r *http.Request, period string) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	UserID := vars[UserIDParam]
	if UserID == "" {
		models.JSONError(h.Log, w, http.StatusBadRequest, "user_id is required")
		return
	}

	dateStr := r.URL.Query().Get(DateParam)

	dateTime, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		models.JSONError(h.Log, w, http.StatusBadRequest, "invalid date format")
		return
	}

	var events []*models.Event
	switch period {
	case Day:
		events, err = h.App.ListEventsByDay(r.Context(), UserID, dateTime)
	case Month:
		events, err = h.App.ListEventsByMonth(r.Context(), UserID, dateTime)
	case Week:
		events, err = h.App.ListEventsByWeek(r.Context(), UserID, dateTime)
	}

	if err != nil {
		models.JSONError(h.Log, w, http.StatusInternalServerError, "internal")
		h.Log.Error("list events", logger.Err(err))
		return
	}

	if err := json.NewEncoder(w).Encode(events); err != nil {
		h.Log.Error("encode events", logger.Err(err))
	}
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/models"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
)

const (
	Day   = "day"
	Week  = "week"
	Month = "month"
)

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		models.JSONError(h.Log, w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// TODO: add validate not nil title and endsAt !before startsAt

	res, err := h.App.CreateEvent(r.Context(), &event)
	if err != nil {
		http.Error(w, "create event error", http.StatusInternalServerError)
		h.Log.Error("create post internal", logger.Err(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	// TODO: мб вынести в отдельную функцию
	if err := json.NewEncoder(w).Encode(res); err != nil {
		models.JSONError(h.Log, w, http.StatusInternalServerError, "internal")
		h.Log.Error("encode JSON", logger.Err(err))
		return
	}
}

func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {

}
func (h *Handlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {

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

func (h *Handlers) listEvents(w http.ResponseWriter, r *http.Request, date string) {
	w.Header().Set("Content-Type", "application/json")

}

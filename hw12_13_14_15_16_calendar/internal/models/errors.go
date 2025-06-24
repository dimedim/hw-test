package models

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
)

var (
	ErrEventNotExists = errors.New("event not exists")
	ErrEventNotSet    = errors.New("event not set")
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func JSONError(log logger.Logger, w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Message: msg}); err != nil {
		log.Error("JSONError encode error", logger.Err(err))
	}
}

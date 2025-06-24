package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	mware "github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/server/middleware"
	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/pkg/logger"
	"github.com/gorilla/mux"
)

type Handlers interface {
	Hello(w http.ResponseWriter, r *http.Request)
	CreateEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	ListEventsByDay(w http.ResponseWriter, r *http.Request)
	ListEventsByWeek(w http.ResponseWriter, r *http.Request)
	ListEventsByMonth(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	Cfg        *config.Config
	Log        logger.Logger
	Router     *mux.Router
	HTTPServer *http.Server
	Handlers   Handlers
}

func NewServer(cfg *config.Config, log logger.Logger, router *mux.Router, handlers Handlers) *Server {
	addr := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	return &Server{
		Cfg: cfg,
		Log: log,
		HTTPServer: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  cfg.HTTP.Timeout,
			WriteTimeout: cfg.HTTP.Timeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		},
		Router:   router,
		Handlers: handlers,
	}
}

func (s *Server) Start() error {
	s.Log.Info("starting HTTP server", "addr", "http://"+s.HTTPServer.Addr)
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.Log.Info("stopping HTTP server", "addr", s.HTTPServer.Addr)
	return s.HTTPServer.Shutdown(ctx)
}

func (s *Server) GetRouter() http.Handler {
	return s.Router
}

func (s *Server) RegisterRoutes() {
	s.Router.Use(mware.PanicRecover(s.Log))

	s.Router.HandleFunc("/", s.Handlers.Hello).Methods(http.MethodGet)

	events := s.Router.PathPrefix("/events").Subrouter()
	events.HandleFunc("", s.Handlers.CreateEvent).Methods(http.MethodPost)
	events.HandleFunc("/{event_id}", s.Handlers.UpdateEvent).Methods(http.MethodPatch, http.MethodPut)
	events.HandleFunc("/{event_id}", s.Handlers.DeleteEvent).Methods(http.MethodDelete)

	// ? GET /events/day/user_id?date=2025-06-10
	events.HandleFunc("/day/{user_id}", s.Handlers.ListEventsByDay).Methods(http.MethodGet)
	events.HandleFunc("/week/{user_id}", s.Handlers.ListEventsByWeek).Methods(http.MethodGet)
	events.HandleFunc("/month/{user_id}", s.Handlers.ListEventsByMonth).Methods(http.MethodGet)
}

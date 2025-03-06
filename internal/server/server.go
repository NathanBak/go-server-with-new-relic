package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"

	newrelic "github.com/newrelic/go-agent/v3/newrelic"
)

// A Server is a web server that extends http.Server.  It should be created with the New() function,
// started with the Server.ListenAndServe() function, and can be cleanly closed with the
// Server.Shutdown() function.
type Server struct {
	http.Server
	log     Logger
	cfg     Config
	storage Storage
}

// New creates, configures, and returns a new server instance.
func New(cfg Config, app *newrelic.Application) (*Server, error) {

	s := &Server{
		Server: http.Server{
			Addr:         fmt.Sprintf(":%v", cfg.Port),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg:     cfg,
		log:     cfg.Logger,
		storage: cfg.Storage,
	}

	if s.storage == nil {
		return nil, errors.New("no storage provided in config")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Use(nrgorilla.Middleware(app))
	for _, route := range s.routes() {
		wrappedHandler := s.requestWrapper(route.HandlerFunc, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(wrappedHandler)
	}
	s.Server.Handler = router

	return s, nil
}

func (s *Server) ListenAndServe() error {
	msg := fmt.Sprintf("Starting server on port %d", s.cfg.Port)
	s.log.Info(msg)
	return s.Server.ListenAndServe()
}

// ServeHttp routes and handles a request.  It is typically used for testing purposes.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

// Shutdown cleanly shutsdown the server and blocks until complete (or timing out).
func (s *Server) Shutdown(ctx context.Context) error {
	msg := "Shutting down server"
	s.log.Info(msg)
	return s.Server.Shutdown(ctx)
}

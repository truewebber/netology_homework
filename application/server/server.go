package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	server struct {
		httpServer *http.Server
	}
)

func New(addr string, router *mux.Router) *server {
	return &server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}

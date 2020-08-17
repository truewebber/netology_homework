package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/truewebber/netology_homework/log"
)

func NewRouter(logger log.Logger) *mux.Router {
	h := New(logger)

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(h.NotFound)

	rApiSlow := r.PathPrefix(ApiSlowPath).Subrouter()
	rApiSlow.Use(h.ValidateSlowMiddleware)
	rApiSlow.HandleFunc("", h.ApiSlow).Methods(http.MethodPost)

	return r
}

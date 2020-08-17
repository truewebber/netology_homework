package handler

import (
	"encoding/json"
	"net/http"

	"github.com/truewebber/netology_homework/log"
)

type (
	Handler struct {
		logger log.Logger
	}
)

func New(logger log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) write(w http.ResponseWriter, obj interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(obj)
	if err != nil {
		h.logger.Error("Error write response", "error", err.Error())
	}
}

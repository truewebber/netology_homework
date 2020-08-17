package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/truewebber/netology_homework/slow"
)

const (
	MaxTimeoutValue = 5000 * time.Millisecond

	ApiSlowPath = "/api/slow"
)

var (
	timeoutTooLongError = newErrorResponse("timeout too long")
)

func (h *Handler) ApiSlow(w http.ResponseWriter, r *http.Request) {
	s := new(slow.Slow)

	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		h.logger.Warn("error decode slow request", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	duration := time.Duration(s.Timeout) * time.Millisecond
	time.Sleep(duration)

	h.logger.Info("slow", "time", duration.String())

	w.WriteHeader(http.StatusOK)
	h.write(w, newResponse("ok"))
}

func (h *Handler) ValidateSlowMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		r = r.WithContext(ctx)
		syncW := NewSyncWriter(w)

		go func() {
			next.ServeHTTP(syncW, r)

			cancel()
		}()

		select {
		case <-time.After(MaxTimeoutValue):
			syncW.WriteHeader(http.StatusBadRequest)
			h.write(syncW, timeoutTooLongError)
		case <-ctx.Done():
		}
	})
}

package handler

import (
	"net/http"
	"sync"
)

type (
	syncWriter struct {
		w            http.ResponseWriter
		m            sync.Mutex
		bodyClosed   bool
		statusClosed bool
	}
)

func NewSyncWriter(w http.ResponseWriter) *syncWriter {
	return &syncWriter{w: w}
}

func (sw *syncWriter) Header() http.Header {
	return sw.w.Header()
}

func (sw *syncWriter) Write(bb []byte) (int, error) {
	sw.m.Lock()
	defer sw.m.Unlock()

	if sw.bodyClosed {
		return 0, nil
	}

	sw.bodyClosed = true

	return sw.w.Write(bb)
}

func (sw *syncWriter) WriteHeader(statusCode int) {
	sw.m.Lock()
	defer sw.m.Unlock()

	if sw.statusClosed {
		return
	}

	sw.statusClosed = true

	sw.w.WriteHeader(statusCode)
}

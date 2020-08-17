package handler

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/truewebber/netology_homework/log"
)

func TestHandler_ApiSlow(t *testing.T) {
	type (
		fields struct {
			logger log.Logger
		}
		args struct {
			w *mockWriter
			r *http.Request
		}
	)

	nopLogger := log.NewNop()

	writerOne := newMockWriter(http.StatusOK, []byte("{\"status\":\"ok\"}\n"))
	reqOne := newPostRequest("{\"timeout\": 1000}")

	writerTwo := newMockWriter(http.StatusBadRequest, nil)
	reqTwo := newPostRequest("{\"timeout\": \"1000\"}")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mockWriter
	}{
		{
			name:   "Timeout 1000, return OK",
			fields: fields{logger: nopLogger},
			args:   args{w: new(mockWriter), r: reqOne},
			want:   writerOne,
		},
		{
			name:   "Timeout 1000, return bad request",
			fields: fields{logger: nopLogger},
			args:   args{w: new(mockWriter), r: reqTwo},
			want:   writerTwo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				logger: tt.fields.logger,
			}
			h.ApiSlow(tt.args.w, tt.args.r)

			if !reflect.DeepEqual(tt.args.w, tt.want) {
				t.Errorf("ApiSlow(), writer is = %v, want %v", tt.args.w, tt.want)
			}
		})
	}
}

func newPostRequest(body string) *http.Request {
	return newRequest(http.MethodPost, body)
}

func newRequest(method string, body string) *http.Request {
	reader := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, "http://tosi.bosi/api/slow", reader)
	if err != nil {
		panic(err)
	}

	return req
}

// MOCK

type (
	mockWriter struct {
		bytes []byte
		code  int
	}
)

func newMockWriter(code int, bytes []byte) *mockWriter {
	return &mockWriter{
		bytes: bytes,
		code:  code,
	}
}

func (m *mockWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockWriter) Write(bb []byte) (int, error) {
	m.bytes = bb

	return 0, nil
}

func (m *mockWriter) WriteHeader(statusCode int) {
	m.code = statusCode
}

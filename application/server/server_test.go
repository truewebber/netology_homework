package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"testing"

	"golang.org/x/net/context"

	"github.com/truewebber/netology_homework/handler"
	"github.com/truewebber/netology_homework/log"
)

const (
	addr = "localhost:9999"
)

func Test_Server(t *testing.T) {
	type fields struct {
		httpServer *http.Server
	}

	tests := []struct {
		name         string
		fields       fields
		request      *http.Request
		wantRespBody []byte
		wantRespCode int
		wantErr      bool
	}{
		{
			name:         "Success test, timeout 1000",
			fields:       fields{httpServer: newServer(addr)},
			request:      newPostRequest("/api/slow", "{\"timeout\":1000}"),
			wantRespBody: []byte("{\"status\":\"ok\"}\n"),
			wantRespCode: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "Success test, timeout 0",
			fields:       fields{httpServer: newServer(addr)},
			request:      newPostRequest("/api/slow", "{\"timeout\":0}"),
			wantRespBody: []byte("{\"status\":\"ok\"}\n"),
			wantRespCode: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "Success test, timeout 4990",
			fields:       fields{httpServer: newServer(addr)},
			request:      newPostRequest("/api/slow", "{\"timeout\":4990}"),
			wantRespBody: []byte("{\"status\":\"ok\"}\n"),
			wantRespCode: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "Fail test, timeout 5005",
			fields:       fields{httpServer: newServer(addr)},
			request:      newPostRequest("/api/slow", "{\"timeout\":5005}"),
			wantRespBody: []byte("{\"error\":\"timeout too long\"}\n"),
			wantRespCode: http.StatusBadRequest,
			wantErr:      false,
		},
		{
			name:         "Fail test, method GET",
			fields:       fields{httpServer: newServer(addr)},
			request:      newRequest(http.MethodGet, "/api/slow", ""),
			wantRespBody: []byte(""),
			wantRespCode: http.StatusMethodNotAllowed,
			wantErr:      false,
		},
		{
			name:         "Fail test, wrong path",
			fields:       fields{httpServer: newServer(addr)},
			request:      newPostRequest("/api/slowww", "{\"timeout\":5005}"),
			wantRespBody: []byte(""),
			wantRespCode: http.StatusNotFound,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{httpServer: tt.fields.httpServer}

			ctx, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ctx.Done()

				if err := s.Shutdown(); err != nil {
					panic(err)
				}
			}()

			go func() {
				if err := s.Start(); err != http.ErrServerClosed {
					panic(err)
				}
			}()

			respBody, respCode, respErr := proceedReq(tt.request)

			if (respErr != nil) != tt.wantErr {
				t.Errorf("Request server, error = %v, wantErr %v", respErr, tt.wantErr)
			}

			if tt.wantRespCode != respCode || !reflect.DeepEqual(tt.wantRespBody, respBody) {
				t.Errorf("Request server, response is = [%d]%v, want [%d]%v", respCode, respBody,
					tt.wantRespCode, tt.wantRespBody)
			}

			cancel()
			wg.Wait()
		})
	}
}

func proceedReq(req *http.Request) ([]byte, int, error) {
	resp, respErr := http.DefaultClient.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if respErr != nil {
		return nil, 0, respErr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

func newPostRequest(path, body string) *http.Request {
	return newRequest(http.MethodPost, path, body)
}

func newRequest(method, path, body string) *http.Request {
	url := fmt.Sprintf("http://%s%s", addr, path)

	reader := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		panic(err)
	}

	return req
}

func newServer(addr string) *http.Server {
	logger := log.NewNop()
	router := handler.NewRouter(logger)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

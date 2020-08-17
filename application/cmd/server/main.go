package main

import (
	"github.com/truewebber/netology_homework/handler"
	"github.com/truewebber/netology_homework/log"
	"github.com/truewebber/netology_homework/server"
)

func main() {
	logger := log.NewZap()
	defer func() {
		if err := logger.Close(); err != nil {
			println("error close logger", err.Error())
		}
	}()

	r := handler.NewRouter(logger)
	s := server.New("localhost:9999", r)

	if err := s.Start(); err != nil {
		logger.Error("error start server", "error", err.Error())
	}
}

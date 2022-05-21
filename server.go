package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type server struct {
	redis  RedisWrapper
	router *httprouter.Router
}

func newServer(conf *config) *server {
	return &server{
		redis:  newRedisWrap(conf),
		router: httprouter.New(),
	}
}

func (s *server) run() error {
	defer s.redis.Close()

	s.router.GET(routeKey, s.handleGet)
	s.router.POST(routeKeyValue, s.handleSet)
	s.router.DELETE(routeKey, s.handleDel)

	srv := http.Server{
		Addr:    ":8080",
		Handler: requestLogger(authWrap(s.router)),
	}

	log.Printf("server listening on %s", srv.Addr)

	return srv.ListenAndServe()
}

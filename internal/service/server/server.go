package server

import (
	"github.com/dwnGnL/ddos-pow/config"
)

type server struct {
	conf *config.Config
}

func (s server) Ping() string {
	return "Pong..."
}

func New(conf *config.Config) *server {
	return &server{
		conf: conf,
	}
}

package service

import (
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/service/server"
)

type ServiceImpl struct {
	conf   *config.Config
	server Server
}

type Server interface {
	Ping() string
}

type Option func(*ServiceImpl)

func New(conf *config.Config, opts ...Option) *ServiceImpl {

	s := ServiceImpl{
		conf:   conf,
		server: server.New(conf),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func (s ServiceImpl) GetServer() Server {
	return s.server
}

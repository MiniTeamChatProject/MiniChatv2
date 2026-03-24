//go:build wireinject
// +build wireinject

package main

import (
	"room/internal/conf"
	"room/internal/server"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		server.NewGRPCServer,
		server.NewHTTPServer,
		newApp,
	)
	return &kratos.App{}, nil, nil
}

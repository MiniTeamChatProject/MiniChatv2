//go:build wireinject
// +build wireinject

package main

import (
	"room/internal/biz"
	"room/internal/conf"
	"room/internal/data"
	"room/internal/server"
	"room/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	)
	return &kratos.App{}, nil, nil
}

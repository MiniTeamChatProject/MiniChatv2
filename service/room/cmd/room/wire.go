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

// wireServerProvider 从 Bootstrap 中提取 Server 配置
func wireServerProvider(bc *conf.Bootstrap) *conf.Server {
	return bc.Server
}

// wireDataProvider 从 Bootstrap 中提取 Data 配置
func wireDataProvider(bc *conf.Bootstrap) *conf.Data {
	return bc.Data
}

func wireApp(*conf.Bootstrap, log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		wireServerProvider,
		wireDataProvider,
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		newApp,
	)
	return &kratos.App{}, nil, nil
}

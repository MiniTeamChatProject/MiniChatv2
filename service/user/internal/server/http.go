//path: ./user/service/http.go 

package server

import (
	v1 "user/api/helloworld/v1"
	"user/internal/conf"
	"user/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
// 注意参数：注入 *service.RegistrationService
func NewHTTPServer(c *conf.Server, greeter *service.RegistrationService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	
	// 注册服务：把我们的实现注册到 HTTP 服务器上
	// 注意函数名变了：RegisterRegistrationHTTPServer
	v1.RegisterRegistrationHTTPServer(srv, greeter)
	
	return srv
}

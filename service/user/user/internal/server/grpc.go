//path: ./registration/internal/service/grpc.go 

package server

import (
	// 注意这里引入的别名是 v1
	v1 "registration/api/helloworld/v1"
	"registration/internal/conf"
	"registration/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
// 注意参数：这里我们要注入 *service.RegistrationService
func NewGRPCServer(c *conf.Server, greeter *service.RegistrationService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	
	// 注册服务：把我们的实现注册到 gRPC 服务器上
	// 注意函数名变了：RegisterRegistrationServer (这是 Proto 生成的)
	v1.RegisterRegistrationServer(srv, greeter)
	
	return srv
}

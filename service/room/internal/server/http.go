package server

import (
	v1 "room/api/room/v1"
	"room/internal/conf"
	"room/internal/service"

	"github.com/go-kratos-ecosystem/components/v2/middleware/cors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, roomSvc *service.RoomService, wsSvc *service.WebSocketService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			cors.Cors(
				cors.AllowedOrigins("*"),
				cors.AllowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS"),
				cors.AllowedHeaders("Content-Type", "Authorization", "X-Requested-With"),
			),
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
	v1.RegisterRoomServiceHTTPServer(srv, roomSvc)

	// Register Swagger UI - must be registered before other routes
	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)

	// Register WebSocket route
	srv.Route("/ws").GET("/room/{room_id}", wsSvc.HandleWebSocket)

	// Handle OPTIONS requests for CORS preflight
	srv.Route("/").OPTIONS("/*", func(ctx http.Context) error {
		return ctx.JSON(200, nil)
	})

	return srv
}

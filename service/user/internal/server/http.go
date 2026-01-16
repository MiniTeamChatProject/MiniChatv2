package server

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"strings"

	v1 "registration/api/helloworld/v1"
	"registration/internal/conf"
	"registration/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("MySecretKey_0721")

// ======================================================
// Swagger OpenAPI embed
// ======================================================

//go:embed openapi.yaml
var swaggerFS embed.FS

// ======================================================
// CORS Middleware
// ======================================================
func CORS() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*khttp.Transport); ok {
					h := ht.ReplyHeader()
					h.Set("Access-Control-Allow-Origin", "*")
					h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

					if ht.Request().Method == http.MethodOptions {
						return nil, nil
					}
				}
			}
			return handler(ctx, req)
		}
	}
}

// ======================================================
// Auth Middleware (JWT)
// ======================================================
func AuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			// 放行 Swagger 相关路径
			if ht, ok := tr.(*khttp.Transport); ok {
				path := ht.Request().URL.Path
				if strings.HasPrefix(path, "/q/") {
					return handler(ctx, req)
				}
			}

			// 白名单操作（注册 / 登录）
			if isWhiteList(tr.Operation()) {
				return handler(ctx, req)
			}

			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return nil, errors.New("missing Authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return nil, errors.New("invalid Authorization format")
			}

			tokenStr := parts[1]
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return nil, errors.New("invalid or expired token")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if uid, ok := claims["user_id"].(float64); ok {
					ctx = context.WithValue(ctx, "user_id", int64(uid))
				}
			}

			return handler(ctx, req)
		}
	}
}

func isWhiteList(operation string) bool {
	return strings.Contains(operation, "Register") ||
		strings.Contains(operation, "Login")
}

// ======================================================
// Swagger UI HTML
// ======================================================
func swaggerUI(w http.ResponseWriter, _ *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>User Service API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
window.onload = () => {
  SwaggerUIBundle({
    url: "/q/openapi.yaml",
    dom_id: "#swagger-ui"
  });
};
</script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

// ======================================================
// HTTP Server
// ======================================================
func NewHTTPServer(
	c *conf.Server,
	greeter *service.RegistrationService,
	logger log.Logger,
) *khttp.Server {

	opts := []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			CORS(),
			AuthMiddleware(),
		),
	}

	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := khttp.NewServer(opts...)

	// 注册业务接口
	v1.RegisterRegistrationHTTPServer(srv, greeter)

	// Swagger UI
	srv.HandleFunc("/q/swagger", swaggerUI)

	// Swagger OpenAPI YAML（embed）
	srv.HandleFunc("/q/openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		data, err := swaggerFS.ReadFile("openapi.yaml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		_, _ = w.Write(data)
	})

	return srv
}


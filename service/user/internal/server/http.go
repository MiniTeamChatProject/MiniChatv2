package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	// ⚠️ 注意：如果你的 go.mod 是 user，这里改成 user/...
	// 如果是 registration，保持 registration/...
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

// ==========================================
// 1. CORS 中间件
// ==========================================
func CORS() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*khttp.Transport); ok {
					ht.ReplyHeader().Set("Access-Control-Allow-Origin", "*")
					ht.ReplyHeader().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					ht.ReplyHeader().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
					if ht.Request().Method == "OPTIONS" {
						return nil, nil
					}
				}
			}
			return handler(ctx, req)
		}
	}
}

// ==========================================
// 2. 鉴权中间件
// ==========================================
func AuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}
			operation := tr.Operation()
			if isWhiteList(operation) {
				return handler(ctx, req)
			}

			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return nil, errors.New("未提供 Token")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return nil, errors.New("Token 格式错误")
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return nil, errors.New("Token 无效或已过期")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userIdFloat, ok := claims["user_id"].(float64); ok {
					ctx = context.WithValue(ctx, "user_id", int64(userIdFloat))
				}
			}
			return handler(ctx, req)
		}
	}
}

func isWhiteList(operation string) bool {
	// 白名单：注册、登录、以及 Swagger 相关资源
	if strings.Contains(operation, "Register") || 
	   strings.Contains(operation, "Login") ||
	   strings.Contains(operation, "/q/") { // 放行 Swagger
		return true
	}
	return false
}

// ==========================================
// 3. Swagger UI 页面
// ==========================================
func swaggerUI(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>User Service API</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/q/openapi.yaml', 
        dom_id: '#swagger-ui',
      });
    };
  </script>
</body>
</html>`
	w.Write([]byte(html))
}

// ==========================================
// 4. HTTP Server 初始化
// ==========================================
func NewHTTPServer(c *conf.Server, greeter *service.RegistrationService, logger log.Logger) *khttp.Server {
	var opts = []khttp.ServerOption{
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
	v1.RegisterRegistrationHTTPServer(srv, greeter)

	// --- 注册 Swagger ---
	
	// 1. 提供 UI 页面
	srv.HandleFunc("/q/swagger", swaggerUI)

	// 2. 提供 YAML 文件下载 (传统文件服务方式)
	srv.HandleFunc("/q/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		// 这里的路径是相对于运行目录 (项目根目录) 的
		filePath := "api/helloworld/v1/openapi.yaml"
		
		// 检查文件是否存在 (为了调试)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			pwd, _ := os.Getwd()
			// 如果在浏览器看到这个错误，说明路径不对
			http.Error(w, "File not found at: "+pwd+"/"+filePath, 404)
			return
		}
		http.ServeFile(w, r, filePath)
	})

	return srv
}

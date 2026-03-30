package server

import (
	"aicode/internal/app/core"
	"aicode/internal/middleware"
	"aicode/internal/router"
	"aicode/pkg/http"
	"context"
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// HTTPServer 应用层 HTTP 服务器封装
// 依赖 core.App (基础设施)，在 New 时完成中间件和路由注册
// 由 main.go 组装并管理生命周期
type HTTPServer struct {
	Engine            *echo.Echo
	Core              *core.App
	MiddlewareManager *middleware.Manager
}

// NewHTTPServer 创建 HTTP 服务器
// 接收 core.App，配置中间件并注册所有路由
func NewHTTPServer(c *core.App, cfg *http.Config) *HTTPServer {
	e := http.New(cfg)
	e.HideBanner = true
	e.HidePort = true
	e.Debug = c.Config.HTTP.EnableDebug

	// 超时配置
	e.Server.ReadTimeout = c.Config.HTTP.ReadTimeout
	e.Server.WriteTimeout = c.Config.HTTP.WriteTimeout
	e.Server.IdleTimeout = c.Config.HTTP.IdleTimeout
	e.Server.MaxHeaderBytes = 1 << 20

	// 注册中间件 (来自 pkg/http 的自定义中间件，通过 New 创建的 echo 实例已经注册了)
	// 这里重新注册，因为我们是自己创建 echo 实例
	// TODO: 后续可以提取中间件注册为独立函数复用
	e.Use(echoMiddleware.RequestID())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))

	// 初始化中间件管理器，绑定 core
	mwManager := middleware.NewManager(c)

	// 注册应用路由
	r := router.New(e, c, mwManager)
	r.Register()

	return &HTTPServer{
		Engine:            e,
		Core:              c,
		MiddlewareManager: mwManager,
	}
}

// Start 启动 HTTP 服务器（阻塞）
func (s *HTTPServer) Start() error {
	// 组合 addr
	addr := fmt.Sprintf("%s:%d", s.Core.Config.HTTP.Host, s.Core.Config.HTTP.Port)
	slog.Info("http server started success. addr is " + addr)

	// 打印路由
	http.PrintRoutes(s.Engine)

	return s.Engine.Start(addr)
}

func (s *HTTPServer) Shutdown() {
	// 优雅关闭 Echo 服务器
	ctx, cancel := context.WithTimeout(context.Background(), s.Core.Config.HTTP.ShutdownTimeout)
	defer cancel()
	if err := s.Engine.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed.", "err", err.Error())
	}
	slog.Info("server shutdown success")
}

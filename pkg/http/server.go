package http

import (
	"aicode/info"
	"aicode/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
)

// Server HTTP服务器封装（精简版）
type Server struct {
	echo   *echo.Echo
	config *Config
}

var log = logger.C("httpserver")

// New 创建Echo服务器
func New(cfg *Config) *Server {

	// 创建Echo实例
	e := echo.New()
	// 隐藏Banner
	e.HideBanner = true
	e.HidePort = true
	e.Debug = cfg.EnableDebug
	// 直接使用 echo 内置 Server（无需手动创建 http.Server）
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Millisecond
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Millisecond
	e.Server.IdleTimeout = time.Duration(cfg.IdleTimeout) * time.Millisecond
	// 1mb 转为 int
	e.Server.MaxHeaderBytes = 1 << 20

	// 创建服务器实例
	server := &Server{
		echo:   e, // Echo实例
		config: cfg,
	}

	// 注册中间件
	server.registerMiddleware()

	// 注册健康检查
	e.GET("/health", server.healthCheck)
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	log.Info("http server started success. addr is " + cfg.Host + ":" + fmt.Sprintf("%d", cfg.Port))

	return server
}

// Engine 为了方便外部访问 echo 实例
func (s *Server) Engine() *echo.Echo {
	return s.echo
}

// RegisterRoutes 注册路由
func (s *Server) RegisterRoutes(registerFunc func(e *echo.Echo)) {
	registerFunc(s.echo)
}

// registerMiddleware 注册中间件
func (s *Server) registerMiddleware() {
	// 原生中间件
	// 添加RequestID
	s.echo.Use(echoMiddleware.RequestID())
	// 添加 prometheus 中间件
	s.echo.Use(echoprometheus.NewMiddleware(info.AppInfo.Name)) // adds middleware to gather metrics
	s.echo.GET("/metrics", echoprometheus.NewHandler())         // adds route to serve gathered metrics
	// 添加 CORS 中间件
	s.echo.Use(echoMiddleware.CORS())
	// 添加限流中间件
	s.echo.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// 自定义的中间件
	// recover 中间件
	s.echo.Use(Recover())
	// 请求日志中间件
	s.echo.Use(Logger())
	// 错误处理中间件
	s.echo.HTTPErrorHandler = ErrorHandler()

}

// healthCheck 健康检查
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"service": "http-server",
	})
}

// Start 后台启动，已经包含了优雅关闭
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		errCh <- s.echo.Start(addr)
	}()

	const banner = `
========================================
 Service Started
 URL: http://127.0.0.1:%d
========================================
`
	log.Info("http server started success. addr is " + addr)
	log.Info(fmt.Sprintf(banner, s.config.Port))
	// 打印路由表
	PrintRoutes(s.echo)

	// 阻塞等待
	sig := <-quit
	log.Info("received signal", "signal", sig.String())
	s.Shutdown()
	return <-errCh
}

func (s *Server) Shutdown() {
	log.Info("http server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()
	err := s.echo.Shutdown(ctx)
	if err != nil {
		log.Error("http server echo shutdown error", "error", err)
		return
	}
	log.Info("http server stopped")
}

// PrintRoutes 格式化并打印所有 Echo 路由
func PrintRoutes(e *echo.Echo) {
	// 获取所有路由
	routes := e.Routes()

	// 按照 Path 进行排序
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path != routes[j].Path {
			return routes[i].Path < routes[j].Path
		}
		return routes[i].Method < routes[j].Method
	})

	// 初始化 tabwriter
	// 参数说明：输出目标, 最小单元格宽度, 制表符宽度, 填充空格数, 填充字符, 标志
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "\n [ROUTE TABLE]")
	fmt.Fprintln(w, " METHOD\t PATH\t HANDLER")
	fmt.Fprintln(w, " ------\t ----\t -------")

	for _, r := range routes {
		// 打印路由信息
		line := fmt.Sprintf(" %s\t %s\t %s", r.Method, r.Path, r.Name)
		fmt.Fprintln(w, line)
	}

	w.Flush() // 必须调用 Flush 才能写入 stdout
	fmt.Println()
}

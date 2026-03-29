package http

import (
	"aicode/info"
	"fmt"
	"net/http"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// New 创建Echo服务器
func New(cfg *Config) *echo.Echo {
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
	// 1 << 20 == 1MB
	e.Server.MaxHeaderBytes = 1 << 20

	// 注册中间件
	registerMiddleware(e, cfg)

	// 注册健康检查
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return e
}

// registerMiddleware 注册中间件
func registerMiddleware(e *echo.Echo, cfg *Config) {
	// 原生中间件
	// 添加RequestID
	e.Use(echoMiddleware.RequestID())
	// 添加 prometheus 中间件
	e.Use(echoprometheus.NewMiddleware(info.AppInfo.Name)) // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler())         // adds route to serve gathered metrics
	// 添加 CORS 中间件
	// e.Use(echoMiddleware.CORS())
	// 添加限流中间件
	// e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// 自定义的中间件
	// recover 中间件
	e.Use(Recover())
	// 请求日志中间件
	e.Use(Logger())
	// 添加 CORS 中间件
	e.Use(CORS(cfg.CORS))
	// 添加限流中间件
	e.Use(RateLimit(cfg.RateLimit))

	// 错误处理中间件
	e.HTTPErrorHandler = ErrorHandler()
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

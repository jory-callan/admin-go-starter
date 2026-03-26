package main

import (
	"aicode/internal/router"
	"aicode/pkg/app"
	"aicode/pkg/config"
	"aicode/pkg/db"
	"aicode/pkg/logger"
	"aicode/pkg/redis"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "aicode",
		Short: "Go Web Framework with RBAC",
		Run:   run,
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file path")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	// 1. 加载配置
	conf := config.Load(configFile)

	// 2. 初始化日志（设置全局默认）
	logger.New(conf)

	// 3. 初始化数据库
	database := db.New(conf)

	// 4. 初始化 Redis
	rdb := redis.New(conf)

	// 6. 创建 App 结构体
	application := &app.App{
		DB:    database,
		Redis: rdb,
		Conf:  conf,
	}

	// 6. 创建 Echo 实例
	e := echo.New()

	// 7. 注册路由和依赖装配
	router.RegisterRoutes(e, application)

	// 8. 启动服务
	port := conf.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	slog.Info("server starting", "port", port)

	if err := e.Start(":" + port); err != nil {
		slog.Error("failed to start server", "error", err)
		panic(err)
	}
}

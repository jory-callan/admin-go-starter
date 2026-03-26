package main

import (
	"aicode/app"
	"aicode/config"
	"aicode/internal"
	"aicode/internal/router"
	pkglogger "aicode/pkg/logger"
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	rootCmd := &cobra.Command{
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
	// 1. 加载配置 (Default + Unmarshal 智能合并)
	cfg, err := config.Load(configFile)
	if err != nil {
		panic(err)
	}

	// 2. 初始化 Logger（设置全局 slog 默认）
	log := pkglogger.New(cfg.Log)

	// 3. 创建 App 运行时
	application := app.New(cfg, log)

	// 4. 初始化基础设施 (DB, Redis...)
	if err := application.Init(); err != nil {
		log.Error("init infrastructure failed", "error", err)
		panic(err)
	}

	// 5. 数据库迁移
	internal.Migrate(application.DB)

	// 6. 创建 HTTP Server
	e := echo.New()
	router.RegisterRoutes(e, application)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	log.Info("server starting", "addr", addr)

	// 7. 启动服务 (非阻塞)
	go func() {
		if err := e.Start(addr); err != nil {
			log.Error("server stopped", "error", err)
		}
	}()

	// 8. 等待信号并优雅退出
	application.WaitForSignal(context.Background(), func(ctx context.Context) error {
		return e.Shutdown(ctx)
	})
}

package main

import (
	"aicode/app"
	"aicode/config"
	"aicode/info"
	"aicode/internal"
	"aicode/internal/router"
	pkghttp "aicode/pkg/http"
	pkglogger "aicode/pkg/logger"
	"context"
	"fmt"

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
	// 1. 设置应用信息（全局变量，非配置文件）
	info.AppInfo.Name = "aicode"
	info.AppInfo.Version = "0.1.0"
	info.AppInfo.Desc = "Go Web Framework with RBAC"

	// 打印 Banner
	fmt.Print(info.AppInfo.PrintBanner())

	// 2. 加载配置 (DefaultConfig + Unmarshal 智能合并)
	cfg, err := config.Load(configFile)
	if err != nil {
		panic(err)
	}

	// 3. 初始化 Logger（设置全局 slog 默认）
	log := pkglogger.New(cfg.Log, info.AppInfo.Name)

	// 4. 创建 App 运行时
	application := app.New(cfg, log)

	// 5. 初始化基础设施 (DB, Redis, JWT...)
	application.Start()

	// 6. 数据库迁移
	internal.Migrate(application.DB)

	// 7. 创建 HTTP 服务器
	server := pkghttp.New(&cfg.HTTP)

	// 8. 注册路由
	router.RegisterRoutes(server.Engine(), application)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	log.Info("server starting", "version", info.AppInfo.Version, "addr", addr)

	// 9. 启动服务（非阻塞）
	go func() {
		if err := server.Start(); err != nil {
			log.Error("server stopped", "error", err)
		}
	}()

	// 10. 等待信号并优雅退出
	application.WaitForSignal(context.Background(), func(ctx context.Context) error {
		return application.Shutdown()
	})
}

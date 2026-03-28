package main

import (
	"context"
	"fmt"

	"aicode/config"
	"aicode/info"
	"aicode/internal"
	"aicode/internal/app/core"
	"aicode/internal/app/server"
	pkglogger "aicode/pkg/logger"

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
	// 1. 设置应用信息
	info.AppInfo.Name = "aicode"
	info.AppInfo.Version = "0.1.0"
	info.AppInfo.Desc = "Go Web Framework with RBAC"

	// 打印 Banner
	fmt.Print(info.AppInfo.PrintBanner())

	// 2. 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		panic(err)
	}

	// 3. 初始化 Logger
	log := pkglogger.New(cfg.Log, info.AppInfo.Name)

	// 4. 创建 Core (基础设施核心)
	appCore := core.New(cfg, log)

	// 5. 初始化基础设施 (DB, Redis, JWT...)
	appCore.Start()

	// 6. 数据库迁移
	internal.Migrate(appCore.DB)

	// 7. 创建 HTTP Server (依赖 Core，注册路由)
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	httpSrv := server.NewHTTPServer(appCore, addr)

	// 8. 启动 HTTP Server (非阻塞)
	go func() {
		if err := httpSrv.Start(); err != nil {
			log.Error("HTTP server crashed", "error", err)
		}
	}()

	log.Info("server started", "version", info.AppInfo.Version, "addr", addr)

	// 9. 等待信号并优雅退出 (Core 控制信号监听和关闭流程)
	appCore.WaitForSignal(context.Background(), func(ctx context.Context) error {
		return httpSrv.Shutdown(ctx)
	})
}

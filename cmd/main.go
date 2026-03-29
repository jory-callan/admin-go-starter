package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"aicode/config"
	"aicode/info"
	"aicode/internal/app/core"
	"aicode/internal/app/server"
	"aicode/internal/migration"
	"aicode/pkg/logger"

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
	// 设置应用信息
	info.AppInfo.Name = "aicode"
	info.AppInfo.Version = "0.1.0"
	info.AppInfo.Desc = "Go Web Framework with RBAC"

	// 打印 Banner
	fmt.Print(info.AppInfo.PrintBanner())

	// 加载配置
	cfg := config.Load(configFile)

	// 初始化 Logger
	log := logger.New(cfg.Log, info.AppInfo.Name)

	// 创建 Core (基础设施核心)
	appCore := core.New(cfg, log)

	// 初始化基础设施 (DB, Redis, JWT...)
	appCore.Start()

	// 数据库迁移
	migration.Migrate(appCore.DB)

	// 创建 HTTP Server (依赖 Core，注册路由)
	httpSrv := server.NewHTTPServer(appCore, &cfg.HTTP)

	// 启动 HTTP Server (非阻塞)
	go func() {
		if err := httpSrv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server crashed", "error", err)
			panic(err) // 仅对真实错误 panic
		}
	}()

	// 等待信号优雅退出
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	// 关闭 HTTP Server
	httpSrv.Shutdown()
	// 关闭 基础设施
	appCore.Shutdown()
}

package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WaitForSignal 阻塞等待系统信号（SIGINT / SIGTERM），然后执行优雅关闭
func (a *App) WaitForSignal(ctx context.Context, shutdown func(ctx context.Context) error) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	a.Log.Info("received signal, shutting down...", "signal", sig.String())

	// 给予 10 秒的优雅关闭窗口
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 1. 调用外部 shutdown 回调（如 echo.Shutdown）
	if shutdown != nil {
		if err := shutdown(shutdownCtx); err != nil {
			a.Log.Error("server shutdown error", "error", err)
		}
	}

	// 2. 关闭所有基础设施
	if err := a.Shutdown(); err != nil {
		a.Log.Error("app shutdown error", "error", err)
	}

	a.Log.Info("graceful shutdown completed")
}

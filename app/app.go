package app

import (
	"aicode/config"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// App 是框架的核心运行时，包含所有基础设施实例
// 通过结构体字段直接访问全局对象，如 app.DB, app.Redis
type App struct {
	mu sync.RWMutex
	// 内部关闭钩子（逆序执行）
	closers []func() error

	// 配置
	Config *config.AppConfig
	// 全局 Logger
	Log   *slog.Logger
	DB    *gorm.DB      // 默认主数据库 (database)
	Redis *redis.Client // 默认 Redis (redis)
	// === JWT ===
	// JWT 已通过 pkg/jwt.Init() 初始化为全局状态

}

// New 创建 App 实例（不启动任何基础设施）
func New(cfg *config.AppConfig, log *slog.Logger) *App {
	return &App{
		Config:  cfg,
		Log:     log,
		closers: make([]func() error, 0),
	}
}

// Init 初始化所有基础设施（数据库、Redis、JWT 等）
// 调用后可通过 a.DB / a.Redis 直接使用默认实例
func (a *App) Start() {
	a.initDatabases()
	a.initRedis()
	a.initJWT()
}

// registerCloser 注册关闭钩子
func (a *App) registerCloser(fn func() error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.closers = append(a.closers, fn)
}

// Shutdown 优雅关闭：逆序执行所有 closers
func (a *App) Shutdown() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i := len(a.closers) - 1; i >= 0; i-- {
		if err := a.closers[i](); err != nil {
			a.Log.Error("shutdown error", "index", i, "error", err)
		}
	}
	return nil
}

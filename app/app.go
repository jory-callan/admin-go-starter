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

	// 强类型配置
	Config *config.AppConfig

	// 全局 Logger（pkg/logger.New() 返回的 slog.Logger）
	Log *slog.Logger

	// === 数据库实例 ===
	DB        *gorm.DB            // 默认数据库 (database.default)
	Databases map[string]*gorm.DB // 命名数据库实例

	// === Redis 实例 ===
	Redis    *redis.Client            // 默认 Redis (redis.default)
	RedisMap map[string]*redis.Client // 命名 Redis 实例

	// 内部关闭钩子（逆序执行）
	closers []func() error
}

// New 创建 App 实例（不启动任何基础设施）
func New(cfg *config.AppConfig, log *slog.Logger) *App {
	return &App{
		Config:    cfg,
		Log:       log,
		Databases: make(map[string]*gorm.DB),
		RedisMap:  make(map[string]*redis.Client),
		closers:   make([]func() error, 0),
	}
}

// Init 初始化所有基础设施（数据库、Redis 等）
// 调用后可通过 a.DB / a.Redis 直接使用默认实例
func (a *App) Init() error {
	if err := a.initDatabases(); err != nil {
		return err
	}
	if err := a.initRedis(); err != nil {
		return err
	}
	return nil
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

	// 逆序关闭
	for i := len(a.closers) - 1; i >= 0; i-- {
		if err := a.closers[i](); err != nil {
			a.Log.Error("shutdown error", "index", i, "error", err)
		}
	}
	return nil
}

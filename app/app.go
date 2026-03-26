package app

import (
	"database/sql"
	"sync"

	"github.com/redis/go-redis/v9"
	// 引入你的配置包
	"aicode/pkg/config"
)

// App 是框架的核心运行时
type App struct {
	mu sync.RWMutex

	// 配置
	Config *config.AppConfig

	// === 基础设施字段 (显式声明，享受类型安全) ===

	// 默认数据库
	DB *sql.DB

	// 扩展数据库
	LogDB   *sql.DB
	UserDB  *sql.DB
	OrderDB *sql.DB

	// 默认 Redis
	Redis *redis.Client

	// 扩展 Redis
	CacheRedis *redis.Client

	// 其他组件...

	// 内部关闭钩子
	closers []func() error
}

// New 创建新实例
func New(cfg *config.AppConfig) *App {
	return &App{
		Config:  cfg,
		closers: make([]func() error, 0),
	}
}

// Run 启动流程
func (a *App) Run() error {
	// 1. 初始化所有组件 (调用分散在各个文件中的 init 方法)
	if err := a.initDatabases(); err != nil {
		return err
	}
	if err := a.initRedis(); err != nil {
		return err
	}

	// 2. 监听信号并阻塞
	// ... (信号处理逻辑)

	return nil
}

// Shutdown 优雅关闭
func (a *App) Shutdown() error {
	// 逆序执行 closers
	// ...
	return nil
}

// registerCloser 内部辅助方法，供 db.go, redis.go 使用
func (a *App) registerCloser(fn func() error) {
	a.closers = append(a.closers, fn)
}

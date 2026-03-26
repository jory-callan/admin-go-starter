package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// New 创建全局默认 logger，并设置为 slog 默认实例
// 调用后各组件可直接使用 slog.Info() / slog.Error() 等包级函数
func New(conf *viper.Viper) {
	cfg := DefaultConfig()

	if conf.IsSet("logger.level") {
		cfg.Level = parseLevel(conf.GetString("logger.level"))
	}
	if conf.IsSet("logger.format") {
		cfg.Format = conf.GetString("logger.format")
	}

	logger := slog.New(newHandler(os.Stdout, cfg))
	logger = logger.With(slog.String("app", "app"))
	slog.SetDefault(logger)
}

// C 创建带 component 标记的子 logger
// 用法: log := logger.C("db")  然后 log.Info("xxx")
func C(component string) *slog.Logger {
	return slog.Default().With(slog.String("component", component))
}

// newHandler 根据配置创建 handler
func newHandler(w io.Writer, cfg Config) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	switch strings.ToLower(cfg.Format) {
	case "text":
		return slog.NewTextHandler(w, opts)
	default:
		return slog.NewJSONHandler(w, opts)
	}
}

// parseLevel 将字符串转为 slog.Level
func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

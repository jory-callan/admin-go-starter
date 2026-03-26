package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

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

// newHandler 根据配置创建 slog.Handler
func newHandler(w io.Writer, cfg Config) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
	}

	switch strings.ToLower(cfg.Format) {
	case "text":
		return slog.NewTextHandler(w, opts)
	default:
		return slog.NewJSONHandler(w, opts)
	}
}

// New 根据 Config 创建 *slog.Logger 并设为全局默认
// appName 会作为每条日志的 "app" 字段
func New(cfg Config, appName string) *slog.Logger {
	var writer io.Writer

	// 优先使用 file_path，其次使用 output
	filePath := cfg.FilePath
	if filePath == "" && cfg.Output != "" && cfg.Output != "stdout" {
		filePath = cfg.Output
	}

	if filePath != "" {
		// 输出到文件（支持日志轮转）
		dir := filepath.Dir(filePath)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create log directory %s: %v\n", dir, err)
				writer = os.Stdout
			} else {
				writer = &lumberjack.Logger{
					Filename:   filePath,
					MaxSize:    cfg.MaxSize,
					MaxBackups: cfg.MaxBackups,
					MaxAge:     cfg.MaxAge,
					Compress:   cfg.Compress,
				}
			}
		} else {
			writer = &lumberjack.Logger{
				Filename:   filePath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}
		}
	} else {
		writer = os.Stdout
	}

	h := newHandler(writer, cfg)
	l := slog.New(h).With(slog.String("app", appName))
	slog.SetDefault(l)
	return l
}

// C 创建带 component 标记的子 logger
func C(component string) *slog.Logger {
	return slog.Default().With(slog.String("component", component))
}

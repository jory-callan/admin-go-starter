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

// Config 日志配置
type Config struct {
	Level      string `mapstructure:"level" yaml:"level"`             // debug, info, warn, error
	Format     string `mapstructure:"format" yaml:"format"`           // json, text
	Output     string `mapstructure:"output" yaml:"output"`           // stdout 或文件路径
	FilePath   string `mapstructure:"file_path" yaml:"file_path"`     // 日志文件路径
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size"`       // 单文件最大 MB
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups"` // 保留旧文件数
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age"`         // 旧文件保留天数
	Compress   bool   `mapstructure:"compress" yaml:"compress"`       // 是否压缩
}

// GetDefault 返回日志默认配置
func GetDefault() Config {
	return Config{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
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
// 返回的 logger 已携带 app 标记
func New(cfg Config) *slog.Logger {
	var writer io.Writer

	switch strings.ToLower(cfg.Output) {
	case "stdout", "":
		writer = os.Stdout
	default:
		// 输出到文件（支持日志轮转）
		dir := filepath.Dir(cfg.Output)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create log directory %s: %v\n", dir, err)
				writer = os.Stdout
			} else {
				filePath := cfg.Output
				if filePath == "" && cfg.FilePath != "" {
					filePath = cfg.FilePath
				}
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
				Filename:   cfg.Output,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}
		}
	}

	h := newHandler(writer, cfg)
	l := slog.New(h).With(slog.String("app", "aicode"))
	slog.SetDefault(l)
	return l
}

// C 创建带 component 标记的子 logger
func C(component string) *slog.Logger {
	return slog.Default().With(slog.String("component", component))
}

package logger

import "log/slog"

// Config 日志配置
type Config struct {
	Level  slog.Level `json:"level"`  // 日志级别: debug=-4, info=0, warn=4, error=8
	Format string     `json:"format"` // 输出格式: json / text
}

// DefaultConfig 默认日志配置
func DefaultConfig() Config {
	return Config{
		Level:  slog.LevelInfo,
		Format: "json",
	}
}

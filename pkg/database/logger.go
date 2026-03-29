package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm"
	gormlogx "gorm.io/gorm/logger"
)

// Silent (1) >= Error (2) >= Warn (3) >= Info (4)。

// gormLogger GORM日志适配器
type gormLogger struct {
	level         gormlogx.LogLevel
	slowThreshold time.Duration
}

// newGormLogger 创建GORM日志适配器
func newGormLogger(level string, slowThreshold time.Duration) gormlogx.Interface {
	var gormLevel gormlogx.LogLevel
	switch strings.ToLower(level) {
	case "silent":
		gormLevel = gormlogx.Silent
	case "error":
		gormLevel = gormlogx.Error
	case "warn":
		gormLevel = gormlogx.Warn
	case "info":
		gormLevel = gormlogx.Info
	default:
		gormLevel = gormlogx.Info
	}
	return &gormLogger{
		level:         gormLevel,
		slowThreshold: slowThreshold,
	}
}

func (l *gormLogger) LogMode(level gormlogx.LogLevel) gormlogx.Interface {
	l.level = level
	return l
}
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogx.Info {
		slog.Info(fmt.Sprintf(msg, data...))
	}
}
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogx.Warn {
		slog.Warn(fmt.Sprintf(msg, data...))
	}
}
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogx.Error {
		slog.Error(fmt.Sprintf(msg, data...))
	}
}
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 如果日志级别低于Silent，则不输出日志
	if l.level <= gormlogx.Silent {
		return
	}
	// 计算执行时间
	elapsed := time.Since(begin)
	// 获取SQL和行数
	sql, rows := fc()
	// 构建日志字段
	fields := fmt.Sprintf("sql=%s, rows=%d, cost=%s", sql, rows, elapsed)
	// 根据错误和执行时间输出日志
	switch {
	// 错误日志, 排除记录未找到
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.level >= gormlogx.Error:
		fields = fields + fmt.Sprintf("error=%s", err)
		slog.Error("gorm error --> " + fields)
	// 慢查询日志, 优先匹配
	case elapsed > l.slowThreshold && l.slowThreshold > 0 && l.level >= gormlogx.Warn:
		fields = fields + fmt.Sprintf("threshold=%s", l.slowThreshold)
		slog.Warn("slow query --> " + fields)
	//
	case l.level >= gormlogx.Info:
		slog.Info("gorm query --> " + fields)
	}
}

func (l *gormLogger) detectCallerSkip() int {
	//find the first caller, the one is this file:lineNumber
	var buf [10]uintptr
	n := runtime.Callers(0, buf[:])
	frames := runtime.CallersFrames(buf[:n])
	for i := 0; i < n; i++ {
		frame, _ := frames.Next()
		fmt.Printf("[%d] %s:%d\n", i, filepath.Base(frame.File), frame.Line)
	}

	// Record a test log and capture its PC
	var testPC uintptr
	func() {
		var pcs [1]uintptr
		runtime.Callers(1, pcs[:]) // skip this func
		testPC = pcs[0]
	}()
	// Now simulate a log call and see what skip gives us testPC
	for skip := 0; skip <= 10; skip++ {
		var pcs [1]uintptr
		n := runtime.Callers(skip, pcs[:])
		if n > 0 && pcs[0] == testPC {
			return skip
		}
	}
	return 1 // fallback
}

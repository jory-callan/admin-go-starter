package core

import (
	"fmt"
	"log/slog"

	"aicode/pkg/jwt"
)

// initJWT 初始化 JWT 全局参数
func (a *App) initJWT() {
	var msg string
	jwt.Init(a.Config.JWT)
	msg = fmt.Sprintf("jwt initialized, issuer: %s, expires: %s", a.Config.JWT.Issuer, a.Config.JWT.Expires)
	slog.Info(msg)
}

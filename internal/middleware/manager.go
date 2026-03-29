package middleware

import "aicode/internal/app/core"

type Manager struct {
	core *core.App
}

func NewManager(core *core.App) *Manager {
	return &Manager{core: core}
}

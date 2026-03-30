package config

import (
	"aicode/pkg/database"
	pkghttp "aicode/pkg/http"
	pkgjwt "aicode/pkg/jwt"
	pkglogger "aicode/pkg/logger"
)

// AppConfig 应用全局配置（强类型）
// 所有子配置均引用 pkg 层各自的 Config 结构体，统一由各包的 DefaultConfig() 提供默认值
type AppConfig struct {
	Log      pkglogger.Config `mapstructure:"log" yaml:"log"`
	HTTP     pkghttp.Config   `mapstructure:"http" yaml:"http"`
	Database database.Config  `mapstructure:"database" yaml:"database"`
	JWT      pkgjwt.Config    `mapstructure:"jwt" yaml:"jwt"`
}

// DefaultConfig 返回全局默认配置
// 各子配置的默认值由 pkg 层各自的 DefaultConfig() 提供
func DefaultConfig() AppConfig {
	return AppConfig{
		Log:      pkglogger.DefaultConfig(),
		HTTP:     pkghttp.DefaultConfig(),
		Database: database.DefaultConfig(),
		JWT:      pkgjwt.DefaultConfig(),
	}
}

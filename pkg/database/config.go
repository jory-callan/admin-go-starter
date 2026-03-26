package database

// Config 数据库配置
//
// YAML 配置示例:
//
//	database:
//	  driver: "mysql"                    # 驱动类型: mysql, postgres, sqlite3
//	  dsn: "root:password@tcp(127.0.0.1:3306)/main_db?charset=utf8mb4&parseTime=True"
//	  max_open_conns: 100               # 最大打开连接数
//	  max_idle_conns: 10                # 最大空闲连接数
//	  conn_max_lifetime: 1800           # 连接最大存活时间(秒)
//	  conn_max_idle_time: 300           # 连接最大空闲时间(秒)
//	  log_level: "warn"                 # gorm 日志级别: silent, error, warn, info
type Config struct {
	Driver          string `mapstructure:"driver" yaml:"driver" json:"driver"`                         // mysql, postgres, sqlite
	DSN             string `mapstructure:"dsn" yaml:"dsn" json:"dsn"`                                 // 连接串
	MaxOpenConns    int    `mapstructure:"max_open_conns" yaml:"max_open_conns" json:"max_open_conns"` // 最大打开连接数
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns" json:"max_idle_conns"` // 最大空闲连接数
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime" json:"conn_max_lifetime"` // 连接最大存活时间(秒)
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time" json:"conn_max_idle_time"` // 连接最大空闲时间(秒)
	LogLevel        string `mapstructure:"log_level" yaml:"log_level" json:"log_level"`             // gorm 日志级别
}

// DefaultConfig 返回数据库默认配置
func DefaultConfig() Config {
	return Config{
		Driver:          "sqlite",
		DSN:             "demo.sqlite.db",
		MaxOpenConns:    50,
		MaxIdleConns:    10,
		ConnMaxLifetime: 1800, // 30 分钟
		ConnMaxIdleTime: 300,  // 5 分钟
		LogLevel:        "warn",
	}
}

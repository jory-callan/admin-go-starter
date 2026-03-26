package redis

// Config Redis 配置
//
// YAML 配置示例:
//
//	redis:
//	  addr: "127.0.0.1:6379"       # Redis 地址 host:port
//	  password: ""                 # 密码（无密码则留空）
//	  db: 0                        # 数据库索引
//	  pool_size: 100               # 连接池大小
//	  min_idle_conns: 10           # 最小空闲连接数
//	  dial_timeout: 5              # 连接超时(秒)
//	  read_timeout: 3              # 读超时(秒)
//	  write_timeout: 3             # 写超时(秒)
type Config struct {
	Addr         string `mapstructure:"addr" yaml:"addr" json:"addr"`                          // Redis 地址 host:port
	Password     string `mapstructure:"password" yaml:"password" json:"password"`              // 密码
	DB           int    `mapstructure:"db" yaml:"db" json:"db"`                                // 数据库索引
	PoolSize     int    `mapstructure:"pool_size" yaml:"pool_size" json:"pool_size"`           // 连接池大小
	MinIdleConns int    `mapstructure:"min_idle_conns" yaml:"min_idle_conns" json:"min_idle_conns"` // 最小空闲连接数
	DialTimeout  int    `mapstructure:"dial_timeout" yaml:"dial_timeout" json:"dial_timeout"`  // 连接超时(秒)
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout" json:"read_timeout"`  // 读超时(秒)
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout" json:"write_timeout"` // 写超时(秒)
}

// DefaultConfig 返回 Redis 默认配置
func DefaultConfig() Config {
	return Config{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		DialTimeout:  5,
		ReadTimeout:  3,
		WriteTimeout: 3,
	}
}

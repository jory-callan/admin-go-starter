package logger

// Config 日志配置
//
// YAML 配置示例:
//
//	log:
//	  level: "info"              # 日志级别: debug, info, warn, error, fatal
//	  format: "json"             # 日志格式: json, text
//	  output: "stdout"           # 输出目标: stdout, 或文件路径（如 /var/log/app.log）
//	  file_path: "/var/log/app.log"  # 日志文件路径，填了就有文件，没填就没有文件
//	  max_size: 100              # 单个日志文件最大大小(MB)
//	  max_backups: 5             # 保留的旧日志文件数量
//	  max_age: 7                 # 旧日志文件保留天数(天)
//	  compress: true             # 是否压缩旧日志文件
type Config struct {
	Level      string `mapstructure:"level" yaml:"level" json:"level"`                 // debug, info, warn, error
	Format     string `mapstructure:"format" yaml:"format" json:"format"`              // json, text
	Output     string `mapstructure:"output" yaml:"output" json:"output"`              // stdout 或文件路径
	FilePath   string `mapstructure:"file_path" yaml:"file_path" json:"file_path"`    // 日志文件路径（与 output 二选一）
	MaxSize    int    `mapstructure:"max_size" yaml:"max_size" json:"max_size"`        // 单文件最大 MB
	MaxBackups int    `mapstructure:"max_backups" yaml:"max_backups" json:"max_backups"` // 保留旧文件数
	MaxAge     int    `mapstructure:"max_age" yaml:"max_age" json:"max_age"`           // 旧文件保留天数
	Compress   bool   `mapstructure:"compress" yaml:"compress" json:"compress"`        // 是否压缩
}

// DefaultConfig 返回日志默认配置
func DefaultConfig() Config {
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

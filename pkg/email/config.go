package email

// Config 邮件服务配置
//
// YAML 配置示例:
//
//	email:
//	  smtp_host: "smtp.example.com"     # SMTP 服务器地址
//	  smtp_port: 587                     # SMTP 端口 (25, 465, 587)
//	  username: "noreply@example.com"    # 发件人邮箱
//	  password: "your-smtp-password"     # SMTP 密码或授权码
//	  from_name: "My App"                # 发件人名称
//	  from_email: "noreply@example.com"  # 默认发件人邮箱
//	  use_tls: true                      # 是否启用 TLS
//	  timeout: 10                        # 连接超时时间 (秒)
type Config struct {
	SMTPHost  string `mapstructure:"smtp_host" yaml:"smtp_host" json:"smtp_host"`    // SMTP 服务器地址
	SMTPPort  int    `mapstructure:"smtp_port" yaml:"smtp_port" json:"smtp_port"`    // SMTP 端口
	Username  string `mapstructure:"username" yaml:"username" json:"username"`       // 发件人邮箱
	Password  string `mapstructure:"password" yaml:"password" json:"password"`       // SMTP 密码或授权码
	FromName  string `mapstructure:"from_name" yaml:"from_name" json:"from_name"`    // 发件人名称
	FromEmail string `mapstructure:"from_email" yaml:"from_email" json:"from_email"` // 默认发件人邮箱
	UseTLS    bool   `mapstructure:"use_tls" yaml:"use_tls" json:"use_tls"`          // 是否启用 TLS
	Timeout   int    `mapstructure:"timeout" yaml:"timeout" json:"timeout"`          // 连接超时时间 (秒)
}

// DefaultConfig 返回邮件服务默认配置
func DefaultConfig() Config {
	return Config{
		SMTPHost:  "smtp.example.com",
		SMTPPort:  587,
		Username:  "noreply@example.com",
		Password:  "",
		FromName:  "My App",
		FromEmail: "noreply@example.com",
		UseTLS:    true,
		Timeout:   10,
	}
}

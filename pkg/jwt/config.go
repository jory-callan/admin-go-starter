package jwt

// Config JWT 配置
//
// YAML 配置示例:
//
//	jwt:
//	  secret: "your-32-char-secret-key-here"  # 密钥（生产环境必须修改）
//	  expires: "24h"                           # Token 过期时间
//	  issuer: "myapp"                          # 签发者
type Config struct {
	Secret  string `mapstructure:"secret" yaml:"secret" json:"secret"`     // 密钥
	Expires string `mapstructure:"expires" yaml:"expires" json:"expires"`   // 过期时间 (如 "24h", "7d")
	Issuer  string `mapstructure:"issuer" yaml:"issuer" json:"issuer"`      // 签发者
}

// DefaultConfig 返回 JWT 默认配置
func DefaultConfig() Config {
	return Config{
		Secret:  "change-me-in-production-32chars!",
		Expires: "24h",
		Issuer:  "aicode",
	}
}

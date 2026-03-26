package rabbitmq

// Config RabbitMQ 配置
//
// YAML 配置示例:
//
//	rabbitmq:
//	  url: "amqp://guest:guest@localhost:5672/"  # RabbitMQ 连接 URL
//	  exchange: "app_events"                      # 默认交换机
//	  queue: "app_queue"                          # 默认队列
type Config struct {
	URL      string `mapstructure:"url" yaml:"url" json:"url"`           // RabbitMQ 连接 URL
	Exchange string `mapstructure:"exchange" yaml:"exchange" json:"exchange"` // 默认交换机
	Queue    string `mapstructure:"queue" yaml:"queue" json:"queue"`       // 默认队列
}

// DefaultConfig 返回 RabbitMQ 默认配置
func DefaultConfig() Config {
	return Config{
		URL:      "amqp://guest:guest@localhost:5672/",
		Exchange: "app_events",
		Queue:    "app_queue",
	}
}

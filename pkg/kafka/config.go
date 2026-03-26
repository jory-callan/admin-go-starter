package kafka

// Config Kafka 配置
//
// YAML 配置示例:
//
//	kafka:
//	  brokers: ["127.0.0.1:9092"]  # Kafka Broker 列表
//	  topic: "app_events"           # 默认 Topic
//	  group_id: "app_consumer"      # 消费者组 ID
//	  partition: 0                  # 默认分区
//	  offset: "latest"              # 消费偏移量: earliest, latest
type Config struct {
	Brokers  []string `mapstructure:"brokers" yaml:"brokers" json:"brokers"` // Kafka Broker 列表
	Topic    string   `mapstructure:"topic" yaml:"topic" json:"topic"`       // 默认 Topic
	GroupID  string   `mapstructure:"group_id" yaml:"group_id" json:"group_id"` // 消费者组 ID
	Partition int     `mapstructure:"partition" yaml:"partition" json:"partition"` // 默认分区
	Offset   string   `mapstructure:"offset" yaml:"offset" json:"offset"`     // 消费偏移量: earliest, latest
}

// DefaultConfig 返回 Kafka 默认配置
func DefaultConfig() Config {
	return Config{
		Brokers:  []string{"127.0.0.1:9092"},
		Topic:    "app_events",
		GroupID:  "app_consumer",
		Partition: 0,
		Offset:   "latest",
	}
}

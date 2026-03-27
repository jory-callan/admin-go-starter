package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

// Client Kafka 客户端（管理生产者和消费者）
type Client struct {
	cfg       Config
	producer  *kafka.Writer
	consumers map[string]*kafka.Reader
	mu        sync.RWMutex
}

// NewClient 创建 Kafka 客户端
func NewClient(cfg Config) (*Client, error) {
	client := &Client{
		cfg:       cfg,
		consumers: make(map[string]*kafka.Reader),
	}

	// 初始化生产者
	client.producer = &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  cfg.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return client, nil
}

// Producer 返回生产者实例
func (c *Client) Producer() *kafka.Writer {
	return c.producer
}

// GetConsumer 获取或创建指定 Topic 的消费者
func (c *Client) GetConsumer(topic string, groupID string) *kafka.Reader {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s:%s", topic, groupID)
	if reader, ok := c.consumers[key]; ok {
		return reader
	}

	// 创建新的消费者
	offset := kafka.FirstOffset
	if c.cfg.Offset == "latest" {
		offset = kafka.LastOffset
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.cfg.Brokers,
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     0,    // 阻塞等待
		StartOffset: offset,
	})

	c.consumers[key] = reader
	return reader
}

// Close 关闭所有资源
func (c *Client) Close() error {
	var errs []error

	// 关闭生产者
	if c.producer != nil {
		if err := c.producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close producer: %w", err))
		}
	}

	// 关闭所有消费者
	c.mu.Lock()
	for _, reader := range c.consumers {
		if err := reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close consumer: %w", err))
		}
	}
	c.consumers = make(map[string]*kafka.Reader)
	c.mu.Unlock()

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}
	return nil
}

// HealthCheck 检查 Kafka 连接健康状态
func (c *Client) HealthCheck(ctx context.Context) error {
	conn, err := kafka.DialLeader(ctx, "tcp", c.cfg.Brokers[0], c.cfg.Topic, 0)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	_, err = conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("read partitions: %w", err)
	}

	return nil
}

package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client RabbitMQ 客户端
type Client struct {
	cfg          Config
	conn         *amqp.Connection
	channel      *amqp.Channel
	mu           sync.RWMutex
	declarations map[string]bool // 记录已声明的队列
}

// NewClient 创建 RabbitMQ 客户端
func NewClient(cfg Config) (*Client, error) {
	client := &Client{
		cfg:          cfg,
		declarations: make(map[string]bool),
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// connect 建立连接
func (c *Client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := amqp.Dial(c.cfg.URL)
	if err != nil {
		return fmt.Errorf("dial rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("open channel: %w", err)
	}

	c.conn = conn
	c.channel = ch

	return nil
}

// Channel 返回当前 channel
func (c *Client) Channel() *amqp.Channel {
	return c.channel
}

// EnsureQueue 确保队列存在（幂等）
func (c *Client) EnsureQueue(queueName string, durable bool) error {
	c.mu.RLock()
	if _, exists := c.declarations[queueName]; exists {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.channel.QueueDeclare(
		queueName,
		durable,
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue %s: %w", queueName, err)
	}

	c.declarations[queueName] = true
	return nil
}

// Publish 发布消息
func (c *Client) Publish(exchange, key string, body []byte, headers map[string]interface{}) error {
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     amqp.Table(headers),
	}

	if err := c.channel.Publish(exchange, key, false, false, msg); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}

// PublishJSON 发布 JSON 消息
func (c *Client) PublishJSON(exchange, key string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	return c.Publish(exchange, key, body, nil)
}

// Consume 消费消息
func (c *Client) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		queue,
		consumer,
		autoAck,
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume from queue %s: %w", queue, err)
	}

	return msgs, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close channel: %w", err))
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}
	return nil
}

// HealthCheck 检查连接健康状态
func (c *Client) HealthCheck(ctx context.Context) error {
	select {
	case <-c.conn.NotifyClose(make(chan *amqp.Error)):
		return fmt.Errorf("connection closed")
	default:
		return nil
	}
}

// Reconnect 重新连接
func (c *Client) Reconnect() error {
	c.Close()
	return c.connect()
}

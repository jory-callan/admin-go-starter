package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	ErrConsumerClosed = errors.New("consumer closed")
)

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	Topic   string
	GroupID string
	Handler func(ctx context.Context, msg Message) error
	Workers int // 并发处理 worker 数量
}

// Consumer 消费者封装
type Consumer struct {
	client *Client
	reader *kafka.Reader
	cfg    ConsumerConfig
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	closed bool
	mu     sync.RWMutex
}

// NewConsumer 创建消费者
func NewConsumer(client *Client, cfg ConsumerConfig) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())

	consumer := &Consumer{
		client: client,
		reader: client.GetConsumer(cfg.Topic, cfg.GroupID),
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	if consumer.cfg.Workers <= 0 {
		consumer.cfg.Workers = 1
	}

	return consumer
}

// Start 启动消费者
func (c *Consumer) Start() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrConsumerClosed
	}
	c.mu.Unlock()

	// 启动多个 worker 并发处理消息
	for i := 0; i < c.cfg.Workers; i++ {
		c.wg.Add(1)
		go c.worker(i)
	}

	return nil
}

// worker 消费消息的工作协程
func (c *Consumer) worker(id int) {
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.reader.FetchMessage(c.ctx)
			if err != nil {
				if c.ctx.Err() != nil {
					return
				}
				continue
			}

			// 转换为 Message 结构
			message := Message{
				Key:     string(msg.Key),
				Value:   msg.Value,
				Headers: make(map[string]string),
			}

			// 解析 headers
			for _, h := range msg.Headers {
				message.Headers[h.Key] = string(h.Value)
			}

			// 处理消息
			if err := c.cfg.Handler(c.ctx, message); err != nil {
				// 处理失败，可以选择重新入队或记录日志
				continue
			}

			// 提交 offset
			if err := c.reader.CommitMessages(c.ctx, msg); err != nil {
				// 提交失败，下次会重新消费
			}
		}
	}
}

// Stop 停止消费者
func (c *Consumer) Stop() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrConsumerClosed
	}
	c.closed = true
	c.mu.Unlock()

	// 取消上下文
	c.cancel()

	// 等待所有 worker 退出
	c.wg.Wait()

	return nil
}

// ReceiveMessage 接收单条消息（适用于手动控制消费场景）
func (c *Consumer) ReceiveMessage(ctx context.Context) (*Message, error) {
	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}

	message := &Message{
		Key:     string(msg.Key),
		Value:   msg.Value,
		Headers: make(map[string]string),
	}

	// 解析 headers
	for _, h := range msg.Headers {
		message.Headers[h.Key] = string(h.Value)
	}

	return message, nil
}

// CommitMessage 提交消息 offset
func (c *Consumer) CommitMessage(ctx context.Context, msg kafka.Message) error {
	return c.reader.CommitMessages(ctx, msg)
}

// UnmarshalValue 解析消息值为指定结构
func (m *Message) UnmarshalValue(v interface{}) error {
	return json.Unmarshal(m.Value, v)
}

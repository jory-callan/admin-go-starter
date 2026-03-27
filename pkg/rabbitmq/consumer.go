package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrConsumerClosed = errors.New("consumer closed")
)

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	Queue         string
	Consumer      string
	AutoAck       bool
	PrefetchCount int // 预取数量
	Handler       func(ctx context.Context, msg *Message) error
	Workers       int // 并发处理 worker 数量
}

// Consumer 消费者封装
type Consumer struct {
	client *Client
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

	// 确保队列存在
	if err := c.client.EnsureQueue(c.cfg.Queue, true); err != nil {
		return err
	}

	// 设置 QoS
	if c.cfg.PrefetchCount > 0 {
		if err := c.client.channel.Qos(
			c.cfg.PrefetchCount,
			0,
			false,
		); err != nil {
			return fmt.Errorf("set qos: %w", err)
		}
	}

	// 获取消息通道
	msgs, err := c.client.Consume(c.cfg.Queue, c.cfg.Consumer, c.cfg.AutoAck)
	if err != nil {
		return err
	}

	// 启动多个 worker 并发处理消息
	for i := 0; i < c.cfg.Workers; i++ {
		c.wg.Add(1)
		go c.worker(i, msgs)
	}

	return nil
}

// worker 消费消息的工作协程
func (c *Consumer) worker(id int, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				return
			}

			msg := ToDelivery(d)

			// 处理消息
			if err := c.cfg.Handler(c.ctx, msg); err != nil {
				// 处理失败，如果非自动 ack，可以 nack 重新入队
				if !c.cfg.AutoAck {
					d.Nack(false, true) // requeue=true
				}
				continue
			}

			// 确认消息
			if !c.cfg.AutoAck {
				d.Ack(false)
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
	// 注意：这个方法需要配合 Get 使用，而不是 Consume
	msg, ok, err := c.client.channel.Get(c.cfg.Queue, !c.cfg.AutoAck)
	if !ok || err != nil {
		return nil, fmt.Errorf("get message: %w", err)
	}

	return ToDelivery(msg), nil
}

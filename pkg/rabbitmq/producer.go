package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Message RabbitMQ 消息结构
type Message struct {
	DeliveryTag   uint64
	Body          []byte
	Headers       map[string]interface{}
	RoutingKey    string
	ReplyTo       string
	CorrelationID string
}

// Producer 生产者封装
type Producer struct {
	client   *Client
	exchange string
}

// NewProducer 创建生产者
func NewProducer(client *Client, exchange string) *Producer {
	return &Producer{
		client:   client,
		exchange: exchange,
	}
}

// Send 发送消息
func (p *Producer) Send(ctx context.Context, routingKey string, body []byte, headers map[string]interface{}) error {
	return p.client.Publish(p.exchange, routingKey, body, headers)
}

// SendJSON 发送 JSON 消息
func (p *Producer) SendJSON(ctx context.Context, routingKey string, data interface{}) error {
	return p.client.PublishJSON(p.exchange, routingKey, data)
}

// SendWithConfirm 发送消息并等待确认
func (p *Producer) SendWithConfirm(ctx context.Context, routingKey string, body []byte) error {
	// 启用 confirm 模式
	if err := p.client.channel.Confirm(false); err != nil {
		return fmt.Errorf("enable confirm mode: %w", err)
	}

	confirms := p.client.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	if err := p.client.channel.Publish(p.exchange, routingKey, false, false, msg); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	// 等待确认
	select {
	case confirmed := <-confirms:
		if !confirmed.Ack {
			return fmt.Errorf("message not acknowledged")
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// UnmarshalMessage 解析消息体
func (m *Message) UnmarshalMessage(v interface{}) error {
	return json.Unmarshal(m.Body, v)
}

// ToDelivery 从 amqp.Delivery 转换为 Message
func ToDelivery(d amqp.Delivery) *Message {
	return &Message{
		DeliveryTag:   d.DeliveryTag,
		Body:          d.Body,
		Headers:       d.Headers,
		RoutingKey:    d.RoutingKey,
		ReplyTo:       d.ReplyTo,
		CorrelationID: d.CorrelationId,
	}
}

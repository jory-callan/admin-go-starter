package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Message Kafka 消息结构
type Message struct {
	Key     string            `json:"key"`
	Value   []byte            `json:"value"`
	Headers map[string]string `json:"headers"`
}

// Producer 生产者封装
type Producer struct {
	client *Client
}

// NewProducer 创建生产者
func NewProducer(client *Client) *Producer {
	return &Producer{client: client}
}

// Send 发送消息（同步）
func (p *Producer) Send(ctx context.Context, topic string, msg Message) error {
	writer := p.client.producer
	writer.Topic = topic

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Key),
		Value: msg.Value,
	}

	// 添加 headers
	if len(msg.Headers) > 0 {
		headers := make([]kafka.Header, 0, len(msg.Headers))
		for k, v := range msg.Headers {
			headers = append(headers, kafka.Header{
				Key:   k,
				Value: []byte(v),
			})
		}
		kafkaMsg.Headers = headers
	}

	if err := writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

// SendJSON 发送 JSON 消息
func (p *Producer) SendJSON(ctx context.Context, topic string, key string, data interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	return p.Send(ctx, topic, Message{
		Key:   key,
		Value: value,
	})
}

// SendBatch 批量发送消息
func (p *Producer) SendBatch(ctx context.Context, topic string, messages []Message) error {
	writer := p.client.producer
	writer.Topic = topic

	kafkaMsgs := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		kafkaMsgs[i] = kafka.Message{
			Key:   []byte(msg.Key),
			Value: msg.Value,
		}

		// 添加 headers
		if len(msg.Headers) > 0 {
			headers := make([]kafka.Header, 0, len(msg.Headers))
			for k, v := range msg.Headers {
				headers = append(headers, kafka.Header{
					Key:   k,
					Value: []byte(v),
				})
			}
			kafkaMsgs[i].Headers = headers
		}
	}

	if err := writer.WriteMessages(ctx, kafkaMsgs...); err != nil {
		return fmt.Errorf("send batch messages: %w", err)
	}

	return nil
}

package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerAddr string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		writer: w,
	}
}

func (p *Producer) SendMessage(ctx context.Context, message []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Value: message,
	})

	if err != nil {
		return fmt.Errorf("kafka.SendMessage: failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("kafka.Close: failed to close producer: %w", err)
	}

	return nil
}

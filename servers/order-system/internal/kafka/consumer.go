package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokerAddr, topic, groutpID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    topic,
		GroupID:  groutpID,
		MaxBytes: 10e6,
	})

	return &Consumer{reader: r}
}

func (c *Consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("FetchMessage: failed to fetch message: %w", err)
	}

	return msg, nil
}

func (c *Consumer) CommitMessage(ctx context.Context, msg kafka.Message) error {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("CommitMessage: failed to commit message: %w", err)
	}

	return nil
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("kafka.close: %w", err)
	}

	return nil
}

package producer

import (
	"context"
	"encoding/json"
	"errors"
	"shorter/internal/events"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewKafkaProducer(brokers []string, topic string, logger *zap.Logger) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           10 * time.Second,
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{
		writer: writer,
		logger: logger,
	}
}

func (p *KafkaProducer) SendClick(ctx context.Context, event *events.ClickEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: value,
		Time:  event.Timestamp,
	}

	const retries = 3
	for range retries {
		err := p.writer.WriteMessages(ctx, msg)
		if errors.Is(err, kafka.UnknownTopicOrPartition) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}
		if err != nil {
			p.logger.Error("failed to write message to kafka", zap.Error(err), zap.String("alias", event.Alias))
			return err
		}
		break
	}

	p.logger.Info("click event send to kafka", zap.String("alias", event.Alias), zap.String("topic", p.writer.Topic))

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

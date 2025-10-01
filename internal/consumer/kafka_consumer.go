package consumer

import (
	"context"
	"encoding/json"
	"shorter/internal/enricher"
	"shorter/internal/metrics"
	"shorter/internal/repository"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	reader          *kafka.Reader
	enricher        enricher.Enricher
	repo            *repository.AnalyticsRepository
	logger          *zap.Logger
	workerCount     int
	shutdownTimeout time.Duration
}

func NewKafkaConsumer(
	brokers []string,
	topic string,
	groupId string,
	enricher enricher.Enricher,
	repo *repository.AnalyticsRepository,
	logger *zap.Logger,
) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupId,
			MinBytes: 10e3,
			MaxBytes: 10e6,
			MaxWait:  1 * time.Second,
		}),
		enricher:        enricher,
		repo:            repo,
		logger:          logger,
		workerCount:     3,
		shutdownTimeout: 30 * time.Second,
	}
}

func (c *KafkaConsumer) Start(parentCtx context.Context) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	for i := 0; i < c.workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.worker(ctx, id)
		}(i)
	}

	// wait done of context
	<-parentCtx.Done()

	// cancel ctx for workers
	cancel()

	// wait workers stop
	wg.Wait()

	return c.reader.Close()
}

func (c *KafkaConsumer) worker(ctx context.Context, _ int) {
	//defer recover()
	for {
		// read msg
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() == nil { // if context alive
				c.logger.Error("worker read error", zap.Error(err))
			}
			return
		}

		var event enricher.ClickTask
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.logger.Error("fail to unmarshall worker task", zap.Error(err))
			continue
		}

		enriched, err := c.enricher.Enrich(ctx, &event)
		if err != nil {
			c.logger.Error("enrich fail", zap.Error(err), zap.String("alias", event.Alias))
			continue
		}

		if err := c.repo.Save(ctx, enriched); err != nil {
			c.logger.Error("fail to save enrich click", zap.Error(err))
			continue
		}
		metrics.EventsProcessed.Inc()
		c.logger.Info("enrich click saved", zap.String("alias", event.Alias))

	}
}

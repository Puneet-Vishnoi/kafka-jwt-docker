package kafka

import (
	"context"
	"log"
	"time"

	"github.com/Puneet-Vishnoi/kafka-simple/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader      *kafka.Reader
	maxRetry    int
	retryBackoff time.Duration
}

func NewKafkaConsumer(cfg config.AppConfig) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{cfg.Brokers},
			GroupID:        cfg.GroupID,
			Topic:          cfg.TopicName,
			MinBytes:       1,
			MaxBytes:       10e6,
			CommitInterval: 0,
	})

	return &KafkaConsumer{
			reader: reader,
			maxRetry: cfg.MaxRetry,
			retryBackoff: cfg.RetryInterval,
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}

func (kc *KafkaConsumer) Consume(ctx context.Context, handler func(kafka.Message) error, dlq *kafka.Writer) {
	for {
			m, err := kc.reader.FetchMessage(ctx)
			if err != nil {
					if ctx.Err() != nil {
							log.Println("Consumer context canceled, exiting")
							return
					}
					log.Printf("Fetch error: %v", err)
					continue
			}

			success := false
			for attempt := 1; attempt <= kc.maxRetry; attempt++ {
					if err := handler(m); err != nil {
							log.Printf("Attempt %d: Failed to process offset=%d: %v", attempt, m.Offset, err)
							time.Sleep(kc.retryBackoff)
							continue
					}
					success = true
					break
			}

			if success {
					if err := kc.reader.CommitMessages(ctx, m); err != nil {
							log.Printf("Commit failed for offset=%d: %v", m.Offset, err)
					}
			} else {
					SendToDLQ(ctx, dlq, m, err, kc.maxRetry)
			}
	}
}
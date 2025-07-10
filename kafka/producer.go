package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Puneet-Vishnoi/kafka-simple/config"
	kafka "github.com/segmentio/kafka-go"
)

// KafkaProducer encapsulates the Kafka writer and related configs.
type KafkaProducer struct {
	writer *kafka.Writer
	maxRetry     int
	retryBackoff time.Duration
}

// NewKafkaProducer initializes a new Kafka producer with the provided brokers and topic.
func NewKafkaProducer(cfg config.AppConfig) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          []string{cfg.Brokers},
		Topic:            cfg.TopicName,
		RequiredAcks:     int(kafka.RequireAll),
		Balancer:         &kafka.RoundRobin{},
		BatchSize:        100,
		BatchTimeout:     1 * time.Second,
		Async:            true,
		QueueCapacity:    1000,
		CompressionCodec: kafka.Snappy.Codec(),
	})

	return &KafkaProducer{
		writer:       writer,
		maxRetry:     cfg.MaxRetry,
		retryBackoff: cfg.RetryInterval,
}
}

// Close gracefully shuts down the Kafka producer.
func (kp *KafkaProducer) Close() error {
	if kp.writer != nil {
		return kp.writer.Close()
	}
	return nil
}

// Produce sends a message to the Kafka topic with retries.
func (kp *KafkaProducer) Produce(ctx context.Context, value []byte) error {
	var err error
	for attempt := 1; attempt <= kp.maxRetry; attempt++ {
		err = kp.writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", time.Now().UnixNano())),
			Value: value,
		})
		if err == nil {
			log.Printf("Produced: %s", string(value))
			return nil
		}
		log.Printf("Attempt %d: Producer error: %v", attempt, err)
		time.Sleep(kp.retryBackoff)
	}
	return fmt.Errorf("failed to produce message after %d attempts: %w", kp.maxRetry, err)
}



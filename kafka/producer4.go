// mannual retry

package kafka

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Puneet-Vishnoi/kafka-simple/config"
	kafka "github.com/segmentio/kafka-go"
)

// KafkaProducer with DLQ support
type KafkaProducer struct {
	writer       *kafka.Writer
	dlqWriter    *kafka.Writer
	maxRetry     int
	retryBackoff time.Duration
}

// NewKafkaProducer initializes the producer and DLQ writer
func NewKafkaProducer(cfg config.AppConfig) *KafkaProducer {
	mainWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers),
		Topic:                  cfg.TopicName,
		RequiredAcks:           kafka.RequireAll,
		Balancer:               &kafka.RoundRobin{},
		Compression:            kafka.Snappy,
		BatchSize:              100,
		BatchTimeout:           1 * time.Second,
		AllowAutoTopicCreation: true,
	}

	dlqWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers),
		Topic:                  cfg.DLQTopic,
		RequiredAcks:           kafka.RequireAll,
		Balancer:               &kafka.RoundRobin{},
		AllowAutoTopicCreation: true,
		Compression:            kafka.Snappy,
	}

	return &KafkaProducer{
		writer:       mainWriter,
		dlqWriter:    dlqWriter,
		maxRetry:     cfg.MaxRetry,
		retryBackoff: cfg.RetryInterval,
	}
}

// Close both main and DLQ writers
func (kp *KafkaProducer) Close() error {
	if kp.writer != nil {
		_ = kp.writer.Close()
	}
	if kp.dlqWriter != nil {
		_ = kp.dlqWriter.Close()
	}
	return nil
}

// Produce sends a message and uses DLQ if retries fail
func (kp *KafkaProducer) Produce(ctx context.Context, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("key-%d", time.Now().UnixNano())),
		Value: value,
	}

	err := writeWithRetry(ctx, kp.writer, msg, kp.maxRetry, kp.retryBackoff)
	if err != nil {
		log.Printf("All retries failed. Sending to DLQ. Reason: %v", err)
		SendToDLQ(ctx, kp.dlqWriter, msg, err, kp.maxRetry)
	}
	return err
}

// writeWithRetry with exponential backoff
func writeWithRetry(ctx context.Context, writer *kafka.Writer, msg kafka.Message, maxRetries int, backoff time.Duration) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = writer.WriteMessages(ctx, msg)
		if err == nil {
			log.Printf("Produced on attempt %d", attempt)
			return nil
		}
		log.Printf("Attempt %d failed: %v", attempt, err)

		select {
		case <-time.After(backoff + jitter(100*time.Millisecond)):
			backoff *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("produce failed after %d attempts: %w", maxRetries, err)
}

// jitter adds randomness to backoff
func jitter(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}

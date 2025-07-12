// Here's your Kafka producer rewritten for high-throughput async mode using a Completion callback.

// This setup is ideal for:

// High message volume (e.g., telemetry, logs, events).

// Applications that want to fire-and-forget but still observe success/failure.

// Environments where latency matters more than guaranteed delivery.

package kafka

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/Puneet-Vishnoi/kafka-simple/config"
// 	kafka "github.com/segmentio/kafka-go"
// )

// // KafkaProducer encapsulates the Kafka writer and related configs.
// type KafkaProducer struct {
// 	writer *kafka.Writer
// }

// // NewKafkaProducer initializes a new async Kafka producer with completion tracking.
// func NewKafkaProducer(cfg config.AppConfig) *KafkaProducer {
// 	writer := &kafka.Writer{
// 		Addr:              kafka.TCP(cfg.Brokers),
// 		Topic:             cfg.TopicName,
// 		RequiredAcks:      kafka.RequireAll,
// 		Balancer:          &kafka.RoundRobin{},
// 		Compression:       kafka.Snappy,
// 		BatchSize:         100,
// 		BatchTimeout:      1 * time.Second,
// 		MaxAttempts:       cfg.MaxRetry,
// 		WriteBackoffMin:   cfg.RetryInterval,
// 		WriteBackoffMax:   2 * cfg.RetryInterval,
// 		WriteTimeout:      10 * time.Second,
// 		Async:             true, // async mode
// 		Completion: func(messages []kafka.Message, err error) {
// 			for _, msg := range messages {
// 				if err != nil {
// 					log.Printf("❌ Failed to send message: %s, err: %v", string(msg.Value), err)
// 				} else {
// 					log.Printf("✅ Successfully sent: %s", string(msg.Value))
// 				}
// 			}
// 		},
// 		AllowAutoTopicCreation: true,
// 	}

// 	return &KafkaProducer{
// 		writer: writer,
// 	}
// }

// // Close gracefully shuts down the Kafka producer.
// func (kp *KafkaProducer) Close() error {
// 	if kp.writer != nil {
// 		return kp.writer.Close()
// 	}
// 	return nil
// }

// // Produce sends a message to Kafka using async mode with completion callback.
// // It returns only validation or queuing errors (rare).
// func (kp *KafkaProducer) Produce(ctx context.Context, value []byte) error {
// 	err := kp.writer.WriteMessages(ctx, kafka.Message{
// 		Key:   []byte(fmt.Sprintf("key-%d", time.Now().UnixNano())),
// 		Value: value,
// 	})

// 	if err != nil {
// 		log.Printf("❗ WriteMessages returned an error before dispatch: %v", err)
// 		return fmt.Errorf("early producer error: %w", err)
// 	}

// 	// No need to log success here — completion callback handles it
// 	return nil
// }

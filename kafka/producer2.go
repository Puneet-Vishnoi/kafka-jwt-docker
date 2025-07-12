//Production-Grade Version (Properly Using Built-in Retries)
// You should do either:

// ✅ Use synchronous mode with built-in retries (recommended for safety), or

// ✅ Use asynchronous mode with a Completion handler (for high-throughput).

// 🔧 Here’s a clean production-grade version using sync mode and built-in retry logic:

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

// // NewKafkaProducer initializes a new Kafka producer with production-grade settings.
// func NewKafkaProducer(cfg config.AppConfig) *KafkaProducer {
// 	writer := &kafka.Writer{
// 		Addr:              kafka.TCP(cfg.Brokers),
// 		Topic:             cfg.TopicName,
// 		RequiredAcks:      kafka.RequireAll,
// 		Balancer:          &kafka.RoundRobin{},
// 		Compression:       kafka.Snappy,
// 		BatchSize:         100,
// 		BatchTimeout:      1 * time.Second,
// 		MaxAttempts:       cfg.MaxRetry,                // built-in retry
// 		WriteBackoffMin:   cfg.RetryInterval,           // backoff between retries
// 		WriteBackoffMax:   2 * cfg.RetryInterval,
// 		WriteTimeout:      10 * time.Second,
// 		AllowAutoTopicCreation: true,
// 		Async:             false,                       // <-- sync mode: ensures retry is meaningful
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

// // Produce sends a message to Kafka using built-in retry logic.
// func (kp *KafkaProducer) Produce(ctx context.Context, value []byte) error {
// 	err := kp.writer.WriteMessages(ctx, kafka.Message{
// 		Key:   []byte(fmt.Sprintf("key-%d", time.Now().UnixNano())),
// 		Value: value,
// 	})
// 	if err != nil {
// 		log.Printf("❌ Failed to produce: %v", err)
// 		return fmt.Errorf("produce failed: %w", err)
// 	}

// 	log.Printf("✅ Produced: %s", string(value))
// 	return nil
// }

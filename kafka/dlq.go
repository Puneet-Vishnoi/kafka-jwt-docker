package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/Puneet-Vishnoi/kafka-simple/config"
	"github.com/segmentio/kafka-go"
)


// // EnsureTopicExists checks if a Kafka topic exists; if not, it creates it.
// func EnsureTopicExists(brokers []string, topic string) error {
// 	conn, err := kafka.Dial("tcp", brokers[0])
// 	if err != nil {
// 		return fmt.Errorf("failed to dial kafka broker: %w", err)
// 	}
// 	defer conn.Close()

// 	controller, err := conn.Controller()
// 	if err != nil {
// 		return fmt.Errorf("failed to get controller: %w", err)
// 	}

// 	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
// 	if err != nil {
// 		return fmt.Errorf("failed to dial controller: %w", err)
// 	}
// 	defer controllerConn.Close()

// 	partitions, err := controllerConn.ReadPartitions(topic)
// 	if err == nil && len(partitions) > 0 {
// 		// Topic exists
// 		log.Printf("Topic %q already exists", topic)
// 		return nil
// 	}

// 	// Create topic if it does not exist
// 	topicConfig := kafka.TopicConfig{
// 		Topic:             topic,
// 		NumPartitions:     1,
// 		ReplicationFactor: 1,
// 	}

// 	if err := controllerConn.CreateTopics(topicConfig); err != nil {
// 		return fmt.Errorf("failed to create topic %s: %w", topic, err)
// 	}

// 	log.Printf("Created topic %q", topic)
// 	return nil
// }

func InitDLQWriter(cfg config.AppConfig) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers),
		Topic:    cfg.DLQTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

func SendToDLQ(ctx context.Context, writer *kafka.Writer, original kafka.Message, processingErr error, retries int) {
	dlqMsg := kafka.Message{
		Key:   original.Key,
		Value: original.Value,
		Headers: []kafka.Header{
			{Key: "original_topic", Value: []byte(original.Topic)},
			{Key: "original_partition", Value: []byte(fmt.Sprint(original.Partition))},
			{Key: "original_offset", Value: []byte(fmt.Sprint(original.Offset))},
			{Key: "retries", Value: []byte(fmt.Sprint(retries))},
			{Key: "error", Value: []byte(processingErr.Error())},
		},
	}
	if err := writer.WriteMessages(ctx, dlqMsg); err != nil {
		log.Printf("Failed to write to DLQ: %v", err)
	} else {
		log.Printf("Sent to DLQ: offset=%d reason=%v", original.Offset, processingErr)
	}
}

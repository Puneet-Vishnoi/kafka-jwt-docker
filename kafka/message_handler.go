package kafka

import (
	"log"

	"github.com/segmentio/kafka-go"
)

func ProcessMessage(m kafka.Message) error {
	log.Printf("Processing: topic=%s partition=%d offset=%d key=%s value=%s",
		m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

	// Simulate processing (can add failure simulation here)
	return nil
}

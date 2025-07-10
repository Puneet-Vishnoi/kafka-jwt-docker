package config

import (
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	TopicName        string
	DLQTopic         string
	Brokers          string
	GroupID          string
	MaxRetry         int
	RetryInterval    time.Duration
}

func LoadConfig() AppConfig {
	maxRetry, _ := strconv.Atoi(os.Getenv("KAFKA_MAX_RETRY"))
	retryInterval, _ := time.ParseDuration(os.Getenv("KAFKA_RETRY_INTERVAL"))

	return AppConfig{
		TopicName:     os.Getenv("KAFKA_TOPIC"),
		DLQTopic:      os.Getenv("KAFKA_DLQ_TOPIC"),
		Brokers:       os.Getenv("KAFKA_BROKERS"),
		GroupID:       os.Getenv("KAFKA_GROUP_ID"),
		MaxRetry:      maxRetry,
		RetryInterval: retryInterval,
	}
}

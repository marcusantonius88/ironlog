package infra

import (
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer wraps kafka.Writer
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireAll),
	})
	return &KafkaProducer{writer: writer}
}

// PublishEvent publishes an event to Kafka
func (kp *KafkaProducer) PublishEvent(key string, event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: payload,
	}

	err = kp.writer.WriteMessages(nil, message)
	if err != nil {
		log.Printf("failed to publish event: %v", err)
		return err
	}

	return nil
}

// Close closes the Kafka writer
func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}

// KafkaConsumer wraps kafka.Reader
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:                brokers,
		Topic:                  topic,
		GroupID:                groupID,
		StartOffset:            kafka.LastOffset,
		CommitInterval:         0,
		PartitionWatchInterval: 0,
	})
	return &KafkaConsumer{reader: reader}
}

// Close closes the Kafka reader
func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}

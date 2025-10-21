package kafka

import (
	"context"
	"fmt"
	"time"

	"yourapp/pkg/config"
	"yourapp/pkg/logger"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	producer *ckafka.Producer
	consumer *ckafka.Consumer
)

// Init initializes the Kafka connection
func Init(ctx context.Context, cfg config.KafkaConfig) error {
	// Initialize producer
	if err := initProducer(cfg); err != nil {
		return fmt.Errorf("failed to initialize Kafka producer: %w", err)
	}

	// Initialize consumer
	if err := initConsumer(cfg); err != nil {
		return fmt.Errorf("failed to initialize Kafka consumer: %w", err)
	}

	return nil
}

// initProducer initializes the Kafka producer
func initProducer(cfg config.KafkaConfig) error {
	conf := &ckafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	// Optional security/SASL configs
	if cfg.SecurityProtocol != "" {
		_ = conf.SetKey("security.protocol", cfg.SecurityProtocol)
	}
	if cfg.SASLMechanism != "" {
		_ = conf.SetKey("sasl.mechanism", cfg.SASLMechanism)
	}
	if cfg.Username != "" {
		_ = conf.SetKey("sasl.username", cfg.Username)
	}
	if cfg.Password != "" {
		_ = conf.SetKey("sasl.password", cfg.Password)
	}
	if cfg.SessionTimeout > 0 {
		_ = conf.SetKey("session.timeout.ms", int(cfg.SessionTimeout/time.Millisecond))
	}
	if cfg.HeartbeatInterval > 0 {
		_ = conf.SetKey("heartbeat.interval.ms", int(cfg.HeartbeatInterval/time.Millisecond))
	}

	var err error
	producer, err = ckafka.NewProducer(conf)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return nil
}

// initConsumer initializes the Kafka consumer
func initConsumer(cfg config.KafkaConfig) error {
	conf := &ckafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		// No group.id to allow manual partition assignment (Assign)
	}

	// Optional security/SASL configs
	if cfg.SecurityProtocol != "" {
		_ = conf.SetKey("security.protocol", cfg.SecurityProtocol)
	}
	if cfg.SASLMechanism != "" {
		_ = conf.SetKey("sasl.mechanism", cfg.SASLMechanism)
	}
	if cfg.Username != "" {
		_ = conf.SetKey("sasl.username", cfg.Username)
	}
	if cfg.Password != "" {
		_ = conf.SetKey("sasl.password", cfg.Password)
	}
	if cfg.SessionTimeout > 0 {
		_ = conf.SetKey("session.timeout.ms", int(cfg.SessionTimeout/time.Millisecond))
	}
	if cfg.HeartbeatInterval > 0 {
		_ = conf.SetKey("heartbeat.interval.ms", int(cfg.HeartbeatInterval/time.Millisecond))
	}

	var err error
	consumer, err = ckafka.NewConsumer(conf)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return nil
}

// GetProducer returns the Kafka producer
func GetProducer() *ckafka.Producer {
	return producer
}

// GetConsumer returns the Kafka consumer
func GetConsumer() *ckafka.Consumer {
	return consumer
}

// Close closes the Kafka connections
func Close() error {
	var err error

	if producer != nil {
		producer.Close()
	}

	if consumer != nil {
		if closeErr := consumer.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("failed to close consumer: %w, previous error: %v", closeErr, err)
			} else {
				err = fmt.Errorf("failed to close consumer: %w", closeErr)
			}
		}
	}

	return err
}

// PublishMessage publishes a message to a topic
func PublishMessage(ctx context.Context, topic, key string, message []byte) error {
	if producer == nil {
		return fmt.Errorf("Kafka producer not initialized")
	}

	deliveryChan := make(chan ckafka.Event, 1)
	defer close(deliveryChan)

	msg := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Key:            []byte(key),
		Value:          message,
	}

	if err := producer.Produce(msg, deliveryChan); err != nil {
		return fmt.Errorf("failed to produce message to topic %s: %w", topic, err)
	}

	select {
	case e := <-deliveryChan:
		m := e.(*ckafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
		}
		logger.Infof("Message sent to topic %s, partition %d, offset %v", topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// ConsumeMessages consumes messages from a topic
func ConsumeMessages(ctx context.Context, topic string, handler func(*ckafka.Message) error) error {
	if consumer == nil {
		return fmt.Errorf("Kafka consumer not initialized")
	}

	// Manually assign to partition 0, starting from latest (similar to previous behavior)
	if err := consumer.Assign([]ckafka.TopicPartition{{Topic: &topic, Partition: 0, Offset: ckafka.OffsetEnd}}); err != nil {
		return fmt.Errorf("failed to assign consumer to topic %s: %w", topic, err)
	}
	defer func() { _ = consumer.Unassign() }()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			e := consumer.Poll(100)
			if e == nil {
				continue
			}
			switch ev := e.(type) {
			case *ckafka.Message:
				if err := handler(ev); err != nil {
					logger.Errorf("Error processing message: %v", err)
				}
			case ckafka.Error:
				logger.Errorf("Consumer error: %v", ev)
			default:
				// ignore other events (stats, etc.)
			}
		}
	}
}

// Health checks the health of the Kafka connection
func Health(ctx context.Context) error {
	if producer == nil && consumer == nil {
		return fmt.Errorf("Kafka client not initialized")
	}

	// Try to get metadata via available client
	if producer != nil {
		_, err := producer.GetMetadata(nil, false, int((5*time.Second)/time.Millisecond))
		return err
	}
	if consumer != nil {
		_, err := consumer.GetMetadata(nil, false, int((5*time.Second)/time.Millisecond))
		return err
	}

	return fmt.Errorf("Kafka client not initialized")
}

// CreateTopic creates a new topic
func CreateTopic(ctx context.Context, topicName string, numPartitions int32, replicationFactor int16) error {
	// This would typically be done through Kafka admin API
	// For now, we'll just log that the topic should be created
	fmt.Printf("Topic %s should be created with %d partitions and replication factor %d\n",
		topicName, numPartitions, replicationFactor)
	return nil
}

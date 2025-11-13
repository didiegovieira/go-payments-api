package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, key string, message interface{}) error
	Close() error
}

type publisher struct {
	writer  *kafka.Writer
	brokers []string
}

func NewPublisher(brokers []string) Publisher {
	log.Printf("🔧 Initializing Kafka publisher with brokers: %v", brokers)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	}

	pub := &publisher{
		writer:  writer,
		brokers: brokers,
	}

	if err := pub.createTopicIfNotExists("payment.events"); err != nil {
		log.Printf("⚠️  Warning: Failed to create topic: %v", err)
	}

	return pub
}

func (p *publisher) createTopicIfNotExists(topic string) error {
	conn, err := kafka.Dial("tcp", p.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial leader: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to dial controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		// Ignorar erro se o tópico já existe
		if err.Error() == "Topic with this name already exists" {
			log.Printf("✅ Topic %s already exists", topic)
			return nil
		}
		return fmt.Errorf("failed to create topic: %w", err)
	}

	log.Printf("✅ Topic %s created successfully", topic)
	return nil
}

func (p *publisher) Publish(ctx context.Context, topic string, key string, message interface{}) error {
	log.Printf("📨 Attempting to publish message - Topic: %s, Key: %s", topic, key)

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Failed to marshal message: %v", err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	log.Printf("📦 Message marshaled - Size: %d bytes", len(data))

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("❌ Failed to write message to Kafka: %v", err)
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("✅ Message published successfully - Topic: %s, Key: %s", topic, key)
	return nil
}

func (p *publisher) Close() error {
	log.Printf("🔒 Closing Kafka publisher")
	return p.writer.Close()
}

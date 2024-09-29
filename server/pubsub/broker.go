package pubsub

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Конфігурація для підключення до Kafka
var kafkaBrokerAddress = "localhost:9092" // Замініть на адресу вашого Kafka сервера
var kafkaTopic = "experiment-events"

// Producer структура для надсилання повідомлень до Kafka
type KafkaProducer struct {
	writer *kafka.Writer
}

// Створення нового Kafka Producer
func NewKafkaProducer() *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokerAddress),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{writer: writer}
}

// Відправка повідомлення до Kafka
func (p *KafkaProducer) PublishMessage(event string) error {
	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(time.Now().Format(time.RFC3339)),
			Value: []byte(event),
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Printf("Message published: %s", event)
	}
	return err
}

// Закриття Producer
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

// Consumer структура для читання повідомлень з Kafka
type KafkaConsumer struct {
	reader *kafka.Reader
}

// Створення нового Kafka Consumer
func NewKafkaConsumer() *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokerAddress},
		Topic:   kafkaTopic,
		GroupID: "experiment-consumer-group",
	})

	return &KafkaConsumer{reader: reader}
}

// Читання повідомлення з Kafka
func (c *KafkaConsumer) SubscribeMessages() {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}
		log.Printf("Received message: %s", string(msg.Value))
	}
}

// Закриття Consumer
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

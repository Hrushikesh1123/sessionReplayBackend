package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func ProduceMessage(broker, topic string) {
	fmt.Println("Starting Kafka producer...")
	writer := kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	for i := 1; i <= 5; i++ {
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", i)),
			Value: []byte(fmt.Sprintf("Message number %d", i)),
		}

		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			log.Fatalf("Failed to write message: %v", err)
		}
		fmt.Printf("Message sent: %s\n", msg.Value)
		time.Sleep(1 * time.Second)
	}
}

func ConsumeMessages(broker, topic, groupID string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})

	defer reader.Close()

	fmt.Println("Starting Kafka consumer...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}
		fmt.Printf("Received: %s => %s\n", string(msg.Key), string(msg.Value))
	}
}

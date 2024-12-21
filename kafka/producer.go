package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func (kc *KafkaClient) ProduceMessage(topic, key, value string) error {
	kc.mu.Lock()
	writer, ok := kc.producers[topic]
	kc.mu.Unlock()

	if !ok {
		return fmt.Errorf("no producer for topic %s", topic)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	if err := writer.WriteMessages(context.Background(), msg); err != nil {
		log.Printf("[Producer Error] %v", err)
		return err
	}

	log.Printf("[Produced] %s => %s (topic: %s)", key, value, topic)
	return nil
}

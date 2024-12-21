package kafka

import (
	"fmt"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	producers map[string]*kafka.Writer // one producer per topic
	consumers map[string]*kafka.Reader // one consumer per (topic + groupID)
	mu        sync.Mutex
}

type KafkaInitParams struct {
	Brokers []string
	Topics  []string
	Groups  []string
}

// NewClient sets up producers and consumers for each topic/group.
func NewClient(params KafkaInitParams) *KafkaClient {
	c := &KafkaClient{
		producers: make(map[string]*kafka.Writer),
		consumers: make(map[string]*kafka.Reader),
	}

	// Producer per topic
	for _, topic := range params.Topics {
		writer := &kafka.Writer{
			Addr:     kafka.TCP(params.Brokers...), // single or multiple brokers
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
		c.producers[topic] = writer
		log.Printf("[KafkaClient] Producer created for topic: %s", topic)
	}

	// Consumer for each (topic, groupID) pair
	// If you have more topics than groups, or vice versa, handle that logic as you see fit
	for i, topic := range params.Topics {
		groupID := "default-group"
		if i < len(params.Groups) {
			groupID = params.Groups[i]
		}
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: params.Brokers, // single or multiple brokers
			Topic:   topic,
			GroupID: groupID,
		})
		key := fmt.Sprintf("%s|%s", topic, groupID)
		c.consumers[key] = reader
		log.Printf("[KafkaClient] Consumer created for topic: %s, group: %s", topic, groupID)
	}

	return c
}

func (kc *KafkaClient) Close() {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	for topic, writer := range kc.producers {
		if err := writer.Close(); err != nil {
			log.Printf("[Close Error] Producer for %s: %v", topic, err)
		}
	}
	for key, reader := range kc.consumers {
		if err := reader.Close(); err != nil {
			log.Printf("[Close Error] Consumer for %s: %v", key, err)
		}
	}
	log.Println("[KafkaClient] All producers & consumers closed.")
}

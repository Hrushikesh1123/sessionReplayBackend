package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
)

func (kc *KafkaClient) StartAllConsumers(ctx context.Context, wg *sync.WaitGroup) {
	for key := range kc.consumers {
		wg.Add(1)
		go func(rKey string) {
			defer wg.Done()
			log.Printf("[Consumer] Starting for %s", rKey)

			for {
				msg, err := kc.consumers[rKey].ReadMessage(ctx)
				if err != nil {
					select {
					case <-ctx.Done():
						log.Printf("[Consumer] Shutting down for %s", rKey)
						return
					default:
						log.Printf("[Consumer Error] %s: %v", rKey, err)
						return
					}
				}
				fmt.Printf("[Consumed] Key=%s, Value=%s, Consumer=%s\n",
					string(msg.Key),
					string(msg.Value),
					rKey,
				)
			}
		}(key)
	}
}

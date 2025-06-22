package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	topic := "coffee_orders"
	msgCount := 0

	// Connect to Kafka consumer
	consumer, err := ConnectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println("❌ Failed to close consumer:", err)
		}
	}()

	fmt.Println("☕ Coffee Order Consumer Started")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				fmt.Println("❌ Consumer error:", err)
			case msg := <-partitionConsumer.Messages():
				msgCount++
				order := string(msg.Value)
				fmt.Printf("📦 Received order #%d | Topic(%s) | Message(%s)\n", msgCount, msg.Topic, order)
				fmt.Printf("☕ Brewing Coffee for order: %s\n", order)
			case <-sigchan:
				fmt.Println("🛑 Interruption detected, shutting down...")
				doneCh <- struct{}{}
				return
			}
		}
	}()

	<-doneCh
	fmt.Println("✅ Processed", msgCount, "messages")
}

// ConnectConsumer connects to a Kafka cluster and returns a consumer instance
func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}

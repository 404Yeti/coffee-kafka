package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

// Order represents a coffee order
type Order struct {
	CustomerName string `json:"customer_name"`
	CoffeeType   string `json:"coffee_type"`
}

func main() {
	http.HandleFunc("/order", placeOrder)
	log.Println("☕ Order service started on :3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}

// ConnectProducer creates a Kafka sync producer
func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

// PushOrderToQueue publishes the order to Kafka
func PushOrderToQueue(topic string, order []byte) error {
	brokers := []string{"localhost:9092"}

	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(order),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("✅ Order published to topic(%s), partition(%d), offset(%d)\n", topic, partition, offset)
	return nil
}

// placeOrder handles incoming HTTP POST requests to /order
func placeOrder(w http.ResponseWriter, r *http.Request) {
	// ✅ CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// ✅ Handle preflight (OPTIONS request)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Println("❌ Error decoding order:", err)
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		log.Println("❌ Error marshaling order:", err)
		http.Error(w, "Error processing order", http.StatusInternalServerError)
		return
	}

	err = PushOrderToQueue("coffee_orders", orderInBytes)
	if err != nil {
		log.Println("❌ Kafka push error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     fmt.Sprintf("Order for %s placed successfully!", order.CustomerName),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("❌ Response encoding error:", err)
		http.Error(w, "Error placing order", http.StatusInternalServerError)
	}
}


☕ Brewster: The Coffee Commander (Kafka + Go)

This is a fun microservice project that simulates a coffee ordering system using Kafka, Go, and Docker. It's designed to help you learn how event-driven architecture works by connecting a Producer (order form) with a Consumer (coffee maker) through a Kafka message queue.

💡 What You'll Learn

Kafka concepts: topics, producers, consumers

How to write Go services that produce and consume Kafka messages

How to use Docker Compose to spin up Kafka + UI

Handling CORS and JSON APIs in Go

🧵 Tech Stack

Go (Producer & Worker)

Kafka + Zookeeper via Docker

Kafka UI for visualizing messages

Postman / curl or HTML form for testing

📂 Project Structure

coffee/
├── Producer/              # HTTP server that sends orders to Kafka
│   └── main.go
├── Worker/                # Kafka consumer that prints out orders
│   └── main.go
├── brewster.html          # (Optional) Fun browser UI to send orders
├── docker-compose.yaml    # Kafka, Zookeeper, Kafka UI
├── .gitignore
└── README.md

🔧 Setup & Run

1. Start Kafka

docker compose up -d

This runs Kafka, Zookeeper, and Kafka UI on:

Kafka: localhost:9092

Kafka UI: http://localhost:8080

2. Run the Producer (HTTP Server)

cd Producer
go run main.go

This starts an HTTP server at http://localhost:3000.

3. Run the Worker (Consumer)

cd Worker
go run main.go

The worker listens for new messages on the coffee_orders topic and logs them.

📢 Send an Order

Option 1: Postman or curl

curl -X POST http://localhost:3000/order \
  -H "Content-Type: application/json" \
  -d '{"customer_name": "Robbie", "coffee_type": "Cold Brew"}'

Option 2: Open the HTML page (Optional)

You can also test it via the fun browser UI:

# From root of project
python3 -m http.server 8081

Then open:

http://localhost:8081/brewster.html

📊 Kafka UI

Visit:

http://localhost:8080

View topics, partitions, consumers

Inspect the coffee_orders topic and messages

🔍 Example Output

From Worker Terminal

️ Coffee Order Consumer Started
📦 Received order #1 | Topic(coffee_orders) | Message({"customer_name":"Robbie","coffee_type":"Cold Brew"})
🥜 Brewing Coffee for order: Robbie

🚩 Dev Notes

CORS enabled in Producer for HTML testing

Messages are not persisted after restart (in-memory only)

docker compose down -v will reset Kafka topics

🎉 Future Ideas

Store orders in a DB (e.g., Postgres)

Add a delivery service as a second consumer

Deploy to cloud with Docker Compose or K8s

🌐 Author

Made with Go, Kafka, and Coffee by 404Yeti ☕

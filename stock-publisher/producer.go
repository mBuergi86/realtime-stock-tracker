package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomPrice(r *rand.Rand) float64 {
	return r.Float64()*500 + 50
}

type StockEvent struct {
	Company   string  `json:"company"`
	EventType string  `json:"eventType"`
	Price     float64 `json:"price"`
}

func stockPublisher(url, queueName, stock string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventTypes := []string{"buy", "sell"}
	src := rand.NewSource(time.Now().UnixNano())
	localRand := rand.New(src)

	tickerIntervaleValue := getEnvWithDefault("TICKER_INTERVAL", "1000")
	tickerInterval, err := strconv.Atoi(tickerIntervaleValue)
	failOnError(err, "Failed to parse ticker interval")

	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Millisecond)
	for range ticker.C {
		eventType := eventTypes[localRand.Intn(len(eventTypes))]
		price := randomPrice(localRand)
		stockEvent := StockEvent{
			Company:   stock,
			EventType: eventType,
			Price:     price,
		}

		jsonBody, err := json.Marshal(stockEvent)
		failOnError(err, "Failed to marshal event")

		err = channel.PublishWithContext(
			ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        jsonBody,
			})
		if err != nil {
			log.Panicf("Failed to publish a message: %s", err)
		}

		log.Printf(" [x]Sent: %s", jsonBody)

	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	rabbitMQConnectionURL := getEnvWithDefault("RABBITMQ_CONNECTION_URL", "amqp://stockmarket:supersecret123@localhost:5672/")

	stocks := []string{"MSFT", "TSLA", "AAPL"}

	for _, stock := range stocks {
		stockPublisher(rabbitMQConnectionURL, "Stock Publisher", stock)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

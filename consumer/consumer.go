package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func RoundToTwoDigits(num float64) float64 {
	return math.Round(num*100) / 100
}

func main() {
	rabbitMQConnectionURL := getEnvWithDefault("RABBITMQ_CONNECTION_URL", "amqp://stockmarket:supersecret123@localhost:5672/")
	conn, err := amqp.Dial(rabbitMQConnectionURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Stock Publisher", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Parse the received JSON into a struct
			var stock map[string]interface{}

			err := json.Unmarshal(d.Body, &stock)
			if err != nil {
				failOnError(err, "Failed to parse JSON")
				continue
			}

			companyName := stock["company"].(string)
			price := RoundToTwoDigits(stock["price"].(float64))

			mongo_uri := getEnvWithDefault("MONGO_URI", "mongodb://localhost:27017")

			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongo_uri))

			if err != nil {
				failOnError(err, "Failed to connect to MongoDB")
				continue
			}

			_, err = client.Database("stockmarket").Collection("stocks").InsertOne(context.Background(), bson.M{
				"company":  companyName,
				"avgPrice": price,
			})

			if err != nil {
				failOnError(err, "Failed to insert document")
				continue
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

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
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type StockEvent struct {
	Company   string  `json:"company"`
	EventType string  `json:"eventType"`
	Price     float64 `json:"price"`
}

// StockConsumer consumes stock events from RabbitMQ and writes them to MongoDB
func StockConsumer(url, queueName string, mongoClient *mongo.Client) {

	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
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
			// Process the message
			log.Printf("Received a message: %s", d.Body)
			var event StockEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Error unmarshalling JSON: %v", err)
				continue
			}

			companyName := event.Company
			price := RoundToTwoDigits(event.Price)

			// Process the event
			if event.EventType == "buy" {
				log.Printf("Processing buy event: %v", event)
			} else if event.EventType == "sell" {
				log.Printf("Processing sell event: %v", event)
			} else {
				log.Printf("Unknown event type: %s", event.EventType)
			}

			_, err = mongoClient.Database("stockmarket").Collection("stocks").InsertOne(context.Background(), bson.M{
				"company":  companyName,
				"avgPrice": price,
			})

			if err != nil {
				failOnError(err, "Failed to insert document")
				continue
			}
		}
	}()
	<-forever
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
	rabbitMQConnectionURL := getEnvWithDefault("RABBITMQ_CONNECTION_URL", "amqp://stockmarket:supersecret123@127.0.0.1:5672/")
	mongoURI := getEnvWithDefault("MONGO_URI", "mongodb://host.docker.internal:27017,host.docker.internal:27018,host.docker.internal:27019/?replicaSet=rs0")
	// cretentials := options.Credential{
	// 	Username: getEnvWithDefault("MONGO_USERNAME", "stockmarket"),
	// 	Password: getEnvWithDefault("MONGO_PASSWORD", "supersecret123"),
	// }
	queueName := "TSLA" // This should match the producer's queue name ☝️

	wc := writeconcern.New(writeconcern.W(1))

	clientOptions := options.Client().ApplyURI(mongoURI).SetWriteConcern(wc)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	failOnError(err, "Failed to connect to MongoDB")

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			failOnError(err, "Failed to disconnect from MongoDB")
		}
	}()

	go StockConsumer(rabbitMQConnectionURL, queueName, mongoClient)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

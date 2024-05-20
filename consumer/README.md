# Stock Consumer

## Overview

The Stock Consumer is a Go application that connects to RabbitMQ to receive stock price updates and stores them in a MongoDB database. This application is a crucial part of the real-time stock tracker system, ensuring that stock data is collected and persisted for further analysis and display.

## Prerequisites

- Docker
- Docker Compose
- Go (for local development and building)

## Setup Instructions

### 1. Clone the Repository

```sh
git clone https://github.com/mBuergi86/realtime-stock-tracker.git
cd realtime-stock-tracker/consumer
```

### 2. Environment Variables

Set the following environment variables in a `.env` file or your shell:

```sh
RABBITMQ_CONNECTION_URL=amqp://youruser:yourpassword@localhost:5672/
MONGODB_URI=mongodb://youruser:yourpassword@localhost:27017/
MONGODB_DATABASE=stockmarket
MONGODB_COLLECTION=stocks
```

### 3. Build and Run the Consumer

Using Docker Compose:

```sh
docker-compose up -d
```

### 4. Local Development

To run the consumer locally without Docker, make sure you have Go installed. Then, use the following commands:

```sh
go mod download
go run main.go
```

## Code Explanation

The main parts of the consumer code are as follows:

### Connecting to RabbitMQ

The consumer establishes a connection to RabbitMQ using the provided connection URL.

```go
conn, err := amqp.Dial("RABBITMQ_CONNECTION_URL")
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
```

### Consuming Messages

The consumer listens for messages on the specified queue and processes each message by parsing the JSON payload.

```go
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

go func() {
    for d := range msgs {
        log.Printf("Received a message: %s", d.Body)

        var stock map[string]interface{}
        err := json.Unmarshal(d.Body, &stock)
        if err != nil {
            failOnError(err, "Failed to parse JSON")
            continue
        }

        companyName := stock["company"].(string)
        price := math.Round(stock["price"].(float64)*100) / 100

        client := db.GetMongoClient()

        _, err = client.Collection.InsertOne(context.Background(), bson.M{
            "company": companyName,
            "price":   price,
        })

        if err != nil {
            failOnError(err, "Failed to insert document")
            continue
        }
    }
}()
```

### Connecting to MongoDB

The consumer connects to MongoDB and inserts the parsed stock data into the specified collection.

```go
client := db.GetMongoClient()
_, err = client.Collection.InsertOne(context.Background(), bson.M{
    "company": companyName,
    "price":   price,
})
```

## Error Handling

Errors are handled using the `failOnError` function, which logs the error message and terminates the application if an error occurs.

```go
func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}
```

## License

This project is licensed under the MIT License.

## Contributions

Contributions are welcome! Please open an issue or submit a pull request.

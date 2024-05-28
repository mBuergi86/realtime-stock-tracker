# First stage: Build the Go binary
FROM golang:1.22.3-alpine AS builder

# Set environment variables
ENV GO111MODULE=on

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o aapl_consumer aapl_consumer.go
RUN go build -o msft_consumer msft_consumer.go
RUN go build -o tsla_consumer tsla_consumer.go

# Second stage: Create a minimal image for running the Go binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/aapl_consumer .
COPY --from=builder /app/msft_consumer .
COPY --from=builder /app/tsla_consumer .

# Set environment variables (common for all)
ENV RABBITMQ_CONNECTION_URL=amqp://stockmarket:supersecret123@rabbitmq:5672/
ENV MONGODB_URI=mongodb://mongodb:27017,mongodb:27018,mongodb:27019/?replicaSet=rs0
# ENV MONGODB_USER=stockmarket
# ENV MONGODB_PASSWORD=supersecret123
ENV MONGODB_DATABASE=stockmarket
ENV MONGODB_COLLECTION=stocks

# Command to run the binary
CMD ["./consumer"]

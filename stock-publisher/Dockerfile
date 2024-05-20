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
RUN go build -o producer main.go

# Second stage: Create a minimal image for running the Go binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/producer .

# Set environment variables
ENV RABBITMQ_CONNECTION_URL=amqp://stockmarket:supersecret123@rabbitmq:5672/

# Command to run the binary
CMD ["./producer"]

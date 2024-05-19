# Use the official Golang image as the base
FROM golang:1.22.3 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy all files from the current directory to the working directory in the container
COPY . .

# Download Go module dependencies
RUN go mod download

# Build the Go application with optimization flags
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /producer 

# Set the entry point to the built binary

# Use a minimal base image for the final stage
CMD ["/producer"]

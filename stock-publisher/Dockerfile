# # Use the official Golang image as the build stage
# FROM golang:alpine as builder

# # Set environment variables for the build
# ENV APP_NAME="rmq-producer" CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# # Set the working directory inside the container
# WORKDIR /build

# # Copy all files from the current directory to the working directory in the container
# COPY . .

# # Download Go module dependencies
# RUN go mod download

# # Build the Go application with optimization flags
# # -ldflags="-s -w": strip debug information to reduce binary size
# RUN go build -ldflags="-s -w" -o ${APP_NAME} ./main.go

# # Use a minimal base image for the final stage
# FROM scratch

# # Copy the built binary from the builder stage to the final image
# COPY --from=builder /build/${APP_NAME} /${APP_NAME}

# # Set the entry point to the built binary
# ENTRYPOINT ["/rmq-producer"]
# 

# Use the official Golang image as the base
FROM golang:1.22.3 as builder

# Label maintainers and contributors for the image and project
LABEL maintainers="Markus"
LABEL contributors="Markus"

# Set the working directory inside the container
WORKDIR /app

# Copy all files from the current directory to the working directory in the container
COPY . .

# Download Go module dependencies
RUN go mod download

# Build the Go application
RUN go build -o . 

# Use a minimal base image for the final stage
FROM ubuntu:latest

# Copy the built binary from the builder stage to the final image
COPY --from=builder /app/main /app/main

# Set the entry point to the built binary
ENTRYPOINT ["/app/main"]

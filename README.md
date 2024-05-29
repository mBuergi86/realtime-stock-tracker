# Real-Time Stock Tracker

## Overview

This project is a real-time stock tracking system that integrates multiple services to provide live updates on stock prices. The system is built using Go for the backend services, RabbitMQ for message brokering, and MongoDB for data storage. Additionally, there is a web-based frontend and backend to display live stock prices.

## Components

1. **Stock Publisher**
2. **Stock Consumer**
3. **RabbitMQ**
4. **MongoDB**
5. **Stock Live View (Frontend & Backend)**
6. **NGINX**

### Stock Publisher

The Stock Publisher service is responsible for generating stock price events and publishing them to a RabbitMQ queue.

#### Key Features:

- Generates random stock prices for a predefined set of companies.
- Publishes stock price updates to a RabbitMQ queue.

### Stock Consumer

The Stock Consumer service listens to the RabbitMQ queue, processes incoming stock price events, and stores them in MongoDB.

#### Key Features:

- Consumes stock price events from RabbitMQ.
- Stores stock price data in MongoDB.

### RabbitMQ

RabbitMQ is used as the message broker to handle communication between the stock publisher and consumer services.

### MongoDB

MongoDB is used to store stock price data consumed by the stock consumer service.

### Stock Live View (Frontend & Backend)

The frontend service displays live stock prices in a web interface.

### NGINX

NGINX is configured as a load balancer to distribute traffic between multiple frontend instances.

#### Key Features:

- Displays real-time stock prices.
- Uses Tailwind CSS for styling.

## Getting Started

### Prerequisites

- Docker
- Docker Compose
- Go

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mBuergi86/realtime-stock-tracker.git
   cd realtime-stock-tracker
   ```

2. Set up environment variables:

   ```bash
   cp .env.example .env
   ```

3. Build and start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

### Running the Services

- The stock publisher and consumer services will start automatically with Docker Compose.
- Access the RabbitMQ management interface at `http://localhost:15672`.
- Access the MongoDB instance at `mongodb://host.docker.internal:27017`.
- The frontend can be accessed at `http://localhost`.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for review.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

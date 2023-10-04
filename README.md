Certainly! Here's a sample README file for your "Order-Book" project:

---

# Order-Book Project

The Order-Book project is a server application designed to manage and serve order book data. This README provides an overview of the project, its key components, and how to get started.

## Table of Contents

- [Introduction](#introduction)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Usage](#usage)

## Introduction

The Order-Book project is built on Go (Golang) and is designed to provide an efficient and reliable system for managing order book data. It serves as a foundation for tracking and managing buy and sell orders for various financial instruments.

### Key Features

- **Server Application**: The project includes a Go-based server application responsible for handling HTTP requests and managing order book data.

- **Database Integration**: It utilizes a PostgreSQL database to store order book data.

- **Kafka Integration**: The project also integrates with Apache Kafka for real-time data processing. Kafka is used for consuming and processing incoming order data.

## Project Structure

The project's main function, defined in `main.go`, initializes and starts the server, sets up the database connection, and initializes Kafka consumer(s).

```go
func main() {
    fmt.Println("start server ...")

    // Initialize the database
    db.InitDatabase()

    // Initialize Kafka consumers
    go kafka.InitKafka()

    // Create HTTP router and define API endpoints
    r := mux.NewRouter()
    r.HandleFunc("/OrderBook", api.HandleOrderBook).Methods("GET")

    // Start the HTTP server
    log.Fatal(http.ListenAndServe("localhost:8080", r))
}
```

## Prerequisites

Before you start, ensure you have the following prerequisites:

- Go (Golang): You'll need Go installed on your system to build and run the project.

- PostgreSQL: make sure it's installed and properly configured.

- Apache Kafka : ensure Kafka is installed and set up.

## Getting Started

1. **Clone the Repository**: Clone this repository to your local machine:

   ```bash
   git clone https://github.com/your-username/order-book.git
   ```

2. **Install Dependencies**: Navigate to the project directory and install any required Go dependencies:

   ```bash
   cd order-book
   go mod download
   ```

3. **Configure Database and Kafka**:  configure the database connection and Kafka settings in the respective configuration files.

4. **Build the Project**: Build the Go project:

   ```bash
   go build
   ```

## Usage

- **Running the Server**: Start the server by running the compiled binary:

  ```bash
  ./order-book
  ```

- **Accessing APIs**: Once the server is running, you can access the available APIs using HTTP requests. ( server_address/OrderBook?limit={limitNumber}(optional)&symbol={symbolName} )
            

- **Database Interaction**: The server interacts with the configured database to store and retrieve order book data. Ensure your psql database is up and running.

- **Kafka Integration**: Kafka is used for real-time data processing. Make sure Kafka is configured correctly for your use case. The server is subscribing to the Kafka "Register" topic on port 9092 so send your messages on this topic

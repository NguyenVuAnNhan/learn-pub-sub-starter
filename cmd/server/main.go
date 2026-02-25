package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting Peril server...")

	const connectionString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")
	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	<- sigCh

	fmt.Println("Shutting down Peril server...")
}

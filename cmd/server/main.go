package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"os/signal"
	"syscall"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func main() {
	fmt.Println("Starting Peril server...")

	const connectionString = "amqp://guest:guest@127.0.0.1:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	fmt.Println("Channel opened")
	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	state := routing.PlayingState{
		IsPaused: true,
	}
	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, state)
	if err != nil {
		panic(err)
	}
	fmt.Println("Published initial paused state")

	<- sigCh

	fmt.Println("Shutting down Peril server...")
}

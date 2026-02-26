package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
	"fmt"
	"errors"
)

func main() {
	fmt.Println("Starting Peril client...")

	const connectionString = "amqp://guest:guest@127.0.0.1:5672/"
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err)
		return
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", username),
		routing.PauseKey,
		pubsub.TransientQueue,
		handlerPause(gameState),
	)

	if err != nil {
		fmt.Printf("failed to subscribe to pause messages: %v\n", err)
		return
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn": {
			err := gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}
		}
		case "move": {
			command, err := gameState.CommandMove(words)
			fmt.Printf("Command: %+v\n", command)
			if err != nil {
				fmt.Println(err)
			}
		}
		case "status": {
			gameState.CommandStatus()
		}
		case "help":
			gamelogic.PrintClientHelp()
		case "spam": {
			fmt.Println("Spamming not allowed yet!")
		}
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			errors.New("unknown command")
		}
	}
	fmt.Println("Shutting down Peril client...")
}

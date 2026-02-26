package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
	"fmt"
	"errors"
	"log"
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

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
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

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s",routing.ArmyMovesPrefix, username),
		fmt.Sprintf("%s.*", routing.ArmyMovesPrefix),
		pubsub.TransientQueue,
		handlerMove(gameState, publishCh),
	)

	if err != nil {
		fmt.Printf("failed to subscribe to army move messages: %v\n", err)
		return
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		"war",
		fmt.Sprintf("%s.*", routing.WarRecognitionsPrefix),
		pubsub.DurableQueue,
		handlerWar(gameState),
	)

	if err != nil {
		fmt.Printf("failed to subscribe to war messages: %v\n", err)
		return
	}

	// REPL loop
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
			if err != nil {
				fmt.Println(err)
			}
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username), command)
			if err != nil {
				fmt.Printf("Failed to publish move: %v\n", err)
			} else {
				fmt.Printf("Published move: %+v\n", command)
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

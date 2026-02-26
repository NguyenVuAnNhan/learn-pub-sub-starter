package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
	"fmt"
	"errors"
)

type SimpleQueueType int
const (
	DurableQueue SimpleQueueType = iota
	TransientQueue
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

	_, _, err = DeclareAndBind(
		conn,
		routing.ExchangePerilDirect, 
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey, 
		TransientQueue,
	)

	if err != nil {
		panic(err)
	}

	gameState := gamelogic.NewGameState(username)

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

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to open channel: %w", err)
	}

	durable := queueType == DurableQueue

	queue, err := ch.QueueDeclare(
		queueName,
		durable,
		!durable, // delete when unused
		!durable, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		queue.Name,
		key,
		exchange,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		return nil, amqp.Queue{}, fmt.Errorf("failed to bind queue: %w", err)
	}

	return ch, queue, nil
}

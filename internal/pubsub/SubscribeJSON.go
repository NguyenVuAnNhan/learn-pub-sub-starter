package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"fmt"
	"encoding/json"
)

func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T),
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)

	if err != nil {
		return fmt.Errorf("failed to declare and bind: %w", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		false,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func() {
		for msg := range msgs {
			var data T
			err := json.Unmarshal(msg.Body, &data)
			if err != nil {
				fmt.Printf("failed to unmarshal message: %v\n", err)
				msg.Nack(false, false) // reject the message and don't requeue
				continue
			}
			handler(data)
			msg.Ack(false) // acknowledge the message
		}
	}()

	return nil
}
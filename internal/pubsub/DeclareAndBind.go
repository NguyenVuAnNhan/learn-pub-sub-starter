package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"fmt"
)

type SimpleQueueType int
const (
	DurableQueue SimpleQueueType = iota
	TransientQueue
)

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
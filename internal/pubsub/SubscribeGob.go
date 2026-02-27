package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"fmt"
	"encoding/gob"
	"bytes"
)

func SubscribeGob[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T) AckMode,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)

	if err != nil {
		return fmt.Errorf("failed to declare and bind: %w", err)
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
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
			err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&data)
			if err != nil {
				fmt.Printf("failed to decode message: %v\n", err)
				msg.Nack(false, false) // reject the message and don't requeue
				continue
			}
			ackMode := handler(data)
			switch ackMode {
			case Ack:
				fmt.Printf("Acknowledging message: %s\n", string(msg.Body))
				msg.Ack(false)
			case NackRequeue:
				fmt.Printf("Nacking and requeuing message: %s\n", string(msg.Body))
				msg.Nack(false, true)
			case NackDiscard:
				fmt.Printf("Nacking and discarding message: %s\n", string(msg.Body))
				msg.Nack(false, false)
			}
		}
	}()

	return nil
}
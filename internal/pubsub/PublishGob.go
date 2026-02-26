package pubsub

import (
	"encoding/gob"
	amqp "github.com/rabbitmq/amqp091-go"
	"bytes"
)

func PublishGob[T any](
	ch *amqp.Channel,
	exchange,
	key string,
	msg T,
) AckMode {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		return NackRequeue
	}

	err := ch.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        buf.Bytes(),
		},
	)

	if err != nil {
		return NackRequeue
	}

	return Ack
}
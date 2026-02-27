package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func publishGameLog(
	ch *amqp.Channel,
	exchange string,
	username string,
	log routing.GameLog,
) error {
	routingKey := routing.GameLogSlug + "." + username
	return pubsub.PublishGob(ch, exchange, routingKey, log)
}
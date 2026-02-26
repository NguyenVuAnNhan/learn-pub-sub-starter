package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	
	"fmt"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckMode {
	return func(move gamelogic.ArmyMove) pubsub.AckMode {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		if moveOutcome == gamelogic.MoveOutcomeSamePlayer {
			return pubsub.NackDiscard
		}
		if moveOutcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		}
		if moveOutcome == gamelogic.MoveOutcomeMakeWar {
			recog := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}
			pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, move.Player.Username),
				recog,
			)
			return pubsub.NackRequeue
		}
		return pubsub.NackDiscard
	}
}
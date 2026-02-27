package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"errors"
	"fmt"
	"time"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckMode {
	return func(recog gamelogic.RecognitionOfWar) pubsub.AckMode {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(recog)
		if outcome == gamelogic.WarOutcomeNotInvolved {
			return pubsub.NackRequeue
		} else if outcome == gamelogic.WarOutcomeNoUnits {
			return pubsub.NackDiscard
		} else if outcome == gamelogic.WarOutcomeOpponentWon {
			gameLog := routing.GameLog{
				CurrentTime: time.Now(),
				Message: fmt.Sprintf("%s won a war against %s", winner, loser),
				Username: gs.Player.Username,
			}
			err := publishGameLog(ch, routing.ExchangePerilTopic, gs.Player.Username, gameLog)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeYouWon {
			gameLog := routing.GameLog{
				CurrentTime: time.Now(),
				Message: fmt.Sprintf("%s won a war against %s", winner, loser),
				Username: gs.Player.Username,
			}
			err := publishGameLog(ch, routing.ExchangePerilTopic, gs.Player.Username, gameLog)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeYouWon {
			gameLog := routing.GameLog{
				CurrentTime: time.Now(),
				Message: fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser),
				Username: gs.Player.Username,
			}
			err := publishGameLog(ch, routing.ExchangePerilTopic, gs.Player.Username, gameLog)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		errors.New("unknown war outcome")
		return pubsub.NackDiscard
	}
}
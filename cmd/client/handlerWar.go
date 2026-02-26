package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"fmt"
	"errors"
)

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckMode {
	return func(recog gamelogic.RecognitionOfWar) pubsub.AckMode {
		defer fmt.Print("> ")
		outcome, _, _ := gs.HandleWar(recog)
		if outcome == gamelogic.WarOutcomeNotInvolved {
			return pubsub.NackRequeue
		} else if outcome == gamelogic.WarOutcomeNoUnits {
			return pubsub.NackDiscard
		} else if outcome == gamelogic.WarOutcomeOpponentWon {
			return pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeYouWon {
			return pubsub.Ack
		} else if outcome == gamelogic.WarOutcomeDraw {
			return pubsub.Ack
		}
		errors.New("unknown war outcome")
		return pubsub.NackDiscard
	}
}
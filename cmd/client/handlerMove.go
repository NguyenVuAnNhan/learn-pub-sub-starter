package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"fmt"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckMode {
	return func(move gamelogic.ArmyMove) pubsub.AckMode {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		if moveOutcome == gamelogic.MoveOutcomeSamePlayer {
			return pubsub.NackDiscard
		}
		if moveOutcome == gamelogic.MoveOutcomeMakeWar || moveOutcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		}
		return pubsub.NackDiscard
	}
}
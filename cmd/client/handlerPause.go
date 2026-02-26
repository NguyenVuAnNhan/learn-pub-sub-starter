package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"fmt"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckMode {
	return func(state routing.PlayingState) pubsub.AckMode {
		defer fmt.Print("> ")
		gs.HandlePause(state)
		return pubsub.Ack
	}
}
package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"fmt"
)

func handlerLog() func(routing.GameLog) pubsub.AckMode {
	return func(log routing.GameLog) pubsub.AckMode {
		defer fmt.Print("> ")
		gamelogic.WriteLog(log)
		return pubsub.Ack
	}
}
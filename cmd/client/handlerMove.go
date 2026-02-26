package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"fmt"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(move gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(move)
	}
}
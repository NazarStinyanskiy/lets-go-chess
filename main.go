package main

import (
	"lets-go-chess/game"
	"lets-go-chess/ui"
)

func main() {
	game.CreateField()
	ui.DrawBoard(game.Field)
}

package main

import (
	"lets-go-chess/cli_ui"
	"lets-go-chess/game"
)

func main() {
	game.CreateField()
	chooseMode(0)
}

func chooseMode(mode int) {
	switch mode {
	case 0:
		cli_ui.StartGame()
	}
}

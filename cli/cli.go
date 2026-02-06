package cli

import (
	"fmt"
	"lets-go-chess/game"
)

func StartGame() {
	g := game.StartGame()
	var move string
	for {
		game.DrawConsoleBoard(g.Field)
		fmt.Print("Enter your move: ")
		fmt.Scan(&move)

		if move == "exit" {
			break
		}

		runes := []rune(move)
		fromY := runes[0] - 96
		fromX := runes[1] - 48
		toY := runes[2] - 96
		toX := runes[3] - 48

		situation, moveErr := g.NextMove(game.Position{X: int(fromX), Y: int(fromY)}, game.Position{X: int(toX), Y: int(toY)})
		if moveErr != nil {
			fmt.Print("\033[31m", moveErr, "\033[0m\n")
			continue
		}
		switch situation {
		case game.Continue:
			continue
		case game.Check:
			fmt.Println("Check!")
			continue
		case game.Checkmate:
			fmt.Print("Checkmate!")
			break
		case game.Stalemate:
			fmt.Print("Stalemate!")
			break
		}
	}
}

package cli_ui

import (
	"fmt"
	"lets-go-chess/game"
)

func StartGame() {
	playerWhite := game.Player{IsWhite: true}
	playerBlack := game.Player{IsWhite: false}
	var from string
	var to string
	isWhiteMove := true
	for {
		drawBoard(game.Field)
		fmt.Print("Enter your 'from' move: ")
		fmt.Scan(&from)
		fmt.Print("Enter your 'to' move: ")
		fmt.Scan(&to)
		fromRunes := []rune(from)
		fromY := fromRunes[0] - 96
		fromX := fromRunes[1] - 48
		toRunes := []byte(to)
		toY := toRunes[0] - 96
		toX := toRunes[1] - 48
		var moveErr error
		if isWhiteMove {
			moveErr = game.Move(game.Position{X: int(fromX), Y: int(fromY)}, game.Position{X: int(toX), Y: int(toY)}, &playerWhite)
		} else {
			moveErr = game.Move(game.Position{X: int(fromX), Y: int(fromY)}, game.Position{X: int(toX), Y: int(toY)}, &playerBlack)
		}
		if moveErr == nil {
			isWhiteMove = !isWhiteMove
		} else {
			fmt.Println(moveErr)
		}
	}
}

func drawBoard(field game.Board) {
	for y := 0; y <= 8; y++ {
		for x := 0; x <= 8; x++ {
			if x == 0 || y == 0 {
				drawNumber(x, y)
			}
			if y > 0 && x > 0 {
				drawCell(field.Cells[game.Position{X: x, Y: y}])
			}
		}
		fmt.Println()
	}
}

func drawNumber(x, y int) {
	if x == 0 && y == 0 {
		fmt.Print("   ")
		return
	}
	if y == 0 {
		fmt.Printf(" %v ", x)
		return
	}
	if x == 0 {
		fmt.Printf(" %c ", 'a'+y-1)
		return
	}
}

func drawCell(figure *game.Figure) {
	if figure == nil || figure.Mover == nil {
		fmt.Print("   ")
		return
	}
	t := fmt.Sprintf("%T", figure.Mover)
	if figure.IsWhite {
		drawWhiteFigure(t)
	} else {
		drawBlackFigure(t)
	}

}

func drawBlackFigure(t string) {
	switch t {
	case "game.Pawn":
		fmt.Print(" \u2659 ")
	case "game.Rook":
		fmt.Print(" \u2656 ")
	case "game.Knight":
		fmt.Print(" \u2658 ")
	case "game.Bishop":
		fmt.Print(" \u2657 ")
	case "game.Queen":
		fmt.Print(" \u2655 ")
	case "game.King":
		fmt.Print(" \u2654 ")
	}
}

func drawWhiteFigure(t string) {
	switch t {
	case "game.Pawn":
		fmt.Print(" \u265F ")
	case "game.Rook":
		fmt.Print(" \u265C ")
	case "game.Knight":
		fmt.Print(" \u265E ")
	case "game.Bishop":
		fmt.Print(" \u265D ")
	case "game.Queen":
		fmt.Print(" \u265B ")
	case "game.King":
		fmt.Print(" \u265A ")
	}
}

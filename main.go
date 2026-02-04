package main

import (
	"lets-go-chess/cli"
)

func main() {
	chooseMode(0)
}

func chooseMode(mode int) {
	switch mode {
	case 0:
		cli.StartGame()
	}
}

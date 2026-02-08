package main

import (
	"lets-go-chess/cli"
	"lets-go-chess/server"
)

func main() {
	chooseMode(1)
}

func chooseMode(mode int) {
	switch mode {
	case 0:
		cli.StartGame()
	case 1:
		server.StartServer()
	}
}

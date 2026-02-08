package storage

import "lets-go-chess/game"

var storage = make(map[int]*game.Game)
var nextGameId = 1

func GetGameById(id int) *game.Game {
	return storage[id]
}

func SetGame(game *game.Game) int {
	defer func() {
		nextGameId++
	}()
	storage[nextGameId] = game
	return nextGameId
}

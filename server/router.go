package server

import (
	"encoding/json"
	"fmt"
	"lets-go-chess/game"
	"lets-go-chess/storage"
	"log"
	"net/http"
	"strings"
	"time"
)

type gameResponse struct {
	Board     [][]string     `json:"board"`
	GameId    int            `json:"gameId,omitempty"`
	Situation game.Situation `json:"situation,omitempty"`
	IsWhite   bool           `json:"isWhite"`
}

type moveRequest struct {
	GameId int `json:"gameId"`
	FromX  int `json:"fromX"`
	FromY  int `json:"fromY"`
	ToX    int `json:"toX"`
	ToY    int `json:"toY"`
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /startGame", corsMiddleware(startGame))
	mux.HandleFunc("POST /move", corsMiddleware(move))
	mux.HandleFunc("OPTIONS /move", corsMiddleware(nil))

	server := http.Server{
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if handler != nil {
			handler(w, r)
		}
	}
}

func startGame(w http.ResponseWriter, r *http.Request) {
	g := game.StartGame()
	gameId := storage.SetGame(g)
	resp := &gameResponse{}
	resp.GameId = gameId
	resp.IsWhite = true
	resp.Board = convertBoard(g)
	marshal, err := json.Marshal(resp)
	if err != nil {
		log.Print("Error marshalling response", err)
	}
	_, err = w.Write(marshal)
	if err != nil {
		log.Print("Error writing response", err)
	}
}

func move(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req moveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Error unmarshalling request", err)
		return
	}

	g := storage.GetGameById(req.GameId)
	situation, err := g.NextMove(game.Position{X: req.FromX, Y: req.FromY}, game.Position{X: req.ToX, Y: req.ToY})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := &gameResponse{}
	resp.Situation = situation
	resp.IsWhite = g.IsWhiteMove
	resp.Board = convertBoard(g)
	marshal, err := json.Marshal(resp)
	if err != nil {
		log.Print("Error marshalling response", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(marshal)
	if err != nil {
		log.Print("Error writing response", err)
	}
}

func convertBoard(g *game.Game) [][]string {
	board := make([][]string, 8)
	for i := 0; i < 8; i++ {
		board[i] = make([]string, 8)
		for j := 0; j < 8; j++ {
			figure := g.Field.Cells[game.Position{X: j + 1, Y: i + 1}]
			board[i][j] = figureToLetter(figure)
		}
	}
	return board
}

func figureToLetter(figure *game.Figure) string {
	if figure == nil {
		return ""
	}
	t := fmt.Sprintf("%T", figure.Mover)
	var f string
	switch t {
	case "game.Pawn":
		f = "p"
	case "game.Rook":
		f = "r"
	case "game.Knight":
		f = "n"
	case "game.Bishop":
		f = "b"
	case "game.Queen":
		f = "q"
	case "game.King":
		f = "k"
	default:
		return ""
	}
	if figure.IsWhite {
		f = strings.ToUpper(f)
	}
	return f
}

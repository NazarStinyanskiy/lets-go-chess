package game

import "errors"

type Player struct {
	IsWhite bool
	Situation
}

type Board struct {
	Cells map[Position]*Figure
}

var (
	Field          Board
	PlayerWhite    *Player
	PlayerBlack    *Player
	isWhiteMove    = true
	enPassantWhite *Figure
	enPassantBlack *Figure
)

var (
	InvalidFrom        = errors.New("invalid from")
	ToOutOfBounds      = errors.New("to out of bounds")
	MoveRulesViolation = errors.New("move rules violation")
	WrongColor         = errors.New("wrong color")
)

func CreateField() {
	Field = Board{make(map[Position]*Figure)}
	for x := 1; x <= 8; x++ {
		for y := 1; y <= 8; y++ {
			if y == 1 || y == 2 || y == 7 || y == 8 {
				Field.Cells[Position{X: x, Y: y}] = createFigure(x, y, y == 1 || y == 2)
			} else {
				Field.Cells[Position{X: x, Y: y}] = createEmptyCell()
			}
		}
	}
}

func NextMove(from, to Position) (Situation, error) {
	var player *Player
	if isWhiteMove {
		player = PlayerWhite
	} else {
		player = PlayerBlack
	}
	situation, err := move(from, to, player)
	if err == nil {
		isWhiteMove = !isWhiteMove
	}
	return situation, err
}

func move(from, to Position, player *Player) (Situation, error) {
	figure := Field.Cells[from]
	if figure == nil {
		return Continue, InvalidFrom
	}
	if player.IsWhite != figure.IsWhite {
		return Continue, WrongColor
	}
	if to.X > 8 || to.X < 1 || to.Y > 8 || to.Y < 1 {
		return Continue, ToOutOfBounds
	}
	if player.IsWhite && enPassantWhite != nil {
		enPassantWhite.IsVulnerableForEnPassant = false
		enPassantWhite = nil
	}
	if !player.IsWhite && enPassantBlack != nil {
		enPassantBlack.IsVulnerableForEnPassant = false
		enPassantBlack = nil
	}
	canMove, moveDetails := figure.canMove(Field, from, to, player.Situation)
	if !canMove {
		return Continue, MoveRulesViolation
	}
	if moveDetails == ReadyForEnPassant {
		figure.IsVulnerableForEnPassant = true
		if isWhiteMove {
			enPassantWhite = figure
		} else {
			enPassantBlack = figure
		}
	}
	Field = figure.move(Field, from, to, moveDetails)
	figure.HasMoved = true
	situation := analyzeSituation(Field, player.IsWhite, player.Situation)
	if player.IsWhite {
		PlayerBlack.Situation = situation
	} else {
		PlayerWhite.Situation = situation
	}
	return situation, nil
}

func createFigure(x, y int, isWhite bool) *Figure {
	f := &Figure{IsWhite: isWhite}
	pawn := Pawn{}
	rook := Rook{}
	knight := Knight{}
	bishop := Bishop{}
	queen := Queen{bishop: bishop, rook: rook}
	king := King{}
	if y == 2 || y == 7 {
		f.Mover = pawn
	}
	if (x == 1 || x == 8) && (y == 1 || y == 8) {
		f.Mover = rook
	}
	if (x == 2 || x == 7) && (y == 1 || y == 8) {
		f.Mover = knight
	}
	if (x == 3 || x == 6) && (y == 1 || y == 8) {
		f.Mover = bishop
	}
	if x == 4 && (y == 1 || y == 8) {
		f.Mover = queen
	}
	if x == 5 && (y == 1 || y == 8) {
		f.Mover = king
	}
	return f
}

func createEmptyCell() *Figure {
	return nil
}

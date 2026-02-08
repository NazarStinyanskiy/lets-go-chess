package game

import "errors"

type Game struct {
	Field          Board
	PlayerWhite    *Player
	PlayerBlack    *Player
	IsWhiteMove    bool
	enPassantWhite *Figure
	enPassantBlack *Figure
}

type Player struct {
	IsWhite bool
	Situation
}

type Board struct {
	Cells map[Position]*Figure
}

var (
	InvalidFrom        = errors.New("invalid from")
	ToOutOfBounds      = errors.New("to out of bounds")
	MoveRulesViolation = errors.New("move rules violation")
	WrongColor         = errors.New("wrong color")
)

func StartGame() *Game {
	return &Game{
		Field:       createField(),
		PlayerWhite: &Player{IsWhite: true, Situation: Continue},
		PlayerBlack: &Player{IsWhite: false, Situation: Continue},
		IsWhiteMove: true,
	}
}

func (g *Game) NextMove(from, to Position) (Situation, error) {
	var player *Player
	if g.IsWhiteMove {
		player = g.PlayerWhite
	} else {
		player = g.PlayerBlack
	}
	situation, err := g.move(from, to, player)
	if err == nil {
		g.IsWhiteMove = !g.IsWhiteMove
	}
	return situation, err
}

func (g *Game) move(from, to Position, player *Player) (Situation, error) {
	figure := g.Field.Cells[from]
	if figure == nil {
		return Continue, InvalidFrom
	}
	if player.IsWhite != figure.IsWhite {
		return Continue, WrongColor
	}
	if to.X > 8 || to.X < 1 || to.Y > 8 || to.Y < 1 {
		return Continue, ToOutOfBounds
	}
	if player.IsWhite && g.enPassantWhite != nil {
		g.enPassantWhite.IsVulnerableForEnPassant = false
		g.enPassantWhite = nil
	}
	if !player.IsWhite && g.enPassantBlack != nil {
		g.enPassantBlack.IsVulnerableForEnPassant = false
		g.enPassantBlack = nil
	}
	canMove, moveDetails := figure.canMove(g.Field, from, to, player.Situation)
	if !canMove {
		return Continue, MoveRulesViolation
	}
	if moveDetails == ReadyForEnPassant {
		figure.IsVulnerableForEnPassant = true
		if g.IsWhiteMove {
			g.enPassantWhite = figure
		} else {
			g.enPassantBlack = figure
		}
	}
	g.Field = figure.move(g.Field, from, to, moveDetails)
	figure.HasMoved = true
	situation := analyzeSituation(g.Field, player.IsWhite, player.Situation)
	if player.IsWhite {
		g.PlayerBlack.Situation = situation
	} else {
		g.PlayerWhite.Situation = situation
	}
	return situation, nil
}

func createField() Board {
	field := Board{make(map[Position]*Figure)}
	for x := 1; x <= 8; x++ {
		for y := 1; y <= 8; y++ {
			if y == 1 || y == 2 || y == 7 || y == 8 {
				field.Cells[Position{X: x, Y: y}] = createFigure(x, y, y == 1 || y == 2)
			} else {
				field.Cells[Position{X: x, Y: y}] = createEmptyCell()
			}
		}
	}
	return field
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

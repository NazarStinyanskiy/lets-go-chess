package game

import "errors"

type Player struct {
	IsWhite bool
}

type Position struct {
	X int
	Y int
}

type Board struct {
	Cells map[Position]*Figure
}

var Field Board
var enPassantWhite *Figure
var enPassantBlack *Figure

var (
	InvalidFrom        = errors.New("Invalid from")
	ToOutOfBounds      = errors.New("To out of bounds")
	MoveRulesViolation = errors.New("Move rules violation")
	WrongColor         = errors.New("Wrong color")
)

func CreateField() {
	Field = Board{make(map[Position]*Figure)}
	for x := 1; x <= 8; x++ {
		for y := 1; y <= 8; y++ {
			if y == 1 || y == 2 || y == 7 || y == 8 {
				createFigure(x, y, y == 1 || y == 2)
			} else {
				createEmptyCell(x, y)
			}
		}
	}
}

func Move(from, to Position, player *Player) error {
	figure := Field.Cells[from]
	if figure == nil {
		return InvalidFrom
	}
	if player.IsWhite != figure.IsWhite {
		return WrongColor
	}
	if to.X > 8 || to.X < 1 || to.Y > 8 || to.Y < 1 {
		return ToOutOfBounds
	}
	if player.IsWhite && enPassantWhite != nil {
		enPassantWhite.IsVulnerableForEnPassant = false
		enPassantWhite = nil
	}
	if !player.IsWhite && enPassantBlack != nil {
		enPassantBlack.IsVulnerableForEnPassant = false
		enPassantBlack = nil
	}
	canMove, specialMove := figure.canMove(from, to)
	if !canMove {
		return MoveRulesViolation
	}
	figure.move(from, to, specialMove)
	return nil
}

func createFigure(x, y int, isWhite bool) {
	f := &Figure{IsWhite: isWhite}
	pawn := Pawn{}
	rook := Rook{}
	knight := Knight{}
	bishop := Bishop{}
	queen := Queen{bishop: bishop, rook: rook}
	king := King{}
	if y == 2 || y == 7 {
		f.Mover = pawn
		Field.Cells[Position{x, y}] = f
		return
	}

	if (x == 1 || x == 8) && (y == 1 || y == 8) {
		f.Mover = rook
		Field.Cells[Position{x, y}] = f
		return
	}

	if (x == 2 || x == 7) && (y == 1 || y == 8) {
		f.Mover = knight
		Field.Cells[Position{x, y}] = f
		return
	}

	if (x == 3 || x == 6) && (y == 1 || y == 8) {
		f.Mover = bishop
		Field.Cells[Position{x, y}] = f
		return
	}

	if x == 4 && (y == 1 || y == 8) {
		f.Mover = queen
		Field.Cells[Position{x, y}] = f
		return
	}

	if x == 5 && (y == 1 || y == 8) {
		f.Mover = king
		Field.Cells[Position{x, y}] = f
		return
	}
}

func createEmptyCell(x, y int) {
	Field.Cells[Position{x, y}] = nil
}

package game

import "errors"

var Field Board

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
	if Field.Cells[from] == nil {
		return InvalidFrom
	}
	if player.IsWhite != Field.Cells[from].IsWhite {
		return WrongColor
	}
	if to.X > 8 || to.X < 1 || to.Y > 8 || to.Y < 1 {
		return ToOutOfBounds
	}
	if !Field.Cells[from].move(from, to, &Field) {
		return MoveRulesViolation
	}
	return nil
}

func createFigure(x, y int, isWhite bool) {
	f := &Figure{IsWhite: isWhite}
	if y == 2 || y == 7 {
		f.Mover = Pawn{}
		Field.Cells[Position{x, y}] = f
		return
	}

	if (x == 1 || x == 8) && (y == 1 || y == 8) {
		f.Mover = Rook{}
		Field.Cells[Position{x, y}] = f
		return
	}

	if (x == 2 || x == 7) && (y == 1 || y == 8) {
		f.Mover = Knight{}
		Field.Cells[Position{x, y}] = f
		return
	}

	if (x == 3 || x == 6) && (y == 1 || y == 8) {
		f.Mover = Bishop{}
		Field.Cells[Position{x, y}] = f
		return
	}

	if x == 4 && (y == 1 || y == 8) {
		f.Mover = Queen{}
		Field.Cells[Position{x, y}] = f
		return
	}

	if x == 5 && (y == 1 || y == 8) {
		f.Mover = King{}
		Field.Cells[Position{x, y}] = f
		return
	}
}

func createEmptyCell(x, y int) {
	Field.Cells[Position{x, y}] = nil
}

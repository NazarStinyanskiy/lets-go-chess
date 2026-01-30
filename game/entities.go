package game

import "math"

type Position struct {
	X int
	Y int
}

type Figure struct {
	IsWhite  bool
	HasMoved bool
	Mover
}

type Mover interface {
	move(from, to Position, board *Board) bool
}

type Board struct {
	Cells map[Position]*Figure
}

/** Figure implementations **/

type King struct{}

func (k King) move(from, to Position, board *Board) bool {
	toIsEmpty := board.Cells[to] == nil
	isFightingEnemy := board.Cells[to] != nil && board.Cells[to].IsWhite != board.Cells[from].IsWhite
	deltaX := to.X - from.X
	deltaY := to.Y - from.Y

	if deltaX <= 1 && deltaX >= -1 || deltaY <= 1 && deltaY >= -1 {
		if isFightingEnemy || toIsEmpty {
			Field.Cells[to] = Field.Cells[from]
			Field.Cells[from] = nil
			return true
		}
	}

	if board.Cells[from].IsWhite == board.Cells[to].IsWhite {
		//TODO: castling
		return false
	}
	return false
}

type Queen struct {
	bishop Bishop
	rook   Rook
}

func (q Queen) move(from, to Position, board *Board) bool {
	return q.bishop.move(from, to, board) || q.rook.move(from, to, board)
}

type Rook struct{}

func (r Rook) move(from, to Position, board *Board) (ok bool) {
	defer func() {
		if ok {
			Field.Cells[to] = Field.Cells[from]
			Field.Cells[from] = nil
		}
	}()
	toIsEmpty := board.Cells[to] == nil
	isFightingEnemy := board.Cells[to] != nil && board.Cells[to].IsWhite != board.Cells[from].IsWhite
	deltaX := to.X - from.X
	deltaY := to.Y - from.Y

	if deltaX != 0 && deltaY != 0 {
		return false
	}
	if deltaX != 0 {
		minX := int(math.Min(float64(from.X), float64(to.X))) + 1
		maxX := int(math.Max(float64(from.X), float64(to.X)))
		for x := minX; x < maxX; x++ {
			if board.Cells[Position{X: x, Y: from.Y}] != nil {
				return false
			}
		}
		if isFightingEnemy || toIsEmpty {
			return true
		}
	}
	if deltaY != 0 {
		minY := int(math.Min(float64(from.Y), float64(to.Y))) + 1
		maxY := int(math.Max(float64(from.Y), float64(to.Y)))
		for y := minY; y < maxY; y++ {
			if board.Cells[Position{X: from.X, Y: y}] != nil {
				return false
			}
		}
		if isFightingEnemy || toIsEmpty {
			return true
		}
	}
	return false
}

type Bishop struct{}

func (b Bishop) move(from, to Position, board *Board) (ok bool) {
	defer func() {
		if ok {
			Field.Cells[to] = Field.Cells[from]
			Field.Cells[from] = nil
		}
	}()
	toIsEmpty := board.Cells[to] == nil
	isFightingEnemy := board.Cells[to] != nil && board.Cells[to].IsWhite != board.Cells[from].IsWhite
	deltaX := to.X - from.X
	deltaY := to.Y - from.Y

	if math.Abs(float64(deltaX)) == math.Abs(float64(deltaY)) {
		x := from.X
		y := from.Y
		for x != to.X && y != to.Y {
			if board.Cells[Position{X: x, Y: y}] != nil {
				return false
			}
			if to.X > from.X {
				x++
			} else {
				x--
			}
			if to.Y > from.Y {
				y++
			} else {
				y--
			}
		}
		if isFightingEnemy || toIsEmpty {
			return true
		}
	}
	return false
}

type Knight struct{}

func (k Knight) move(from, to Position, board *Board) (ok bool) {
	defer func() {
		if ok {
			Field.Cells[to] = Field.Cells[from]
			Field.Cells[from] = nil
		}
	}()
	toIsEmpty := board.Cells[to] == nil
	isFightingEnemy := board.Cells[to] != nil && board.Cells[to].IsWhite != board.Cells[from].IsWhite
	deltaX := to.X - from.X
	deltaY := to.Y - from.Y

	if isFightingEnemy || toIsEmpty {
		if (deltaX == 2 || deltaX == -2) && (deltaY == 1 || deltaY == -1) {
			return true
		}
		if (deltaY == 2 || deltaY == -2) && (deltaX == 1 || deltaX == -1) {
			return true
		}
	}
	return false
}

type Pawn struct{}

func (p Pawn) move(from, to Position, board *Board) (ok bool) {
	defer func() {
		if ok {
			Field.Cells[to] = Field.Cells[from]
			Field.Cells[from] = nil
		}
	}()
	figure := board.Cells[from]
	isWhite := figure.IsWhite
	encounter := board.Cells[to] != nil
	isFightingEnemy := encounter && board.Cells[to].IsWhite != isWhite
	deltaX := to.X - from.X
	deltaY := to.Y - from.Y
	if isWhite {
		// move strictly forward without encounter
		if deltaY == 1 && deltaX == 0 && !encounter {
			return true
		}
		// fight left or right
		if (deltaX == 1 || deltaX == -1) && deltaY == 1 && isFightingEnemy {
			return true
		}
		// first move
		if (!figure.HasMoved && deltaY == 2 && deltaX == 0 && !encounter && board.Cells[Position{to.X, to.Y - 1}] == nil) {
			figure.HasMoved = true
			return true
		}
		//TODO: en passant
	} else {
		// move strictly forward without encounter
		if deltaY == -1 && deltaX == 0 && !encounter {
			return true
		}
		// fight left or right
		if (deltaX == 1 || deltaX == -1) && deltaY == -1 && isFightingEnemy {
			return true
		}
		// first move
		if (!figure.HasMoved && deltaY == -2 && deltaX == 0 && !encounter && board.Cells[Position{to.X, to.Y + 1}] == nil) {
			figure.HasMoved = true
			return true
		}
		//TODO: en passant
	}
	return false
}

/** Players **/

type Player struct {
	IsWhite bool
}

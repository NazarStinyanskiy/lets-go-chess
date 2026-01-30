package game

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
	if to.X != from.X || to.X != from.X+1 || to.X != from.X-1 && to.Y != from.Y || to.Y != from.Y+1 || to.Y != from.Y-1 {
		return false
	}
	if board.Cells[from].IsWhite == board.Cells[to].IsWhite {
		//TODO: castling
		return false
	}
	Field.Cells[to] = Field.Cells[from]
	Field.Cells[from] = nil
	return true
}

type Queen struct{}

func (q Queen) move(from, to Position, board *Board) bool {
	return false
}

type Rook struct{}

func (r Rook) move(from, to Position, board *Board) bool {
	return false
}

type Bishop struct{}

func (b Bishop) move(from, to Position, board *Board) bool {
	return false
}

type Knight struct{}

func (k Knight) move(from, to Position, board *Board) bool {
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

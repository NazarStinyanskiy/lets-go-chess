package game

import "math"

type MoveDetails int

const (
	None MoveDetails = iota
	ReadyForEnPassant
	EnPassant
	ShortCastling
	LongCastling
)

type Situation int

const (
	Continue Situation = iota
	Check
	Checkmate
	Stalemate
)

type Position struct {
	X int
	Y int
}

type Figure struct {
	IsWhite                  bool
	HasMoved                 bool
	IsVulnerableForEnPassant bool
	Mover
}

type Mover interface {
	canMove(field Board, from, to Position, situation Situation) (bool, MoveDetails)
	move(field Board, from, to Position, move MoveDetails) Board
}

type King struct{}

func (k King) canMove(field Board, from, to Position, situation Situation) (ok bool, move MoveDetails) {
	figure := field.Cells[from]
	deltaX, deltaY := getDelta(from, to)
	defer func() {
		if ok {
			board := k.move(field, from, to, move)
			kingPos := findKing(board, figure.IsWhite)
			ok = !isFigureInThreat(board, kingPos, situation)
		}
	}()
	isRookAndDidNotMove := func(pos Position) bool {
		f := field.Cells[pos]
		if f != nil {
			if _, isRook := f.Mover.(Rook); isRook && !f.HasMoved {
				return true
			}
		}
		return false
	}
	if deltaX <= 1 && deltaX >= -1 && deltaY <= 1 && deltaY >= -1 {
		if isFightingEnemy(field, from, to) || isCellEmpty(field, to) {
			return true, None
		}
	}
	if figure.IsWhite {
		if deltaX == 2 && deltaY == 0 && !figure.HasMoved && isRookAndDidNotMove(Position{X: 8, Y: 1}) {
			if field.Cells[Position{X: 6, Y: 1}] == nil && field.Cells[Position{X: 7, Y: 1}] == nil {
				prevPos := Position{X: to.X - 1, Y: to.Y}
				if !isFigureInThreat(k.move(field, from, prevPos, move), prevPos, situation) {
					return true, ShortCastling
				}
			}
		}
		if deltaX == -2 && deltaY == 0 && !figure.HasMoved && isRookAndDidNotMove(Position{X: 1, Y: 1}) {
			if field.Cells[Position{X: 4, Y: 1}] == nil && field.Cells[Position{X: 3, Y: 1}] == nil && field.Cells[Position{X: 2, Y: 1}] == nil {
				prevPos := Position{X: to.X + 1, Y: to.Y}
				if !isFigureInThreat(k.move(field, from, prevPos, move), prevPos, situation) {
					return true, LongCastling
				}
			}
		}
	} else {
		if deltaX == 2 && deltaY == 0 && !figure.HasMoved && isRookAndDidNotMove(Position{X: 8, Y: 8}) {
			if field.Cells[Position{X: 6, Y: 8}] == nil && field.Cells[Position{X: 7, Y: 8}] == nil {
				prevPos := Position{X: to.X - 1, Y: to.Y}
				if !isFigureInThreat(k.move(field, from, prevPos, move), prevPos, situation) {
					return true, ShortCastling
				}
			}
		}
		if deltaX == -2 && deltaY == 0 && !figure.HasMoved && isRookAndDidNotMove(Position{X: 1, Y: 8}) {
			if field.Cells[Position{X: 4, Y: 8}] == nil && field.Cells[Position{X: 3, Y: 8}] == nil && field.Cells[Position{X: 2, Y: 8}] == nil {
				prevPos := Position{X: to.X + 1, Y: to.Y}
				if !isFigureInThreat(k.move(field, from, prevPos, move), prevPos, situation) {
					return true, LongCastling
				}
			}
		}
	}
	return false, None
}

func (k King) move(field Board, from, to Position, move MoveDetails) Board {
	field = copyField(field)
	figure := field.Cells[from]
	replace(&field, from, to)
	if move == ShortCastling {
		if figure.IsWhite {
			field.Cells[Position{X: 6, Y: 1}] = field.Cells[Position{X: 8, Y: 1}]
			field.Cells[Position{X: 8, Y: 1}] = nil
		} else {
			field.Cells[Position{X: 6, Y: 8}] = field.Cells[Position{X: 8, Y: 8}]
			field.Cells[Position{X: 8, Y: 8}] = nil
		}
	}
	if move == LongCastling {
		if figure.IsWhite {
			field.Cells[Position{X: 4, Y: 1}] = field.Cells[Position{X: 1, Y: 1}]
			field.Cells[Position{X: 1, Y: 1}] = nil
		} else {
			field.Cells[Position{X: 4, Y: 8}] = field.Cells[Position{X: 1, Y: 8}]
			field.Cells[Position{X: 1, Y: 8}] = nil
		}
	}
	return field
}

type Queen struct {
	bishop Bishop
	rook   Rook
}

func (q Queen) canMove(field Board, from, to Position, situation Situation) (ok bool, move MoveDetails) {
	defer func() {
		if ok {
			board := q.move(field, from, to, move)
			kingPos := findKing(board, field.Cells[from].IsWhite)
			ok = !isFigureInThreat(board, kingPos, situation)
		}
	}()
	bishopOk, _ := q.bishop.canMove(field, from, to, situation)
	rookOk, _ := q.rook.canMove(field, from, to, situation)
	return bishopOk || rookOk, None
}

func (q Queen) move(field Board, from, to Position, _ MoveDetails) Board {
	field = copyField(field)
	replace(&field, from, to)
	return field
}

type Rook struct{}

func (r Rook) canMove(field Board, from, to Position, situation Situation) (ok bool, move MoveDetails) {
	deltaX, deltaY := getDelta(from, to)
	defer func() {
		if ok {
			board := r.move(field, from, to, move)
			kingPos := findKing(board, field.Cells[from].IsWhite)
			ok = !isFigureInThreat(board, kingPos, situation)
		}
	}()
	if deltaX != 0 && deltaY != 0 {
		return false, None
	}
	if deltaX != 0 {
		minX := int(math.Min(float64(from.X), float64(to.X))) + 1
		maxX := int(math.Max(float64(from.X), float64(to.X)))
		for x := minX; x < maxX; x++ {
			if field.Cells[Position{X: x, Y: from.Y}] != nil {
				return false, None
			}
		}
		if isFightingEnemy(field, from, to) || isCellEmpty(field, to) {
			return true, None
		}
	}
	if deltaY != 0 {
		minY := int(math.Min(float64(from.Y), float64(to.Y))) + 1
		maxY := int(math.Max(float64(from.Y), float64(to.Y)))
		for y := minY; y < maxY; y++ {
			if field.Cells[Position{X: from.X, Y: y}] != nil {
				return false, None
			}
		}
		if isFightingEnemy(field, from, to) || isCellEmpty(field, to) {
			return true, None
		}
	}
	return false, None
}

func (r Rook) move(field Board, from, to Position, _ MoveDetails) Board {
	field = copyField(field)
	replace(&field, from, to)
	return field
}

type Bishop struct{}

func (b Bishop) canMove(field Board, from, to Position, situation Situation) (ok bool, move MoveDetails) {
	deltaX, deltaY := getDelta(from, to)
	defer func() {
		if ok {
			board := b.move(field, from, to, move)
			kingPos := findKing(board, field.Cells[from].IsWhite)
			ok = !isFigureInThreat(board, kingPos, situation)
		}
	}()
	if math.Abs(float64(deltaX)) == math.Abs(float64(deltaY)) {
		x := from.X
		y := from.Y
		for {
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
			if x == to.X && y == to.Y {
				if isFightingEnemy(field, from, to) || isCellEmpty(field, to) {
					return true, None
				}
				return false, None
			}
			if field.Cells[Position{X: x, Y: y}] != nil {
				return false, None
			}
		}
	}
	return false, None
}

func (b Bishop) move(field Board, from, to Position, _ MoveDetails) Board {
	field = copyField(field)
	replace(&field, from, to)
	return field
}

type Knight struct{}

func (k Knight) canMove(field Board, from, to Position, situation Situation) (ok bool, move MoveDetails) {
	deltaX, deltaY := getDelta(from, to)
	defer func() {
		if ok {
			board := k.move(field, from, to, move)
			kingPos := findKing(board, field.Cells[from].IsWhite)
			ok = !isFigureInThreat(board, kingPos, situation)
		}
	}()
	if isFightingEnemy(field, from, to) || isCellEmpty(field, to) {
		if (deltaX == 2 || deltaX == -2) && (deltaY == 1 || deltaY == -1) {
			return true, None
		}
		if (deltaY == 2 || deltaY == -2) && (deltaX == 1 || deltaX == -1) {
			return true, None
		}
	}
	return false, None
}

func (k Knight) move(field Board, from, to Position, _ MoveDetails) Board {
	field = copyField(field)
	replace(&field, from, to)
	return field
}

type Pawn struct{}

func (p Pawn) canMove(field Board, from, to Position, situation Situation) (ok bool, move MoveDetails) {
	figure := field.Cells[from]
	deltaX, deltaY := getDelta(from, to)
	defer func() {
		if ok {
			board := p.move(field, from, to, move)
			kingPos := findKing(board, figure.IsWhite)
			ok = !isFigureInThreat(board, kingPos, situation)
		}
	}()
	if figure.IsWhite {
		if deltaY == 1 && deltaX == 0 && isCellEmpty(field, to) {
			return true, None
		}
		if (deltaX == 1 || deltaX == -1) && deltaY == 1 && isFightingEnemy(field, from, to) {
			return true, None
		}
		if (!figure.HasMoved && deltaY == 2 && deltaX == 0 && isCellEmpty(field, to) && field.Cells[Position{to.X, to.Y - 1}] == nil) {
			return true, ReadyForEnPassant
		}
		enPassantEnemy := field.Cells[Position{to.X, to.Y - 1}]
		if math.Abs(float64(deltaX)) == 1 && deltaY == 1 && enPassantEnemy != nil && !enPassantEnemy.IsWhite && enPassantEnemy.IsVulnerableForEnPassant {
			return true, EnPassant
		}
	} else {
		if deltaY == -1 && deltaX == 0 && isCellEmpty(field, to) {
			return true, None
		}
		if (deltaX == 1 || deltaX == -1) && deltaY == -1 && isFightingEnemy(field, from, to) {
			return true, None
		}
		if (!figure.HasMoved && deltaY == -2 && deltaX == 0 && isCellEmpty(field, to) && field.Cells[Position{to.X, to.Y + 1}] == nil) {
			return true, ReadyForEnPassant
		}
		enPassantEnemy := field.Cells[Position{to.X, to.Y + 1}]
		if math.Abs(float64(deltaX)) == 1 && deltaY == -1 && enPassantEnemy != nil && enPassantEnemy.IsWhite && enPassantEnemy.IsVulnerableForEnPassant {
			return true, EnPassant
		}
	}
	return false, None
}

func (p Pawn) move(field Board, from, to Position, move MoveDetails) Board {
	field = copyField(field)
	figure := field.Cells[from]
	replace(&field, from, to)
	if move == EnPassant {
		if figure.IsWhite {
			field.Cells[Position{X: to.X, Y: to.Y - 1}] = nil
		} else {
			field.Cells[Position{X: to.X, Y: to.Y + 1}] = nil
		}
	}
	return field
}

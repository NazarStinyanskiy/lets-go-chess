package game

import "math"

type Figure struct {
	IsWhite                  bool
	HasMoved                 bool
	IsVulnerableForEnPassant bool
	Mover
}

type MoveDetails int

const (
	None MoveDetails = iota
	EnPassant
	ShortCastling
	LongCastling
)

type Mover interface {
	canMove(from, to Position) (bool, MoveDetails)
	move(from, to Position, move MoveDetails)
}

type King struct{}

func (k King) canMove(from, to Position) (ok bool, move MoveDetails) {
	figure := Field.Cells[from]
	deltaX, deltaY := getDelta(from, to)
	if deltaX <= 1 && deltaX >= -1 && deltaY <= 1 && deltaY >= -1 {
		if (isFightingEnemy(from, to) || isCellEmpty(to)) && !isFigureInThreatAfterMove(from, to) {
			return true, None
		}
	}
	if figure.IsWhite {
		if deltaX == 2 && deltaY == 0 && !figure.HasMoved && !Field.Cells[Position{X: 8, Y: 1}].HasMoved {
			if Field.Cells[Position{X: 6, Y: 1}] == nil && Field.Cells[Position{X: 7, Y: 1}] == nil {
				if !isFigureInThreatAfterMove(from, Position{X: to.X - 1, Y: to.Y}) && !isFigureInThreatAfterMove(from, to) {
					return true, ShortCastling
				}
			}
		}
		if deltaX == -2 && deltaY == 0 && !figure.HasMoved && !Field.Cells[Position{X: 1, Y: 1}].HasMoved {
			if Field.Cells[Position{X: 4, Y: 1}] == nil && Field.Cells[Position{X: 3, Y: 1}] == nil && Field.Cells[Position{X: 2, Y: 1}] == nil {
				if !isFigureInThreatAfterMove(from, Position{X: to.X + 1, Y: to.Y}) && !isFigureInThreatAfterMove(from, to) {
					return true, LongCastling
				}
			}
		}
	} else {
		if deltaX == 2 && deltaY == 0 && !figure.HasMoved && !Field.Cells[Position{X: 8, Y: 8}].HasMoved {
			if Field.Cells[Position{X: 6, Y: 8}] == nil && Field.Cells[Position{X: 7, Y: 8}] == nil {
				if !isFigureInThreatAfterMove(from, Position{X: to.X - 1, Y: to.Y}) && !isFigureInThreatAfterMove(from, to) {
					return true, ShortCastling
				}
			}
		}
		if deltaX == -2 && deltaY == 0 && !figure.HasMoved && !Field.Cells[Position{X: 1, Y: 8}].HasMoved {
			if Field.Cells[Position{X: 4, Y: 8}] == nil && Field.Cells[Position{X: 3, Y: 8}] == nil && Field.Cells[Position{X: 2, Y: 8}] == nil {
				if !isFigureInThreatAfterMove(from, Position{X: to.X + 1, Y: to.Y}) && !isFigureInThreatAfterMove(from, to) {
					return true, LongCastling
				}
			}
		}
	}
	return false, None
}

func (k King) move(from, to Position, move MoveDetails) {
	figure := Field.Cells[from]
	figure.HasMoved = true
	replace(from, to)
	if move == ShortCastling {
		if figure.IsWhite {
			Field.Cells[Position{X: 6, Y: 1}] = Field.Cells[Position{X: 8, Y: 1}]
			Field.Cells[Position{X: 8, Y: 1}] = nil
		} else {
			Field.Cells[Position{X: 6, Y: 8}] = Field.Cells[Position{X: 8, Y: 8}]
			Field.Cells[Position{X: 8, Y: 8}] = nil
		}
	}
	if move == LongCastling {
		if figure.IsWhite {
			Field.Cells[Position{X: 4, Y: 1}] = Field.Cells[Position{X: 1, Y: 1}]
			Field.Cells[Position{X: 1, Y: 1}] = nil
		} else {
			Field.Cells[Position{X: 4, Y: 8}] = Field.Cells[Position{X: 1, Y: 8}]
			Field.Cells[Position{X: 1, Y: 8}] = nil
		}
	}
}

type Queen struct {
	bishop Bishop
	rook   Rook
}

func (q Queen) canMove(from, to Position) (bool, MoveDetails) {
	bishopOk, _ := q.bishop.canMove(from, to)
	rookOk, _ := q.rook.canMove(from, to)
	return bishopOk || rookOk, None
}

func (q Queen) move(from, to Position, move MoveDetails) {
	replace(from, to)
}

type Rook struct{}

func (r Rook) canMove(from, to Position) (ok bool, move MoveDetails) {
	deltaX, deltaY := getDelta(from, to)

	if deltaX != 0 && deltaY != 0 {
		return false, None
	}
	if deltaX != 0 {
		minX := int(math.Min(float64(from.X), float64(to.X))) + 1
		maxX := int(math.Max(float64(from.X), float64(to.X)))
		for x := minX; x < maxX; x++ {
			if Field.Cells[Position{X: x, Y: from.Y}] != nil {
				return false, None
			}
		}
		if isFightingEnemy(from, to) || isCellEmpty(to) {
			return true, None
		}
	}
	if deltaY != 0 {
		minY := int(math.Min(float64(from.Y), float64(to.Y))) + 1
		maxY := int(math.Max(float64(from.Y), float64(to.Y)))
		for y := minY; y < maxY; y++ {
			if Field.Cells[Position{X: from.X, Y: y}] != nil {
				return false, None
			}
		}
		if isFightingEnemy(from, to) || isCellEmpty(to) {
			return true, None
		}
	}
	return false, None
}

func (r Rook) move(from, to Position, move MoveDetails) {
	Field.Cells[from].HasMoved = true
	replace(from, to)
}

type Bishop struct{}

func (b Bishop) canMove(from, to Position) (ok bool, move MoveDetails) {
	deltaX, deltaY := getDelta(from, to)
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
				if isFightingEnemy(from, to) || isCellEmpty(to) {
					return true, None
				}
				return false, None
			}
			if Field.Cells[Position{X: x, Y: y}] != nil {
				return false, None
			}
		}
	}
	return false, None
}

func (b Bishop) move(from, to Position, move MoveDetails) {
	replace(from, to)
}

type Knight struct{}

func (k Knight) canMove(from, to Position) (ok bool, move MoveDetails) {
	deltaX, deltaY := getDelta(from, to)
	if isFightingEnemy(from, to) || isCellEmpty(to) {
		if (deltaX == 2 || deltaX == -2) && (deltaY == 1 || deltaY == -1) {
			return true, None
		}
		if (deltaY == 2 || deltaY == -2) && (deltaX == 1 || deltaX == -1) {
			return true, None
		}
	}
	return false, None
}

func (k Knight) move(from, to Position, move MoveDetails) {
	replace(from, to)
}

type Pawn struct{}

func (p Pawn) canMove(from, to Position) (ok bool, move MoveDetails) {
	figure := Field.Cells[from]
	deltaX, deltaY := getDelta(from, to)
	if figure.IsWhite {
		if deltaY == 1 && deltaX == 0 && isCellEmpty(to) {
			return true, None
		}
		if (deltaX == 1 || deltaX == -1) && deltaY == 1 && isFightingEnemy(from, to) {
			return true, None
		}
		if (!figure.HasMoved && deltaY == 2 && deltaX == 0 && isCellEmpty(to) && Field.Cells[Position{to.X, to.Y - 1}] == nil) {
			figure.IsVulnerableForEnPassant = true
			enPassantWhite = figure
			return true, None
		}
		enPassantEnemy := Field.Cells[Position{to.X, to.Y - 1}]
		if math.Abs(float64(deltaX)) == 1 && deltaY == 1 && enPassantEnemy != nil && !enPassantEnemy.IsWhite && enPassantEnemy.IsVulnerableForEnPassant {
			return true, EnPassant
		}
	} else {
		if deltaY == -1 && deltaX == 0 && isCellEmpty(to) {
			return true, None
		}
		if (deltaX == 1 || deltaX == -1) && deltaY == -1 && isFightingEnemy(from, to) {
			return true, None
		}
		if (!figure.HasMoved && deltaY == -2 && deltaX == 0 && isCellEmpty(to) && Field.Cells[Position{to.X, to.Y + 1}] == nil) {
			figure.IsVulnerableForEnPassant = true
			enPassantBlack = figure
			return true, None
		}
		enPassantEnemy := Field.Cells[Position{to.X, to.Y + 1}]
		if math.Abs(float64(deltaX)) == 1 && deltaY == -1 && enPassantEnemy != nil && enPassantEnemy.IsWhite && enPassantEnemy.IsVulnerableForEnPassant {
			return true, EnPassant
		}
	}
	return false, None
}

func (p Pawn) move(from, to Position, move MoveDetails) {
	figure := Field.Cells[from]
	figure.HasMoved = true
	replace(from, to)
	if move == EnPassant {
		if figure.IsWhite {
			Field.Cells[Position{X: to.X, Y: to.Y - 1}] = nil
		} else {
			Field.Cells[Position{X: to.X, Y: to.Y + 1}] = nil
		}
	}
}

func replace(from, to Position) {
	Field.Cells[to] = Field.Cells[from]
	Field.Cells[from] = nil
}

func isCellEmpty(pos Position) bool {
	return Field.Cells[pos] == nil
}

func isFightingEnemy(from, to Position) bool {
	return Field.Cells[to] != nil && Field.Cells[to].IsWhite != Field.Cells[from].IsWhite
}

func getDelta(from, to Position) (x int, y int) {
	return to.X - from.X, to.Y - from.Y
}

func isFigureInThreatNow(pos Position, cells map[Position]*Figure) bool {
	for figurePos, figure := range cells {
		if figure == nil || figurePos == pos {
			continue
		}
		canMove, _ := figure.canMove(figurePos, pos)
		if canMove {
			return true
		}
	}
	return false
}

func isFigureInThreatAfterMove(from, to Position) bool {
	var futureField = make(map[Position]*Figure, len(Field.Cells))
	for figurePos, figure := range Field.Cells {
		futureField[figurePos] = figure
	}
	futureField[to] = futureField[from]
	futureField[from] = nil
	return isFigureInThreatNow(to, futureField)
}

package game

import "math"

type Figure struct {
	IsWhite  bool
	HasMoved bool
	Mover
}

type Mover interface {
	move(from, to Position) bool
}

type King struct{}

func (k King) move(from, to Position) bool {
	deltaX, deltaY := getDelta(from, to)
	if deltaX <= 1 && deltaX >= -1 || deltaY <= 1 && deltaY >= -1 {
		if isFightingEnemy(from, to) || isCellEmpty(to) {
			replace(from, to)
			return true
		}
	}
	if Field.Cells[from].IsWhite == Field.Cells[to].IsWhite {
		//TODO: castling
		return false
	}
	return false
}

type Queen struct {
	bishop Bishop
	rook   Rook
}

func (q Queen) move(from, to Position) bool {
	return q.bishop.move(from, to) || q.rook.move(from, to)
}

type Rook struct{}

func (r Rook) move(from, to Position) (ok bool) {
	defer func() {
		if ok {
			replace(from, to)
		}
	}()
	deltaX, deltaY := getDelta(from, to)

	if deltaX != 0 && deltaY != 0 {
		return false
	}
	if deltaX != 0 {
		minX := int(math.Min(float64(from.X), float64(to.X))) + 1
		maxX := int(math.Max(float64(from.X), float64(to.X)))
		for x := minX; x < maxX; x++ {
			if Field.Cells[Position{X: x, Y: from.Y}] != nil {
				return false
			}
		}
		if isFightingEnemy(from, to) || isCellEmpty(to) {
			return true
		}
	}
	if deltaY != 0 {
		minY := int(math.Min(float64(from.Y), float64(to.Y))) + 1
		maxY := int(math.Max(float64(from.Y), float64(to.Y)))
		for y := minY; y < maxY; y++ {
			if Field.Cells[Position{X: from.X, Y: y}] != nil {
				return false
			}
		}
		if isFightingEnemy(from, to) || isCellEmpty(to) {
			return true
		}
	}
	return false
}

type Bishop struct{}

func (b Bishop) move(from, to Position) (ok bool) {
	defer func() {
		if ok {
			replace(from, to)
		}
	}()
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
					return true
				}
				return false
			}
			if Field.Cells[Position{X: x, Y: y}] != nil {
				return false
			}
		}
	}
	return false
}

type Knight struct{}

func (k Knight) move(from, to Position) (ok bool) {
	defer func() {
		if ok {
			replace(from, to)
		}
	}()
	deltaX, deltaY := getDelta(from, to)

	if isFightingEnemy(from, to) || isCellEmpty(to) {
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

func (p Pawn) move(from, to Position) (ok bool) {

	defer func() {
		if ok {
			replace(from, to)
		}
	}()
	figure := Field.Cells[from]
	deltaX, deltaY := getDelta(from, to)
	if figure.IsWhite {
		// move strictly forward without encounter
		if deltaY == 1 && deltaX == 0 && isCellEmpty(to) {
			return true
		}
		// fight left or right
		if (deltaX == 1 || deltaX == -1) && deltaY == 1 && isFightingEnemy(from, to) {
			return true
		}
		// first move
		if (!figure.HasMoved && deltaY == 2 && deltaX == 0 && isCellEmpty(to) && Field.Cells[Position{to.X, to.Y - 1}] == nil) {
			figure.HasMoved = true
			return true
		}
		//TODO: en passant
	} else {
		// move strictly forward without encounter
		if deltaY == -1 && deltaX == 0 && isCellEmpty(to) {
			return true
		}
		// fight left or right
		if (deltaX == 1 || deltaX == -1) && deltaY == -1 && isFightingEnemy(from, to) {
			return true
		}
		// first move
		if (!figure.HasMoved && deltaY == -2 && deltaX == 0 && isCellEmpty(to) && Field.Cells[Position{to.X, to.Y + 1}] == nil) {
			figure.HasMoved = true
			return true
		}
		//TODO: en passant
	}
	return false
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

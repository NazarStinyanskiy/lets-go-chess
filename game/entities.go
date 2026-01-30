package game

type Position struct {
	X int
	Y int
}

type Figure struct {
	IsWhite bool
	Mover
}

type Mover interface {
	move(from, to Position, board Board) bool
}

type Board struct {
	Cells map[Position]*Figure
}

/** Figure implementations **/

type King struct{}

func (k King) move(from, to Position, board Board) bool {
	return false
}

type Queen struct{}

func (q Queen) move(from, to Position, board Board) bool {
	return false
}

type Rook struct{}

func (r Rook) move(from, to Position, board Board) bool {
	return false
}

type Bishop struct{}

func (b Bishop) move(from, to Position, board Board) bool {
	return false
}

type Knight struct{}

func (k Knight) move(from, to Position, board Board) bool {
	return false
}

type Pawn struct{}

func (p Pawn) move(from, to Position, board Board) bool {
	return false
}

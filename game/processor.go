package game

var Field Board

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

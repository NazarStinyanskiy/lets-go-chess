package game

func replace(field *Board, from, to Position) {
	field.Cells[to] = field.Cells[from]
	field.Cells[from] = nil
}

func copyField(field Board) Board {
	var c = make(map[Position]*Figure, len(field.Cells))
	for figurePos, figure := range field.Cells {
		c[figurePos] = figure
	}
	return Board{Cells: c}
}

func isCellEmpty(field Board, pos Position) bool {
	return field.Cells[pos] == nil
}

func isFightingEnemy(field Board, from, to Position) bool {
	return field.Cells[to] != nil && field.Cells[to].IsWhite != field.Cells[from].IsWhite
}

func getDelta(from, to Position) (x int, y int) {
	return to.X - from.X, to.Y - from.Y
}

func findKing(field Board, isWhite bool) Position {
	for figurePos, figure := range field.Cells {
		if figure == nil {
			continue
		}
		_, ok := figure.Mover.(King)
		if ok && figure.IsWhite == isWhite {
			return figurePos
		}
	}
	return Position{}
}

func isFigureInThreat(field Board, pos Position, situation Situation) bool {
	for figurePos, figure := range field.Cells {
		if figure == nil || figurePos == pos {
			continue
		}
		canMove, _ := figure.canMove(field, figurePos, pos, situation)
		if canMove {
			return true
		}
	}
	return false
}

func analyzeSituation(field Board, playerIsWhite bool, situation Situation) Situation {
	enemyKingPos := findKing(field, !playerIsWhite)
	if isFigureInThreat(field, enemyKingPos, situation) {
		if isCheckmate(field, playerIsWhite, enemyKingPos, situation) {
			return Checkmate
		}
		return Check
	}
	if isStalemate(field, playerIsWhite, situation) {
		return Stalemate
	}
	return Continue
}

func findEnemyPositions(field Board, playerIsWhite bool) []Position {
	var enemyPositions []Position
	for position, figure := range field.Cells {
		if figure == nil || figure.IsWhite == playerIsWhite {
			continue
		}
		enemyPositions = append(enemyPositions, position)
	}
	return enemyPositions
}

func isCheckmate(field Board, playerIsWhite bool, enemyKingPos Position, situation Situation) bool {
	for _, enemyPosition := range findEnemyPositions(field, playerIsWhite) {
		for fieldPos := range field.Cells {
			if enemyPosition == fieldPos {
				continue
			}
			canMove, details := field.Cells[enemyPosition].canMove(field, enemyPosition, fieldPos, situation)
			if canMove {
				filedAfterMove := field.Cells[enemyPosition].move(field, enemyPosition, fieldPos, details)
				if !isFigureInThreat(filedAfterMove, enemyKingPos, situation) {
					return false
				}
			}
		}
	}
	return false
}

func isStalemate(field Board, playerIsWhite bool, situation Situation) bool {
	for _, enemyPosition := range findEnemyPositions(field, playerIsWhite) {
		for fieldPos := range field.Cells {
			if enemyPosition == fieldPos {
				continue
			}
			canMove, _ := field.Cells[enemyPosition].canMove(field, enemyPosition, fieldPos, situation)
			if canMove {
				return false
			}
		}
	}
	return true
}

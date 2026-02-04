package game

import (
	"errors"
	"testing"
)

func TestNextMoveSimple(t *testing.T) {
	g1 := StartGame()
	g2 := StartGame()
	g2.isWhiteMove = false
	tests := []struct {
		from, to   Position
		g          *Game
		eField     Board
		eSituation Situation
		eError     error
	}{
		{from: Position{1, 2}, to: Position{1, 3}, g: g1, eField: *whitePawnMovedField(g1), eSituation: Continue, eError: nil},
		{from: Position{1, 7}, to: Position{1, 6}, g: g2, eField: *blackPawnMovedField(g2), eSituation: Continue, eError: nil},
	}
	for _, test := range tests {
		situation, err := test.g.NextMove(test.from, test.to)
		isExpected, wrongPos, wrongFigure := isAllFiguresExpected(test.g.Field, test.eField)
		if isExpected && situation == test.eSituation && errors.Is(err, test.eError) {
			continue
		}
		if wrongFigure != nil {
			t.Errorf("NextMove(%v, %v) wrong figure in wrong place: %v, %v)", test.from, test.to, wrongPos, wrongFigure)
		} else {
			t.Errorf("NextMove(%v, %v) expected situation: %v, error: %v, got situation: %v, error: %v",
				test.from, test.to, test.eSituation, test.eError, situation, err)
		}
	}
}

func isAllFiguresExpected(actual, expected Board) (bool, Position, *Figure) {
	for position, figure := range expected.Cells {
		if actual.Cells[position] != figure {
			return false, position, figure
		}
	}
	return true, Position{}, nil
}

func whitePawnMovedField(g *Game) *Board {
	field := copyField(g.Field)
	field.Cells[Position{1, 3}] = field.Cells[Position{1, 2}]
	field.Cells[Position{1, 2}] = nil
	return &field
}

func blackPawnMovedField(g *Game) *Board {
	field := copyField(g.Field)
	field.Cells[Position{1, 6}] = field.Cells[Position{1, 7}]
	field.Cells[Position{1, 7}] = nil
	return &field
}

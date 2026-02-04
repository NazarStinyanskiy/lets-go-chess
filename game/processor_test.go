package game

import (
	"errors"
	"testing"
)

type TestCase struct {
	from, to   Position
	g          *Game
	eField     Board
	eSituation Situation
	eError     error
}

func TestNextMoveSimplePawnMove(t *testing.T) {
	var tests []TestCase
	for x := 1; x <= 8; x++ {
		g := StartGame()
		tests = append(tests, TestCase{from: Position{x, 2}, to: Position{x, 3}, g: g, eField: *whitePawnMovedField(g, x, 1), eSituation: Continue, eError: nil})
		g2 := StartGame()
		g2.isWhiteMove = false
		tests = append(tests, TestCase{from: Position{x, 7}, to: Position{x, 6}, g: g2, eField: *blackPawnMovedField(g2, x, -1), eSituation: Continue, eError: nil})
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

func TestNextMoveLongPawnMove(t *testing.T) {
	var tests []TestCase
	for x := 1; x <= 8; x++ {
		g := StartGame()
		tests = append(tests, TestCase{from: Position{x, 2}, to: Position{x, 4}, g: g, eField: *whitePawnMovedField(g, x, 2), eSituation: Continue, eError: nil})
		g2 := StartGame()
		g2.isWhiteMove = false
		tests = append(tests, TestCase{from: Position{x, 7}, to: Position{x, 5}, g: g2, eField: *blackPawnMovedField(g2, x, -2), eSituation: Continue, eError: nil})
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

func whitePawnMovedField(g *Game, x, deltaY int) *Board {
	field := copyField(g.Field)
	field.Cells[Position{x, 2 + deltaY}] = field.Cells[Position{x, 2}]
	field.Cells[Position{x, 2}] = nil
	return &field
}

func blackPawnMovedField(g *Game, x, deltaY int) *Board {
	field := copyField(g.Field)
	field.Cells[Position{x, 7 + deltaY}] = field.Cells[Position{x, 7}]
	field.Cells[Position{x, 7}] = nil
	return &field
}

func createCustomField(figures map[Position]*Figure) Board {
	field := Board{make(map[Position]*Figure)}
	for x := 1; x <= 8; x++ {
		for y := 1; y <= 8; y++ {
			pos := Position{x, y}
			if figures[pos] != nil {
				field.Cells[pos] = figures[pos]
			} else {
				field.Cells[pos] = nil
			}
		}
	}
	return field
}

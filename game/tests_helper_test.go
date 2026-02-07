package game

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type OneMoveTestCase struct {
	from, to   Position
	g          *Game
	eField     Board
	eSituation Situation
	eError     error
}

func test(tests []OneMoveTestCase, t *testing.T) {
	for id, test := range tests {
		t.Run(fmt.Sprintf("OneMoveTestCase_%d", id), func(t *testing.T) {
			fmt.Println("______________________________")
			t.Log("Starting test case")
			DrawConsoleBoard(test.g.Field)
			fmt.Println()
			situation, err := test.g.NextMove(test.from, test.to)
			DrawConsoleBoard(test.g.Field)
			isExpected, wrongPos, wrongFigure := isAllFiguresExpected(test.g.Field, test.eField)
			if isExpected && situation == test.eSituation && errors.Is(err, test.eError) {
				t.Log("Test case successfully finished")
				fmt.Println("______________________________")
				return
			}
			if wrongFigure != nil {
				t.Errorf("NextMove(%v, %v) wrong figure in wrong place: %v, %v\n)", test.from, test.to, wrongPos, wrongFigure)
			} else {
				t.Errorf("NextMove(%v, %v) expected situation: %v, error: %v, got situation: %v, error: %v\n",
					test.from, test.to, test.eSituation, test.eError, situation, err)
			}
			fmt.Println()
		})
	}
}

func isAllFiguresExpected(actual, expected Board) (bool, Position, *Figure) {
	for position, figure := range expected.Cells {
		if !isFiguresEqual(actual.Cells[position], figure) {
			return false, position, figure
		}
	}
	return true, Position{}, nil
}

func isFiguresEqual(actual, expected *Figure) bool {
	if actual == nil || expected == nil {
		return actual == expected
	}
	return actual.IsWhite == expected.IsWhite &&
		actual.HasMoved == expected.HasMoved &&
		actual.IsVulnerableForEnPassant == expected.IsVulnerableForEnPassant &&
		reflect.TypeOf(actual.Mover) == reflect.TypeOf(expected.Mover)
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

package game

import (
	"testing"
)

func TestNextMoveSimpleMovePawn(t *testing.T) {
	var tests []OneMoveTestCase
	for x := 1; x <= 8; x++ {
		g := StartGame()
		tests = append(tests, OneMoveTestCase{from: Position{x, 2}, to: Position{x, 3}, g: g, eField: *pawnMovedField(g, Position{x, 2}, Position{x, 3}), eSituation: Continue, eError: nil})
		g2 := StartGame()
		g2.IsWhiteMove = false
		tests = append(tests, OneMoveTestCase{from: Position{x, 7}, to: Position{x, 6}, g: g2, eField: *pawnMovedField(g2, Position{x, 7}, Position{x, 6}), eSituation: Continue, eError: nil})
	}
	test(tests, t)
}

func TestNextMoveLongMovePawn(t *testing.T) {
	var tests []OneMoveTestCase
	for x := 1; x <= 8; x++ {
		g := StartGame()
		tests = append(tests, OneMoveTestCase{from: Position{x, 2}, to: Position{x, 4}, g: g, eField: *pawnMovedField(g, Position{x, 2}, Position{x, 4}), eSituation: Continue, eError: nil})
		g2 := StartGame()
		g2.IsWhiteMove = false
		tests = append(tests, OneMoveTestCase{from: Position{x, 7}, to: Position{x, 5}, g: g2, eField: *pawnMovedField(g2, Position{x, 7}, Position{x, 5}), eSituation: Continue, eError: nil})
	}
	test(tests, t)
}

func TestNextMoveBeatPawn(t *testing.T) {
	var tests []OneMoveTestCase
	g := StartGame()
	g.Field = createCustomField(map[Position]*Figure{
		Position{1, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{5, 4}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{6, 5}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
	})
	g2 := StartGame()
	g2.IsWhiteMove = false
	g2.Field = createCustomField(map[Position]*Figure{
		Position{1, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{5, 4}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{6, 5}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
	})

	tests = append(tests, OneMoveTestCase{
		from:       Position{5, 4},
		to:         Position{6, 5},
		g:          g,
		eField:     *pawnMovedField(g, Position{5, 4}, Position{6, 5}),
		eSituation: Continue,
		eError:     nil,
	})
	tests = append(tests, OneMoveTestCase{
		from:       Position{6, 5},
		to:         Position{5, 4},
		g:          g2,
		eField:     *pawnMovedField(g2, Position{6, 5}, Position{5, 4}),
		eSituation: Continue,
		eError:     nil,
	})
	test(tests, t)
}

func TestNextMoveEnPassantPawn(t *testing.T) {
	g := StartGame()
	g.Field = createCustomField(map[Position]*Figure{
		Position{1, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{1, 2}: {IsWhite: true, HasMoved: false, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{2, 4}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
	})

	eFinalField := createCustomField(map[Position]*Figure{
		Position{1, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{1, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
	})

	tests := []OneMoveTestCase{
		{
			from:       Position{1, 2},
			to:         Position{1, 4},
			g:          g,
			eField:     *pawnMovedField(g, Position{1, 2}, Position{1, 4}),
			eSituation: Continue,
			eError:     nil,
		},
		{
			from:       Position{2, 4},
			to:         Position{1, 3},
			g:          g,
			eField:     eFinalField,
			eSituation: Continue,
			eError:     nil,
		},
	}
	test(tests, t)
}

func TestNextMoveImpossiblePawn(t *testing.T) {
	g := StartGame()
	g.Field = createCustomField(map[Position]*Figure{
		Position{1, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{1, 2}: {IsWhite: true, HasMoved: false, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{1, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
	})
	eField := copyField(g.Field)

	tests := []OneMoveTestCase{
		{
			from:       Position{1, 2},
			to:         Position{1, 3},
			g:          g,
			eField:     eField,
			eSituation: Continue,
			eError:     MoveRulesViolation,
		},
		{
			from:       Position{1, 2},
			to:         Position{1, 4},
			g:          g,
			eField:     eField,
			eSituation: Continue,
			eError:     MoveRulesViolation,
		},
	}
	test(tests, t)
}

func pawnMovedField(g *Game, from, to Position) *Board {
	field := copyField(g.Field)
	field.Cells[to] = field.Cells[from]
	field.Cells[from] = nil
	return &field
}

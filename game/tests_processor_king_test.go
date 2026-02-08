package game

import "testing"

func TestNextMoveSimpleKing(t *testing.T) {
	g := StartGame()
	g.Field = createCustomField(map[Position]*Figure{
		Position{2, 2}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
	})
	eXField := createCustomField(map[Position]*Figure{
		Position{3, 2}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
	})
	eYField := createCustomField(map[Position]*Figure{
		Position{3, 2}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 7}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
	})
	eXYField := createCustomField(map[Position]*Figure{
		Position{4, 3}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 7}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
	})

	tests := []OneMoveTestCase{
		{
			from:       Position{2, 2},
			to:         Position{3, 2},
			g:          g,
			eField:     eXField,
			eSituation: Continue,
			eError:     nil,
		},
		{
			from:       Position{8, 8},
			to:         Position{8, 7},
			g:          g,
			eField:     eYField,
			eSituation: Continue,
			eError:     nil,
		},
		{
			from:       Position{3, 2},
			to:         Position{4, 3},
			g:          g,
			eField:     eXYField,
			eSituation: Continue,
			eError:     nil,
		},
	}
	test(tests, t)
}

func TestNextMoveCheckKing(t *testing.T) {
	g := StartGame()
	g.IsWhiteMove = false
	g.Field = createCustomField(map[Position]*Figure{
		Position{2, 2}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{3, 4}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
	})
	eCheckField := createCustomField(map[Position]*Figure{
		Position{2, 2}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{3, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
	})

	tests := []OneMoveTestCase{
		{
			from:       Position{3, 4},
			to:         Position{3, 3},
			g:          g,
			eField:     eCheckField,
			eSituation: Check,
			eError:     nil,
		},
	}
	test(tests, t)
}

func TestNextMoveImpossibleMoveKing(t *testing.T) {
	g := StartGame()
	g.Field = createCustomField(map[Position]*Figure{
		Position{2, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{3, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
	})
	eField := createCustomField(map[Position]*Figure{
		Position{2, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{3, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
	})

	tests := []OneMoveTestCase{
		{
			from:       Position{2, 1},
			to:         Position{2, 2},
			g:          g,
			eField:     eField,
			eSituation: Continue,
			eError:     MoveRulesViolation,
		},
	}
	test(tests, t)
}

func TestNextMoveCheckmateKing(t *testing.T) {
	g := StartGame()
	g.IsWhiteMove = false
	g.Field = createCustomField(map[Position]*Figure{
		Position{2, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{3, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{1, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Rook{}},
		Position{8, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Queen{}},
	})
	eCheckmateField := createCustomField(map[Position]*Figure{
		Position{2, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{3, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{1, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Rook{}},
		Position{8, 1}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Queen{}},
	})

	tests := []OneMoveTestCase{
		{
			from:       Position{8, 8},
			to:         Position{8, 1},
			g:          g,
			eField:     eCheckmateField,
			eSituation: Checkmate,
			eError:     nil,
		},
	}
	test(tests, t)
}

func TestNextMoveStalemateKing(t *testing.T) {
	g := StartGame()
	g.IsWhiteMove = false
	g.Field = createCustomField(map[Position]*Figure{
		Position{1, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{4, 4}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{4, 5}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{1, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{1, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Rook{}},
	})
	eStalemateField := createCustomField(map[Position]*Figure{
		Position{1, 1}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{4, 4}: {IsWhite: true, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{4, 5}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Pawn{}},
		Position{1, 3}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: King{}},
		Position{2, 8}: {IsWhite: false, HasMoved: true, IsVulnerableForEnPassant: false, Mover: Rook{}},
	})

	tests := []OneMoveTestCase{
		{
			from:       Position{1, 8},
			to:         Position{2, 8},
			g:          g,
			eField:     eStalemateField,
			eSituation: Stalemate,
			eError:     nil,
		},
	}
	test(tests, t)
}

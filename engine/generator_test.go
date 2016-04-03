package engine

import "testing"

func TestGenerateMovesForDefaultBoardPosition(t *testing.T) {

	expected := []Move{
		Move{From: A2, To: A3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: A2, To: A4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: B2, To: B3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: B2, To: B4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: C2, To: C3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: C2, To: C4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: D2, To: D3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: D2, To: D4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: E2, To: E3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: E2, To: E4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: F2, To: F3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: F2, To: F4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: G2, To: G3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: G2, To: G4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: H2, To: H3, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: H2, To: H4, MovedPiece: WhitePawn, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: B1, To: A3, MovedPiece: WhiteKnight, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: B1, To: C3, MovedPiece: WhiteKnight, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: G1, To: F3, MovedPiece: WhiteKnight, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: G1, To: H3, MovedPiece: WhiteKnight, Special: moveOrdinary, Content: Empty, Promoted: Empty},
	}

	doTestMovesForFEN(defaultFEN, expected, t)

}

func TestGenerateMovesForCheckPositionToCreateAllEscapeMovesForKing(t *testing.T) {
	fen := "8/7p/1R2k1p1/3pp1P1/7P/7r/8/5K2 b - - 3 39"

	expected := []Move{
		Move{From: E6, To: E7, MovedPiece: BlackKing, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: E6, To: D7, MovedPiece: BlackKing, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: E6, To: F7, MovedPiece: BlackKing, Special: moveOrdinary, Content: Empty, Promoted: Empty},
		Move{From: E6, To: F5, MovedPiece: BlackKing, Special: moveOrdinary, Content: Empty, Promoted: Empty},
	}

	doTestMovesForFEN(fen, expected, t)
}

func TestCheckSimpleShouldBeFalseForStartingPosition(t *testing.T) {
	doTestCheckSimple(defaultFEN, false, t)
}

func TestCheckSimpleShouldBeFalseForNonCheckPosition(t *testing.T) {
	doTestCheckSimple("rnbqkbnr/3p1ppp/p3p3/1pp3B1/3PP3/2N2N2/PPP2PPP/R2QKB1R b KQkq - 1 5", false, t)
}

func TestCheckSimpleShouldBeTrueForCheckPositionFromRook(t *testing.T) {
	doTestCheckSimple("k7/8/8/8/8/8/8/RK6 b - - 0 2", true, t)
}

func TestAttackCheckSimpleForSingleRookOnFile(t *testing.T) {
	doTestCheckSimple("r7/8/8/8/8/8/8/K7 w - - 0 1", true, t)
}

func TestAttackCheckSimpleForSingleRookOnRank(t *testing.T) {
	doTestCheckSimple("r6K/8/8/8/8/8/8/8 w - - 0 1", true, t)
}

/* helper */

func doTestCheckSimple(fen string, expected bool, t *testing.T) {
	board, _ := parseFEN(fen)
	if NewGenerator(board).CheckSimple() != expected {
		t.Errorf("Expected CheckSimple() to be %t but was %t for board\n%s\n",
			expected, !expected, formatBoard(board))
	}
}

func doTestMovesForFEN(fen string, expected []Move, t *testing.T) {
	board, _ := parseFEN(fen)
	gen := NewGenerator(board)

	actual := gen.GenerateMoves()

	if len(actual) != len(expected) {
		t.Errorf("Expected %d moves but generated %d\n", len(expected), len(actual))
	}

	for _, e := range expected {
		if !contains(actual, e) {
			t.Errorf("Expected move %s was not created\n", e.String())
		}
	}
}

func contains(moves []Move, move Move) bool {

	for _, m := range moves {
		if m.From == move.From && m.To == move.To && m.MovedPiece == move.MovedPiece && m.Content == move.Content && m.Special == move.Special && m.Promoted == move.Promoted {
			return true
		}
	}

	return false
}

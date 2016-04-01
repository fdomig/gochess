package engine

import (
	"testing"
	"time"
)

func TestMateInOneFromRook(t *testing.T) {
	e := Move{From: G3, To: G8, MovedPiece: WhiteRook}
	doTestBestMoveForFEN("r3k3/2R5/4p2p/4Pp1P/8/5KR1/8/8 w - - 16 70", e, t)
}

func TestMateInOneFromQueen(t *testing.T) {
	e := Move{From: B6, To: B8, MovedPiece: WhiteQueen}
	doTestBestMoveForFEN("k7/P7/1Q6/8/8/8/8/K7 w - - 1 0", e, t)
}

func TestQueenPromotion(t *testing.T) {
	e := Move{From: A7, To: A8, MovedPiece: WhitePawn, Promoted: WhiteQueen, Special: movePromotion}
	doTestBestMoveForFEN("7k/P7/8/8/8/8/8/K7 w - - 1 0", e, t)
}

func doTestBestMoveForFEN(fen string, e Move, t *testing.T) {
	b := NewBoard(fen)

	searchVerbose = false
	searchMaxTime = 100 * time.Millisecond

	a := Search(b)

	if a.From != e.From || a.To != e.To || a.MovedPiece != e.MovedPiece ||
		a.Content != e.Content || a.Promoted != e.Promoted || a.Special != e.Special {

		t.Errorf("Expected %s but found %s\n%s\n", e.String(), a.String(), b.String())
	}
}

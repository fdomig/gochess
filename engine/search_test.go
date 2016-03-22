package engine

import "testing"

func TestBestMoveForEndgamePosition(t *testing.T) {

	board := NewBoard("4k3/2R5/4p2b/4Pb1P/8/5K2/8/8 w - - 16 70")

	best := Search(board)

	if best.From == C7 && best.To == G7 {
		t.Errorf("Expected best move to not be: %s", best.String())
	}

}

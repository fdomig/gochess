package engine

import "testing"

func TestEvaluateStartingPosition(t *testing.T) {
	doTestEvalForFEN(defaultFEN, t, 0)
}

func TestEvaluateOnePawnStartingPosition(t *testing.T) {
	doTestEvalForFEN("8/pppppppp/8/8/8/8/8/8 b - - 0 1", t, 810)
}

func doTestEvalForFEN(fen string, t *testing.T, e int) {
	b, _ := parseFEN(fen)

	a := Evaluate(b)

	if a != e {
		t.Errorf("Expected %d but got %d\n", e, a)
	}

}

package engine

const (
	pawnValue   = 100
	knightValue = 325
	bishopValue = 325
	rookValue   = 500
	queenValue  = 1050
	kingValue   = 40000

	evalPenaltyDoublePawn    = -8
	evalBonusCasteling       = 16
	evalBonusCheck           = 50
	evalBonusEndgamePawnMove = 50
	evalEndGameLevel         = 1500
	evalMateSearchLevel      = 600
	evalKingSafteyDivisor    = 3100

	scoreMate = 24000
	scoreDraw = 0
)

var (
	flipTable = []int{
		112, 113, 114, 115, 116, 117, 118, 119, 0, 0, 0, 0, 0, 0, 0, 0,
		96, 97, 98, 99, 100, 101, 102, 103, 0, 0, 0, 0, 0, 0, 0, 0,
		80, 81, 82, 83, 84, 85, 86, 87, 0, 0, 0, 0, 0, 0, 0, 0,
		64, 65, 66, 67, 68, 69, 70, 71, 0, 0, 0, 0, 0, 0, 0, 0,
		48, 49, 50, 51, 52, 53, 54, 55, 0, 0, 0, 0, 0, 0, 0, 0,
		32, 33, 34, 35, 36, 37, 38, 39, 0, 0, 0, 0, 0, 0, 0, 0,
		16, 17, 18, 19, 20, 21, 22, 23, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	pawnTable = []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, -20, -20, 10, 10, 5, 0, 0, 0, 0, 0, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 20, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		5, 5, 10, 25, 25, 10, 5, 5, 0, 0, 0, 0, 0, 0, 0, 0,
		10, 10, 20, 30, 30, 20, 10, 10, 0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	knightTable = []int{
		-50, -40, -30, -30, -30, -30, -40, -50, 0, 0, 0, 0, 0, 0, 0, 0,
		-40, -20, 0, 5, 5, 0, -20, -40, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, 0, 10, 15, 15, 10, 0, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, 5, 15, 20, 20, 15, 5, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, 0, 15, 20, 20, 15, 0, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, 5, 10, 15, 15, 10, 5, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-40, -20, 0, 0, 0, 0, -20, -40, 0, 0, 0, 0, 0, 0, 0, 0,
		-50, -40, -30, -30, -30, -30, -40, -50, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	bishopTable = []int{
		-20, -10, -10, -10, -10, -10, -10, -20, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 5, 0, 0, 0, 0, 5, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 10, 10, 10, 10, 10, 10, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 10, 10, 10, 10, 0, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 5, 5, 10, 10, 5, 5, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 5, 10, 10, 5, 0, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 0, 0, 0, 0, 0, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-20, -10, -10, -10, -10, -10, -10, -20, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	rookTable = []int{
		0, 0, 0, 5, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5, 0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5, 0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5, 0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5, 0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5, 0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	queenTable = []int{
		-20, -10, -10, -5, -5, -10, -10, -20, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 0, 0, 0, 5, 0, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 5, 5, 5, 5, 5, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 5, 5, 5, 5, 0, -5, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 5, 5, 5, 5, 0, -5, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 5, 5, 5, 5, 0, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 0, 0, 0, 0, 0, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-20, -10, -10, -5, -5, -10, -10, -20, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	kingTableMiddle = []int{
		20, 30, 10, 0, 0, 10, 30, 20, 0, 0, 0, 0, 0, 0, 0, 0,
		20, 20, 0, 0, 0, 0, 20, 20, 0, 0, 0, 0, 0, 0, 0, 0,
		-10, -20, -20, -20, -20, -20, -20, -10, 0, 0, 0, 0, 0, 0, 0, 0,
		-20, -30, -30, -40, -40, -30, -30, -20, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -40, -40, -50, -50, -40, -40, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -40, -40, -50, -50, -40, -40, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -40, -40, -50, -50, -40, -40, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -40, -40, -50, -50, -40, -40, -30, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	kingTableEnd = []int{
		-50, -30, -30, -30, -30, -30, -30, -50, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -30, 0, 0, 0, 0, -30, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -10, 20, 30, 30, 20, -10, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -10, 30, 40, 40, 30, -10, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -10, 30, 40, 40, 30, -10, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -10, 20, 30, 30, 20, -10, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-30, -20, -10, 0, 0, -10, -20, -30, 0, 0, 0, 0, 0, 0, 0, 0,
		-50, -40, -30, -20, -20, -30, -40, -50, 0, 0, 0, 0, 0, 0, 0, 0,
	}
)

// Evaluate the score of a given board
func Evaluate(b *Board) int {

	scoreWhite := 0
	scoreBlack := 0

	materialWhite := 0
	materialBlack := 0

	for rank := int8(0); rank < size; rank++ {
		for file := int8(0); file < size; file++ {
			sq := square(rank, file)
			switch b.data[sq] {
			case WhitePawn:
				materialWhite += pawnValue
				scoreWhite += evaluatePawn(b, sq)
			case WhiteKnight:
				materialWhite += knightValue
				scoreWhite += evaluateKnight(b, sq)
			case WhiteBishop:
				materialWhite += bishopValue
				scoreWhite += evaluateBishop(b, sq)
			case WhiteRook:
				materialWhite += rookValue
				scoreWhite += evaluateRook(b, sq)
			case WhiteQueen:
				materialWhite += queenValue
				scoreWhite += evaluateQueen(b, sq)
			case BlackPawn:
				materialBlack += pawnValue
				scoreBlack += evaluatePawn(b, sq)
			case BlackKnight:
				materialBlack += knightValue
				scoreBlack += evaluateKnight(b, sq)
			case BlackBishop:
				materialBlack += bishopValue
				scoreBlack += evaluateBishop(b, sq)
			case BlackRook:
				materialBlack += rookValue
				scoreBlack += evaluateRook(b, sq)
			case BlackQueen:
				materialBlack += queenValue
				scoreBlack += evaluateQueen(b, sq)
			}
		}
	}

	// mate level?
	if materialWhite <= evalMateSearchLevel || materialBlack <= evalMateSearchLevel {
		generator := NewGenerator(b)

		if generator.CheckSimple() {
			// if white is to move, black just made a check move
			if b.sideToMove == White {
				scoreBlack += evalBonusCheck
			} else {
				scoreWhite += evalBonusCheck
			}
		}
	}

	// evaluate kings
	scoreWhite += evaluateKing(b, int8(b.whiteKingPosition), materialWhite, materialBlack)
	scoreBlack += evaluateKing(b, int8(b.blackKingPosition), materialWhite, materialBlack)

	// special moves

	// casteling bonus

	scoreWhite += materialWhite
	scoreBlack += materialBlack

	return int(b.sideToMove) * (scoreWhite - scoreBlack)
}

func evaluateKing(b *Board, sq int8, materialWhite int, materialBlack int) int {
	score := 0

	// White
	if b.data[sq] > 0 {
		if materialWhite > evalEndGameLevel {
			score += kingTableMiddle[sq]
		} else {
			score += kingTableEnd[sq]
		}
		// king saftey
		score *= materialBlack
		score /= evalKingSafteyDivisor

	} else {
		if materialBlack > evalEndGameLevel {
			score += kingTableMiddle[flipTable[sq]]
		} else {
			score += kingTableEnd[flipTable[sq]]
		}
		// king saftey
		score *= materialWhite
		score /= evalKingSafteyDivisor
	}

	return score
}

func evaluatePawn(b *Board, sq int8) int {

	if b.data[sq] > 0 {
		score := pawnTable[sq]
		if sq >= nextRank && b.data[sq-nextRank] == WhitePawn {
			score += evalPenaltyDoublePawn
		}

		return score
	}

	score := pawnTable[flipTable[sq]]
	if sq+nextRank < boardSize && sq+nextRank >= 0 && b.data[sq+nextRank] == BlackPawn {
		score += evalPenaltyDoublePawn
	}

	return score
}

func evaluateKnight(b *Board, sq int8) int {
	if b.data[sq] > 0 {
		return knightTable[sq]
	}
	return knightTable[flipTable[sq]]
}

func evaluateBishop(b *Board, sq int8) int {
	if b.data[sq] > 0 {
		return bishopTable[sq]
	}
	return bishopTable[flipTable[sq]]
}

func evaluateRook(b *Board, sq int8) int {
	if b.data[sq] > 0 {
		return rookTable[sq]
	}
	return rookTable[flipTable[sq]]
}

func evaluateQueen(b *Board, sq int8) int {
	if b.data[sq] > 0 {
		return queenTable[sq]
	}
	return queenTable[flipTable[sq]]
}

package engine

import "time"

var (
	searchVerbose = true
	searchMaxTime = 16 * time.Second
)

const (
	searchMaxDepth  = 20
	searchMaxPly    = 128
	searchEvalStart = 50000
)

type pvSearch struct {
	board         *Board
	checkedNodes  int64
	path          [searchMaxPly][searchMaxPly]Move
	pathLength    [searchMaxPly]int
	bestMoves     [searchMaxDepth]Move
	bestMovesPlys [searchMaxDepth]int
	bestScores    [searchMaxDepth]int
	stopByTime    bool
	stopTime      time.Time
	followPv      bool
	ply           int
}

// Search finds the best available move
func Search(board *Board) Move {

	// TODO book

	startTime := time.Now()

	pv := pvSearch{}
	pv.stopTime = startTime.Add(searchMaxTime)
	pv.board = &Board{}
	*pv.board = *board
	pv.board.ply = 0

	printSearchHead()

	foundMate := false
	depth := 1

	for ; depth < searchMaxDepth && !pv.stopByTime; depth++ {
		pv.followPv = true
		score := pv.alphaBeta(depth, -searchEvalStart, searchEvalStart)

		pv.bestMoves[depth] = pv.path[0][0]
		pv.bestMovesPlys[depth] = pv.pathLength[0]
		pv.bestScores[depth] = score

		printSearchLevel(&pv, depth, score, startTime)

		if score >= scoreMate || score <= -scoreMate {
			foundMate = true
			break
		}

	}

	best := Move{}

	if pv.stopByTime {
		best = pv.bestMoves[depth-2]

	} else if foundMate {
		best = pv.bestMoves[depth]

	} else {
		best = pv.bestMoves[depth-1]
	}

	printSearchResult(&pv, startTime)

	return best
}

func (pv *pvSearch) alphaBeta(depth, alpha, beta int) int {
	if depth == 0 {
		return pv.quiescence(alpha, beta)
	}
	pv.checkedNodes++

	// check time all 4096 nodes
	if pv.checkedNodes%4095 == 0 {
		if time.Now().After(pv.stopTime) {
			pv.stopByTime = true
			return 0
		}
	}

	// TODO: index out of range
	if len(pv.pathLength) <= pv.board.ply {
		return 0
	}
	pv.pathLength[pv.board.ply] = pv.board.ply

	generator := Generator{board: pv.board}
	moves := generator.GenerateMoves()

	if generator.kingUnderCheck {
		depth++
	}

	// repetition
	if pv.board.ply > 0 && pv.board.repetitions() >= 3 {
		return scoreDraw
	}

	if pv.followPv {
		moves = pv.sortPv(moves)
	}

	playedMove := false
	score := 0
	pvSearch := true

	for _, move := range moves {
		pv.board.MakeMove(move)
		playedMove = true

		if pvSearch {
			score = -pv.alphaBeta(depth-1, -beta, -alpha)
		} else {
			score = -pv.alphaBeta(depth-1, -alpha-1, -alpha)
			if score > alpha && score < beta {
				score = -pv.alphaBeta(depth-1, -beta, -alpha)
			}
		}
		pv.board.UndoMove()

		if pv.stopByTime {
			return 0
		}

		if score > alpha {
			if score >= beta {
				return score
			}
			alpha = score
			pvSearch = false

			pv.path[pv.board.ply][pv.board.ply] = move
			for j := pv.board.ply + 1; j < pv.pathLength[pv.board.ply+1]; j++ {
				pv.path[pv.board.ply][j] = pv.path[pv.board.ply+1][j]
			}
			pv.pathLength[pv.board.ply] = pv.pathLength[pv.board.ply+1]
		}

	}

	if !playedMove {
		if generator.kingUnderCheck {
			return -(scoreMate + pv.board.ply)
		}
		return scoreDraw
	}

	// fifty moves rule
	if pv.board.halfMoveClock >= 100 {
		return scoreDraw
	}

	return alpha
}

func (pv *pvSearch) quiescence(alpha, beta int) int {

	pv.checkedNodes++

	// check time all 4096 nodes
	if pv.checkedNodes%4095 == 0 {
		if time.Now().After(pv.stopTime) {
			pv.stopByTime = true
			return 0
		}
	}

	pv.pathLength[pv.board.ply] = pv.board.ply

	eval := Evaluate(pv.board)

	if eval >= beta {
		return beta
	}

	if eval > alpha {
		alpha = eval
	}

	generator := Generator{board: pv.board}

	for _, move := range generator.GenerateMoves() {

		// only check capture moves
		// TODO: should be optimized from the generator!
		if move.Content == Empty {
			continue
		}

		pv.board.MakeMove(move)
		score := -pv.quiescence(-beta, -alpha)
		pv.board.UndoMove()
		if score > alpha {
			if score >= beta {
				return beta
			}
			alpha = score

			// store new, better alpha node in the path
			pv.path[pv.board.ply][pv.board.ply] = move
			for j := pv.board.ply + 1; j < pv.pathLength[pv.board.ply+1]; j++ {
				pv.path[pv.board.ply][j] = pv.path[pv.board.ply+1][j]
			}
			pv.pathLength[pv.board.ply] = pv.pathLength[pv.board.ply+1]
		}
	}

	return alpha
}

func (pv *pvSearch) sortPv(moves []Move) []Move {
	pv.followPv = false
	for i := 0; i < len(moves); i++ {
		if moves[i] == pv.path[0][pv.board.ply] {
			pv.followPv = true
			tmp := moves[0]
			moves[0] = moves[i]
			moves[i] = tmp
			return moves
		}
	}
	return moves
}

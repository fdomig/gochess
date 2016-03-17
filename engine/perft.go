package engine

import "time"

var (
	position1FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	position1Table = []PerftData{
		PerftData{depth: 0, nodes: 1, captures: 0, enPassants: 0, castles: 0, promotions: 0, checks: 0, mates: 0},
		PerftData{depth: 1, nodes: 20, captures: 0, enPassants: 0, castles: 0, promotions: 0, checks: 0, mates: 0},
		PerftData{depth: 2, nodes: 400, captures: 0, enPassants: 0, castles: 0, promotions: 0, checks: 0, mates: 0},
		PerftData{depth: 3, nodes: 8902, captures: 34, enPassants: 0, castles: 0, promotions: 0, checks: 12, mates: 0},
		PerftData{depth: 4, nodes: 197281, captures: 1576, enPassants: 0, castles: 0, promotions: 0, checks: 469, mates: 8},
		PerftData{depth: 5, nodes: 4865609, captures: 82719, enPassants: 258, castles: 0, promotions: 0, checks: 27351, mates: 347},
		PerftData{depth: 6, nodes: 119060324, captures: 2812008, enPassants: 5248, castles: 0, promotions: 0, checks: 809099, mates: 10828},
	}

	position2FEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -"

	position2Table = []PerftData{
		PerftData{depth: 0, nodes: 1},
		PerftData{depth: 1, nodes: 48, captures: 8, enPassants: 0, castles: 2, promotions: 0, checks: 0, mates: 0},
		PerftData{depth: 2, nodes: 2039, captures: 351, enPassants: 1, castles: 91, promotions: 0, checks: 3, mates: 0},
		PerftData{depth: 3, nodes: 97862, captures: 17102, enPassants: 45, castles: 3162, promotions: 0, checks: 993, mates: 1},
		PerftData{depth: 4, nodes: 4085603, captures: 757163, enPassants: 1929, castles: 128013, promotions: 15172, checks: 25523, mates: 43},
		PerftData{depth: 5, nodes: 193690690, captures: 35043416, enPassants: 73365, castles: 4993637, promotions: 8392, checks: 3309887, mates: 30171},
	}
)

// PerftData aggregates performance test data in a structure
type PerftData struct {
	depth      int
	nodes      int64
	captures   int64
	enPassants int64
	castles    int64
	promotions int64
	checks     int64
	mates      int64
	elapsed    time.Duration
}

// Perft runs a performance test against a given FEN and expected results
func Perft(fen string, expected []PerftData) {
	printPerftData(NewBoard(fen), expected)
}

func perft(depth int, board *Board) PerftData {

	data := PerftData{depth: depth}
	generator := NewGenerator(board)

	start := time.Now()

	if depth == 0 {
		data.depth = 0
		data.nodes = 1
		return data
	}

	moves := generator.GenerateMoves()

	if len(moves) == 0 {
		data.mates++
	}

	if generator.kingUnderCheck {
		data.checks++
	}

	for _, move := range moves {
		board.MakeMove(move)

		res := perft(depth-1, board)
		data.nodes += res.nodes
		data.captures += res.captures
		data.enPassants += res.enPassants
		data.castles += res.castles
		data.promotions += res.promotions
		data.checks += res.checks
		data.mates += res.mates

		switch move.Special {
		case moveCastelingShort:
			data.castles++
		case moveCastelingLong:
			data.castles++
		case movePromotion:
			data.promotions++
		case moveEnPassant:
			data.enPassants++
		}

		if move.Content != Empty {
			data.captures++
		}

		switch board.status {
		case statusCheck:
			data.checks++
		case statusBlackMates:
			data.mates++
		case statusWhiteMates:
			data.mates++
		}

		board.UndoMove()
	}

	data.elapsed = time.Since(start)

	return data
}

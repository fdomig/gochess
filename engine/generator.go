package engine

import "fmt"

var (
	nextRank int8 = 16
	nextFile int8 = 1

	moveUp        = nextRank
	moveDown      = -nextRank
	moveLeft      = -nextFile
	moveRight     = nextFile
	moveUpLeft    = moveUp + moveLeft
	moveUpRight   = moveUp + moveRight
	moveDownLeft  = moveDown + moveLeft
	moveDownRight = moveDown + moveRight

	deltaAll    = []int8{moveUp, moveDown, moveLeft, moveRight, moveUpLeft, moveUpRight, moveDownLeft, moveDownRight}
	deltaKnight = []int8{moveUp*2 + moveLeft, moveUp*2 + moveRight, moveRight*2 + moveUp, moveRight*2 + moveDown, moveDown*2 + moveRight, moveDown*2 + moveLeft, moveLeft*2 + moveUp, moveLeft*2 + moveDown}
	deltaRook   = []int8{moveUp, moveDown, moveLeft, moveRight}
	deltaBishop = []int8{moveUpLeft, moveUpRight, moveDownLeft, moveDownRight}
	deltaQueen  = deltaAll
	deltaKing   = deltaAll

	deltaWhitePawn = []int8{moveUp, moveUpLeft, moveUpRight}       // first move for pawns has to be forward!
	deltaBlackPawn = []int8{moveDown, moveDownLeft, moveDownRight} // first move for pawns has to be forward!

	whitePawnStartPos int8 = 1 // rank 2
	blackPawnStartPos int8 = 6 // rank 7

	castleShortDistanceRook int8 = 3
	castleLongDistanceRook  int8 = 4

	whiteKingStartSquare = E1
	blackKingStartSquare = E8
	whiteRookShortSquare = A1
	whiteRookLongSquare  = H1
	blackRookShortSquare = A8
	blackRookLongSquare  = H8
)

// squares that need to be empty and not under check for castling
var (
	shortWhiteSquares = []Square{F1, G1}
	longWhiteSquares  = []Square{B1, C1, D1}
	shortBlackSquares = []Square{F8, G8}
	longBlackSquares  = []Square{B8, C8, D8}
)

const (
	boardSize int8 = 120
	size      int8 = 8
)

// Generator creates possible moves for a given board position
type Generator struct {
	board                  *Board
	legalEnding            []bool
	legalDelta             []int8
	lastMoveSquare         Square
	moves                  []Move
	kingSquare             int8
	kingUnderCheck         bool
	kingUnderCheckByKnight int8
	isCheckMate            bool
}

// NewGenerator creates a new generator for a given board
func NewGenerator(board *Board) *Generator {
	g := new(Generator)
	g.board = board

	return g
}

// GenerateMoves creates a list of possible moves
func (g *Generator) GenerateMoves() []Move {

	g.reset()

	// generation

	threats := g.numberOfCheckThreats()

	g.generateKingMoves(g.kingSquare)

	// only legal move is moving the king ... else we have a mate
	if threats > 1 {
		g.sortMoves()
		return g.moves
	}

	if threats == 1 && g.kingUnderCheckByKnight != Empty {
		g.generateCaptureMovesForOpponentKnight(g.kingUnderCheckByKnight)
		g.sortMoves()
		return g.moves
	}

	if threats == 0 {
		g.generateCastlingMoves()
	}

	// do generation

	for rank := int8(0); rank < size; rank++ {
		for file := int8(0); file < size; file++ {
			square := square(rank, file)
			piece := g.board.data[square] * g.board.sideToMove

			if piece > 0 {
				switch piece {
				case Pawn:
					g.generateMovesPawn(square)
				case Rook:
					g.generateGenericMoves(square, deltaRook, false)
				case Bishop:
					g.generateGenericMoves(square, deltaBishop, false)
				case Queen:
					g.generateGenericMoves(square, deltaQueen, false)
				case Knight:
					g.generateGenericMoves(square, deltaKnight, true)
				}
			}

		}
	}

	g.sortMoves()
	return g.moves
}

// CheckSimple finds possible check attacks
func (g *Generator) CheckSimple() bool {

	// check for knight attacks first
	for _, delta := range deltaKnight {
		to := g.kingSquare + int8(delta)
		if g.board.legalSquare(to) && g.board.data[to]*opponent(g.board.sideToMove) == Knight {
			return true
		}
	}

	// check for all other pieces
	for _, delta := range deltaAll {
		depth := 0
		for to := g.kingSquare + delta; g.board.legalSquare(to); to += delta {
			found := g.board.data[to] * opponent(g.board.sideToMove)

			// empty square, go ahead
			if found == Empty {
				continue
			}

			if found == Pawn || found == King {
				if depth == 1 {
					if g.attackPosible(to, opposite(delta)) {
						return true
					}
				}
				break
			}

			// stop direction if own piece is found
			if found < 0 {
				break
			}

			// found opponent
			if g.attackPosible(to, delta) {
				return true
			}
			break
		}
	}

	return false
}

func (g *Generator) reset() {
	g.moves = make([]Move, 0, 48)

	g.kingUnderCheck = false
	g.kingUnderCheckByKnight = Empty

	g.kingSquare = int8(g.board.whiteKingPosition)
	if g.board.sideToMove == Black {
		g.kingSquare = int8(g.board.blackKingPosition)
	}

	g.legalEnding = make([]bool, boardSize)
	g.legalDelta = make([]int8, boardSize)

	if len(g.board.history) > 0 {
		g.lastMoveSquare = g.board.history[len(g.board.history)-1].move.To
	}
}

func (g *Generator) sortMoves() {
	lastCapture := make([]Move, 0, len(g.moves))
	sorted := make([]Move, 0, len(g.moves))
	captured := make([]Move, 0, len(g.moves))
	promotions := make([]Move, 0, len(g.moves))
	castelings := make([]Move, 0, len(g.moves))
	ordinary := make([]Move, 0, len(g.moves))

	for _, move := range g.moves {
		switch {
		case move.To == g.lastMoveSquare:
			lastCapture = append(lastCapture, move)
		case move.Content != Empty:
			captured = append(captured, move)
		case move.Special == movePromotion:
			promotions = append(promotions, move)
		case move.Special == moveCastelingShort || move.Special == moveCastelingLong:
			castelings = append(castelings, move)
		default:
			ordinary = append(ordinary, move)
		}
	}

	// 1. last moved piece capture
	sorted = append(sorted, lastCapture...)

	// 2. capture moves
	sorted = append(sorted, captured...)

	// 3.promotion moves
	sorted = append(sorted, promotions...)

	// 4.castling moves
	sorted = append(sorted, castelings...)

	// 5. normal moves
	sorted = append(sorted, ordinary...)

	if len(sorted) != len(g.moves) {
		fmt.Printf("len(g.moves): %d, len(sorted): %d\n", len(g.moves), len(sorted))
		panic("Sorted moves are missing moves!")
	}

	g.moves = sorted
}

func (g *Generator) addMove(move Move) {
	g.moves = append(g.moves, move)
}

func (g *Generator) generateKingMoves(square int8) {
	to := int8(0)
	for _, delta := range deltaKing {
		to = square + delta
		if g.board.legalSquare(to) {
			move := g.createMove(square, to)
			if move.MovedPiece == Empty {
				continue
			}

			if len(g.findThreats(Square(to), g.board.sideToMove, true)) == 0 {
				g.addMove(move)
			}
		}
	}

}

func (g *Generator) generateCaptureMovesForOpponentKnight(square int8) {
	threats := g.findThreats(Square(square), opponent(g.board.sideToMove), false)

	if len(threats) >= 0 {
		for _, threat := range threats {
			g.addMove(g.createMove(threat, square))
		}
	}
}

func (g *Generator) generateGenericMoves(from int8, delta []int8, singleStep bool) {
	for _, d := range delta {
		for to := from + d; g.board.legalSquare(to); to += d {
			move := g.createMove(from, to)

			// we are facing the same color; skip
			if move.MovedPiece == Empty {
				break
			}

			if g.legalDelta[from] == 0 || (g.legalDelta[from] == d || g.legalDelta[from] == opposite(d)) {
				if !g.kingUnderCheck || g.legalEnding[to] {
					g.addMove(move)
				}
			} else {
				break
			}

			// this move was single step only or we captured a piece; no need to go further then
			if singleStep || move.Content != Empty {
				break
			}
		}
	}
}

func (g *Generator) createMove(from, to int8) Move {
	move := Move{From: Square(from), To: Square(to)}
	move.Promoted = Empty
	move.Content = g.board.data[to]

	// piece of same color on to square
	if g.board.data[to]*g.board.data[from] > 0 {
		move.MovedPiece = Empty
		return move
	}

	move.MovedPiece = g.board.data[from]
	move.From = Square(from)
	move.To = Square(to)

	switch {
	case move.MovedPiece == WhiteKing || move.MovedPiece == WhiteRook:
		move.Promoted = g.board.whiteCastle
	case move.MovedPiece == BlackKing || move.MovedPiece == BlackRook:
		move.Promoted = g.board.blackCastle
	}

	return move
}

func (g *Generator) generateMovesPawn(from int8) {

	startPos := false
	delta := deltaWhitePawn

	switch g.board.data[from] {
	case WhitePawn:
		startPos = rank(from) == whitePawnStartPos
	case BlackPawn:
		startPos = rank(from) == blackPawnStartPos
		delta = deltaBlackPawn
	}

	for _, d := range delta {
		to := from + d

		if g.board.legalSquare(to) {
			// first check if it's has the direction to go and then if it is a legal ending
			if g.legalDelta[from] == 0 || (g.legalDelta[from] == d || g.legalDelta[from] == opposite(d)) {
				g.generateSingleMovePawn(from, to, startPos)
			}
		}
	}
}

func (g *Generator) generateSingleMovePawn(from int8, to int8, startingPos bool) {
	move := Move{}

	move.MovedPiece = g.board.data[from]
	move.Content = g.board.data[to]
	move.From = Square(from)
	move.To = Square(to)
	move.Promoted = Empty

	enPassantRemovesThreat := false

	if g.kingUnderCheck && g.board.enPassant == Square(to) {
		enPassantRemovesThreat = g.legalEnding[to+(moveDown*g.board.sideToMove)]
	}

	if (!g.kingUnderCheck || g.legalEnding[to]) || enPassantRemovesThreat {
		// move left/right
		if file(from) != file(to) {

			// en passant?
			if move.To == g.board.enPassant {
				move.Special = moveEnPassant
				move.Content = -move.MovedPiece

			} else if g.board.data[from]*g.board.data[to] >= 0 {
				// must be opposite pawn
				return
			}

		} else if g.board.data[to] != Empty {
			// moving forward requires an empty field
			return
		}

		// promotions (queen or rook are only relevant pieces to check)
		if rank(to)%7 == 0 {
			move.Promoted = move.MovedPiece * Queen
			move.Special = movePromotion
			g.addMove(move)

			move.Promoted = move.MovedPiece * Rook
			move.Special = movePromotion
			g.addMove(move)

		} else {
			g.addMove(move)
		}
	}

	// starting position; do two steps if on same file; do not care about direction
	if startingPos && file(from) == file(to) && g.board.data[to] == Empty {
		move.To = Square(to + to - from)
		if g.board.data[move.To] == Empty && (!g.kingUnderCheck || g.legalEnding[move.To]) {
			g.addMove(move)
		}
	}

}

func (g *Generator) generateCastlingMoves() {
	// assume king is not under check
	switch g.board.sideToMove {
	case White:
		if g.canCastle(g.board.sideToMove, castleShort) {
			g.addMove(Move{From: E1, To: G1, Content: Empty, MovedPiece: WhiteKing, Special: moveCastelingShort})
		}
		if g.canCastle(g.board.sideToMove, castleLong) {
			g.addMove(Move{From: E1, To: C1, Content: Empty, MovedPiece: WhiteKing, Special: moveCastelingLong})
		}
	case Black:
		if g.canCastle(g.board.sideToMove, castleShort) {
			g.addMove(Move{From: E8, To: G8, Content: Empty, MovedPiece: BlackKing, Special: moveCastelingShort})
		}
		if g.canCastle(g.board.sideToMove, castleLong) {
			g.addMove(Move{From: E8, To: C8, Content: Empty, MovedPiece: BlackKing, Special: moveCastelingLong})
		}
	}
}

func (g *Generator) canCastle(color int8, dir int8) bool {
	switch color {
	case White:
		if g.board.whiteCastle&dir == dir {
			if dir == castleShort {
				return g.board.isEmpty(F1, G1) && len(g.findThreats(F1, color, false)) == 0 && len(g.findThreats(G1, color, false)) == 0
			}
			return g.board.isEmpty(B1, C1, D1) && len(g.findThreats(C1, color, false)) == 0 && len(g.findThreats(D1, color, false)) == 0
		}
	case Black:
		if g.board.blackCastle&dir == dir {
			if dir == castleShort {
				return g.board.isEmpty(F8, G8) && len(g.findThreats(F8, color, false)) == 0 && len(g.findThreats(G8, color, false)) == 0
			}
			return g.board.isEmpty(B8, C8, D8) && len(g.findThreats(C8, color, false)) == 0 && len(g.findThreats(D8, color, false)) == 0
		}
	}

	return false
}

func (g *Generator) findThreats(square Square, sideToMove int8, skipKing bool) []int8 {

	threats := []int8{}

	// knights need some special treatment
	for _, delta := range deltaKnight {
		to := int8(square) + delta
		if g.board.legalSquare(to) && g.board.data[to]*opponent(sideToMove) == Knight {
			threats = append(threats, to)
		}
	}

	// check all possible directions
	for _, delta := range deltaAll {
		depth := 0
		for to := int8(square) + delta; g.board.legalSquare(to); to += delta {
			content := g.board.data[to] * sideToMove

			// if there is a pawn/king on next square, check for attack
			if content == Pawn || content == King {
				if depth == 1 {
					if g.attackPosible(to, opposite(delta)) {
						threats = append(threats, to)
					}
				}
				break
			}

			// skip all empty squares
			if content == Empty {
				continue
			}

			if skipKing && content == King {
				continue
			}

			// own piece
			if content > 0 {
				break
			}

			// found opponent piece
			if g.attackPosible(to, delta) {
				threats = append(threats, to)
			}

			// no point in checking further
			break
		}
	}

	return threats
}

// find the threats for the current color's king possition
func (g *Generator) numberOfCheckThreats() int {
	color := g.board.sideToMove
	threats := 0
	depth := 0
	squareOfGuardingPiece := Invalid

	// check if there is a threat from an opponents knight (could only be one)
	for _, delta := range deltaKnight {
		newSquare := g.kingSquare + delta

		if g.board.legalSquare(newSquare) && g.board.data[newSquare]*opponent(color) == Knight {
			g.kingUnderCheck = true
			g.kingUnderCheckByKnight = int8(newSquare)
			threats++
			break
		}
	}

	// checking all possible directions
	for _, delta := range deltaAll {
		depth = 0
		squareOfGuardingPiece = Invalid

		for newSquare := g.kingSquare + delta; g.board.legalSquare(newSquare); newSquare += delta {
			depth++

			// check found piece
			// positive: opponent
			// zero:     empty
			// negative: same color
			piece := g.board.data[newSquare] * opponent(color)

			if piece == Empty {
				continue
			}

			// pawn on the next square, check if attack possible
			if piece == Pawn || piece == King {
				if depth == 1 {
					if g.attackPosible(newSquare, opposite(delta)) {
						g.kingUnderCheck = true
						threats++
						g.legalEnding[newSquare] = true
					}
				}
				break
			}

			// own piece -> we stop but might have a protecting piece
			if piece < 0 {
				if squareOfGuardingPiece == Invalid {
					squareOfGuardingPiece = Square(newSquare)
					continue
				} else {
					break
				}
			}

			// we found an opponents piece; might be a threat for the current square
			if g.attackPosible(newSquare, delta) {
				if squareOfGuardingPiece == Invalid {
					threats++
					g.kingUnderCheck = true
					g.legalEnding[newSquare] = true

					// set legal endings; looks hack-ish: it is!
					for i := g.kingSquare + delta; i != newSquare; i += delta {
						g.legalEnding[i] = true
					}
				} else {
					g.legalDelta[uint8(squareOfGuardingPiece)] = delta
				}
			}
		}
	}

	return threats
}

// check wether a piece can move towards a specific direction
// does not work for knight obviously
func (g *Generator) hasDirection(delta []int8, direction int8) bool {
	for _, d := range delta {
		if d == direction {
			return true
		}
	}

	return false
}

// check for attacks from a diven square in a particular direction
// does not work for knight obviously
func (g *Generator) attackPosible(from int8, direction int8) bool {
	absPiece := g.board.data[from]
	if absPiece < 0 {
		absPiece = -absPiece
	}

	switch absPiece {
	case King:
		return g.hasDirection(deltaKing, direction)
	case Queen:
		return g.hasDirection(deltaQueen, direction)
	case Rook:
		return g.hasDirection(deltaRook, direction)
	case Bishop:
		return g.hasDirection(deltaBishop, direction)
	case Pawn:
		// remove the first move of pawns (move foreward = no attack move!)
		switch g.board.data[from] {
		case WhitePawn:
			return g.hasDirection(deltaWhitePawn[1:], direction)
		case BlackPawn:
			return g.hasDirection(deltaBlackPawn[1:], direction)
		}
	}

	return false
}

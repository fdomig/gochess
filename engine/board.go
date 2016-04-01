package engine

import "fmt"

const (
	statusNormal     = 0
	statusCheck      = 1
	statusWhiteMates = 2
	statusBlackMates = 3
	statusStaleMate  = 4
	statusDraw       = 5
	statusWhiteWins  = 6
	statusBlackWins  = 7
)

// HistoryItem chess move and flags
type HistoryItem struct {
	move          Move
	whiteCastle   int8
	blackCastle   int8
	enPassant     Square
	halfMoveClock int
	hash          int64
}

// Board represents a chessboard
type Board struct {
	history           []HistoryItem
	halfMoveClock     int
	fullMoves         int
	ply               int
	data              [boardSize]int8
	sideToMove        int8
	enPassant         Square
	whiteCastle       int8
	blackCastle       int8
	whiteKingPosition Square
	blackKingPosition Square
	status            int
	zobristTable      *ZobristTable
	currentHash       int64
}

// NewBoard creates a new chessboard from given fen
func NewBoard(fen string) *Board {
	b, err := parseFEN(fen)

	if err != nil {
		fmt.Printf("invalid FEN: \"%s\"\n", fen)
	}

	b.zobristTable = NewZobristTable()
	b.currentHash = b.generateHash()

	return b
}

func (b *Board) legalSquare(square int8) bool {
	// the magic of this 0x88 board representation
	return !(uint8(square)&0x88 != 0)
}

// MakeMove does a move on the board
func (b *Board) MakeMove(m Move) {

	historyItem := HistoryItem{
		move:          m,
		whiteCastle:   b.whiteCastle,
		blackCastle:   b.blackCastle,
		enPassant:     b.enPassant,
		halfMoveClock: b.halfMoveClock,
	}

	b.halfMoveClock++
	b.enPassant = Invalid

	if b.sideToMove == Black {
		b.fullMoves++
	}

	switch m.Special {
	case moveOrdinary:
		b.data[m.From] = Empty
		b.data[m.To] = m.MovedPiece

		if m.Content != Empty {
			b.halfMoveClock = 0
		}

		switch m.MovedPiece {
		case WhiteKing:
			b.whiteKingPosition = m.To
			if m.From == whiteKingStartSquare {
				b.whiteCastle = castleNone
			}
		case BlackKing:
			b.blackKingPosition = m.To
			if m.From == blackKingStartSquare {
				b.blackCastle = castleNone
			}
		case WhiteRook:
			if m.From == whiteRookShortSquare {
				b.whiteCastle &= ^castleShort
			} else if m.From == whiteRookLongSquare {
				b.whiteCastle &= ^castleLong
			}
		case BlackRook:
			if m.From == blackRookShortSquare {
				b.blackCastle &= ^castleShort
			} else if m.From == whiteRookLongSquare {
				b.blackCastle &= ^castleLong
			}
		case WhitePawn:
			b.halfMoveClock = 0
			steps := rank(int8(m.To)) - rank(int8(m.From))

			if steps > 1 {
				b.enPassant = Square(int8(m.From) + nextRank)
			}
		case BlackPawn:
			b.halfMoveClock = 0
			steps := rank(int8(m.From)) - rank(int8(m.To))

			if steps > 1 {
				b.enPassant = Square(int8(m.From) - nextRank)
			}
		}

	case moveCastelingShort:
		// king
		b.data[m.To] = m.MovedPiece
		b.data[m.From] = Empty
		// rook
		rookPos := int8(m.From) + castleShortDistanceRook*nextFile
		b.data[int8(m.From)+nextFile] = b.data[rookPos]
		b.data[rookPos] = Empty
		if m.MovedPiece == WhiteKing {
			b.whiteCastle = castleNone
			b.whiteKingPosition = m.To
		} else {
			b.blackCastle = castleNone
			b.blackKingPosition = m.To
		}
	case moveCastelingLong:
		// king
		b.data[m.To] = m.MovedPiece
		b.data[m.From] = Empty
		// rook
		rookPos := int8(m.From) - castleLongDistanceRook*nextFile
		b.data[int8(m.From)-nextFile] = b.data[rookPos]
		b.data[rookPos] = Empty
		if m.MovedPiece == WhiteKing {
			b.whiteCastle = castleNone
			b.whiteKingPosition = m.To
		} else {
			b.blackCastle = castleNone
			b.blackKingPosition = m.To
		}
	case movePromotion:
		b.data[m.From] = Empty
		b.data[m.To] = m.Promoted
		b.halfMoveClock = 0
	case moveEnPassant:
		b.data[m.From] = Empty
		b.data[m.To] = m.MovedPiece
		b.data[int8(m.To)-m.MovedPiece*nextRank] = Empty
		b.halfMoveClock = 0
	}

	b.sideToMove = opponent(b.sideToMove)
	b.ply++

	historyItem.hash = b.currentHash
	b.history = append(b.history, historyItem)

	b.updateHash(m)
}

// UndoMove undoes the last move on the board
func (b *Board) UndoMove() {
	if len(b.history) < 1 {
		fmt.Println("could not undo move")
		return
	}

	historyItem := b.history[len(b.history)-1]
	b.history = b.history[0 : len(b.history)-1]

	// update previous flags
	b.whiteCastle = historyItem.whiteCastle
	b.blackCastle = historyItem.blackCastle
	b.enPassant = historyItem.enPassant
	b.halfMoveClock = historyItem.halfMoveClock

	m := historyItem.move

	// XXX: might to be done at the end?
	b.updateHash(m)

	switch {
	case m.Special == moveOrdinary || m.Special == movePromotion:
		b.data[m.To] = m.Content
		b.data[m.From] = m.MovedPiece
		switch m.MovedPiece {
		case WhiteKing:
			b.whiteKingPosition = m.From
		case BlackKing:
			b.blackKingPosition = m.From
		}
	case m.Special == moveCastelingShort:
		b.data[m.From] = m.MovedPiece
		b.data[int8(m.From)+castleShortDistanceRook*nextFile] = b.data[int8(m.From)+nextFile]
		b.data[m.To] = Empty
		b.data[int8(m.From)+nextFile] = Empty
		switch m.MovedPiece {
		case WhiteKing:
			b.whiteKingPosition = m.From
		case BlackKing:
			b.blackKingPosition = m.From
		}
	case m.Special == moveCastelingLong:
		b.data[m.From] = m.MovedPiece
		b.data[int8(m.From)-castleLongDistanceRook*nextFile] = b.data[int8(m.From)-nextFile]
		b.data[m.To] = Empty
		b.data[int8(m.From)-nextFile] = Empty
		switch m.MovedPiece {
		case WhiteKing:
			b.whiteKingPosition = m.From
		case BlackKing:
			b.blackKingPosition = m.From
		}
	case m.Special == moveEnPassant:
		b.data[m.From] = m.MovedPiece
		b.data[m.To] = Empty
		b.data[int8(m.To)-m.MovedPiece*nextRank] = m.Content
	}

	b.sideToMove = opponent(b.sideToMove)
	b.ply--

	if b.sideToMove == Black {
		b.fullMoves--
	}

}

func (b *Board) isEmpty(squares ...Square) bool {
	for _, s := range squares {
		if b.data[uint8(s)] != Empty {
			return false
		}
	}
	return true
}

func (b *Board) repetitions() int {
	r := 0
	first := len(b.history) - b.halfMoveClock
	if first >= 0 {
		for i := first; i < len(b.history)-1; i++ {
			if b.history[i].hash == b.currentHash {
				r++
			}
		}
	}
	return r
}

func (b *Board) updateHash(m Move) {
	key := b.currentHash
	color := 0
	if b.sideToMove == Black {
		color = 1
	}
	piece := abs(m.MovedPiece) - 1

	key ^= b.zobristTable.hashPieces[piece][color][int8(m.From)]
	key ^= b.zobristTable.hashPieces[piece][color][int8(m.To)]

	if m.Content != Empty {
		key ^= b.zobristTable.hashPieces[abs(m.Content)-1][color][int8(m.To)]
	}

	switch m.Special {
	case moveCastelingLong:
		if color == 0 {
			key ^= b.zobristTable.hashCastelingWhite[castleLong]
		} else {
			key ^= b.zobristTable.hashCastelingBlack[castleLong]
		}
	case moveCastelingShort:
		if color == 0 {
			key ^= b.zobristTable.hashCastelingWhite[castleShort]
		} else {
			key ^= b.zobristTable.hashCastelingBlack[castleLong]
		}
	case moveEnPassant:
		key ^= b.zobristTable.hashEnPassant[int8(m.To)-int8(m.MovedPiece)*nextRank]
	case movePromotion:
		key ^= b.zobristTable.hashPromotion[abs(m.Promoted)-1]
	}

	b.currentHash = key
}

func (b *Board) generateHash() int64 {
	key := int64(0)

	for square := int8(0); square < boardSize; square++ {
		piece := b.data[square]
		if piece > Empty {
			key ^= b.zobristTable.hashPieces[piece-1][0][square]
		} else if piece < Empty {
			key ^= b.zobristTable.hashPieces[-piece-1][1][square]
		}
	}

	key ^= b.zobristTable.hashCastelingWhite[b.whiteCastle]
	key ^= b.zobristTable.hashCastelingBlack[b.blackCastle]
	if b.sideToMove == Black {
		key ^= b.zobristTable.hashSide
	}

	return key
}

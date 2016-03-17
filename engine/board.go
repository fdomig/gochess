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
}

// NewBoard creates a new chessboard from given fen
func NewBoard(fen string) *Board {
	b, err := parseFEN(fen)

	if err != nil {
		fmt.Printf("invalid FEN: \"%s\"\n", fen)
	}

	return b
}

func (b *Board) legalSquare(square int8) bool {
	// the magic of this 0x88 board representation
	return !(uint8(square)&0x88 != 0)
}

// MakeMove does a move on the board
func (b *Board) MakeMove(m Move) {

	historyItem := HistoryItem{
		move: m, whiteCastle: b.whiteCastle,
		blackCastle: b.blackCastle, enPassant: b.enPassant, halfMoveClock: b.halfMoveClock,
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
		b.data[int8(m.From)+nextFile] = b.data[rookPos]
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

	b.history = append(b.history, historyItem)
	b.ply++

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

	return r
}

func (b *Board) String() string {
	files := "   a  b  c  d  e  f  g  h"
	str := fmt.Sprintf("%s\n", files)

	for rank := int8(7); rank >= 0; rank-- {
		var s = fmt.Sprintf("%d  ", rank+1)
		for file := int8(0); file < 8; file++ {
			s += fmt.Sprintf("%s  ", symbol(b.data[square(rank, file)]))
		}
		if rank == 4 {
			color := "white"
			if b.sideToMove == Black {
				color = "black"
			}
			s += fmt.Sprintf("\t(%d) %s's move", b.fullMoves, color)
		}
		str += fmt.Sprintf("%s\n", s)
	}
	lastMove := ""
	if len(b.history) > 0 {
		lastMove = b.history[len(b.history)-1].move.String()
	}
	str += fmt.Sprintf("%s\t%s\t", files, lastMove)

	switch b.status {
	case statusCheck:
		str += "Check!"
	case statusDraw:
		str += "Draw!"
	case statusWhiteMates:
		str += "Mate! White wins."
	case statusBlackMates:
		str += "Mate! Black wins."
	case statusStaleMate:
		str += "Stale mate!"
	}

	str += "\n"

	return str
}

package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

	// defaultFEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -"

	// defaultFEN = "k7/8/K7/8/7p/7P/8/8 w - - 0 1" // only kings
	// defaultFEN = "8/7K/5k2/7P/8/8/8/8 w - - 7 15" // repetitions
	// defaultFEN = "6k1/8/6KP/8/8/8/8/8 w - - 7 15" // repetitions

	// defaultFEN = "r7/8/8/2k5/5K2/8/8/7R - - 49 1" // repetitions

	// defaultFEN = "r1bqkb1r/ppp1pppp/1nnp4/4P3/2P5/2N2N2/PP1PQPPP/R1B1KB1R b KQkq - 1 6" // leads to a thrid rook for black
	// defaultFEN = "rnbqkbnr/3p1ppp/p3p3/1pp3B1/3PP3/2N2N2/PPP2PPP/R2QKB1R b KQkq - 1 5" // thinks it's check but isn't
)

func generateFEN(board *Board) string {
	fen := ""

	// pieces
	count := 0
	for rank := int8(size - 1); rank >= 0; rank-- {
		for file := int8(0); file < size; file++ {
			if board.data[square(rank, file)] != Empty {
				if count > 0 {
					fen += fmt.Sprintf("%d", count)
					count = 0
				}
				fen += pieceString(board.data[square(rank, file)])
			} else {
				count++
			}
		}
		if count > 0 {
			fen += fmt.Sprintf("%d", count)
			count = 0
		}
		if rank != 0 {
			fen += "/"
		}
	}

	fen += " "

	// side to move
	if board.sideToMove == White {
		fen += "w"
	} else {
		fen += "b"
	}

	fen += " "

	// casteling

	if board.whiteCastle&castleShort != 0 {
		fen += "K"
	}

	if board.whiteCastle&castleLong != 0 {
		fen += "Q"
	}

	if board.blackCastle&castleShort != 0 {
		fen += "k"
	}

	if board.blackCastle&castleLong != 0 {
		fen += "q"
	}

	if (board.whiteCastle<<2 | board.blackCastle) == 0 {
		fen += "-"
	}

	fen += " "

	// en passant square

	if board.enPassant != Invalid {
		fen += strings.ToLower(SquareMap[board.enPassant])
	} else {
		fen += "-"
	}

	fen += " "

	// fifty moves count

	fen += fmt.Sprintf("%d", board.halfMoveClock)

	fen += " "

	fen += fmt.Sprintf("%d", board.fullMoves)

	return fen
}

func parseFEN(fen string) (*Board, error) {

	board := Board{}
	board.enPassant = Invalid
	board.fullMoves = 1

	for i := boardSize - 1; i >= 0; i-- {
		board.data[i] = Empty
	}

	parts := strings.Split(fen, " ")

	if len(parts) < 4 {
		return &board, errors.New("invalid FEN")
	}

	// parts[0]: piece placement
	i := 0
	for j := 0; j < len(parts[0]); j++ {
		switch parts[0][j] {
		case 'p':
			board.data[board64square[i]] = BlackPawn
		case 'r':
			board.data[board64square[i]] = BlackRook
		case 'n':
			board.data[board64square[i]] = BlackKnight
		case 'b':
			board.data[board64square[i]] = BlackBishop
		case 'q':
			board.data[board64square[i]] = BlackQueen
		case 'k':
			board.data[board64square[i]] = BlackKing
			board.blackKingPosition = Square(board64square[i])
		case 'P':
			board.data[board64square[i]] = WhitePawn
		case 'R':
			board.data[board64square[i]] = WhiteRook
		case 'N':
			board.data[board64square[i]] = WhiteKnight
		case 'B':
			board.data[board64square[i]] = WhiteBishop
		case 'Q':
			board.data[board64square[i]] = WhiteQueen
		case 'K':
			board.data[board64square[i]] = WhiteKing
			board.whiteKingPosition = Square(board64square[i])
		case '/':
			i--
		case '1':
			// noop
		case '2':
			i++
		case '3':
			i += 2
		case '4':
			i += 3
		case '5':
			i += 4
		case '6':
			i += 5
		case '7':
			i += 6
		case '8':
			i += 7
		default:
			return &board, errors.New("invalid FEN")
		}

		i++
	}

	// parts[1]: active color
	if parts[1] == "w" {
		board.sideToMove = White
	} else {
		board.sideToMove = Black
	}

	// parts[2]:casteling availability
	if strings.Contains(parts[2], "k") {
		board.blackCastle |= castleShort
	}

	if strings.Contains(parts[2], "q") {
		board.blackCastle |= castleLong
	}

	if strings.Contains(parts[2], "K") {
		board.whiteCastle |= castleShort
	}

	if strings.Contains(parts[2], "Q") {
		board.whiteCastle |= castleLong
	}

	// parts[3]: en passant target square
	if parts[3] != "-" {
		board.enPassant = SquareLookup[parts[3]]
	}

	// optionals:

	// parts[4]: halfmove clock (fifty move rule)
	if len(parts) >= 5 {
		board.halfMoveClock, _ = strconv.Atoi(parts[4])
	}

	// parts[5]: fullmove clock
	if len(parts) >= 6 {
		fullMoves, _ := strconv.Atoi(parts[5])
		board.fullMoves = fullMoves
	}

	return &board, nil
}

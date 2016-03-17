package engine

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	moveOrdinary       int8 = 0
	moveCastelingShort int8 = 1
	moveCastelingLong  int8 = 2
	movePromotion      int8 = 3
	moveEnPassant      int8 = 4

	castleNone  int8 = 0
	castleLong  int8 = 1
	castleShort int8 = 2
)

// Move on the board representation
type Move struct {
	From       Square
	To         Square
	Special    int8
	MovedPiece int8
	Content    int8
	Promoted   int8
}

func (m Move) String() string {
	capture := ""

	if m.Special == moveCastelingLong {
		return "O-O-O"

	} else if m.Special == moveCastelingShort {
		return "O-O"
	}

	if m.Content != Empty {
		capture = "x"
	}

	return fmt.Sprintf("%s%s%s", SquareMap[m.From], capture, SquareMap[m.To])
}

func createMove(str string) (Move, error) {

	// TODO castling/promotion

	if m, _ := regexp.MatchString("^[a-h][1-8][a-h][1-8]$", str); !m {
		return Move{}, errors.New("invalid move")
	}

	from := str[:2]
	to := str[2:]

	return Move{From: SquareLookup[from], To: SquareLookup[to]}, nil
}

func printMoves(moves []Move) {
	str := fmt.Sprintf("%d available moves:\n", len(moves))
	for i, move := range moves {
		captured := ""
		if move.Content != Empty {
			captured = fmt.Sprintf(" {%s}", symbol(move.Content))
		}

		if move.Special == moveCastelingLong {
			str += fmt.Sprintf("%s: O-O-O\t", symbol(move.MovedPiece))

		} else if move.Special == moveCastelingShort {
			str += fmt.Sprintf("%s: O-O  \t", symbol(move.MovedPiece))

		} else {
			str += fmt.Sprintf("%s: %s..%s%s\t", symbol(move.MovedPiece), SquareMap[move.From], SquareMap[move.To], captured)
		}

		if i%2 != 0 {
			str += "\n"
		}
	}
	fmt.Printf("%s\n", str)
}

package engine

import "math/rand"

const (
	numPieces     = 6
	numColors     = 2
	numCastelings = 4 // none, short, long, short&long
)

type ZobristTable struct {
	hashPieces         [numPieces][numColors][boardSize]int64
	hashEnPassant      [boardSize]int64
	hashCastelingBlack [numCastelings]int64
	hashCastelingWhite [numCastelings]int64
	hashPromotion      [numPieces]int64
	hashSide           int64
	hashCurrent        int64
}

func NewZobristTable() *ZobristTable {

	rand.Seed(4711)

	z := ZobristTable{}

	// pieces
	for piece := 0; piece < numPieces; piece++ {
		for color := 0; color < numColors; color++ {
			for square := int8(0); square < boardSize; square++ {
				z.hashPieces[piece][color][square] = hashRand()
			}
		}
	}
	// en passant
	for square := int8(0); square < boardSize; square++ {
		z.hashEnPassant[square] = hashRand()
	}
	// castling options
	for i := int8(0); i <= castleLong; i++ {
		z.hashCastelingBlack[i] = hashRand()
		z.hashCastelingWhite[i] = hashRand()
	}

	// promotion pieces
	for piece := int8(0); piece < numPieces; piece++ {
		z.hashPromotion[piece] = hashRand()
	}

	// side
	z.hashSide = hashRand()

	return &z
}

func hashRand() int64 {
	return rand.Int63()
}

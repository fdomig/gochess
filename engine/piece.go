package engine

const (
	White int8 = 1
	Black int8 = -1

	Empty int8 = 0

	Pawn   int8 = 1
	Knight int8 = 2
	Bishop int8 = 3
	Rook   int8 = 4
	Queen  int8 = 5
	King   int8 = 6

	WhitePawn   int8 = White * Pawn
	WhiteKnight int8 = White * Knight
	WhiteBishop int8 = White * Bishop
	WhiteRook   int8 = White * Rook
	WhiteQueen  int8 = White * Queen
	WhiteKing   int8 = White * King

	BlackPawn   int8 = Black * Pawn
	BlackKnight int8 = Black * Knight
	BlackBishop int8 = Black * Bishop
	BlackRook   int8 = Black * Rook
	BlackQueen  int8 = Black * Queen
	BlackKing   int8 = Black * King
)

var (
	symbols = []string{
		".",
		"P", "N", "B", "R", "Q", "K",
		"p", "n", "b", "r", "q", "k",
	}

	symbolsUnicode = []string{
		".",
		"♙", "♘", "♗", "♖", "♕", "♔",
		"♟", "♞", "♝", "♜", "♛", "♚",
	}
)

func symbol(piece int8) string {
	if piece >= 0 {
		return symbolsUnicode[piece]
	}

	return symbolsUnicode[piece*-1+6]
}

func pieceString(piece int8) string {
	if piece >= 0 {
		return symbols[piece]
	}

	return symbols[piece*-1+6]
}

func opposite(piece int8) int8 {
	return ^piece + 1
}

func opponent(color int8) int8 {
	return (-1 * color)
}

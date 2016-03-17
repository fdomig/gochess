package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Game represents a gochess game
type Game struct {
	board *Board
}

// NewGame creates a new gochess game and returns a reference
func NewGame() *Game {
	g := new(Game)
	g.board = NewBoard(defaultFEN)

	return g
}

// Run a given game
func (g *Game) Run() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")

	for scanner.Scan() {

		in := scanner.Text()

		if in == "quit" || in == "q" {
			break

		} else if in == "moves" || in == "m" {
			gen := NewGenerator(g.board)
			printMoves(gen.GenerateMoves())

		} else if in == "perft" {
			Perft(position1FEN, position1Table)

		} else if in == "perft2" {
			Perft(position2FEN, position2Table)

		} else if in == "new" || in == "n" {
			g.board = NewBoard(defaultFEN)

		} else if in == "fen" || in == "f" {
			fmt.Printf("%s\n", generateFEN(g.board))

		} else if in == "undo" || in == "u" {
			g.board.UndoMove()

		} else if strings.HasPrefix(in, "fen ") {
			g.board = NewBoard(in[3:])

		} else if in == "print" || in == "p" {
			fmt.Printf("%s\n", g.board.String())

		} else if in == "search" || in == "s" {
			Search(g.board)

		} else if in == "do" || in == "d" {
			g.board.MakeMove(Search(g.board))

		} else if in == "eval" || in == "e" {
			fmt.Printf("Score: %d\n", Evaluate(g.board))

		} else if in == "auto" {
			for g.board.status == statusNormal {
				g.board.MakeMove(Search(g.board))
				fmt.Printf("%s\n", g.board.String())
			}

		} else if m, err := createMove(in); err == nil {
			gen := NewGenerator(g.board)
			moves := gen.GenerateMoves()

			found := Move{From: Invalid}

			for _, move := range moves {
				if move.From == m.From && move.To == m.To {
					found = move
					break
				}
			}

			if found.From != Invalid {
				g.board.MakeMove(found)
			} else {
				fmt.Printf("illegal move\n")
			}

		} else {
			fmt.Printf("invalid input\n")
		}

		fmt.Printf("> ")
	}
}

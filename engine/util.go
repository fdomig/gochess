package engine

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func printSearchHead() {
	if !searchVerbose {
		return
	}

	fmt.Printf("ply  score   time   nodes  pv\n")
}

func printSearchLevel(pv *pvSearch, depth, score int, startTime time.Time) {
	if !searchVerbose {
		return
	}

	fmt.Printf("%3d %6s %6s %7s  ", depth,
		formatScore(score), formatDuration(time.Since(startTime)), formatNodesCount(pv.checkedNodes))

	for j := 0; j < pv.pathLength[0]; j++ {
		if pv.board.sideToMove == Black {
			if j == 0 {
				fmt.Printf("%d. ... ", pv.board.fullMoves)
			} else {
				if (j+1)%2 == 0 {
					fmt.Printf("%d. ", pv.board.fullMoves+(j/2+1))
				}
			}
		} else {
			if j%2 == 0 {
				fmt.Printf("%d. ", pv.board.fullMoves+(j/2))
			}
		}
		fmt.Printf("%s ", pv.path[0][j].String())
	}
	fmt.Printf("\n")
}

func printSearchResult(pv *pvSearch, startTime time.Time) {
	if !searchVerbose {
		return
	}

	totalTime := time.Since(startTime) // time is in nanoseconds
	fmt.Printf("%s nodes searched in %s secs (%.1fK nodes/sec)\n",
		formatNodesCount(pv.checkedNodes), formatDuration(totalTime),
		float64(pv.checkedNodes*1000000)/float64(totalTime))
}

func printPerftData(board *Board, expected []PerftData) {
	if !searchVerbose {
		return
	}

	fmt.Printf(color.WhiteString("D   Nodes    Capt.   E.p.   Cast.   Prom.  Checks   Mates   Time\n"))
	for i := 0; i < len(expected); i++ {

		res := perft(i, board)

		fmt.Printf("%d %7s %7s %7s %7s %7s %7s %7s  %5ss\n",
			i,
			formatNodesCount(res.nodes),
			formatNodesCount(res.captures),
			formatNodesCount(res.enPassants),
			formatNodesCount(res.castles),
			formatNodesCount(res.promotions),
			formatNodesCount(res.checks),
			formatNodesCount(res.mates),
			formatDuration(res.elapsed),
		)

		fmt.Printf("  %s %s %s %s %s %s %s\n\n",
			formatPerftEntry(res.nodes, expected[i].nodes),
			formatPerftEntry(res.captures, expected[i].captures),
			formatPerftEntry(res.enPassants, expected[i].enPassants),
			formatPerftEntry(res.castles, expected[i].castles),
			formatPerftEntry(res.promotions, expected[i].promotions),
			formatPerftEntry(res.checks, expected[i].checks),
			formatPerftEntry(res.mates, expected[i].mates))

	}

}

func formatPerftEntry(actual, expected int64) string {

	diff := actual - expected

	if diff != 0 {
		return color.RedString("%7s", formatNodesCount(diff))
	}

	return color.GreenString("      0")
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.2f", float64(d)*float64(1e-9))
}

func formatScore(score int) string {
	if score >= scoreMate {
		return "+Mate"
	} else if score <= -scoreMate {
		return "-Mate"
	}
	return fmt.Sprintf("%.2f", float64(score)/100)
}

func formatNodesCount(nodes int64) string {
	if nodes < 1000 && nodes > -1000 {
		return fmt.Sprintf("%d", nodes)
	} else if nodes < 1000000 && nodes > -1000000 {
		return fmt.Sprintf("%.1fK", float64(nodes)/1000)
	}
	return fmt.Sprintf("%.2fM", float64(nodes)/1000000)
}

func formatBoard(b *Board) string {
	files := "   a  b  c  d  e  f  g  h"
	str := fmt.Sprintf("%s\n", files)

	lastMoveSquare := Invalid

	if len(b.history) > 0 {
		lastMoveSquare = b.history[len(b.history)-1].move.To
	}

	for r := int8(7); r >= 0; r-- {
		var s = fmt.Sprintf("%d  ", r+1)
		for f := int8(0); f < 8; f++ {
			p := symbol(b.data[square(r, f)])

			if p == "." && (r+f)&1 == 0 {
				p = ","
			}

			s += fmt.Sprintf("%s", p)

			if square(r, f) == int8(lastMoveSquare) {
				s += "* "
			} else {
				s += "  "
			}
		}
		if r == 4 {
			color := "white"
			if b.sideToMove == Black {
				color = "black"
			}
			s += fmt.Sprintf("\t(%d) %s's move", b.fullMoves, color)
		}
		if r == 3 {
			c := ""
			if b.whiteCastle&castleShort != 0 {
				c += "K"
			}
			if b.whiteCastle&castleLong != 0 {
				c += "Q"
			}
			if b.blackCastle&castleShort != 0 {
				c += "k"
			}
			if b.blackCastle&castleLong != 0 {
				c += "q"
			}
			if len(c) == 0 {
				c = "-"
			}
			s += fmt.Sprintf("\tCasteling: %s", c)
		}
		if r == 2 {
			gen := NewGenerator(b)
			if gen.CheckSimple() {
				s += fmt.Sprintf("\tCheck!")
			}
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

func abs(v int8) int8 {
	if v >= 0 {
		return v
	}
	return -v
}

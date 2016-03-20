package engine

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func printSearchHead() {
	fmt.Printf("ply  score   time   nodes  pv\n")
}

func printSearchLevel(pv *pvSearch, depth, score int, startTime time.Time) {
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
	totalTime := time.Since(startTime) // time is in nanoseconds
	fmt.Printf("%s nodes searched in %s secs (%.1fK nodes/sec)\n",
		formatNodesCount(pv.checkedNodes), formatDuration(totalTime),
		float64(pv.checkedNodes*1000000)/float64(totalTime))
}

func printPerftData(board *Board, expected []PerftData) {

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

func abs(v int8) int8 {
	if v >= 0 {
		return v
	}
	return -v
}

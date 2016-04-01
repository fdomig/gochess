# Gochess

> A chess engine written in Go.

Well, this is a chess engine highly inspired by [chess-at-nite][chess-at-nite] which I have written several years ago with my friends in Copenhagen, Denmark. Since I am interested in writting an application in Golang I thought of doing chess.

## Installation

```
$ go get github.com/fdomig/gochess
$ godep restore
```

You might need to install `godep` with `github.com/tools/godep` first.

## Run

```
$ gochess
```

You then may just enter a move in the coordinate notation (e.g. `e2e4`). 

## Commands

```
auto, a      let the engine play against itself

do, d        search the best available move and play it

eval, e      displays the current board's score 

fen, f       displays the current board position in the Forsyth Edwards Notation (FEN)

new, n       start a new game

moves, m     show a list of all possible moves

print, p     shows the current board position

quit, q      quits this game and the application

search, s    search the current board position for the best possible move

undo, u      undo the last move

```     


## Ideas

* Use algebraic notation for input and display
* Make multi threaded with channels
* Use zobrist hash vor repetitions, score cache, etc. 

## Contribute 

Feel free to contribute and fix things via GitHub Pull Requests.

[chess-at-nite]: https://github.com/fdomig/chess-at-nite
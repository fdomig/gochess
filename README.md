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

## Ideas

* Make multi threaded with channels
* Use zobrist hash vor repetitions, score cache, etc. 

## Contribute 

Feel free to contribute and fix things via GitHub Pull Requests.

[chess-at-nite]: https://github.com/fdomig/chess-at-nite
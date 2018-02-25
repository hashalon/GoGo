# GoGo
Go game made in Go language

## Play
You can read the rules of Go on [wikipedia](https://en.wikipedia.org/wiki/Rules_of_Go).
Simply click on a position of the Goban to place a stone.
The banner at the bottom show which player can play.

There is no end game yet.

## Compiling
You'll need the library [Ebiten by hajimehoshi](https://github.com/hajimehoshi/ebiten) to compile it.
to build, execute:
```
go build board.go liberty.go main.go vector.go window.go -o gogo
```
to run, execute:
```
go run board.go liberty.go main.go vector.go window.go
```

## Side notes
The program has only been tested on Linux until now.
If you have trouble compiling it on Windows, I cannot provide help, sorry.

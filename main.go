package main

import (
	"fmt"
	tm "github.com/buger/goterm"
)

func main() {
	tm.Clear()
	tm.MoveCursor(0,0)
	fmt.Println("GO-GO running...")
	tm.Flush()
	
	game := MakeGame()
	for true {
		game.Turn()
	}
}

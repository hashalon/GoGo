package main

import "fmt"

func main() {
	fmt.Println("GoGo running...")

	// build a board of size 19x19
	board := MakeBoard(19)

	// make a window with library Ebiten to play the game
	var window Window
	window.Start(&board)

	fmt.Println("GoGo stopped.")
}

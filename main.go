package main

import "fmt"

func main() {
	fmt.Println("GO-GO running...")

	game := MakeGame()
	for {
		game.Turn()
	}
}

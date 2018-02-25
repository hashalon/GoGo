package main

import "fmt"

func main() {
	fmt.Println("GO-GO running...")

	board := MakeBoard(19)
	/*
	for i := 0; i < 19; i+=2 {
	for j := 0; j < 19; j+=2 {
		st := Stone{Vec2{int8(i), int8(j)}, i%4 == j%4}
		board.Place(st)
	}} // */
	var window Window
	window.Init(&board)

}

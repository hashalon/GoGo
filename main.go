package main

import (
	"fmt"
	"github.com/fatih/color"
)

func main() {
	fmt.Println("GO-GO running...")
	/*
	chars := CharSet{ '○', '●',
		'╔', '╤', '╗',
		'╟', '┼', '╢',
		'╚', '╧', '╝',
		     '╦',
		'╠', '╬', '╣',
		     '╩',
		'╪', '╫',
		"ABCDEFGHIJKLMNOPQRST",
		"①②③④⑤⑥⑦⑧⑨⑩⑪⑫⑬⑭⑮⑯⑰⑱⑲⑳" } /*/
	chars := CharSet{ 'O', 'Q',
		'┌', '─', '┐',
		'│', ' ', '│',
		'└', '─', '┘',
			 '┬',
		'├', '┼', '┤',
		     '┴',
		'─', '│',
		"ABCDEFGHIJKLMNOPQRST",
		"abcdefghijklmnopqrst" } //*/
	/*
	colors := ColorSet{
		color.New(color.FgHiCyan),
		color.New(color.FgHiRed),
		color.New(color.FgHiBlack),
		color.New(color.FgMagenta)} /*/
	colors := ColorSet{
		color.New(color.FgHiWhite),
		color.New(color.FgHiRed),
		color.New(color.FgHiBlack),
		color.New(color.FgBlue)} //*/
	display := MakeDisplay(chars, colors, 19)
	board   := MakeBoard()
	board.set[0] = Stone{Vec2{2, 3}, true}
	board.set[1] = Stone{Vec2{4, 5}, false}
	display.Draw(board)
}

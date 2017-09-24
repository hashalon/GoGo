package main

import (
	"fmt"
	"regexp"
	"strconv"
	"github.com/fatih/color"
	tm "github.com/buger/goterm"
)

// Game a board and a display
type Game struct {
	board Board
	display Display
	team, invalid, occupied bool
}

// MakeGame make a new game
func MakeGame() Game {
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
		"ABCDEFGHIJKLMNOPQRST", "" } //*/
	colors := ColorSet{
		color.New(color.FgHiWhite),
		color.New(color.FgHiRed),
		color.New(color.FgHiBlack),
		color.New(color.FgYellow) }
	// start the game with black
	return Game{ MakeBoard(19), MakeDisplay(chars, colors, 19), true, false, false }
}

// Turn manage one turn of the game
func (game *Game) Turn() {
	tm.Clear()
	// display the current state of the game
	game.display.Draw(game.board)
	// retry until we have a valid position
	retry := true
	for retry {
		retry = false
		var pos Vec2
		pos, game.invalid = game.Select()
		if !game.invalid {
			stone := Stone{pos, game.team}
			game.occupied = !game.board.Place(stone)
			if game.occupied { retry = true }
		} else { retry = true }
	}
	// stone placed
	game.team = !game.team
	game.invalid  = false
	game.occupied = false
	tm.Flush()
}

var (
	colregex = regexp.MustCompile(`[A-Z]|[a-z]`)
	rowregex = regexp.MustCompile(`[0-9]+`)
)

// Select selection for the current player, return a error for incorrect input
func (game *Game) Select() (Vec2, bool) {

	// state the error
	if game.invalid {
		fmt.Println(color.YellowString("Position selected is invalid!"))
	} else if game.occupied {
		fmt.Println(color.YellowString("Position is occupied already!"))
	}

	// ask the current player to place a stone
	if game.team {
		game.display.colorB.Printf("Black place a stone: ")
	} else {
		game.display.colorW.Printf("White place a stone: ")
	}

	// read the selected position
	var selection string
	fmt.Scanln(&selection)
	strcol := colregex.FindString(selection)
	strrow := rowregex.FindString(selection)

	// extract coordinates
	if strcol == "" || strrow == "" { return Vec2{-1, -1}, true }

	// parse the coordinates
	x := strcol[0]
	y, err := strconv.ParseInt(strrow, 10, 8)
	if err != nil { return Vec2{-1, -1}, true }
	if x >= 'a' { x -= 'a' } else if x >= 'A' { x -= 'A' }

	// check that the position is inside of the board
	pos := Vec2{int8(x), int8(y-1)}
	if !game.board.Valid(pos) { return Vec2{-1, -1}, true }
	return pos, false
}
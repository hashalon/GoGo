package main

import (
	"time"
	"fmt"
	"regexp"
	"strconv"

	tm "github.com/buger/goterm"
	"github.com/fatih/color"
	term "github.com/nsf/termbox-go"
)

// Game a board and a display
type Game struct {
	board    Board
	display  Display
	selector Vec2
	team     bool
}

// MakeGame make a new game
func MakeGame() Game {
	/*
		chars := CharSet{'○', '●', '×',
			'╔', '╤', '╗',
			'╟', '┼', '╢',
			'╚', '╧', '╝',
			     '╦',
			'╠', '╬', '╣',
			     '╩',
			'╪', '╫',
			"ABCDEFGHIJKLMNOPQRST",
			"①②③④⑤⑥⑦⑧⑨⑩⑪⑫⑬⑭⑮⑯⑰⑱⑲⑳" } // */
	chars := CharSet{'O', 'Q', '#',
		'┌', '─', '┐',
		'│', ' ', '│',
		'└', '─', '┘',
		     '┬',
		'├', '┼', '┤',
		     '┴',
		'─', '│',
		"ABCDEFGHIJKLMNOPQRST", ""} // */
	colors := ColorSet{
		color.New(color.FgHiWhite),
		color.New(color.FgHiRed),
		color.New(color.FgHiBlack),
		color.New(color.BgHiCyan),
		color.New(color.FgYellow)}
	// start the game with black
	return Game{MakeBoard(19), MakeDisplay(chars, colors, 19), Vec2{10, 10}, true}
}

// Turn manage one turn of the game
func (game *Game) Turn() {
	tm.Clear()
	// detect arrow keys input
	for {
		tm.MoveCursor(0, 0)
		// move the selector and check if we try a position
		selected := false
		go func () { selected = game.Interact() }()
		if selected {
			stone := Stone{game.selector, game.team}
			if game.board.Valid(game.selector) &&
			   game.board.Place(stone) { break }
		} else {
			game.display.Draw(game.board, game.selector)
			tm.Flush()
			time.Sleep(time.Second)
		}
	}
	// stone placed
	game.team = !game.team
	game.selector = Vec2{10, 10}
}

// Interact manage key inputs
func (game *Game) Interact() bool {
	switch ev := term.PollEvent(); ev.Type {
	case term.EventKey:
		switch ev.Key {
		case term.KeyArrowUp:
			game.selector.y++
			if game.selector.y > 18 { game.selector.y = 18 }
		case term.KeyArrowDown:
			game.selector.y--
			if game.selector.y <  0 { game.selector.y =  0 }
		case term.KeyArrowLeft:
			game.selector.x--
			if game.selector.x <  0 { game.selector.x =  0 }
		case term.KeyArrowRight:
			game.selector.x++
			if game.selector.y > 18 { game.selector.x = 18 }
		case term.KeyEnter:
			return true
		}
	case term.EventError:
		panic(ev.Err)
	}
	return false
}

var (
	colregex = regexp.MustCompile(`[A-Z]|[a-z]`)
	rowregex = regexp.MustCompile(`[0-9]+`)
	invalid  = false
	occupied = false
)

// TurnText manage one turn of the game
func (game *Game) TurnText() {
	tm.Clear()
	// display the current state of the game
	game.display.Draw(game.board, Vec2{-1, -1})
	// retry until we have a valid position
	for {
		tm.MoveCursor(0, 0)
		var pos Vec2
		pos, invalid = game.SelectText()
		if !invalid {
			stone := Stone{pos, game.team}
			occupied = !game.board.Place(stone)
			if occupied { break }
		}
	}
	// stone placed
	game.team = !game.team
	invalid  = false
	occupied = false
	tm.Flush()
}

// SelectText selection for the current player, return a error for incorrect input
func (game *Game) SelectText() (Vec2, bool) {

	// state the error
	if invalid {
		fmt.Println(color.YellowString("Position selected is invalid!"))
	} else if occupied {
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
	// manage coordinates with upper and lower case characters
	if x >= 'a' { 
		x -= 'a' 
	} else if x >= 'A' {
		x -= 'A' 
	}

	// check that the position is inside of the board
	pos := Vec2{int8(x), int8(y - 1)}
	if !game.board.Valid(pos) { return Vec2{-1, -1}, true }
	return pos, false
}



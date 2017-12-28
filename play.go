package main

// to correct
// https://github.com/nsf/termbox-go
// https://github.com/1984weed/2048-go/blob/master/application.go


import (
	"time"
	"fmt"
	"regexp"
	"strconv"
	"github.com/nsf/termbox-go"
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
		termbox.ColorWhite,
		termbox.ColorRed,
		termbox.ColorBlack,
		termbox.ColorCyan,
		termbox.ColorYellow}
	// start the game with black
	return Game{MakeBoard(19), MakeDisplay(chars, colors, 19), Vec2{10, 10}, true}
}

// Turn manage one turn of the game
func (game *Game) Turn() {
	const defcol = termbox.ColorDefault
	termbox.Clear(defcol, defcol)
	// detect arrow keys input
	for {
		termbox.SetCursor(0, 0)
		// move the selector and check if we try a position
		selected := false
		go func () { selected = game.Interact() }()
		if selected {
			stone := Stone{game.selector, game.team}
			if game.board.Valid(game.selector) &&
			   game.board.Place(stone) { break }
		} else {
			game.display.Draw(game.board, game.selector)
			termbox.Flush()
			time.Sleep(time.Second)
		}
	}
	// stone placed
	game.team = !game.team
	game.selector = Vec2{10, 10}
}

// Interact manage key inputs
func (game *Game) Interact() bool {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp:
			game.selector.y++
			if game.selector.y > 18 { game.selector.y = 18 }
		case termbox.KeyArrowDown:
			game.selector.y--
			if game.selector.y <  0 { game.selector.y =  0 }
		case termbox.KeyArrowLeft:
			game.selector.x--
			if game.selector.x <  0 { game.selector.x =  0 }
		case termbox.KeyArrowRight:
			game.selector.x++
			if game.selector.y > 18 { game.selector.x = 18 }
		case termbox.KeyEnter:
			return true
		}
	case termbox.EventError:
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
	// display the current state of the game
	game.display.Draw(game.board, Vec2{-1, -1})
	// retry until we have a valid position
	for {
		termbox.SetCursor(0, 0)
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
}

// SelectText selection for the current player, return a error for incorrect input
func (game *Game) SelectText() (Vec2, bool) {

	// state the error
	if invalid {
		fmt.Println("Position selected is invalid!")
	} else if occupied {
		fmt.Println("Position is occupied already!")
	}

	// ask the current player to place a stone
	if game.team {
		fmt.Printf("Black place a stone: ")
	} else {
		fmt.Printf("White place a stone: ")
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



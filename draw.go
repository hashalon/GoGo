package main

// to correct
// https://github.com/nsf/termbox-go
// https://github.com/1984weed/2048-go/blob/master/application.go

import (
	"bytes"
	"fmt"
	"github.com/nsf/termbox-go"
)

// CharSet define the characters to use for the display
type CharSet struct {
	white, black, selector,
	tl, tm, tr,
	ml, mm, mr,
	bl, bm, br,
	    tc,
	cl, cc, cr,
	    bc,
	ch, cv rune
	columns, rows string
}

// ColorSet set of colors used to display the game
type ColorSet struct {
	colorW, colorB, selector, background, label termbox.Attribute
}

// Display configuration to display board
type Display struct {
	size                                    uint8
	stoneW, stoneB, selector                rune
	colorW, colorB, colorSelect, background termbox.Attribute
	labelCols                               string
	labelRows                               [20]string
	layout                                  [4][20]rune
}

// MakeDisplay return a new display configuration
func MakeDisplay(chars CharSet, colors ColorSet, size uint8) Display {
	if size > 20 { size = 20 }
	display := Display{size, chars.white, chars.black, chars.selector,
		colors.colorW, colors.colorB, colors.selector, colors.background,
		"", [20]string{}, [4][20]rune{}}
	colorLabel := colors.label.SprintFunc()
	if uint8(len(chars.rows)) >= size {
		display.labelCols = colorLabel(" " + chars.columns[:size])
		rowRunes := []rune(chars.rows)
		for i := uint8(0); i < size; i++ {
			display.labelRows[i] = colorLabel(string(rowRunes[i]))
		}
	} else {
		display.labelCols = colorLabel("  " + chars.columns[:size])
		for i := uint8(0); i < size; i++ {
			var label string
			if i < 9 {
				label = fmt.Sprintf(" %d", i+1)
			} else {
				label = fmt.Sprintf("%d", i+1)
			}
			display.labelRows[i] = colorLabel(label)
		}
	}
	// generate the layout from the size of the board
	display.layout[0][0] = chars.tl
	display.layout[1][0] = chars.ml
	display.layout[2][0] = chars.cl
	display.layout[3][0] = chars.bl
	for i := uint8(1); i < size-1; i++ {
		if i%6 == 3 {
			display.layout[0][i] = chars.tc
			display.layout[1][i] = chars.cv
			display.layout[2][i] = chars.cc
			display.layout[3][i] = chars.bc
		} else {
			display.layout[0][i] = chars.tm
			display.layout[1][i] = chars.mm
			display.layout[2][i] = chars.ch
			display.layout[3][i] = chars.bm
		}
	}
	display.layout[0][size-1] = chars.tr
	display.layout[1][size-1] = chars.mr
	display.layout[2][size-1] = chars.cr
	display.layout[3][size-1] = chars.br
	return display
}

// Draw the board
func (display *Display) Draw(board Board, highlight Vec2) {
	layout := [20][20]rune{}
	// generate the board
	for i := uint8(0); i < display.size; i++ {
		row := 1
		if i == 0 { 
			row = 0 
		} else if i == display.size-1 { 
			row = 3 
		} else if i%6 == 3 {
			row = 2 
		}
		layout[i] = display.layout[row]
	}
	// place the stones
	for _, stone := range board.set {
		stoneRune := display.stoneW
		if stone.team { stoneRune = display.stoneB }
		layout[stone.y][stone.x] = stoneRune
	}
	// add the selector
	if highlight.x > -1 && highlight.x < 19 &&
	   highlight.y > -1 && highlight.y < 19 {
		layout[highlight.x][highlight.y] = display.selector
	}

	// add labels above the board
	fmt.Println(display.labelCols)
	// convert to string, put the colors and display
	for i := uint8(0); i < display.size; i++ {
		var mainBuffer, secondBuffer bytes.Buffer
		mainBuffer.WriteString(display.labelRows[i])
		secondBuffer.WriteRune(layout[i][0]) // push first rune in buffer
		for j := uint8(1); j < display.size; j++ {
			curr := display.typeRune(layout[i][j])
			prev := display.typeRune(layout[i][j-1])
			if curr != prev {
				color := display.selectColor(prev).SprintFunc()
				mainBuffer.WriteString(color(secondBuffer.String()))
				secondBuffer.Reset()
			}
			secondBuffer.WriteRune(layout[i][j])
		}
		color := display.selectColor(layout[i][display.size-1]).SprintFunc()
		mainBuffer.WriteString(color(secondBuffer.String()))
		fmt.Println(mainBuffer.String())
	}
}

func (display *Display) selectColor(r rune) *color.Color {
	switch r {
	case display.selector : return display.colorSelect
	case display.stoneW   : return display.colorW
	case display.stoneB   : return display.colorB
	}
	return display.background
}

func (display *Display) typeRune(r rune) rune {
	if r != display.stoneW &&
	   r != display.stoneB { return 0 }
	return r
}

package main

import (
	"fmt"
	"log"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// size, padding, scaling
const (
	GOBAN_SIZE = 512
	STONE_SIZE =  25.5

	// hardcoded values relative to sprites used
	GOBAN_SCALE = 1.0 /  4.0
	STONE_SCALE = 1.0 / 10.0
	PADDING     = 27

	START_X =              PADDING - STONE_SIZE / 2.0
	START_Y = GOBAN_SIZE - PADDING - STONE_SIZE / 2.0
)

// store the necessities to draw the board
type Window struct {
	// images to draw the board
	gobanImage, blackImage, whiteImage * ebiten.Image
	// draw options for the board
	gobanOpts * ebiten.DrawImageOptions

	// position of the selector when using key inputs
	selectorPos Vec2
	// the board to display
	board * Board
	// team to play
	team bool
}
var window * Window

// initialize the window and run the game
func (win * Window) Init(board * Board) {

	// get the images to display
	var err error
	if win.gobanImage, _, err = ebitenutil.NewImageFromFile(  "goban.png", ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.blackImage, _, err = ebitenutil.NewImageFromFile("stone_b.png", ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.whiteImage, _, err = ebitenutil.NewImageFromFile("stone_w.png", ebiten.FilterNearest); err != nil {log.Fatal(err)}

	// create a single draw option for the board
	win.gobanOpts = &ebiten.DrawImageOptions{}
	win.gobanOpts.GeoM.Scale(GOBAN_SCALE, GOBAN_SCALE) // 4 times smaller
	win.board = board
	win.team  = false
	window = win

	// run the game loop
	if err := ebiten.Run(update, GOBAN_SIZE, GOBAN_SIZE, 1, "GoGo"); err != nil {log.Fatal(err)}
}


func update (screen * ebiten.Image) error {

	/* update game state */

	// get inputs
	if ebiten.IsKeyPressed(ebiten.KeyUp   ) {window.selectorPos.y += 1}
	if ebiten.IsKeyPressed(ebiten.KeyDown ) {window.selectorPos.y -= 1}
	if ebiten.IsKeyPressed(ebiten.KeyLeft ) {window.selectorPos.x -= 1}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {window.selectorPos.x += 1}

	// we prefer interactions with the mouse
	x, y := ebiten.CursorPosition()
	c := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// when we click try to place a stone in the board
	// if it succedes, end turn
	if c {
		posX := (float64(x) - START_X) / STONE_SIZE
		posY := (START_Y - float64(y)) / STONE_SIZE
		stone := Stone{Vec2{int8(posX), int8(posY + 1)}, window.team}
		if window.board.Place(stone) {
			window.team = !window.team // end turn
		}
	}

	// do not force the system to render if it gets too slow
	if ebiten.IsRunningSlowly() { return nil }

	/* render the game */

	// draw a cool pattern
	screen.Fill(color.NRGBA{0x30, 0x30, 0x30, 0xff})

	// draw the goban
	screen.DrawImage(window.gobanImage, window.gobanOpts)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d, FPS: %f", x, y, ebiten.CurrentFPS()))
	if c { ebitenutil.DebugPrint(screen, "\n\nClick") }

	if window.team {ebitenutil.DebugPrint(screen, "\nWHITE")
	} else         {ebitenutil.DebugPrint(screen, "\nBLACK")}

	// for each stone on the board, draw it on the display
	for _, stone := range window.board.set {
		posX := START_X + float64(stone.Vec2.x) * STONE_SIZE
		posY := START_Y - float64(stone.Vec2.y) * STONE_SIZE
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(STONE_SCALE, STONE_SCALE) // ten times smaller
		opts.GeoM.Translate(posX, posY)
		if stone.team {screen.DrawImage(window.whiteImage, opts)
		} else        {screen.DrawImage(window.blackImage, opts)}
	}
	return nil
}

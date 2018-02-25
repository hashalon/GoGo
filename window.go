package main

import (
	//"fmt"
	"log"
	"image"
	//"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// size, padding, scaling
const (
	GOBAN_SIZE  = 608
	STONE_SIZE  = 32
	BANNER_SIZE = 32
)

// store the necessities to draw the board
type Window struct {
	// images to draw the board
	stoneBlackImage, bannerBlackImage,
	stoneWhiteImage, bannerWhiteImage,
	gobanImage * ebiten.Image

	// draw options for the board and banner
	gobanOpts, bannerOpts * ebiten.DrawImageOptions

	// position of the selector when using key inputs
	selectorPos Vec2 // TODO: not used yet
	wasClicked bool // true if the mouse button was pressed already on the frame before

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
	if win.gobanImage,       _, err = ebitenutil.NewImageFromFile("goban.png",       ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.stoneBlackImage,  _, err = ebitenutil.NewImageFromFile("stoneBlack.png",  ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.stoneWhiteImage,  _, err = ebitenutil.NewImageFromFile("stoneWhite.png",  ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.bannerBlackImage, _, err = ebitenutil.NewImageFromFile("bannerBlack.png", ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.bannerWhiteImage, _, err = ebitenutil.NewImageFromFile("bannerWhite.png", ebiten.FilterNearest); err != nil {log.Fatal(err)}

	// if we cannot load the icon, it doesn't matter
	_, icon, err := ebitenutil.NewImageFromFile("icon.png", ebiten.FilterNearest)
	if err == nil {ebiten.SetWindowIcon([]image.Image{icon})}

	// create a single draw option for the board
	win.gobanOpts  = &ebiten.DrawImageOptions{}
	win.bannerOpts = &ebiten.DrawImageOptions{}
	win.bannerOpts.GeoM.Translate(0, GOBAN_SIZE)

	win.board = board
	win.team  = false
	window = win

	// place the selector in the center
	win.selectorPos.x = int8(10)
	win.selectorPos.y = int8(10)
	win.wasClicked = false

	// run the game loop
	if err := ebiten.Run(update, GOBAN_SIZE, GOBAN_SIZE + BANNER_SIZE, 1, "GoGo"); err != nil {log.Fatal(err)}
}

// update selector position and play the game
func (win * Window) Update() {
	// get inputs
	if ebiten.IsKeyPressed(ebiten.KeyUp   ) {win.selectorPos.y += 1}
	if ebiten.IsKeyPressed(ebiten.KeyDown ) {win.selectorPos.y -= 1}
	if ebiten.IsKeyPressed(ebiten.KeyLeft ) {win.selectorPos.x -= 1}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {win.selectorPos.x += 1}

	// we prefer interactions with the mouse
	x, y := ebiten.CursorPosition()
	click := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	// when we click try to place a stone in the board
	// if it succedes, end turn
	if click && !win.wasClicked {
		stone := Stone{win.BoardCoords(x, y), win.team}
		if win.board.Place(stone) {
			win.team = !win.team // end turn
		}
	}

	// keep track of the state of click
	if click {win.wasClicked = true} else {win.wasClicked = false}
}

// used to redraw the board only when necessary
func (win * Window) Redraw(screen * ebiten.Image) {
	// draw the goban and the banner
	screen.DrawImage(win.gobanImage, win.gobanOpts)
	if win.team {screen.DrawImage(win.bannerWhiteImage, win.bannerOpts)
	} else      {screen.DrawImage(win.bannerBlackImage, win.bannerOpts)}

	// for each stone on the board, draw it on the display
	for _, stone := range win.board.set {
		x, y := win.WindowCoords(stone.Vec2)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, y)
		if stone.team {screen.DrawImage(win.stoneWhiteImage, opts)
		} else        {screen.DrawImage(win.stoneBlackImage, opts)}
	}
}

// from window coords, get board coords
// used to click on the board
func (win * Window) BoardCoords(x, y int) Vec2 {
	if x < 0 || y < 0 || x > GOBAN_SIZE || y > GOBAN_SIZE {return InvalidVec()}
	return Vec2{int8(x / STONE_SIZE), int8(y / STONE_SIZE)}
}

// from board coords, get window coords
// used to draw elements on the board
func (win * Window) WindowCoords(pos Vec2) (float64, float64) {
	return float64(pos.x) * STONE_SIZE, float64(pos.y) * STONE_SIZE
}

// function called in loop by ebiten
func update (screen * ebiten.Image) error {

	/* update game state */
	window.Update()

	// do not force the system to render if it gets too slow
	if ebiten.IsRunningSlowly() { return nil }

	/* render the game */
	window.Redraw(screen)

	return nil
}

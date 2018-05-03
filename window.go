package main

import (
	"log"
	"image"
	_ "image/png"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// size, padding, scaling
const (
	GOBAN_SIZE  = 608
	STONE_SIZE  = 32
	BANNER_SIZE = 32
	PASS_SIZE   = 64
)

// store the necessities to draw the board
type Window struct {
	// images to draw the board
	stoneBlackImage, bannerBlackImage, terrBlackImage,
	stoneWhiteImage, bannerWhiteImage, terrWhiteImage,
	bannerEndImage , gobanImage * ebiten.Image

	// draw options for the board and banner
	gobanOpts, bannerOpts * ebiten.DrawImageOptions

	// position of the selector when using key inputs
	selectorPos Vec2 // TODO: not used yet
	wasClicked bool // true if the mouse button was pressed already on the frame before

	// the board to display
	board * Board
	// team to play: false=black ; true=white
	team, hasEnded bool

	// count the number of passes
	blackPasses, whitePasses int8

	// end game containing the set of all territories
	end EndGame
}
var window * Window

// initialize the window and run the game
func (win * Window) Start(board * Board) {

	// get the images to display
	var err error
	if win.gobanImage,       _, err = ebitenutil.NewImageFromFile("goban.png",          ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.stoneBlackImage,  _, err = ebitenutil.NewImageFromFile("stoneBlack.png",     ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.stoneWhiteImage,  _, err = ebitenutil.NewImageFromFile("stoneWhite.png",     ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.bannerBlackImage, _, err = ebitenutil.NewImageFromFile("bannerBlack.png",    ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.bannerWhiteImage, _, err = ebitenutil.NewImageFromFile("bannerWhite.png",    ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.bannerEndImage,   _, err = ebitenutil.NewImageFromFile("bannerEnd.png",      ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.terrBlackImage,   _, err = ebitenutil.NewImageFromFile("territoryBlack.png", ebiten.FilterNearest); err != nil {log.Fatal(err)}
	if win.terrWhiteImage,   _, err = ebitenutil.NewImageFromFile("territoryWhite.png", ebiten.FilterNearest); err != nil {log.Fatal(err)}

	// if we cannot load the icon, it doesn't matter
	_, icon, err := ebitenutil.NewImageFromFile("icon.png", ebiten.FilterNearest)
	if err == nil {ebiten.SetWindowIcon([]image.Image{icon})}

	// create a single draw option for the board
	win.gobanOpts  = &ebiten.DrawImageOptions{}
	win.bannerOpts = &ebiten.DrawImageOptions{}
	win.bannerOpts.GeoM.Translate(0, GOBAN_SIZE)

	win.board = board
	win.team  = false // start with black
	window = win

	win.blackPasses = 0
	win.whitePasses = 0
	win.hasEnded = false

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

		// if the game is over, define the territories
		// if the player pressed the pass button, pass the turn
		// if we clicked on the goban, place a stone
		if win.IsOver() {
			// get the territories if null
			if !win.hasEnded {win.end = MakeTerritories(win.board)}
			
		} else if PassButton(x, y) {
			// increment the right passes counter
			if win.team {win.whitePasses += 1} else {win.blackPasses += 1}
			win.team = !win.team // end turn
		} else {
			stone := Stone{BoardCoords(x, y), win.team}
			if win.board.Place(stone) {
				win.team = !win.team // end turn
			}
		}
	}

	// keep track of the state of click
	win.wasClicked = click
}

// used to redraw the board only when necessary
func (win * Window) Redraw(screen * ebiten.Image) {
	var img * ebiten.Image // shortcut to access the image to print

	// draw the goban
	screen.DrawImage(win.gobanImage, win.gobanOpts)

	// if the game is over display the end banner,
	// otherwise display the banner of the player
	if win.IsOver() {
		img = win.bannerEndImage

		// display the territories too
		for _, terr := range win.end.set {
			noteam := terr.whiteCount == terr.blackCount
			team   := terr.whiteCount >  terr.blackCount

			// pick the image to use based on the team of the territory
			if noteam {continue} else if team {img = win.terrWhiteImage} else {img = win.terrBlackImage}
			// print the territory
			for _, pos := range terr.set {
				x, y := WindowCoords(pos)
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(x, y)
				screen.DrawImage(img, opts)
			}
		}

	} else if win.team {img = win.bannerWhiteImage} else {img = win.bannerBlackImage}
	screen.DrawImage(img, win.bannerOpts)

	// for each stone on the board, draw it on the display
	for _, stone := range win.board.set {
		x, y := WindowCoords(stone.Vec2)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, y)

		// print the stone at the right location
		if stone.team {img = win.stoneWhiteImage} else {img = win.stoneBlackImage}
		screen.DrawImage(img, opts)
	}

	// display the scores as an overlay
	if win.IsOver() {
		// TODO...
	}
}

// from window coords, get board coords
// used to click on the board
func BoardCoords(x, y int) Vec2 {
	if x < 0 || y < 0 || x > GOBAN_SIZE || y > GOBAN_SIZE {return InvalidVec()}
	return Vec2{int8(x / STONE_SIZE), int8(y / STONE_SIZE)}
}

// from board coords, get window coords
// used to draw elements on the board
func WindowCoords(pos Vec2) (float64, float64) {
	return float64(pos.x) * STONE_SIZE, float64(pos.y) * STONE_SIZE
}

// tells if the cursor is hover the pass button
func PassButton(x, y int) bool {
	return y > GOBAN_SIZE && x < PASS_SIZE
}

// tells if the game is over: if both players have passed twice
func (win * Window) IsOver () bool {
	return win.blackPasses >= 2 && win.whitePasses >= 2
}

/* EBITEN */

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

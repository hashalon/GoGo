package main

//import "fmt"

// Stone stone placed in the board
type Stone struct {
	Vec2
	team bool
}

// Board board of the game
type Board struct {
	size uint8
	set  map[uint]Stone
}

// MakeBoard build a new empty board
func MakeBoard(size uint8) Board {
	return Board{size, make(map[uint]Stone)}
}

// FindStone find the stone at the given position
func (board Board) FindStone(pos Vec2) (Stone, bool) {
	return board.findID(pos.ID())
}
func (board Board) findID(id uint) (Stone, bool) {
	stone, found := board.set[id]
	return stone, found
}
func (board Board) remID(id uint) {
	delete(board.set, id)
}

// Valid return true if the position is in the board
func (board Board) Valid(pos Vec2) bool {
	return -1 < pos.x && pos.x < int8(board.size) && -1 < pos.y && pos.y < int8(board.size)
}

// Place place a new stone in the board if possible and capture, true if placed
func (board Board) Place(stone Stone) bool {
	id := stone.ID()

	// check that the position is not occupied
	_, occupied := board.set[id]
	if occupied { return false }

	// check that if we place the stone at that location,
	// we wouldn't close the last liberty of the group
	libs := GatherLiberties(board, stone)
	if len(libs.liberties) <= 0 { return false }

	// add the stone to the board
	board.set[id] = stone

	// for each stones around ours
	traject := Vec2{1, 0}
	for i := 0; i < 4; i++ {
		testPos := stone.Add(traject)
		traject = traject.Rot(1) // rotate
		id2 := testPos.ID()
		other, found := board.findID(id2)
		if !found { continue } // position is empty

		// check the liberties of the surrounding group
		libs2 := MakeLibSet()
		gatherLibs(board, other.Vec2, other.team, false, &libs2)
		if len(libs2.liberties) > 0 { continue } // the group will still have liberties
		for id3 := range libs2.checked { board.remID(id3) } // no more liberties -> capture
	}
	return true
}

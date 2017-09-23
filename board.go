package main

// Stone stone placed in the board
type Stone struct {
	Vec2
	team bool
}

// Board board of the game
type Board struct {
	size int8
	set  map[uint]Stone
}

// MakeBoard build a new empty board
func MakeBoard() Board {
	return Board{19, make(map[uint]Stone)}
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
	return -1 < pos.x && pos.x < board.size && -1 < pos.y && pos.y < board.size
}

// Place place a new stone in the board if possible and capture
func (board Board) Place(stone Stone) bool {
	id := stone.ID()
	_, occupied := board.set[id]
	if occupied {
		return false
	}
	libs := GatherLiberties(board, stone)
	if len(libs.liberties) <= 0 {
		return false
	}
	traject := Vec2{1, 0}
	for i := 0; i < 4; i++ {
		// test position to find a stone
		testPos := stone.Add(traject)
		id := testPos.ID()
		other, found := board.findID(id)
		if !found { // position is empty
			break
		}
		libs := GatherLiberties(board, other)
		if len(libs.liberties) > 0 { // the group will still have liberties
			break
		}
		for id := range libs.checked { // no more liberties -> capture
			board.remID(id)
		}
		traject = traject.Rot(1)
	}
	return true
}

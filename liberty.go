package main

//import "fmt"

// LibSet set of liberties for a group of stones
type LibSet struct {
	liberties map[uint]Vec2
	checked   map[uint]Vec2
}

// MakeLibSet create a new empty libset
func MakeLibSet() LibSet {
	return LibSet{make(map[uint]Vec2), make(map[uint]Vec2)}
}

// GatherLiberties return a set of liberties and a set of stones
func GatherLiberties(board Board, stone Stone) LibSet {
	libs := MakeLibSet()
	gatherLibs(board, stone.Vec2, stone.team, true, &libs)
	return libs
}
func gatherLibs(board Board, pos Vec2, team, checkNext bool, libs *LibSet) {
	id := pos.ID()
	libs.checked[id] = pos // position is checked
	// for each stone around ours
	traject := Vec2{1, 0}
	for i := 0; i < 4; i++ { // test neighbors
		testPos := pos.Add(traject)
		traject = traject.Rot(1)

		// position outside of the board
		if !board.Valid(testPos) { continue }
		id2 := testPos.ID()

		// see if the position has been tested yet
		if _, checked := libs.checked[id2]; !checked {
			// try to find a stone at the position
			other, found := board.findID(id2)

			if !found { // no stone
				libs.liberties[id2] = testPos // add liberty
			} else if team == other.team { // ally
				// find liberties of ally
				gatherLibs(board, testPos, team, checkNext, libs)
			} else { // enemy
				if !checkNext { continue } // already in prediction of next turn

				// find the liberties of the group after insertion
				nextLibs := MakeLibSet()
				gatherLibs(board, other.Vec2, other.team, false, &nextLibs)
				if len(nextLibs.liberties) > 1 { continue } // no problem

				// test if the remaining liberty is actually our stone
				if _, equal := nextLibs.liberties[id]; equal {
					libs.liberties[id] = pos
				}
			}
		}
	}
}

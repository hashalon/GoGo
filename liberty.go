package main

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
	gatherLibs(board, stone.Vec2, stone.team, true, libs)
	return libs
}
func gatherLibs(board Board, pos Vec2, team, checkNext bool, libs LibSet) {
	id := pos.ID()
	libs.checked[id] = pos // position is checked
	traject := Vec2{1, 0}
	for i := 0; i < 4; i++ { // test neighbors
		testPos := pos.Add(traject)
		id2 := testPos.ID()
		posChecked, checked := libs.checked[id2]
		var (
			other Stone
			found = false
		)
		if !checked { // the position hasn't been checked yet
			other, found = board.findID(id)
		}
		if found { // we found a stone at the tested position
			if team != other.team {
				if !checkNext {
					break
				}
				nextLibs := MakeLibSet()
				gatherLibs(board, other.Vec2, other.team, false, nextLibs)
				if len(nextLibs.liberties) > 1 {
					break
				}
				// test if the remaining liberty is actually our stone
				_, equal := nextLibs.liberties[id]
				if equal {
					libs.liberties[id] = pos
				}
			} else { // if it is a stone of our team, test it
				gatherLibs(board, posChecked, team, true, libs)
			}
		} else if board.Valid(testPos) { // add liberty if inside board
			libs.liberties[id2] = testPos
		}
		traject = traject.Rot(1)
	}
}

package main

// set of territories
type EndGame struct {
	// set of each territory identified by its start position
	set map[uint]Territory

	// set of all positions that have been checked
	checked map[uint]Vec2

	// the board to explore
	board * Board
}

// Territory: free cells surrounded by stones
type Territory struct {

	// identify the territory by its top left position
	start Vec2
	
	// count the number of free cells
	count uint

	// free cells of the territory
	set map[uint]Vec2

	// number of stones of each color surrounding the territory
	blackCount, whiteCount int8
}

// MakeLibSet create a new empty libset
func MakeTerritories(board * Board) EndGame {
	end := EndGame{make(map[uint]Territory), make(map[uint]Vec2), board}

	// for each cell of the goban, get territories
	for x := 0; x < int(board.size); x++ {
		for y := 0; y < int(board.size); y++ {
			pos := Vec2{int8(x), int8(y)}

			// skip position that have been already visited
			if _, checked := end.checked[pos.ID()]; checked {continue}

			terr := end.makeTerritory(pos)
			end.set[terr.start.ID()] = terr
		}
	}

	return end
}

// GatherLiberties return a set of liberties and a set of stones
func (end * EndGame) makeTerritory(pos Vec2) Territory {
	terr := Territory{Vec2{0, 0}, 0, make(map[uint]Vec2), 0, 0}
	end.gatherCells(pos, &terr)
	return terr
}

// gather cells recursively
func (end * EndGame) gatherCells(pos Vec2, terr * Territory) {
	// if the position has been already checked => stop
	id := pos.ID()
	if _, checked := end.checked[id]; checked {return}
	end.checked[id] = pos

	stone, present := end.board.FindStone(pos)

	if present { // stone present
		// increment the stone counter
		if stone.team {terr.whiteCount += 1} else {terr.blackCount += 1}

	} else { // no stone
		// add the position to the set
		terr.addPos(pos)

		// for each stone around ours
		traject := Vec2{1, 0}
		for i := 0; i < 4; i++ { // test neighbors
			testPos := pos.Add(traject)
			traject = traject.Rot(1)

			// position outside of the board
			if !end.board.Valid(testPos) { continue }
			
			// gather next position
			end.gatherCells(testPos, terr)
		}
	}
}

// add the position to the set and update the starting point
func (terr * Territory) addPos(pos Vec2) {
	terr.set[pos.ID()] = pos
	terr.count += 1

	up := false
	if pos.x < terr.start.x {up = true} else if pos.y < terr.start.y {up = true}
	
	// change the starting point if necessary
	if up {terr.start = pos}
}
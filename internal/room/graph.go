package room

import (
	"context"
	"errors"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

const errListingRooms string = "error listing rooms from database"

var ErrListingRooms error = errors.New(errListingRooms)

type Node struct {
	North       *Node
	Northeast   *Node
	East        *Node
	Southeast   *Node
	South       *Node
	Southwest   *Node
	West        *Node
	Northwest   *Node
	Title       string
	Description string
	ID          int64
}

var nilNode *Node = nil

func (n *Node) IsExitEmpty(dir string) bool {
	exitNode := n.Exit(dir)
	if exitNode == nil {
		return true
	}
	return exitNode.ID == 0
}

func (n *Node) ExitID(dir string) int64 {
	if n.IsExitEmpty(dir) {
		return 0
	}
	return n.Exit(dir).ID
}

func (n *Node) Exit(dir string) *Node {
	switch dir {
	case DirectionNorth:
		return n.North
	case DirectionNortheast:
		return n.Northeast
	case DirectionEast:
		return n.East
	case DirectionSoutheast:
		return n.Southeast
	case DirectionSouth:
		return n.South
	case DirectionSouthwest:
		return n.Southwest
	case DirectionWest:
		return n.West
	case DirectionNorthwest:
		return n.Northwest
	}
	return nil
}

func (n *Node) GetExitID(dir string) int64 {
	exitID := int64(0)
	if !n.IsExitEmpty(dir) {
		return n.Exit(dir).ID
	}
	return exitID
}

func (n *Node) Bind() fiber.Map {
	bind := fiber.Map{
		"ID":          n.ID,
		"Title":       n.Title,
		"Description": n.Description,
	}

	for _, dir := range DirectionsList {
		// TODO: Change this ID name to Key
		bindID := DirectionBindID(dir)
		bind[bindID] = n.BindGridExit(dir)
	}

	return bind
}

func (n *Node) BindGridExit(dir string) fiber.Map {
	b := n.BindEmptyExit(dir)
	b["ID"] = n.GetExitID(dir)
	return b
}

func MatrixCoordinateForDirection(dir string, row, col int) (int, int) {
	switch dir {
	case DirectionNorth:
		return row - 1, col
	case DirectionNortheast:
		return row - 1, col + 1
	case DirectionEast:
		return row, col + 1
	case DirectionSoutheast:
		return row + 1, col + 1
	case DirectionSouth:
		return row + 1, col
	case DirectionSouthwest:
		return row + 1, col - 1
	case DirectionWest:
		return row, col - 1
	case DirectionNorthwest:
		return row - 1, col - 1
	}
	return 0, 0
}

type BindMatrixParams struct {
	Matrix   [][]fiber.Map
	Priority []int64
	Row      int
	Col      int
	Shallow  bool
}

func (n *Node) BindMatrix(p BindMatrixParams) [][]fiber.Map {
	// TODO: Set up an algorithm for this to write *out* from the center node
	// This would predicate the grid being an odd-sized square: i.e. 5x5, 7x7, etc
	// Bind the center cell to the data of the root node

	priorityMap := PriorityMap(p.Priority)

	if !IsValidMatrixCoordinate(p.Matrix, p.Row, p.Col) {
		return p.Matrix
	}

	if !p.Shallow {
		for _, dir := range DirectionsList {
			row, col := MatrixCoordinateForDirection(dir, p.Row, p.Col)
			if n.IsExitEmpty(dir) {
				continue
			}
			p.Matrix = n.Exit(dir).BindMatrix(BindMatrixParams{
				Matrix:  p.Matrix,
				Row:     row,
				Col:     col,
				Shallow: true,
			})
		}

		for _, dir := range DirectionsList {
			row, col := MatrixCoordinateForDirection(dir, p.Row, p.Col)
			if n.IsExitEmpty(dir) {
				continue
			}
			p.Matrix = n.Exit(dir).BindMatrix(BindMatrixParams{
				Matrix:  p.Matrix,
				Row:     row,
				Col:     col,
				Shallow: false,
			})
		}
	}

	visitedID := p.Matrix[p.Row][p.Col]["ID"].(int64)
	if visitedID == int64(0) || Priority(priorityMap, n.ID, visitedID) {
		p.Matrix[p.Row][p.Col] = n.Bind()
	}

	return p.Matrix
}

func IsValidMatrixCoordinate(matrix [][]fiber.Map, row, col int) bool {
	if row < 0 {
		return false
	}

	if col < 0 {
		return false
	}

	if len(matrix) == 0 {
		return false
	}

	if row > len(matrix)-1 {
		return false
	}

	if col > len(matrix[0])-1 {
		return false
	}

	return true
}

func EmptyBindMatrix(size int) [][]fiber.Map {
	matrix := [][]fiber.Map{}
	for i := 0; i < size; i++ {
		row := []fiber.Map{}
		for j := 0; j < size; j++ {
			row = append(row, fiber.Map{"ID": int64(0)})
		}
		matrix = append(matrix, row)
	}
	return matrix
}

type BuildGraphParams struct {
	Queries  *query.Queries
	Room     *query.Room
	MaxDepth int
	Depth    int
}

func BuildGraph(p BuildGraphParams) (Node, error) {
	maxDepth := p.MaxDepth
	if maxDepth == 0 {
		maxDepth = constant.DefaultRoomGraphDepth
	}
	exitIDs := ExitIDs(p.Room)
	exitRooms, err := p.Queries.ListRoomsByIDs(context.Background(), exitIDs)
	if err != nil {
		return Node{ID: p.Room.ID}, ErrListingRooms
	}
	exitRoomByID := map[int64]query.Room{}
	for _, exitRoom := range exitRooms {
		exitRoomByID[exitRoom.ID] = exitRoom
	}
	nodesByDirection := EmptyNodesByDirection()
	for _, dir := range DirectionsList {
		exitID := ExitID(p.Room, dir)
		if exitID == 0 {
			continue
		} else {
			exitRoom, ok := exitRoomByID[ExitID(p.Room, dir)]
			if !ok {
				continue
			}
			if p.Depth >= maxDepth {
				node := TerminalNode(&exitRoom)
				nodesByDirection[dir] = node
			} else {
				node, err := BuildGraph(BuildGraphParams{
					Queries:  p.Queries,
					Room:     &exitRoom,
					MaxDepth: maxDepth,
					Depth:    p.Depth + 1,
				})
				if err != nil {
					return EmptyGraphNode(), err
				}
				nodesByDirection[dir] = node
			}
		}
	}

	northNode := nodesByDirection[DirectionNorth]
	northeastNode := nodesByDirection[DirectionNortheast]
	eastNode := nodesByDirection[DirectionEast]
	southeastNode := nodesByDirection[DirectionSoutheast]
	southNode := nodesByDirection[DirectionSouth]
	southwestNode := nodesByDirection[DirectionSouthwest]
	westNode := nodesByDirection[DirectionWest]
	northwestNode := nodesByDirection[DirectionNorthwest]
	return Node{
		ID:          p.Room.ID,
		Title:       p.Room.Title,
		Description: p.Room.Description,
		North:       &northNode,
		Northeast:   &northeastNode,
		East:        &eastNode,
		Southeast:   &southeastNode,
		South:       &southNode,
		Southwest:   &southwestNode,
		West:        &westNode,
		Northwest:   &northwestNode,
	}, nil
}

func (n *Node) IsExitTwoWay(en *Node, dir string) bool {
	if !IsDirectionValid(dir) {
		return false
	}

	roomExitID := n.ExitID(dir)
	if roomExitID == 0 {
		return false
	}

	opposite := DirectionOpposite(dir)
	if len(opposite) == 0 {
		return false
	}

	exitRoomExitID := en.ExitID(opposite)
	if exitRoomExitID == 0 {
		return false
	}

	return exitRoomExitID == n.ID
}

func (n *Node) BuildExitRooms() map[string]*Node {
	exitRooms := map[string]*Node{}
	for _, dir := range DirectionsList {
		if n.IsExitEmpty(dir) {
			emptyNode := EmptyGraphNode()
			exitRooms[dir] = &emptyNode
		} else {
			exitRooms[dir] = n.Exit(dir)
		}
	}
	return exitRooms
}

func (n *Node) BindExits() []fiber.Map {
	exits := []fiber.Map{}

	for _, dir := range DirectionsList {
		if n.IsExitEmpty(dir) {
			exits = append(exits, n.BindEmptyExit(dir))
		} else {
			exits = append(exits, n.BindExit(dir))
		}
	}

	return exits
}

func (n *Node) BindEmptyExit(dir string) fiber.Map {
	// TODO: Remove elements from the main map that aren't being used
	return fiber.Map{
		"ID":              0,
		"RoomID":          n.ID,
		"Exit":            dir,
		"ExitLetter":      DirectionLetter(dir),
		"ExitTitle":       DirectionTitle(dir),
		"EditElementID":   ExitEditElementID(dir),
		"SelectElementID": ExitSelectElementID(dir),
		"RoomsPath":       route.Rooms,
		"RoomExitsPath":   route.RoomExitsPath(n.ID),
		"RoomExitPath":    route.RoomExitPath(n.ID, dir),
		"CreateDialog": fiber.Map{
			"Exit":          dir,
			"RoomID":        n.ID,
			"RoomsPath":     route.Rooms,
			"EditElementID": ExitEditElementID(dir),
		},
		"LinkDialog": fiber.Map{
			"Exit":          dir,
			"RoomExitsPath": route.RoomExitsPath(n.ID),
			"EditElementID": ExitEditElementID(dir),
		},
	}
}

func (n *Node) BindExit(dir string) fiber.Map {
	exit := n.BindEmptyExit(dir)
	if n.IsExitEmpty(dir) {
		return exit
	}
	en := n.Exit(dir)
	exit["ID"] = en.ID
	exit["Title"] = en.Title
	exit["Description"] = en.Description
	exit["ExitPath"] = route.RoomPath(en.ID)
	exit["ExitEditPath"] = route.EditRoomPath(en.ID)
	exit["TwoWay"] = n.IsExitTwoWay(en, dir)
	return exit
}

func TerminalNode(room *query.Room) Node {
	northNode := EmptyGraphNode()
	northeastNode := EmptyGraphNode()
	eastNode := EmptyGraphNode()
	southeastNode := EmptyGraphNode()
	southNode := EmptyGraphNode()
	southwestNode := EmptyGraphNode()
	westNode := EmptyGraphNode()
	northwestNode := EmptyGraphNode()
	return Node{
		ID:          room.ID,
		Title:       room.Title,
		Description: room.Description,
		North:       &northNode,
		Northeast:   &northeastNode,
		East:        &eastNode,
		Southeast:   &southeastNode,
		South:       &southNode,
		Southwest:   &southwestNode,
		West:        &westNode,
		Northwest:   &northwestNode,
	}
}

func EmptyNodesByDirection() map[string]Node {
	nodesByDir := map[string]Node{}
	for _, dir := range DirectionsList {
		nodesByDir[dir] = EmptyGraphNode()
	}
	return nodesByDir
}

func EmptyGraphNode() Node {
	return Node{
		ID: 0,
	}
}

// TODO: Clean up and test
// This may be a good use case for a fully-qualified data structure that maintains
// independent pointers to Nodes in the graph by ID
func AnnotateMatrixExits(matrix [][]fiber.Map) [][]fiber.Map {
	matrixIDs := map[int64]bool{}
	for i := 0; i < len(matrix); i++ {
		row := matrix[i]
		for j := 0; j < len(row); j++ {
			room := row[j]
			matrixIDs[room["ID"].(int64)] = true
		}
	}

	for i := 0; i < len(matrix); i++ {
		row := matrix[i]
		for j := 0; j < len(row); j++ {
			room := row[j]

			for _, dir := range DirectionsList {

				bindID := DirectionBindID(dir)
				exit := room[bindID]
				if exit == nil {
					continue
				}
				exitBound := exit.(fiber.Map)
				exitID := exitBound["ID"].(int64)
				if exitID == int64(0) {
					continue
				}

				_, ok := matrixIDs[exitID]
				if ok {
					exitBound["InMatrix"] = true
				}
				row, col := MatrixCoordinateForDirection(dir, i, j)
				if IsValidMatrixCoordinate(matrix, row, col) {
					exitBound["Canonical"] = matrix[row][col]["ID"].(int64) == exitID
					exitBound["OffMatrix"] = false
				} else {
					exitBound["OffMatrix"] = true
				}
			}

		}
	}

	return matrix
}

func PriorityMap(priority []int64) map[int64]int {
	priorityMap := map[int64]int{}
	priorityWeight := 1
	for i := len(priority) - 1; i >= 0; i++ {
		priorityID := priority[i]
		priorityMap[priorityID] = priorityWeight
		priorityWeight++
	}
	return priorityMap
}

func Priority(priorityMap map[int64]int, priorityID, visitedID int64) bool {
	priorityWeight, ok := priorityMap[priorityID]
	if !ok {
		return false
	}
	visitedWeight, ok := priorityMap[visitedID]
	if !ok {
		return true
	}
	return priorityWeight > visitedWeight
}

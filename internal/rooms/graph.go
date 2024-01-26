package rooms

import (
	"context"
	"errors"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
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
	exitNode := n.GetExit(dir)
	if exitNode == nil {
		return true
	}

	return exitNode.ID == 0
}

func (n *Node) ExitID(dir string) int64 {
	if n.IsExitEmpty(dir) {
		return 0
	}
	return n.GetExit(dir).ID
}

func (n *Node) GetExit(dir string) *Node {
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
		return n.GetExit(dir).ID
	}
	return exitID
}

func (n *Node) Bind() fiber.Map {
	return fiber.Map{
		"ID":          n.ID,
		"Title":       n.Title,
		"Description": n.Description,
		"North":       n.GetExitID(DirectionNorth),
		"Northeast":   n.GetExitID(DirectionNortheast),
		"East":        n.GetExitID(DirectionEast),
		"Southeast":   n.GetExitID(DirectionSoutheast),
		"South":       n.GetExitID(DirectionSouth),
		"Southwest":   n.GetExitID(DirectionSouthwest),
		"West":        n.GetExitID(DirectionWest),
		"Northwest":   n.GetExitID(DirectionNorthwest),
	}
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
	Matrix  [][]fiber.Map
	Row     int
	Col     int
	Shallow bool
}

func (n *Node) BindMatrix(p BindMatrixParams) [][]fiber.Map {
	// TODO: Set up an algorithm for this to write *out* from the center node
	// This would predicate the grid being an odd-sized square: i.e. 5x5, 7x7, etc
	// Bind the center cell to the data of the root node

	if !IsValidMatrixCoordinate(p.Matrix, p.Row, p.Col) {
		return p.Matrix
	}

	if p.Matrix[p.Row][p.Col]["ID"] != int64(0) {
		return p.Matrix
	}

	p.Matrix[p.Row][p.Col] = n.Bind()

	if p.Shallow {
		return p.Matrix
	}

	for _, dir := range DirectionsList {
		row, col := MatrixCoordinateForDirection(dir, p.Row, p.Col)
		if n.IsExitEmpty(dir) {
			continue
		}
		p.Matrix = n.GetExit(DirectionNorth).BindMatrix(BindMatrixParams{
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
		p.Matrix = n.GetExit(DirectionNorth).BindMatrix(BindMatrixParams{
			Matrix:  p.Matrix,
			Row:     row,
			Col:     col,
			Shallow: false,
		})
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

func EmptyBindMatrix() [][]fiber.Map {
	matrix := [][]fiber.Map{}
	for i := 0; i < 5; i++ {
		row := []fiber.Map{}
		for j := 0; j < 5; j++ {
			row = append(row, fiber.Map{"ID": int64(0)})
		}
		matrix = append(matrix, row)
	}
	return matrix
}

type BuildGraphParams struct {
	Queries  *queries.Queries
	Room     *queries.Room
	MaxDepth int
	Depth    int
}

func BuildGraph(p BuildGraphParams) (Node, error) {
	exitIDs := ExitIDs(p.Room)
	exitRooms, err := p.Queries.ListRoomsByIDs(context.Background(), exitIDs)
	if err != nil {
		return Node{ID: p.Room.ID}, ErrListingRooms
	}
	exitRoomByID := map[int64]queries.Room{}
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
			if p.Depth >= p.MaxDepth {
				node := TerminalNode(&exitRoom)
				nodesByDirection[dir] = node
			} else {
				node, err := BuildGraph(BuildGraphParams{
					Queries:  p.Queries,
					Room:     &exitRoom,
					MaxDepth: p.MaxDepth,
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
			exitRooms[dir] = n.GetExit(dir)
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
			exits = append(exits, n.BindExit(n.GetExit(dir), dir))
		}
	}

	return exits
}

func (n *Node) BindEmptyExit(dir string) fiber.Map {
	return fiber.Map{
		"ID":              0,
		"RoomID":          n.ID,
		"Exit":            dir,
		"ExitLetter":      DirectionLetter(dir),
		"ExitTitle":       DirectionTitle(dir),
		"EditElementID":   ExitEditElementID(dir),
		"SelectElementID": ExitSelectElementID(dir),
		"RoomsPath":       routes.Rooms,
		"RoomExitsPath":   routes.RoomExitsPath(n.ID),
		"RoomExitPath":    routes.RoomExitPath(n.ID, dir),
	}
}

func (n *Node) BindExit(en *Node, dir string) fiber.Map {
	exit := n.BindEmptyExit(dir)
	exit["ID"] = en.ID
	exit["Title"] = en.Title
	exit["Description"] = en.Description
	exit["ExitPath"] = routes.RoomPath(en.ID)
	exit["ExitEditPath"] = routes.EditRoomPath(en.ID)
	exit["TwoWay"] = n.IsExitTwoWay(en, dir)
	return exit
}

func TerminalNode(room *queries.Room) Node {
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

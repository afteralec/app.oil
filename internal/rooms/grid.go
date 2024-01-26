package rooms

import (
	"context"
	"errors"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
)

const errListingRooms string = "error listing rooms from database"

var ErrListingRooms error = errors.New(errListingRooms)

const (
	GridRowOneElementID   string = "room-grid-row-one"
	GridRowTwoElementID   string = "room-grid-row-two"
	GridRowThreeElementID string = "room-grid-row-three"
	GridRowFourElementID  string = "room-grid-row-four"
	GridRowFiveElementID  string = "room-grid-row-five"
)

var EmptyRoomExitMap fiber.Map = fiber.Map{
	DirectionNorth:     fiber.Map{"ID": int64(0)},
	DirectionNortheast: fiber.Map{"ID": int64(0)},
	DirectionEast:      fiber.Map{"ID": int64(0)},
	DirectionSoutheast: fiber.Map{"ID": int64(0)},
	DirectionSouth:     fiber.Map{"ID": int64(0)},
	DirectionSouthwest: fiber.Map{"ID": int64(0)},
	DirectionWest:      fiber.Map{"ID": int64(0)},
	DirectionNorthwest: fiber.Map{"ID": int64(0)},
}

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

func (n *Node) BindMatrix(matrix [][]fiber.Map, row, col int) [][]fiber.Map {
	// TODO: Set up an algorithm for this to write *out* from the center node
	// This would predicate the grid being an odd-sized square: i.e. 5x5, 7x7, etc
	// Bind the center cell to the data of the root node

	if !IsValidMatrixCoordinate(matrix, row, col) {
		return matrix
	}

	if matrix[row][col]["ID"] != int64(0) {
		return matrix
	}

	matrix[row][col] = n.Bind()

	if !n.IsExitEmpty(DirectionNorth) {
		matrix = n.GetExit(DirectionNorth).BindMatrix(matrix, row-1, col)
	}
	if !n.IsExitEmpty(DirectionNortheast) {
		matrix = n.GetExit(DirectionNortheast).BindMatrix(matrix, row-1, col+1)
	}
	if !n.IsExitEmpty(DirectionEast) {
		matrix = n.GetExit(DirectionEast).BindMatrix(matrix, row, col+1)
	}
	if !n.IsExitEmpty(DirectionSoutheast) {
		matrix = n.GetExit(DirectionSoutheast).BindMatrix(matrix, row+1, col+1)
	}
	if !n.IsExitEmpty(DirectionSouth) {
		matrix = n.GetExit(DirectionSouth).BindMatrix(matrix, row+1, col)
	}
	if !n.IsExitEmpty(DirectionSouthwest) {
		matrix = n.GetExit(DirectionSouthwest).BindMatrix(matrix, row+1, col-1)
	}
	if !n.IsExitEmpty(DirectionWest) {
		matrix = n.GetExit(DirectionWest).BindMatrix(matrix, row, col-1)
	}
	if !n.IsExitEmpty(DirectionNorthwest) {
		matrix = n.GetExit(DirectionNorthwest).BindMatrix(matrix, row-1, col-1)
	}

	return matrix
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
	nodesByDirection := BuildEmptyNodesByDirection()
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
				node := BuildTerminalNode(&exitRoom)
				nodesByDirection[dir] = node
			} else {
				node, err := BuildGraph(BuildGraphParams{
					Queries:  p.Queries,
					Room:     &exitRoom,
					MaxDepth: p.MaxDepth,
					Depth:    p.Depth + 1,
				})
				if err != nil {
					return BuildEmptyGraphNode(), err
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

func BuildTerminalNode(room *queries.Room) Node {
	northNode := BuildEmptyGraphNode()
	northeastNode := BuildEmptyGraphNode()
	eastNode := BuildEmptyGraphNode()
	southeastNode := BuildEmptyGraphNode()
	southNode := BuildEmptyGraphNode()
	southwestNode := BuildEmptyGraphNode()
	westNode := BuildEmptyGraphNode()
	northwestNode := BuildEmptyGraphNode()
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

func BuildEmptyNodesByDirection() map[string]Node {
	nodesByDir := map[string]Node{}
	for _, dir := range DirectionsList {
		nodesByDir[dir] = BuildEmptyGraphNode()
	}
	return nodesByDir
}

func BuildEmptyGraphNode() Node {
	return Node{
		ID: 0,
	}
}

type GridExitMapParams struct {
	Queries  *queries.Queries
	Room     *queries.Room
	Depth    int
	MaxDepth int
}

func GridRoomEmpty() fiber.Map {
	room := fiber.Map{
		"ID": int64(0),
	}
	for _, dir := range DirectionsList {
		room[dir] = fiber.Map{"ID": int64(0)}
	}
	return room
}

func GridRoom(room *queries.Room) fiber.Map {
	return fiber.Map{
		"ID":               room.ID,
		"Title":            room.Title,
		"Description":      room.Description,
		DirectionNorth:     fiber.Map{"ID": room.North},
		DirectionNortheast: fiber.Map{"ID": room.Northeast},
		DirectionEast:      fiber.Map{"ID": room.East},
		DirectionSoutheast: fiber.Map{"ID": room.Southeast},
		DirectionSouth:     fiber.Map{"ID": room.South},
		DirectionSouthwest: fiber.Map{"ID": room.Southwest},
		DirectionWest:      fiber.Map{"ID": room.West},
		DirectionNorthwest: fiber.Map{"ID": room.Northwest},
	}
}

// TODO: Update this to use the ListRoomsByIDs query
// Also, I think this could be reimagined as a graph of nodes
func GridExitMap(p GridExitMapParams) fiber.Map {
	rooms := fiber.Map{
		"ID":          p.Room.ID,
		"Title":       p.Room.Title,
		"Description": p.Room.Description,
	}
	for _, dir := range DirectionsList {
		exitID := ExitID(p.Room, dir)
		if exitID > 0 {
			room, err := p.Queries.GetRoom(context.Background(), exitID)
			if err != nil {
				return GridRoomEmpty()
			}
			if p.Depth >= p.MaxDepth {
				rooms[dir] = GridRoom(&room)
			} else {
				rooms[dir] = GridExitMap(GridExitMapParams{
					Queries:  p.Queries,
					Room:     &room,
					Depth:    p.Depth + 1,
					MaxDepth: p.MaxDepth,
				})
			}
		} else {
			rooms[dir] = GridRoomEmpty()
		}
	}
	return rooms
}

func FilterGridRoom(p map[string]fiber.Map) fiber.Map {
	for dir, room := range p {
		if room[dir].(fiber.Map)["ID"].(int64) > 0 {
			return room[dir].(fiber.Map)
		}
	}

	return GridRoomEmpty()
}

func Grid(exitMap fiber.Map) []fiber.Map {
	northMap := exitMap[DirectionNorth].(fiber.Map)
	northeastMap := exitMap[DirectionNortheast].(fiber.Map)
	eastMap := exitMap[DirectionEast].(fiber.Map)
	southeastMap := exitMap[DirectionSoutheast].(fiber.Map)
	southMap := exitMap[DirectionSouth].(fiber.Map)
	southwestMap := exitMap[DirectionSouthwest].(fiber.Map)
	westMap := exitMap[DirectionWest].(fiber.Map)
	northwestMap := exitMap[DirectionNorthwest].(fiber.Map)

	rowOne := fiber.Map{
		"ElementID": GridRowOneElementID,
		"Rooms": []fiber.Map{
			northwestMap[DirectionNorthwest].(fiber.Map),
			FilterGridRoom(map[string]fiber.Map{
				DirectionNorth:     northwestMap,
				DirectionNorthwest: northMap,
			}),
			FilterGridRoom(map[string]fiber.Map{
				DirectionNortheast: northwestMap,
				DirectionNorth:     northMap,
				DirectionNorthwest: northeastMap,
			}),
			FilterGridRoom(map[string]fiber.Map{
				DirectionNortheast: northMap,
				DirectionNorth:     northeastMap,
			}),
			northeastMap[DirectionNortheast].(fiber.Map),
		},
	}

	rowTwo := fiber.Map{
		"ElementID": GridRowTwoElementID,
		"Rooms": []fiber.Map{
			FilterGridRoom(map[string]fiber.Map{
				DirectionWest:      northwestMap,
				DirectionNorthwest: westMap,
			}),
			northwestMap,
			northMap,
			northeastMap,
			FilterGridRoom(map[string]fiber.Map{
				DirectionEast:      northeastMap,
				DirectionNortheast: eastMap,
			}),
		},
	}

	rowThree := fiber.Map{
		"ElementID": GridRowThreeElementID,
		"Rooms": []fiber.Map{
			FilterGridRoom(map[string]fiber.Map{
				DirectionSouthwest: northwestMap,
				DirectionWest:      westMap,
				DirectionNorthwest: southwestMap,
			}),
			westMap,
			{"Self": true},
			eastMap,
			FilterGridRoom(map[string]fiber.Map{
				DirectionSoutheast: northeastMap,
				DirectionEast:      eastMap,
				DirectionNortheast: southeastMap,
			}),
		},
	}

	rowFour := fiber.Map{
		"ElementID": GridRowFourElementID,
		"Rooms": []fiber.Map{
			FilterGridRoom(map[string]fiber.Map{
				DirectionSouthwest: westMap,
				DirectionWest:      southwestMap,
			}),
			southwestMap,
			southMap,
			southeastMap,
			FilterGridRoom(map[string]fiber.Map{
				DirectionSoutheast: eastMap,
				DirectionEast:      southeastMap,
			}),
		},
	}

	rowFive := fiber.Map{
		"ElementID": GridRowFiveElementID,
		"Rooms": []fiber.Map{
			southwestMap[DirectionSouthwest].(fiber.Map),
			FilterGridRoom(map[string]fiber.Map{
				DirectionSouth:     southwestMap,
				DirectionSouthwest: southMap,
			}),
			FilterGridRoom(map[string]fiber.Map{
				DirectionSouth:     southMap,
				DirectionSoutheast: southwestMap,
				DirectionSouthwest: southeastMap,
			}),
			FilterGridRoom(map[string]fiber.Map{
				DirectionSouth:     southeastMap,
				DirectionSoutheast: southMap,
			}),
			southeastMap[DirectionSoutheast].(fiber.Map),
		},
	}

	return []fiber.Map{
		rowOne,
		rowTwo,
		rowThree,
		rowFour,
		rowFive,
	}
}

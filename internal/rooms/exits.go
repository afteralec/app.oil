package rooms

import (
	"context"
	"errors"
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

var ErrExitIDNotFound error = errors.New("no exit found for that RID")

func LoadExitRooms(q *queries.Queries, room *queries.Room) (map[string]queries.Room, error) {
	exitRooms := make(map[string]queries.Room)
	exitRoomDirections := make(map[int64]string)

	exitIDs := []int64{}
	for _, dir := range DirectionsList {
		exitID := ExitID(room, dir)
		if exitID == 0 {
			continue
		}
		exitIDs = append(exitIDs, exitID)
		exitRoomDirections[exitID] = dir
	}

	records, err := q.ListRoomsByIDs(context.Background(), exitIDs)
	if err != nil {
		return exitRooms, err
	}

	for _, record := range records {
		dir, ok := exitRoomDirections[record.ID]
		if !ok {
			// TODO: This should be a fatal error
			continue
		}
		exitRooms[dir] = record
	}

	return exitRooms, nil
}

func (n *Node) BuildExitRooms() map[string]*Node {
	exitRooms := map[string]*Node{}
	for _, dir := range DirectionsList {
		if n.IsExitEmpty(dir) {
			emptyNode := BuildEmptyGraphNode()
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

func ExitEditElementID(dir string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "edit-room-exits-edit-%s", dir)
	return sb.String()
}

func ExitSelectElementID(dir string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "edit-room-exits-select-%s", dir)
	return sb.String()
}

func ExitIDs(room *queries.Room) []int64 {
	ids := []int64{}
	for _, dir := range DirectionsList {
		ids = append(ids, ExitID(room, dir))
	}
	return ids
}

func ExitID(room *queries.Room, dir string) int64 {
	switch dir {
	case DirectionNorth:
		return room.North
	case DirectionNortheast:
		return room.Northeast
	case DirectionEast:
		return room.East
	case DirectionSoutheast:
		return room.Southeast
	case DirectionSouth:
		return room.South
	case DirectionSouthwest:
		return room.Southwest
	case DirectionWest:
		return room.West
	case DirectionNorthwest:
		return room.Northwest
	default:
		return 0
	}
}

func ExitDirection(room *queries.Room, id int64) (string, error) {
	switch id {
	case room.North:
		return DirectionNorth, nil
	case room.Northeast:
		return DirectionNortheast, nil
	case room.East:
		return DirectionEast, nil
	case room.Southeast:
		return DirectionSoutheast, nil
	case room.South:
		return DirectionSouth, nil
	case room.Southwest:
		return DirectionSouthwest, nil
	case room.West:
		return DirectionWest, nil
	case room.Northwest:
		return DirectionNorthwest, nil
	default:
		return "", ErrExitIDNotFound
	}
}

func IsExitTwoWay(room *queries.Room, exitRoom *queries.Room, dir string) bool {
	if !IsDirectionValid(dir) {
		return false
	}

	roomExitID := ExitID(room, dir)
	if roomExitID == 0 {
		return false
	}

	opposite := DirectionOpposite(dir)
	if len(opposite) == 0 {
		return false
	}

	exitRoomExitID := ExitID(exitRoom, opposite)
	if exitRoomExitID == 0 {
		return false
	}

	return exitRoomExitID == room.ID
}

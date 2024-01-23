package rooms

import (
	"context"
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

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

func BuildExits(room *queries.Room, exitRooms map[string]queries.Room) []fiber.Map {
	exits := []fiber.Map{}

	for _, dir := range DirectionsList {
		exitRoom, ok := exitRooms[dir]
		if !ok {
			exits = append(exits, BuildEmptyExit(room, dir))
			continue
		}

		exit := BuildExit(room, &exitRoom, dir)

		exits = append(exits, exit)
	}

	return exits
}

func BuildEmptyExit(room *queries.Room, dir string) fiber.Map {
	return fiber.Map{
		"ID":              0,
		"RoomID":          room.ID,
		"Exit":            dir,
		"ExitLetter":      DirectionLetter(dir),
		"ExitTitle":       DirectionTitle(dir),
		"EditElementID":   ExitEditElementID(dir),
		"SelectElementID": ExitSelectElementID(dir),
		"RoomsPath":       routes.Rooms,
	}
}

func BuildExit(room *queries.Room, exitRoom *queries.Room, dir string) fiber.Map {
	exit := BuildEmptyExit(room, dir)
	exit["ID"] = exitRoom.ID
	exit["Title"] = exitRoom.Title
	exit["Description"] = exitRoom.Description
	exit["ExitPath"] = routes.RoomPath(exitRoom.ID)
	exit["ExitEditPath"] = routes.EditRoomPath(exitRoom.ID)
	exit["TwoWay"] = IsExitTwoWay(room, exitRoom, dir)
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

package rooms

import (
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

func BuildExits(room *queries.Room) []fiber.Map {
	exits := []fiber.Map{}

	for _, dir := range DirectionsList {
		exits = append(exits, BuildExit(room, dir))
	}

	return exits
}

func BuildExit(room *queries.Room, dir string) fiber.Map {
	id := ExitID(room, dir)
	exit := fiber.Map{
		"ID":         id,
		"RoomID":     room.ID,
		"Exit":       dir,
		"ExitLetter": DirectionLetter(dir),
		"ExitTitle":  DirectionTitle(dir),
		"ElementID":  ExitElementID(dir),
	}

	// TODO: Figure out getting room exit summaries in here
	if id > 0 {
		exit["Title"] = DefaultTitle
		exit["Description"] = DefaultDescription
		exit["RoomPath"] = routes.RoomPath(id)
	}

	return exit
}

func ExitElementID(dir string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "edit-room-exits-edit-%s", dir)
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

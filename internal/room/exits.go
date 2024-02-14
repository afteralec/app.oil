package room

import (
	"errors"
	"fmt"
	"strings"

	"petrichormud.com/app/internal/query"
)

var ErrExitIDNotFound error = errors.New("no exit found for that RID")

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

func ExitIDs(room *query.Room) []int64 {
	ids := []int64{}
	for _, dir := range DirectionsList {
		ids = append(ids, ExitID(room, dir))
	}
	return ids
}

func ExitID(room *query.Room, dir string) int64 {
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

func ExitDirection(room *query.Room, id int64) (string, error) {
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

func IsExitTwoWay(room *query.Room, exitRoom *query.Room, dir string) bool {
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

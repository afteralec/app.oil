package routes

import (
	"fmt"
	"strings"
)

const (
	Rooms        string = "/rooms"
	RoomImages   string = "/rooms/images"
	NewRoomImage string = "/rooms/images/new"
)

func RoomImagePath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d", RoomImages, id)
	return sb.String()
}

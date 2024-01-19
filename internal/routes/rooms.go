package routes

import (
	"fmt"
	"strings"
)

const (
	RoomImages             string = "/rooms/images"
	RoomImagePathParam     string = "/rooms/images/:id"
	NewRoomImage           string = "/rooms/images/new"
	EditRoomImagePathParam string = "/rooms/images/:id/edit"
	Rooms                  string = "/rooms"
)

func RoomImagePath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d", RoomImages, id)
	return sb.String()
}

func EditRoomImagePath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/edit", RoomImages, id)
	return sb.String()
}

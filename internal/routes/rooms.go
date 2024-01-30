package routes

import (
	"fmt"
	"strings"
)

const (
	Rooms                    string = "/rooms"
	RoomPathParam            string = "/rooms/:id"
	NewRoom                  string = "/rooms/new"
	EditRoomPathParam        string = "/rooms/:id/edit"
	RoomGridPathParam        string = "/rooms/:id/grid/:selected"
	RoomExitPathParam        string = "/rooms/:id/:exit"
	RoomExitsPathParam       string = "/rooms/:id/exits"
	RoomTitlePathParam       string = "/rooms/:id/title"
	RoomDescriptionPathParam string = "/rooms/:id/description"
)

func RoomPath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d", Rooms, id)
	return sb.String()
}

func EditRoomPath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/edit", Rooms, id)
	return sb.String()
}

func RoomExitPath(id int64, exit string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/%s", Rooms, id, exit)
	return sb.String()
}

func RoomExitsPath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/exits", Rooms, id)
	return sb.String()
}

func RoomTitlePath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/title", Rooms, id)
	return sb.String()
}

func RoomDescriptionPath(id int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s/%d/description", Rooms, id)
	return sb.String()
}

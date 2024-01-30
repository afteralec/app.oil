package rooms

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
)

func BindSizeRadioGroup(b fiber.Map, room *queries.Room) fiber.Map {
	b["SizeRadioGroup"] = []fiber.Map{
		{
			"ID":       "edit-room-image-size-tiny",
			"Name":     "size",
			"Variable": "size",
			"Value":    "0",
			"Active":   room.Size == 0,
			"Label":    "Tiny",
		},
		{
			"ID":       "edit-room-image-size-small",
			"Name":     "size",
			"Variable": "size",
			"Value":    "1",
			"Active":   room.Size == 1,
			"Label":    "Small",
		},
		{
			"ID":       "edit-room-image-size-medium",
			"Name":     "size",
			"Variable": "size",
			"Value":    "2",
			"Active":   room.Size == 2,
			"Label":    "Medium",
		},
		{
			"ID":       "edit-room-image-size-large",
			"Name":     "size",
			"Variable": "size",
			"Value":    "3",
			"Active":   room.Size == 3,
			"Label":    "Large",
		},
		{
			"ID":       "edit-room-image-size-huge",
			"Name":     "size",
			"Variable": "size",
			"Value":    "4",
			"Active":   room.Size == 4,
			"Label":    "Huge",
		},
	}
	return b
}

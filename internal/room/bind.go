package room

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/query"
)

func BindSizeRadioGroup(b fiber.Map, room *query.Room) fiber.Map {
	b["SizeRadioGroup"] = []bind.Radio{
		{
			ID:       "edit-room-image-size-tiny",
			Name:     "size",
			Variable: "size",
			Value:    "0",
			Label:    "Tiny",
			Active:   room.Size == 0,
		},
		{
			ID:       "edit-room-image-size-small",
			Name:     "size",
			Variable: "size",
			Value:    "1",
			Label:    "Small",
			Active:   room.Size == 1,
		},
		{
			ID:       "edit-room-image-size-medium",
			Name:     "size",
			Variable: "size",
			Value:    "2",
			Label:    "Medium",
			Active:   room.Size == 2,
		},
		{
			ID:       "edit-room-image-size-large",
			Name:     "size",
			Variable: "size",
			Value:    "3",
			Label:    "Large",
			Active:   room.Size == 3,
		},
		{
			ID:       "edit-room-image-size-huge",
			Name:     "size",
			Variable: "size",
			Value:    "4",
			Label:    "Huge",
			Active:   room.Size == 4,
		},
	}
	return b
}

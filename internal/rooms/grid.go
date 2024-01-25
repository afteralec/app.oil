package rooms

import (
	"context"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
)

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

func NewGridFromExitMap(exitMap fiber.Map) []fiber.Map {
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
			{"ID": exitMap["ID"].(int64)},
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

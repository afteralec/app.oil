package room

import "petrichormud.com/app/internal/query"

var Default query.CreateRoomParams = query.CreateRoomParams{
	Title:       "A dark ocean expanse",
	Description: "Dark, roiling waters stretch to the horizon in every direction, torn to white where the wind rakes the tops of the waves. The ocean itself is a deep green that turns black in the troughs, promising blindness to any unfortunate enough to go under. Stinging salt hangs heavy in every breath of the air.",
	Size:        3,
}

const (
	DefaultTitle       string = "A dark ocean expanse"
	DefaultDescription string = "Dark, roiling waters stretch to the horizon in every direction, torn to white where the wind rakes the tops of the waves. The ocean itself is a deep green that turns black in the troughs, promising blindness to any unfortunate enough to go under. Stinging salt hangs heavy in every breath of the air."
	DefaultSize        int32  = 2
)

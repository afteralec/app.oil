package rooms

const (
	DirectionNorth     string = "north"
	DirectionNortheast string = "northeast"
	DirectionEast      string = "east"
	DirectionSoutheast string = "southeast"
	DirectionSouth     string = "south"
	DirectionSouthwest string = "southwest"
	DirectionWest      string = "west"
	DirectionNorthwest string = "northwest"
)

const (
	DirectionNorthOpposite     string = DirectionSouth
	DirectionNortheastOpposite string = DirectionSouthwest
	DirectionEastOpposite      string = DirectionWest
	DirectionSoutheastOpposite string = DirectionNorthwest
	DirectionSouthOpposite     string = DirectionNorth
	DirectionSouthwestOpposite string = DirectionNortheast
	DirectionWestOpposite      string = DirectionEast
	DirectionNorthwestOpposite string = DirectionSoutheast
)

func IsDirectionValid(dir string) bool {
	switch dir {
	case DirectionNorth:
		return true
	case DirectionNortheast:
		return true
	case DirectionEast:
		return true
	case DirectionSoutheast:
		return true
	case DirectionSouth:
		return true
	case DirectionSouthwest:
		return true
	case DirectionWest:
		return true
	case DirectionNorthwest:
		return true
	}
	return false
}

func GetDirectionOpposite(dir string) string {
	switch dir {
	case DirectionNorth:
		return DirectionNorthOpposite
	case DirectionNortheast:
		return DirectionNortheastOpposite
	case DirectionEast:
		return DirectionEastOpposite
	case DirectionSoutheast:
		return DirectionSoutheastOpposite
	case DirectionSouth:
		return DirectionSouthOpposite
	case DirectionSouthwest:
		return DirectionSouthwestOpposite
	case DirectionWest:
		return DirectionWestOpposite
	case DirectionNorthwest:
		return DirectionNorthwestOpposite
	}

	return ""
}

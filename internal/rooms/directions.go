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
	DirectionLetterNorth     string = "n"
	DirectionLetterNortheast string = "ne"
	DirectionLetterEast      string = "e"
	DirectionLetterSoutheast string = "se"
	DirectionLetterSouth     string = "s"
	DirectionLetterSouthwest string = "sw"
	DirectionLetterWest      string = "w"
	DirectionLetterNorthwest string = "nw"
)

const (
	DirectionTitleNorth     string = "North"
	DirectionTitleNortheast string = "Northeast"
	DirectionTitleEast      string = "East"
	DirectionTitleSoutheast string = "Southeast"
	DirectionTitleSouth     string = "South"
	DirectionTitleSouthwest string = "Southwest"
	DirectionTitleWest      string = "West"
	DirectionTitleNorthwest string = "Northwest"
)

var DirectionsList []string = []string{
	DirectionNorth,
	DirectionNortheast,
	DirectionEast,
	DirectionSoutheast,
	DirectionSouth,
	DirectionSouthwest,
	DirectionWest,
	DirectionNorthwest,
}

var Directions map[string]bool = map[string]bool{
	DirectionNorth:     true,
	DirectionNortheast: true,
	DirectionEast:      true,
	DirectionSoutheast: true,
	DirectionSouth:     true,
	DirectionSouthwest: true,
	DirectionWest:      true,
	DirectionNorthwest: true,
}

var DirectionOpposites map[string]string = map[string]string{
	DirectionNorth:     DirectionSouth,
	DirectionNortheast: DirectionSouthwest,
	DirectionEast:      DirectionWest,
	DirectionSoutheast: DirectionNorthwest,
	DirectionSouth:     DirectionNorth,
	DirectionSouthwest: DirectionNortheast,
	DirectionWest:      DirectionEast,
	DirectionNorthwest: DirectionSoutheast,
}

var DirectionLetters map[string]string = map[string]string{
	DirectionNorth:     DirectionLetterNorth,
	DirectionNortheast: DirectionLetterNortheast,
	DirectionEast:      DirectionLetterEast,
	DirectionSoutheast: DirectionLetterSoutheast,
	DirectionSouth:     DirectionLetterSouth,
	DirectionSouthwest: DirectionLetterSouthwest,
	DirectionWest:      DirectionLetterWest,
	DirectionNorthwest: DirectionLetterNorthwest,
}

var DirectionTitles map[string]string = map[string]string{
	DirectionNorth:     DirectionTitleNorth,
	DirectionNortheast: DirectionTitleNortheast,
	DirectionEast:      DirectionTitleEast,
	DirectionSoutheast: DirectionTitleSoutheast,
	DirectionSouth:     DirectionTitleSouth,
	DirectionSouthwest: DirectionTitleSouthwest,
	DirectionWest:      DirectionTitleWest,
	DirectionNorthwest: DirectionTitleNorthwest,
}

func IsDirectionValid(dir string) bool {
	_, ok := Directions[dir]
	return ok
}

func DirectionOpposite(dir string) string {
	opposite, ok := DirectionOpposites[dir]
	if !ok {
		return ""
	}
	return opposite
}

func DirectionLetter(dir string) string {
	letter, ok := DirectionLetters[dir]
	if !ok {
		return ""
	}
	return letter
}

func DirectionTitle(dir string) string {
	title, ok := DirectionTitles[dir]
	if !ok {
		return ""
	}
	return title
}

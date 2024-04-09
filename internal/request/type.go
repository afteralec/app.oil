package request

import "slices"

const TypeCharacterApplication string = "CharacterApplication"

var Types []string = []string{
	TypeCharacterApplication,
}

func IsTypeValid(t string) bool {
	return slices.Contains(Types, t)
}

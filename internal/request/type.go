package request

const (
	TypeCharacterApplication string = "CharacterApplication"
)

func IsTypeValid(t string) bool {
	return t == TypeCharacterApplication
}

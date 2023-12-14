package password

const (
	MinLength = 8
	MaxLength = 255
)

func IsValid(pw string) bool {
	if len(pw) < MinLength {
		return false
	}
	if len(pw) > MaxLength {
		return false
	}
	return true
}

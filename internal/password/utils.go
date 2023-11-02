package password

func Validate(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	if len(pw) > 255 {
		return false
	}
	return true
}

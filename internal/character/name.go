package character

import (
	"regexp"
)

func SanitizeName(n string) string {
	re := regexp.MustCompile("[^a-zA-Z'-]+")
	s := re.ReplaceAllString(n, "")
	return s
}

func IsValidName(n string) bool {
	if len(n) < 4 {
		return false
	}

	if len(n) > 16 {
		return false
	}

	return true
}

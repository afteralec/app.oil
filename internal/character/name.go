package character

import (
	"regexp"

	"petrichormud.com/app/internal/shared"
)

func SanitizeName(n string) string {
	// TODO: Add this regex test to IsNameValid too
	re := regexp.MustCompile("[^a-zA-Z'-]+")
	s := re.ReplaceAllString(n, "")
	return s
}

func IsNameValid(n string) bool {
	if len(n) < shared.MinCharacterNameLength {
		return false
	}

	if len(n) > shared.MaxCharacterNameLength {
		return false
	}

	return true
}

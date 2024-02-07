package character

import (
	"regexp"

	"petrichormud.com/app/internal/constants"
)

func SanitizeName(n string) string {
	// TODO: Add this regex test to IsNameValid too
	re := regexp.MustCompile("[^a-zA-Z'-]+")
	s := re.ReplaceAllString(n, "")
	return s
}

func IsNameValid(n string) bool {
	if len(n) < constants.MinCharacterNameLength {
		return false
	}

	if len(n) > constants.MaxCharacterNameLength {
		return false
	}

	return true
}

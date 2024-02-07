package character

import (
	"regexp"

	"petrichormud.com/app/internal/constants"
)

func SanitizeBackstory(backstory string) string {
	// TODO: Add this regex test to IsBackstoryValid too
	re := regexp.MustCompile("[^\r\na-zA-Z, -.!()]+")
	s := re.ReplaceAllString(backstory, "")
	return s
}

func IsBackstoryValid(backstory string) bool {
	if len(backstory) < constants.MinCharacterBackstoryLength {
		return false
	}

	if len(backstory) > constants.MaxCharacterBackstoryLength {
		return false
	}

	return true
}

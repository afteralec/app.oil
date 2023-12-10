package character

import (
	"regexp"

	"petrichormud.com/app/internal/shared"
)

func SanitizeBackstory(backstory string) string {
	// TODO: Add this regex test to IsBackstoryValid too
	re := regexp.MustCompile("[^\r\na-zA-Z, -.!()]+")
	s := re.ReplaceAllString(backstory, "")
	return s
}

func IsBackstoryValid(backstory string) bool {
	if len(backstory) < shared.MinCharacterBackstoryLength {
		return false
	}

	if len(backstory) > shared.MaxCharacterBackstoryLength {
		return false
	}

	return true
}

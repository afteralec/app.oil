package character

import (
	"regexp"

	"petrichormud.com/app/internal/shared"
)

func SanitizeShortDescription(sdesc string) string {
	// TODO: Add this regex test to IsShortDescriptionValid too
	re := regexp.MustCompile("[^a-zA-Z, -]+")
	s := re.ReplaceAllString(sdesc, "")
	return s
}

func IsShortDescriptionValid(sdesc string) bool {
	if len(sdesc) < shared.MinCharacterShortDescriptionLength {
		return false
	}

	if len(sdesc) > shared.MaxCharacterShortDescriptionLength {
		return false
	}

	return true
}

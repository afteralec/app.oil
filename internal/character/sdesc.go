package character

import (
	"regexp"

	"petrichormud.com/app/internal/constants"
)

func SanitizeShortDescription(sdesc string) string {
	// TODO: Add this regex test to IsShortDescriptionValid too
	re := regexp.MustCompile("[^a-zA-Z, -]+")
	s := re.ReplaceAllString(sdesc, "")
	return s
}

func IsShortDescriptionValid(sdesc string) bool {
	if len(sdesc) < constants.MinCharacterShortDescriptionLength {
		return false
	}

	if len(sdesc) > constants.MaxCharacterShortDescriptionLength {
		return false
	}

	return true
}

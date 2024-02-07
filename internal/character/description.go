package character

import (
	"regexp"

	"petrichormud.com/app/internal/constants"
)

func SanitizeDescription(description string) string {
	// TODO: Add this regex test to IsDescriptionValid too
	re := regexp.MustCompile("[^a-zA-Z, -.!()]+")
	s := re.ReplaceAllString(description, "")
	return s
}

func IsDescriptionValid(description string) bool {
	if len(description) < constants.MinCharacterDescriptionLength {
		return false
	}

	if len(description) > constants.MaxCharacterDescriptionLength {
		return false
	}

	return true
}

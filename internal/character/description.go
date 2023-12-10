package character

import (
	"regexp"

	"petrichormud.com/app/internal/shared"
)

func SanitizeDescription(description string) string {
	// TODO: Add this regex test to IsDescriptionValid too
	re := regexp.MustCompile("[^a-zA-Z, -.!()]+")
	s := re.ReplaceAllString(description, "")
	return s
}

func IsDescriptionValid(description string) bool {
	if len(description) < shared.MinCharacterDescriptionLength {
		return false
	}

	if len(description) > shared.MaxCharacterDescriptionLength {
		return false
	}

	return true
}

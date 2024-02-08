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

const (
	GenderMale      = "Male"
	GenderFemale    = "Female"
	GenderNonBinary = "NonBinary"
)

func SanitizeGender(str string) string {
	if str != GenderMale && str != GenderFemale && str != GenderNonBinary {
		return GenderNonBinary
	}
	return str
}

func IsGenderValid(gender string) bool {
	if gender == GenderNonBinary {
		return true
	}

	if gender == GenderFemale {
		return true
	}

	if gender == GenderMale {
		return true
	}

	return false
}

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

package character

import (
	"regexp"
)

func SanitizeName(n string) string {
	// TODO: Add this regex test to IsNameValid too
	re := regexp.MustCompile("[^a-zA-Z'-]+")
	s := re.ReplaceAllString(n, "")
	return s
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

func SanitizeShortDescription(sdesc string) string {
	// TODO: Add this regex test to IsShortDescriptionValid too
	re := regexp.MustCompile("[^a-zA-Z, -]+")
	s := re.ReplaceAllString(sdesc, "")
	return s
}

func SanitizeDescription(description string) string {
	// TODO: Add this regex test to IsDescriptionValid too
	re := regexp.MustCompile("[^a-zA-Z, -.!()]+")
	s := re.ReplaceAllString(description, "")
	return s
}

func SanitizeBackstory(backstory string) string {
	// TODO: Add this regex test to IsBackstoryValid too
	re := regexp.MustCompile("[^\r\na-zA-Z, -.!()]+")
	s := re.ReplaceAllString(backstory, "")
	return s
}

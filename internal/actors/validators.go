package actors

import "regexp"

const (
	MinimumImageNameLength        int = 4
	MaximumImageNameLength        int = 50
	MinimumShortDescriptionLength int = 8
	MaximumShortDescriptionLength int = 300
	MinimumDescriptionLength      int = 32
	MaximumDescriptionLength      int = 2000
)

func IsImageNameValid(name string) bool {
	if len(name) < MinimumImageNameLength {
		return false
	}

	if len(name) > MaximumImageNameLength {
		return false
	}

	re := regexp.MustCompile("[^a-z-]+")
	return !re.MatchString(name)
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

func IsShortDescriptionValid(sdesc string) bool {
	if len(sdesc) < MinimumShortDescriptionLength {
		return false
	}
	if len(sdesc) > MaximumShortDescriptionLength {
		return false
	}

	re := regexp.MustCompile("[^a-zA-Z, -]+")
	return !re.MatchString(sdesc)
}

func IsDescriptionValid(desc string) bool {
	if len(desc) < MinimumDescriptionLength {
		return false
	}
	if len(desc) > MaximumDescriptionLength {
		return false
	}

	re := regexp.MustCompile("[^a-zA-Z, '-.!()]+")
	return !re.MatchString(desc)
}

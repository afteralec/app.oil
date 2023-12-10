package character

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

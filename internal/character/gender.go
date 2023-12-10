package character

const (
	GenderMale      = "Male"
	GenderFemale    = "Female"
	GenderNonBinary = "NonBinary"
)

func ValidateGender(str string) string {
	if str != GenderMale && str != GenderFemale && str != GenderNonBinary {
		return GenderNonBinary
	}
	return str
}

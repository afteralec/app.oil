package requests

import (
	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/shared"
)

func IsCharacterApplicationValid(app *queries.CharacterApplicationContent) bool {
	if !IsNameValid(app.Name) {
		return false
	}

	if !IsGenderValid(app.Gender) {
		return false
	}

	if !IsShortDescriptionValid(app.ShortDescription) {
		return false
	}

	if !IsDescriptionValid(app.Description) {
		return false
	}

	if !IsBackstoryValid(app.Backstory) {
		return false
	}

	return true
}

func IsNameValid(n string) bool {
	if len(n) < shared.MinCharacterNameLength {
		return false
	}

	if len(n) > shared.MaxCharacterNameLength {
		return false
	}

	return true
}

func IsGenderValid(gender string) bool {
	if gender == constants.GenderNonBinary {
		return true
	}

	if gender == constants.GenderFemale {
		return true
	}

	if gender == constants.GenderMale {
		return true
	}

	return false
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

func IsDescriptionValid(description string) bool {
	if len(description) < shared.MinCharacterDescriptionLength {
		return false
	}

	if len(description) > shared.MaxCharacterDescriptionLength {
		return false
	}

	return true
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

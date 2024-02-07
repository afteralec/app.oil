package requests

import (
	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/queries"
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
	if len(n) < constants.MinCharacterNameLength {
		return false
	}

	if len(n) > constants.MaxCharacterNameLength {
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
	if len(sdesc) < constants.MinCharacterShortDescriptionLength {
		return false
	}

	if len(sdesc) > constants.MaxCharacterShortDescriptionLength {
		return false
	}

	return true
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

func IsBackstoryValid(backstory string) bool {
	if len(backstory) < constants.MinCharacterBackstoryLength {
		return false
	}

	if len(backstory) > constants.MaxCharacterBackstoryLength {
		return false
	}

	return true
}

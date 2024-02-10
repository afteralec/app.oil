package requests

import (
	"petrichormud.com/app/internal/constant"
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
	if len(n) < constant.MinCharacterNameLength {
		return false
	}

	if len(n) > constant.MaxCharacterNameLength {
		return false
	}

	return true
}

func IsGenderValid(gender string) bool {
	if gender == constant.GenderNonBinary {
		return true
	}

	if gender == constant.GenderFemale {
		return true
	}

	if gender == constant.GenderMale {
		return true
	}

	return false
}

func IsShortDescriptionValid(sdesc string) bool {
	if len(sdesc) < constant.MinCharacterShortDescriptionLength {
		return false
	}

	if len(sdesc) > constant.MaxCharacterShortDescriptionLength {
		return false
	}

	return true
}

func IsDescriptionValid(description string) bool {
	if len(description) < constant.MinCharacterDescriptionLength {
		return false
	}

	if len(description) > constant.MaxCharacterDescriptionLength {
		return false
	}

	return true
}

func IsBackstoryValid(backstory string) bool {
	if len(backstory) < constant.MinCharacterBackstoryLength {
		return false
	}

	if len(backstory) > constant.MaxCharacterBackstoryLength {
		return false
	}

	return true
}

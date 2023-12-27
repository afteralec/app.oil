package request

import (
	"petrichormud.com/app/internal/queries"
)

// TODO: Create a cleaner way to do this
func CharacterApplicationGetNextIncompleteField(app *queries.CharacterApplicationContent) (string, error) {
	if len(app.Name) == 0 {
		return FieldName, nil
	}

	if len(app.Gender) == 0 {
		return FieldGender, nil
	}

	if len(app.ShortDescription) == 0 {
		return FieldShortDescription, nil
	}

	if len(app.Description) == 0 {
		return FieldDescription, nil
	}

	if len(app.Backstory) == 0 {
		return FieldBackstory, nil
	}

	return "", ErrNoIncompleteFields
}

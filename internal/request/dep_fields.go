package request

import (
	"errors"

	"petrichormud.com/app/internal/view"
)

// Request fields
const FieldStatus string = "status"

// Content fields

// Character Application fields
const (
	FieldName             string = "name"
	FieldGender           string = "gender"
	FieldShortDescription string = "sdesc"
	FieldDescription      string = "desc"
	FieldBackstory        string = "backstory"
)

// Errors
var ErrNoIncompleteFields error = errors.New("no incomplete fields")

var ViewsByFieldAndType map[string]map[string]string = map[string]map[string]string{
	TypeCharacterApplication: {
		FieldName:             view.CharacterApplicationName,
		FieldGender:           view.CharacterApplicationGender,
		FieldShortDescription: view.CharacterApplicationShortDescription,
		FieldDescription:      view.CharacterApplicationDescription,
		FieldBackstory:        view.CharacterApplicationBackstory,
	},
}

func IsFieldNameValid(t, name string) bool {
	definition, ok := Definitions.Get(t)
	if !ok {
		return false
	}

	return definition.IsFieldNameValid(name)
}

var (
	ErrMalformedUpdateInput error = errors.New("no field matched in input")
	ErrInvalidInput         error = errors.New("field value didn't pass validation")
)

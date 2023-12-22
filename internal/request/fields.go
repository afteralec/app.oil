package request

import "petrichormud.com/app/internal/views"

// Request fields
const FieldStatus string = "status"

// Content fields

// Character Application fields
const (
	FieldName             string = "name"
	FieldGender           string = "gender"
	FieldShortDescription string = "sdesc"
	FieldDescription      string = "description"
	FieldBackstory        string = "backstory"
)

var FieldsByType map[string]map[string]bool = map[string]map[string]bool{
	TypeCharacterApplication: {
		FieldName:             true,
		FieldGender:           true,
		FieldShortDescription: true,
		FieldDescription:      true,
		FieldBackstory:        true,
	},
}

var ViewsByFieldAndType map[string]map[string]string = map[string]map[string]string{
	TypeCharacterApplication: {
		FieldName:             views.CharacterApplicationName,
		FieldGender:           views.CharacterApplicationGender,
		FieldShortDescription: views.CharacterApplicationShortDescription,
		FieldDescription:      views.CharacterApplicationDescription,
		FieldBackstory:        views.CharacterApplicationBackstory,
	},
}

func IsFieldValid(t, field string) bool {
	fieldsByType, ok := FieldsByType[t]
	if !ok {
		return false
	}

	_, ok = fieldsByType[field]
	return ok
}

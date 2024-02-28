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

func IsFieldValid(t, field string) bool {
	fieldsByType, ok := FieldMapsByType[t]
	if !ok {
		return false
	}

	_, ok = fieldsByType[field]
	return ok
}

var (
	ErrMalformedUpdateInput error = errors.New("no field matched in input")
	ErrInvalidInput         error = errors.New("field value didn't pass validation")
)

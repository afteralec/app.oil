package request

import (
	"context"
	"errors"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/views"
)

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

type UpdateInput struct {
	// Request fields
	Status string `form:"status"`
	// Character Application fields
	Name             string `form:"name"`
	Gender           string `form:"gender"`
	ShortDescription string `form:"sdesc"`
	Description      string `form:"desc"`
	Backstory        string `form:"backstory"`
}

func (in *UpdateInput) GetField() (string, error) {
	if len(in.Status) > 0 {
		return FieldStatus, nil
	}

	if len(in.Name) > 0 {
		return FieldName, nil
	}

	if len(in.Gender) > 0 {
		return FieldGender, nil
	}

	if len(in.ShortDescription) > 0 {
		return FieldShortDescription, nil
	}

	if len(in.Description) > 0 {
		return FieldDescription, nil
	}

	if len(in.Backstory) > 0 {
		return FieldBackstory, nil
	}

	return "", errors.New("malformed input")
}

var (
	ErrMalformedUpdateInput error = errors.New("no field matched in input")
	ErrInvalidInput         error = errors.New("field value didn't pass validation")
)

// TODO: Turn this into a map of updaters by field - can create an interface for the Updater
// TODO: Get Validators in front of everything here
func (in *UpdateInput) UpdateField(q *queries.Queries, req *queries.Request, field string) error {
	switch field {
	case FieldName:
		if !IsNameValid(field) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentName(context.Background(), queries.UpdateCharacterApplicationContentNameParams{
			RID:  req.ID,
			Name: in.Name,
		}); err != nil {
			return err
		}
	case FieldGender:
		if !IsGenderValid(field) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentGender(context.Background(), queries.UpdateCharacterApplicationContentGenderParams{
			RID:    req.ID,
			Gender: in.Gender,
		}); err != nil {
			return err
		}
	case FieldShortDescription:
		if !IsShortDescriptionValid(field) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentShortDescription(context.Background(), queries.UpdateCharacterApplicationContentShortDescriptionParams{
			RID:              req.ID,
			ShortDescription: in.ShortDescription,
		}); err != nil {
			return err
		}
	case FieldDescription:
		if !IsDescriptionValid(field) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentDescription(context.Background(), queries.UpdateCharacterApplicationContentDescriptionParams{
			RID:         req.ID,
			Description: in.Description,
		}); err != nil {
			return err
		}
	case FieldBackstory:
		if !IsBackstoryValid(field) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentBackstory(context.Background(), queries.UpdateCharacterApplicationContentBackstoryParams{
			RID:       req.ID,
			Backstory: in.Backstory,
		}); err != nil {
			return err
		}
	default:
		return ErrMalformedUpdateInput
	}

	return nil
}

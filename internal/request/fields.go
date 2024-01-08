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
	FieldDescription      string = "desc"
	FieldBackstory        string = "backstory"
)

// Errors
var ErrNoIncompleteFields error = errors.New("no incomplete fields")

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
	fieldsByType, ok := FieldMapsByType[t]
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
func (in *UpdateInput) UpdateField(pid int64, q *queries.Queries, req *queries.Request, field string) error {
	switch field {
	case FieldName:
		if !IsNameValid(in.Name) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentName(context.Background(), queries.UpdateCharacterApplicationContentNameParams{
			RID:  req.ID,
			Name: in.Name,
		}); err != nil {
			return err
		}

	case FieldGender:
		if !IsGenderValid(in.Gender) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentGender(context.Background(), queries.UpdateCharacterApplicationContentGenderParams{
			RID:    req.ID,
			Gender: in.Gender,
		}); err != nil {
			return err
		}
	case FieldShortDescription:
		if !IsShortDescriptionValid(in.ShortDescription) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentShortDescription(context.Background(), queries.UpdateCharacterApplicationContentShortDescriptionParams{
			RID:              req.ID,
			ShortDescription: in.ShortDescription,
		}); err != nil {
			return err
		}
	case FieldDescription:
		if !IsDescriptionValid(in.Description) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentDescription(context.Background(), queries.UpdateCharacterApplicationContentDescriptionParams{
			RID:         req.ID,
			Description: in.Description,
		}); err != nil {
			return err
		}
	case FieldBackstory:
		if !IsBackstoryValid(in.Backstory) {
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

	if req.Type == TypeCharacterApplication {
		app, err := q.GetCharacterApplicationContentForRequest(context.Background(), req.ID)
		if err != nil {
			return err
		}

		ready := IsCharacterApplicationValid(&app)

		if ready && req.Status == StatusIncomplete {
			if err := q.CreateHistoryForRequestStatusChange(context.Background(), queries.CreateHistoryForRequestStatusChangeParams{
				RID: req.ID,
				PID: pid,
			}); err != nil {
				return err
			}

			if err := q.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
				ID:     req.ID,
				Status: StatusReady,
			}); err != nil {
				return err
			}
		} else if !ready && req.Status == StatusReady {
			if err := q.CreateHistoryForRequestStatusChange(context.Background(), queries.CreateHistoryForRequestStatusChangeParams{
				RID: req.ID,
				PID: pid,
			}); err != nil {
				return err
			}

			if err := q.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
				ID:     req.ID,
				Status: StatusIncomplete,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

package request

import (
	"context"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/definition"
	"petrichormud.com/app/internal/request/field"
)

var FieldsByType map[string]field.Group = map[string]field.Group{
	TypeCharacterApplication: definition.CharacterApplicationFields,
}

type UpdateFieldParams struct {
	Request *query.Request
	Field   *query.RequestField
	Value   string
	PID     int64
}

func UpdateField(q *query.Queries, p UpdateFieldParams) error {
	if !IsTypeValid(p.Request.Type) {
		return ErrInvalidType
	}

	fg, ok := FieldsByType[p.Request.Type]
	if !ok {
		return ErrNoDefinition
	}

	if err := q.UpdateRequestFieldValue(context.Background(), query.UpdateRequestFieldValueParams{
		ID:    p.Field.ID,
		Value: p.Value,
	}); err != nil {
		return err
	}
	if err := q.UpdateRequestFieldStatus(context.Background(), query.UpdateRequestFieldStatusParams{
		ID:     p.Field.ID,
		Status: FieldStatusNotReviewed,
	}); err != nil {
		return err
	}

	// TODO: Split this whole thing out into a tested unit?
	fields, err := q.ListRequestFieldsForRequest(context.Background(), p.Request.ID)
	if err != nil {
		return err
	}
	ready := true
	for _, field := range fields {
		fd, ok := fg.Get(field.Type)
		if !ok {
			// TODO: This means there's a field on a request that doesn't have a definition
			return ErrNoDefinition
		}
		if !fd.IsValid(field.Value) {
			ready = false
		}
	}
	if ready && p.Request.Status == StatusIncomplete {
		if err := UpdateStatus(q, UpdateStatusParams{
			PID:    p.PID,
			RID:    p.Request.ID,
			Status: StatusReady,
		}); err != nil {
			return err
		}
	} else if !ready && p.Request.Status == StatusReady {
		if err := UpdateStatus(q, UpdateStatusParams{
			PID:    p.PID,
			RID:    p.Request.ID,
			Status: StatusIncomplete,
		}); err != nil {
			return err
		}
	}

	return nil
}

type FieldsForOverviewParams = field.ForOverviewParams

func FieldsForOverview(p field.ForOverviewParams) ([]field.ForOverview, error) {
	fields, ok := FieldsByType[p.Request.Type]
	if !ok {
		return []field.ForOverview{}, ErrNoDefinition
	}
	return fields.ForOverview(p), nil
}

func FieldMap(fields []query.RequestField) field.Map {
	m := field.Map{}
	for _, field := range fields {
		m[field.Type] = &field
	}
	return m
}

func NextIncompleteField(t string, fieldmap field.Map) (field.NextIncompleteOutput, error) {
	fields, ok := FieldsByType[t]
	if !ok {
		return field.NextIncompleteOutput{}, ErrNoDefinition
	}
	return fields.NextIncomplete(fieldmap)
}

func NextUnreviewedField(t string, fieldmap field.Map) (field.NextUnreviewedOutput, error) {
	fields, ok := FieldsByType[t]
	if !ok {
		return field.NextUnreviewedOutput{}, ErrNoDefinition
	}
	return fields.NextUnreviewed(fieldmap)
}

func IsFieldTypeValid(t, ft string) bool {
	fields, ok := FieldsByType[t]
	if !ok {
		return false
	}
	_, ok = fields.Map()[ft]
	return ok
}

package request

import (
	"context"

	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/definition"
	"petrichormud.com/app/internal/request/field"
)

var FieldsByType map[string]field.Group = map[string]field.Group{
	TypeCharacterApplication: definition.FieldsCharacterApplication,
}

func GetFieldDefinition(t, ft string) (field.Field, error) {
	fg, ok := FieldsByType[t]
	if !ok {
		return field.Field{}, ErrNoDefinition
	}

	fd, ok := fg.Get(ft)
	if !ok {
		// TODO: Make this an Invalid Field Type error instead
		return field.Field{}, ErrInvalidType
	}

	return fd, nil
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

	if p.Request.Status == StatusIncomplete {
		fields, err := q.ListRequestFieldsForRequest(context.Background(), p.Request.ID)
		if err != nil {
			return err
		}
		rfids := []int64{}
		for _, field := range fields {
			rfids = append(rfids, field.ID)
		}
		subfields, err := q.ListRequestSubfieldsForFields(context.Background(), rfids)
		if err != nil {
			return err
		}
		ready, err := AreFieldsReady(AreFieldsReadyParams{
			FieldGroup: fg,
			Fields:     fields,
			SubFields:  subfields,
			PlayerOnly: true,
		})
		if err != nil {
			return err
		}

		if ready {
			if err := UpdateStatus(q, UpdateStatusParams{
				PID:    p.PID,
				RID:    p.Request.ID,
				Status: StatusReady,
			}); err != nil {
				return err
			}
		}
	}

	if p.Request.Status == StatusReady {
		fields, err := q.ListRequestFieldsForRequest(context.Background(), p.Request.ID)
		if err != nil {
			return err
		}
		rfids := []int64{}
		for _, field := range fields {
			rfids = append(rfids, field.ID)
		}
		subfields, err := q.ListRequestSubfieldsForFields(context.Background(), rfids)
		if err != nil {
			return err
		}
		ready, err := AreFieldsReady(AreFieldsReadyParams{
			FieldGroup: fg,
			Fields:     fields,
			SubFields:  subfields,
			PlayerOnly: true,
		})
		if err != nil {
			return err
		}

		if ready {
			if err := UpdateStatus(q, UpdateStatusParams{
				PID:    p.PID,
				RID:    p.Request.ID,
				Status: StatusIncomplete,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

type AreFieldsReadyParams struct {
	Fields     []query.RequestField
	SubFields  []query.RequestSubfield
	FieldGroup field.Group
	PlayerOnly bool
}

func AreFieldsReady(p AreFieldsReadyParams) (bool, error) {
	sfmap := SubfieldMap(p.SubFields)

	ready := true
	for _, field := range p.Fields {
		fd, ok := p.FieldGroup.Get(field.Type)
		if !ok {
			// TODO: This means there's a field on a request that doesn't have a definition
			return false, ErrNoDefinition
		}

		if p.PlayerOnly && !fd.ForPlayer() {
			continue
		}

		if fd.SubfieldConfig.Require {
			subfields, ok := sfmap[field.ID]
			if ok {
				if len(subfields) < fd.SubfieldConfig.MinValues || len(subfields) > fd.SubfieldConfig.MaxValues {
					ready = false
				}

				for _, subfield := range subfields {
					if !fd.IsValid(subfield.Value) {
						ready = false
					}
				}
			} else {
				// TODO: If there are any Fields that allow 0+ Subfields, this will trigger on 0
				ready = false
			}
			continue
		}

		if !fd.IsValid(field.Value) {
			ready = false
		}
	}

	return ready, nil
}

type FieldsForOverviewParams = field.ForOverviewParams

func FieldsForOverview(e *html.Engine, p field.ForOverviewParams) ([]field.ForOverview, error) {
	fields, ok := FieldsByType[p.Request.Type]
	if !ok {
		return []field.ForOverview{}, ErrNoDefinition
	}
	return fields.ForOverview(e, p), nil
}

func FieldMap(fields []query.RequestField) field.Map {
	m := field.Map{}
	for _, field := range fields {
		m[field.Type] = field
	}
	return m
}

// TODO: Make this SubfieldMap a discrete type?
func SubfieldMap(subfields []query.RequestSubfield) map[int64][]query.RequestSubfield {
	m := map[int64][]query.RequestSubfield{}
	for _, subfield := range subfields {
		_, ok := m[subfield.RFID]
		if !ok {
			m[subfield.RFID] = []query.RequestSubfield{}
		}
		m[subfield.RFID] = append(m[subfield.RFID], subfield)
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

func FieldSubfieldConfig(t, ft string) (field.SubfieldConfig, error) {
	fd, err := GetFieldDefinition(t, ft)
	if err != nil {
		// TODO: Trace this error
		return field.SubfieldConfig{}, ErrNoDefinition
	}
	return fd.SubfieldConfig, nil
}

// TODO: Deprecate this in favor of FieldSubfieldConfig()
func FieldRequiresSubfields(t, ft string) bool {
	fd, err := GetFieldDefinition(t, ft)
	if err != nil {
		// TODO: Trace this error
		return false
	}
	return fd.SubfieldConfig.Require
}
